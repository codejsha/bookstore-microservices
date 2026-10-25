from abc import ABC, abstractmethod
from collections.abc import Mapping
from typing import Any


class EventPublisher(ABC):
    @abstractmethod
    async def publish(self, resource: str, uid: str, action: str, fields: Mapping[str, Any]) -> None: ...
