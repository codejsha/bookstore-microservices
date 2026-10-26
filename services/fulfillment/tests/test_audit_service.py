from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.audit_status import AuditStatus
from internal.domain.model.command.inventory_command import StockAuditCreateCommand, StockAuditItemCommand
from internal.domain.model.error import InvalidTransitionError
from internal.domain.service.inventory_service import AuditService
from tests.conftest import make_audit


@pytest.fixture
def service(audit_repo: MagicMock, event_publisher: MagicMock) -> AuditService:
    return AuditService(audit_repo=audit_repo, event_publisher=event_publisher)


async def test_create_audit_maps_items_and_publishes(
    service: AuditService, audit_repo: MagicMock, event_publisher: MagicMock
) -> None:
    audit = make_audit()
    audit_repo.create.return_value = audit
    edition = uuid4()
    command = StockAuditCreateCommand(
        warehouse_uid=uuid4(),
        items=[StockAuditItemCommand(edition_uid=edition, actual_quantity=3)],
        notes="count",
    )

    result = await service.create_audit(command)

    assert result is audit
    params = audit_repo.create.call_args.args[0]
    assert params.warehouse_uid == command.warehouse_uid
    assert [(i.edition_uid, i.actual_quantity) for i in params.items] == [(edition, 3)]
    assert params.notes == "count"
    assert event_publisher.publish.await_args.args[:3] == ("audit", str(audit.uid), "created")


async def test_complete_audit_missing_is_none(service: AuditService, audit_repo: MagicMock) -> None:
    audit_repo.find_by_uid.return_value = None
    assert await service.complete_audit(uuid4()) is None
    audit_repo.complete.assert_not_awaited()


async def test_complete_audit_completed_raises_invalid_transition(service: AuditService, audit_repo: MagicMock) -> None:
    audit_repo.find_by_uid.return_value = make_audit(status=AuditStatus.COMPLETED)
    with pytest.raises(InvalidTransitionError):
        await service.complete_audit(uuid4())
    audit_repo.complete.assert_not_awaited()


async def test_complete_audit_in_progress_delegates_and_publishes(
    service: AuditService, audit_repo: MagicMock, event_publisher: MagicMock
) -> None:
    audit = make_audit()
    audit_repo.find_by_uid.return_value = audit
    audit_repo.complete.return_value = make_audit(uid=audit.uid, status=AuditStatus.COMPLETED)

    result = await service.complete_audit(audit.uid)

    assert result is not None and result.status == AuditStatus.COMPLETED
    assert event_publisher.publish.await_args.args[:3] == ("audit", str(audit.uid), "completed")


async def test_get_audit_missing_is_none(service: AuditService, audit_repo: MagicMock) -> None:
    audit_repo.find_by_uid.return_value = None
    assert await service.get_audit(uuid4()) is None
