from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.author_role import AuthorRole
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.error import UnknownReferenceError
from internal.domain.model.command.ticket_command import (
    AddCommentCommand,
    CreateTicketCommand,
    UpdateTicketCommand,
    UpdateTicketStatusCommand,
)
from internal.domain.model.option.ticket_option import TicketFilterOption
from internal.domain.service.support_service import SupportService
from tests.conftest import make_comment, make_ticket


@pytest.fixture
def service(
    ticket_repo: MagicMock,
    comment_repo: MagicMock,
    category_repo: MagicMock,
    faq_repo: MagicMock,
) -> SupportService:
    return SupportService(
        ticket_repo=ticket_repo,
        comment_repo=comment_repo,
        category_repo=category_repo,
        faq_repo=faq_repo,
    )


class TestCreateTicket:
    async def test_create_ticket_without_category_persists_open_ticket(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        customer_uid = uuid4()
        command = CreateTicketCommand(
            customer_uid=customer_uid,
            subject="Order issue",
            description="Did not arrive",
            priority=TicketPriority.HIGH,
        )
        result = await service.create_ticket(command)
        assert result.status == TicketStatus.OPEN
        assert result.priority == TicketPriority.HIGH
        assert result.customer_uid == customer_uid
        assert result.category_id is None
        category_repo.find_id_by_uid.assert_not_called()
        ticket_repo.save.assert_called_once()

    async def test_create_ticket_with_category_uid_resolves_category_id(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        category_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 7
        command = CreateTicketCommand(customer_uid=uuid4(), subject="s", description="d", category_uid=category_uid)
        result = await service.create_ticket(command)
        assert result.category_id == 7
        category_repo.find_id_by_uid.assert_called_once_with(category_uid)

    async def test_create_ticket_unknown_category_uid_raises_unknown_reference_error(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        category_repo.find_id_by_uid.return_value = None
        command = CreateTicketCommand(customer_uid=uuid4(), subject="s", description="d", category_uid=uuid4())
        with pytest.raises(UnknownReferenceError):
            await service.create_ticket(command)
        ticket_repo.save.assert_not_called()


class TestGetTicket:
    async def test_get_ticket_delegates_to_repository(self, service: SupportService, ticket_repo: MagicMock) -> None:
        ticket = make_ticket()
        ticket_repo.find_by_uid.return_value = ticket
        assert await service.get_ticket(ticket.uid) is ticket


class TestListTickets:
    async def test_list_tickets_with_option_delegates_to_repository(
        self, service: SupportService, ticket_repo: MagicMock
    ) -> None:
        option = TicketFilterOption()
        ticket_repo.find_all.return_value = ([make_ticket()], 1)
        result, total = await service.list_tickets(option)
        assert total == 1


class TestUpdateTicket:
    async def test_update_ticket_missing_is_none(self, service: SupportService, ticket_repo: MagicMock) -> None:
        ticket_repo.find_by_uid.return_value = None
        assert await service.update_ticket(uuid4(), UpdateTicketCommand(subject="s")) is None
        ticket_repo.save.assert_not_called()

    async def test_update_ticket_partial_fields_updates_only_those(
        self, service: SupportService, ticket_repo: MagicMock
    ) -> None:
        ticket = make_ticket()
        original_description = ticket.description
        ticket_repo.find_by_uid.return_value = ticket
        result = await service.update_ticket(
            ticket.uid,
            UpdateTicketCommand(subject="new subject", priority=TicketPriority.URGENT),
        )
        assert result is not None
        assert result.subject == "new subject"
        assert result.priority == TicketPriority.URGENT
        assert result.description == original_description

    async def test_update_ticket_with_category_uid_resolves_category_id(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        ticket = make_ticket()
        ticket_repo.find_by_uid.return_value = ticket
        category_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 9
        result = await service.update_ticket(ticket.uid, UpdateTicketCommand(category_uid=category_uid))
        assert result is not None
        assert result.category_id == 9

    async def test_update_ticket_unknown_category_uid_raises_unknown_reference_error(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        ticket_repo.find_by_uid.return_value = make_ticket()
        category_repo.find_id_by_uid.return_value = None
        with pytest.raises(UnknownReferenceError):
            await service.update_ticket(uuid4(), UpdateTicketCommand(category_uid=uuid4()))


class TestUpdateTicketStatus:
    async def test_update_ticket_status_missing_is_none(self, service: SupportService, ticket_repo: MagicMock) -> None:
        ticket_repo.find_by_uid.return_value = None
        assert (
            await service.update_ticket_status(uuid4(), UpdateTicketStatusCommand(status=TicketStatus.RESOLVED)) is None
        )

    async def test_update_ticket_status_resolved_sets_resolved_at(
        self, service: SupportService, ticket_repo: MagicMock
    ) -> None:
        ticket = make_ticket(status=TicketStatus.OPEN, resolved_at=None)
        ticket_repo.find_by_uid.return_value = ticket
        result = await service.update_ticket_status(ticket.uid, UpdateTicketStatusCommand(status=TicketStatus.RESOLVED))
        assert result is not None
        assert result.status == TicketStatus.RESOLVED
        assert result.resolved_at is not None

    async def test_update_ticket_status_already_resolved_keeps_original_resolved_at(
        self, service: SupportService, ticket_repo: MagicMock
    ) -> None:
        from datetime import UTC, datetime, timedelta

        original_resolved = datetime.now(UTC) - timedelta(days=1)
        ticket = make_ticket(status=TicketStatus.RESOLVED, resolved_at=original_resolved)
        ticket_repo.find_by_uid.return_value = ticket
        result = await service.update_ticket_status(ticket.uid, UpdateTicketStatusCommand(status=TicketStatus.RESOLVED))
        assert result is not None
        assert result.resolved_at == original_resolved

    async def test_update_ticket_status_not_resolved_leaves_resolved_at_unset(
        self, service: SupportService, ticket_repo: MagicMock
    ) -> None:
        ticket = make_ticket(status=TicketStatus.OPEN, resolved_at=None)
        ticket_repo.find_by_uid.return_value = ticket
        result = await service.update_ticket_status(
            ticket.uid, UpdateTicketStatusCommand(status=TicketStatus.IN_PROGRESS)
        )
        assert result is not None
        assert result.status == TicketStatus.IN_PROGRESS
        assert result.resolved_at is None


class TestAddComment:
    async def test_add_comment_missing_is_none(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        comment_repo: MagicMock,
    ) -> None:
        ticket_repo.find_id_by_uid.return_value = None
        result = await service.add_comment(
            AddCommentCommand(
                ticket_uid=uuid4(),
                author_uid=uuid4(),
                author_role=AuthorRole.CUSTOMER,
                body="hi",
            )
        )
        assert result is None
        comment_repo.save.assert_not_called()

    async def test_add_comment_found_creates_comment_with_its_id(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        comment_repo: MagicMock,
    ) -> None:
        author_uid = uuid4()
        ticket_repo.find_id_by_uid.return_value = 42
        result = await service.add_comment(
            AddCommentCommand(
                ticket_uid=uuid4(),
                author_uid=author_uid,
                author_role=AuthorRole.AGENT,
                body="agent reply",
                internal=True,
            )
        )
        assert result is not None
        assert result.ticket_id == 42
        assert result.author_uid == author_uid
        assert result.author_role == AuthorRole.AGENT
        assert result.body == "agent reply"
        assert result.internal is True


class TestListComments:
    async def test_list_comments_missing_is_none(self, service: SupportService, ticket_repo: MagicMock) -> None:
        ticket_repo.find_id_by_uid.return_value = None
        assert await service.list_comments(uuid4()) is None

    async def test_list_comments_with_include_internal_passes_it_to_repository(
        self,
        service: SupportService,
        ticket_repo: MagicMock,
        comment_repo: MagicMock,
    ) -> None:
        ticket_repo.find_id_by_uid.return_value = 5
        comments = [make_comment(ticket_id=5)]
        comment_repo.find_by_ticket_id.return_value = comments
        result = await service.list_comments(uuid4(), include_internal=False)
        assert result == comments
        comment_repo.find_by_ticket_id.assert_called_once_with(5, include_internal=False)
