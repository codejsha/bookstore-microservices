from collections.abc import Mapping
from typing import Any

import structlog

from internal.application.port.event.event_publisher import EventPublisher

logger = structlog.get_logger()


class NoopEventPublisher(EventPublisher):
    async def publish(self, resource: str, uid: str, action: str, fields: Mapping[str, Any]) -> None:
        logger.debug("event not published (noop publisher)", resource=resource, uid=uid, action=action, **fields)
