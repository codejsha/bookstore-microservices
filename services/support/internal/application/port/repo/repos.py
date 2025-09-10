from abc import ABC, abstractmethod
from collections.abc import Callable
from uuid import UUID

from internal.domain.aggregate.faq_aggregate import FaqAggregate
from internal.domain.aggregate.ticket_aggregate import TicketAggregate
from internal.domain.aggregate.ticket_category_aggregate import TicketCategoryAggregate
from internal.domain.aggregate.ticket_comment_aggregate import TicketCommentAggregate
from internal.domain.model.option.faq_option import FaqSearchOption
from internal.domain.model.option.ticket_option import TicketFilterOption


class TicketRepository(ABC):
    @abstractmethod
    async def save(self, ticket: TicketAggregate) -> TicketAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> TicketAggregate | None: ...

    @abstractmethod
    async def find_id_by_uid(self, uid: UUID) -> int | None: ...

    @abstractmethod
    async def find_all(self, option: TicketFilterOption) -> tuple[list[TicketAggregate], int]: ...

    @abstractmethod
    async def update(
        self,
        uid: UUID,
        mutator: Callable[[TicketAggregate], TicketAggregate],
    ) -> TicketAggregate | None: ...


class TicketCommentRepository(ABC):
    @abstractmethod
    async def save(self, comment: TicketCommentAggregate) -> TicketCommentAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> TicketCommentAggregate | None: ...

    @abstractmethod
    async def find_by_ticket_id(
        self, ticket_id: int, include_internal: bool = True
    ) -> list[TicketCommentAggregate]: ...


class TicketCategoryRepository(ABC):
    @abstractmethod
    async def save(self, category: TicketCategoryAggregate) -> TicketCategoryAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> TicketCategoryAggregate | None: ...

    @abstractmethod
    async def find_id_by_uid(self, uid: UUID) -> int | None: ...

    @abstractmethod
    async def find_all(self) -> list[TicketCategoryAggregate]: ...

    @abstractmethod
    async def delete_by_uid(self, uid: UUID) -> bool: ...


class FaqRepository(ABC):
    @abstractmethod
    async def save(self, faq: FaqAggregate) -> FaqAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> FaqAggregate | None: ...

    @abstractmethod
    async def search(self, option: FaqSearchOption) -> tuple[list[FaqAggregate], int]: ...

    @abstractmethod
    async def update(
        self,
        uid: UUID,
        mutator: Callable[[FaqAggregate], FaqAggregate],
    ) -> FaqAggregate | None: ...

    @abstractmethod
    async def increment_view_count(self, uid: UUID) -> None: ...

    @abstractmethod
    async def delete_by_uid(self, uid: UUID) -> bool: ...
