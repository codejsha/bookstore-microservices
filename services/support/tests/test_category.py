from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.error import ConflictError, UnknownReferenceError
from internal.domain.model.command.category_command import (
    CreateCategoryCommand,
    UpdateCategoryCommand,
)
from internal.domain.service.support_service import SupportService
from tests.conftest import make_category


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


class TestCreateCategory:
    async def test_create_category_without_parent_persists_root_category(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        result = await service.create_category(CreateCategoryCommand(name="Billing", description="d"))
        assert result.name == "Billing"
        assert result.parent_id is None
        category_repo.find_id_by_uid.assert_not_called()
        category_repo.save.assert_called_once()

    async def test_create_category_with_parent_uid_resolves_parent_id(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        parent_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 5
        result = await service.create_category(CreateCategoryCommand(name="Sub", parent_uid=parent_uid))
        assert result.parent_id == 5
        category_repo.find_id_by_uid.assert_called_once_with(parent_uid)

    async def test_create_category_unknown_parent_uid_raises_unknown_reference_error(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category_repo.find_id_by_uid.return_value = None
        with pytest.raises(UnknownReferenceError):
            await service.create_category(CreateCategoryCommand(name="Sub", parent_uid=uuid4()))
        category_repo.save.assert_not_called()

    async def test_create_category_duplicate_name_raises_conflict_error(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category_repo.find_by_name.return_value = make_category(name="Billing")
        with pytest.raises(ConflictError):
            await service.create_category(CreateCategoryCommand(name="Billing"))
        category_repo.save.assert_not_called()


class TestGetCategory:
    async def test_get_category_delegates_to_repository(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category = make_category()
        category_repo.find_by_uid.return_value = category
        assert await service.get_category(category.uid) is category


class TestListCategories:
    async def test_list_categories_delegates_to_repository(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        categories = [make_category()]
        category_repo.find_all.return_value = categories
        assert await service.list_categories() == categories


class TestUpdateCategory:
    async def test_update_category_missing_is_none(self, service: SupportService, category_repo: MagicMock) -> None:
        category_repo.find_by_uid.return_value = None
        assert await service.update_category(uuid4(), UpdateCategoryCommand(name="X")) is None
        category_repo.save.assert_not_called()

    async def test_update_category_partial_fields_updates_only_those(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category = make_category(name="Original", description="d")
        category_repo.find_by_uid.return_value = category
        result = await service.update_category(category.uid, UpdateCategoryCommand(name="New name"))
        assert result is not None
        assert result.name == "New name"
        assert result.description == "d"

    async def test_update_category_with_parent_uid_resolves_parent_id(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category = make_category()
        category_repo.find_by_uid.return_value = category
        parent_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 11
        result = await service.update_category(category.uid, UpdateCategoryCommand(parent_uid=parent_uid))
        assert result is not None
        assert result.parent_id == 11

    async def test_update_category_unknown_parent_uid_raises_unknown_reference_error(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category_repo.find_by_uid.return_value = make_category()
        category_repo.find_id_by_uid.return_value = None
        with pytest.raises(UnknownReferenceError):
            await service.update_category(uuid4(), UpdateCategoryCommand(parent_uid=uuid4()))
        category_repo.save.assert_not_called()

    async def test_update_category_name_taken_by_another_category_raises_conflict_error(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category = make_category(name="Original")
        category_repo.find_by_uid.return_value = category
        category_repo.find_by_name.return_value = make_category(name="Billing")
        with pytest.raises(ConflictError):
            await service.update_category(category.uid, UpdateCategoryCommand(name="Billing"))
        category_repo.save.assert_not_called()

    async def test_update_category_name_unchanged_saves_without_conflict(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category = make_category(name="Billing")
        category_repo.find_by_uid.return_value = category
        category_repo.find_by_name.return_value = category
        result = await service.update_category(category.uid, UpdateCategoryCommand(name="Billing"))
        assert result is not None
        assert result.name == "Billing"


class TestDeleteCategory:
    async def test_delete_category_delegates_to_repository(
        self, service: SupportService, category_repo: MagicMock
    ) -> None:
        category_repo.delete_by_uid.return_value = True
        uid = uuid4()
        assert await service.delete_category(uid) is True
        category_repo.delete_by_uid.assert_called_once_with(uid)
