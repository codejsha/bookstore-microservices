from __future__ import annotations

import os
from pathlib import Path
from typing import Any

import structlog
import yaml
from pydantic.fields import FieldInfo
from pydantic_settings import PydanticBaseSettingsSource

logger = structlog.get_logger()

CONFIG_FILE_ENV = "APP_CONFIG_FILE"
DEFAULT_CONFIG_FILE = "/app/config/application.yaml"


class ConfigFileInvalidError(ValueError): ...


def config_file_path() -> Path:
    return Path(os.environ.get(CONFIG_FILE_ENV, "").strip() or DEFAULT_CONFIG_FILE)


def load_config_file(path: Path) -> dict[str, Any]:
    try:
        content = path.read_text(encoding="utf-8")
    except FileNotFoundError:
        logger.info("config file not found, using env and defaults", path=str(path))
        return {}

    data = yaml.safe_load(content)
    if data is None:
        return {}
    if not isinstance(data, dict):
        raise ConfigFileInvalidError(f"config file {path} must contain a mapping at the top level")
    return _expand_env(data)


def _expand_env(value: Any) -> Any:
    if isinstance(value, str):
        return os.path.expandvars(value)
    if isinstance(value, dict):
        return {key: _expand_env(item) for key, item in value.items()}
    if isinstance(value, list):
        return [_expand_env(item) for item in value]
    return value


class YamlFileSettingsSource(PydanticBaseSettingsSource):
    def __init__(self, settings_cls: type, path: Path | None = None) -> None:
        super().__init__(settings_cls)
        self._data: dict[str, Any] = load_config_file(path or config_file_path())

    def get_field_value(self, field: FieldInfo, field_name: str) -> tuple[Any, str, bool]:
        value = self._data.get(field_name)
        return value, field_name, isinstance(value, (dict, list))

    def __call__(self) -> dict[str, Any]:
        return self._data
