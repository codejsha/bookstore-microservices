from __future__ import annotations

import json
import os
import time
import urllib.request
from typing import Any

import structlog
from pydantic.fields import FieldInfo
from pydantic_settings import PydanticBaseSettingsSource

logger = structlog.get_logger()


class CloudConfigUnavailableError(RuntimeError): ...


_CONFIGSERVER_PREFIX = "configserver:"
_REQUEST_TIMEOUT_SECONDS = 5.0
_MAX_ATTEMPTS = 3
_RETRY_BACKOFF_SECONDS = 1.0


def fetch_cloud_config(app_name: str) -> dict[str, Any]:
    server = os.environ.get("APP_CONFIG_SERVER", "").strip()
    if not server:
        logger.debug("cloud config disabled: APP_CONFIG_SERVER unset", app=app_name)
        return {}

    server = _strip_configserver_prefix(server)
    profile = os.environ.get("APP_CONFIG_PROFILE", "").strip() or "default"
    label = os.environ.get("APP_CONFIG_LABEL", "").strip()
    url = _build_url(server, app_name, profile, label)

    payload = _fetch_with_retry(url, app_name)

    sources = payload.get("propertySources") or []
    merged = _merge_property_sources(sources)
    return _expand_source_map(merged)


def _strip_configserver_prefix(server: str) -> str:
    if server.startswith(_CONFIGSERVER_PREFIX):
        return server[len(_CONFIGSERVER_PREFIX) :]
    return server


def _build_url(base: str, app_name: str, profile: str, label: str) -> str:
    base = base.rstrip("/")
    segments = [app_name.lower(), profile.lower()]
    if label:
        segments.append(label.lower())
    return base + "/" + "/".join(segments)


def _fetch_with_retry(url: str, app_name: str) -> dict[str, Any]:
    last_error: Exception | None = None
    for attempt in range(1, _MAX_ATTEMPTS + 1):
        try:
            with urllib.request.urlopen(url, timeout=_REQUEST_TIMEOUT_SECONDS) as response:
                body = response.read()
            return json.loads(body)
        except Exception as exc:  # noqa: BLE001
            last_error = exc
            logger.debug(
                "cloud config fetch attempt failed",
                app=app_name,
                url=url,
                attempt=attempt,
                error=str(exc),
            )
            if attempt < _MAX_ATTEMPTS:
                time.sleep(_RETRY_BACKOFF_SECONDS * attempt)

    logger.error(
        "cloud config unreachable, refusing to start on built-in defaults",
        app=app_name,
        url=url,
        attempts=_MAX_ATTEMPTS,
        error=str(last_error),
    )
    raise CloudConfigUnavailableError(f"config server unreachable after {_MAX_ATTEMPTS} attempts: {url}: {last_error}")


def _merge_property_sources(sources: list[dict[str, Any]]) -> dict[str, Any]:
    merged: dict[str, Any] = {}
    for entry in sources:
        source = entry.get("source") or {}
        for key, value in source.items():
            if key not in merged:
                merged[key] = value
    return merged


def _expand_source_map(flat: dict[str, Any]) -> dict[str, Any]:
    root: dict[str, Any] = {}
    for full_key, value in flat.items():
        tokens = _parse_key(full_key)
        if not tokens:
            continue
        root = _set_value(root, tokens, value)
    return root


def _parse_key(full_key: str) -> list[tuple[str, Any]]:
    tokens: list[tuple[str, Any]] = []
    for seg in full_key.split("."):
        name, indices = _parse_segment(seg)
        if name:
            tokens.append(("key", name))
        for idx in indices:
            tokens.append(("index", idx))
    return tokens


def _parse_segment(seg: str) -> tuple[str, list[int]]:
    open_idx = seg.find("[")
    if open_idx < 0:
        return seg, []

    name = seg[:open_idx]
    rest = seg[open_idx:]
    indices: list[int] = []
    while rest:
        if not rest.startswith("["):
            return seg, []
        close = rest.find("]")
        if close < 0:
            return seg, []
        inner = rest[1:close]
        if not inner.isdigit():
            return seg, []
        indices.append(int(inner))
        rest = rest[close + 1 :]
    return name, indices


def _set_value(container: Any, tokens: list[tuple[str, Any]], value: Any) -> Any:
    if not tokens:
        return value

    kind, key = tokens[0]
    if kind == "index":
        if not isinstance(container, list):
            container = []
        while len(container) <= key:
            container.append(None)
        container[key] = _set_value(container[key], tokens[1:], value)
        return container

    if not isinstance(container, dict):
        container = {}
    container[key] = _set_value(container.get(key), tokens[1:], value)
    return container


class CloudConfigSettingsSource(PydanticBaseSettingsSource):
    def __init__(self, settings_cls: type, app_name: str) -> None:
        super().__init__(settings_cls)
        self._app_name = app_name
        self._data: dict[str, Any] = fetch_cloud_config(app_name)

    def get_field_value(self, field: FieldInfo, field_name: str) -> tuple[Any, str, bool]:
        value = self._data.get(field_name)
        return value, field_name, isinstance(value, (dict, list))

    def __call__(self) -> dict[str, Any]:
        return self._data
