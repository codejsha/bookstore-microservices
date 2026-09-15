import asyncio
import json
import time
import urllib.parse
import urllib.request
from collections.abc import Awaitable, Callable, Mapping
from dataclasses import dataclass

import structlog

from internal.config.config import TemporalAuthConfig

_logger = structlog.get_logger()

_AUTHORIZATION_HEADER = "authorization"
_BEARER_PREFIX = "Bearer "
_MIN_RETRY_BACKOFF_SECONDS = 1.0
_MAX_RETRY_BACKOFF_SECONDS = 30.0


class TemporalTokenError(RuntimeError): ...


@dataclass(frozen=True)
class AccessToken:
    value: str
    expires_in: float


TokenFetcher = Callable[[], Awaitable[AccessToken]]
MetadataListener = Callable[[Mapping[str, str]], None]


class TemporalTokenProvider:
    def __init__(
        self,
        fetch: TokenFetcher,
        *,
        refresh_ratio: float = 0.75,
        min_backoff: float = _MIN_RETRY_BACKOFF_SECONDS,
        max_backoff: float = _MAX_RETRY_BACKOFF_SECONDS,
        clock: Callable[[], float] = time.monotonic,
    ):
        if not 0.0 < refresh_ratio < 1.0:
            raise ValueError("refresh_ratio must be between 0 and 1")
        self._fetch = fetch
        self._refresh_ratio = refresh_ratio
        self._min_backoff = min_backoff
        self._max_backoff = max_backoff
        self._clock = clock
        self._token: str | None = None
        self._expires_at = 0.0
        self._refresh_at = 0.0
        self._backoff = min_backoff
        self._lock = asyncio.Lock()

    async def initialize(self) -> None:
        async with self._lock:
            self._store(await self._fetch())
            self._backoff = self._min_backoff

    def rpc_metadata(self) -> dict[str, str]:
        if self._token is None:
            raise TemporalTokenError("temporal access token is not initialized")
        return {_AUTHORIZATION_HEADER: _BEARER_PREFIX + self._token}

    def next_refresh_delay(self) -> float:
        return max(self._refresh_at - self._clock(), 0.0)

    async def refresh(self) -> float:
        async with self._lock:
            try:
                self._store(await self._fetch())
            except asyncio.CancelledError:
                raise
            except Exception as exc:  # noqa: BLE001
                delay = self._backoff
                self._backoff = min(self._backoff * 2, self._max_backoff)
                _logger.warning(
                    "temporal access token refresh failed",
                    retry_in_seconds=delay,
                    expires_in_seconds=max(self._expires_at - self._clock(), 0.0),
                    error=str(exc),
                )
                return delay
            self._backoff = self._min_backoff
            return self.next_refresh_delay()

    async def run(self, on_refresh: MetadataListener) -> None:
        delay = self.next_refresh_delay()
        while True:
            await asyncio.sleep(delay)
            previous = self._token
            delay = await self.refresh()
            if self._token != previous:
                on_refresh(self.rpc_metadata())

    def _store(self, token: AccessToken) -> None:
        if not token.value:
            raise TemporalTokenError("temporal access token is empty")
        if token.expires_in <= 0:
            raise TemporalTokenError("temporal access token lifetime must be positive")
        now = self._clock()
        self._token = token.value
        self._expires_at = now + token.expires_in
        self._refresh_at = now + token.expires_in * self._refresh_ratio


class KeycloakTokenFetcher:
    def __init__(self, token_url: str, client_id: str, client_secret: str, timeout_seconds: float):
        self._token_url = token_url
        self._client_id = client_id
        self._client_secret = client_secret
        self._timeout_seconds = timeout_seconds

    async def __call__(self) -> AccessToken:
        return await asyncio.to_thread(self._fetch)

    def _fetch(self) -> AccessToken:
        body = urllib.parse.urlencode(
            {
                "grant_type": "client_credentials",
                "client_id": self._client_id,
                "client_secret": self._client_secret,
            }
        ).encode("utf-8")
        request = urllib.request.Request(
            self._token_url,
            data=body,
            method="POST",
            headers={"Content-Type": "application/x-www-form-urlencoded", "Accept": "application/json"},
        )
        with urllib.request.urlopen(request, timeout=self._timeout_seconds) as response:
            payload = json.loads(response.read().decode("utf-8"))

        access_token = payload.get("access_token")
        expires_in = payload.get("expires_in")
        if not isinstance(access_token, str) or not access_token:
            raise TemporalTokenError("token response has no access_token")
        if isinstance(expires_in, bool) or not isinstance(expires_in, int | float) or expires_in <= 0:
            raise TemporalTokenError("token response has no positive expires_in")
        return AccessToken(value=access_token, expires_in=float(expires_in))


def build_temporal_token_provider(config: TemporalAuthConfig) -> TemporalTokenProvider | None:
    if not config.enabled:
        return None
    missing = [
        name
        for name, value in (
            ("token_url", config.token_url),
            ("client_id", config.client_id),
            ("client_secret", config.client_secret),
        )
        if not value
    ]
    if missing:
        raise ValueError(f"temporal auth is enabled but missing: {', '.join(missing)}")
    return TemporalTokenProvider(
        KeycloakTokenFetcher(
            token_url=config.token_url,
            client_id=config.client_id,
            client_secret=config.client_secret,
            timeout_seconds=config.request_timeout_seconds,
        ),
        refresh_ratio=config.refresh_ratio,
    )
