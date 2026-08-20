from unittest.mock import MagicMock
from uuid import uuid4

import pytest

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
    async def test_persists_root_category(self, service: SupportService, category_repo: MagicMock) -> None:
        result = await service.create_category(CreateCategoryCommand(name="Billing", description="d"))
        assert result.name == "Billing"
        assert result.parent_id is None
        category_repo.find_id_by_uid.assert_not_called()
        category_repo.save.assert_called_once()

    async def test_resolves_parent_uid(self, service: SupportService, category_repo: MagicMock) -> None:
        parent_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 5
        result = await service.create_category(CreateCategoryCommand(name="Sub", parent_uid=parent_uid))
        assert result.parent_id == 5
        category_repo.find_id_by_uid.assert_called_once_with(parent_uid)


class TestGetCategory:
    async def test_delegates_to_repository(self, service: SupportService, category_repo: MagicMock) -> None:
        category = make_category()
        category_repo.find_by_uid.return_value = category
        assert await service.get_category(category.uid) is category


class TestListCategories:
    async def test_delegates_to_repository(self, service: SupportService, category_repo: MagicMock) -> None:
        categories = [make_category()]
        category_repo.find_all.return_value = categories
        assert await service.list_categories() == categories


class TestUpdateCategory:
    async def test_returns_none_when_missing(self, service: SupportService, category_repo: MagicMock) -> None:
        category_repo.find_by_uid.return_value = None
        assert await service.update_category(uuid4(), UpdateCategoryCommand(name="X")) is None
        category_repo.save.assert_not_called()

    async def test_updates_only_provided_fields(self, service: SupportService, category_repo: MagicMock) -> None:
        category = make_category(name="Original", description="d")
        category_repo.find_by_uid.return_value = category
        result = await service.update_category(category.uid, UpdateCategoryCommand(name="New name"))
        assert result is not None
        assert result.name == "New name"
        assert result.description == "d"

    async def test_resolves_parent_uid(self, service: SupportService, category_repo: MagicMock) -> None:
        category = make_category()
        category_repo.find_by_uid.return_value = category
        parent_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 11
        result = await service.update_category(category.uid, UpdateCategoryCommand(parent_uid=parent_uid))
        assert result is not None
        assert result.parent_id == 11


class TestDeleteCategory:
    async def test_delegates_to_repository(self, service: SupportService, category_repo: MagicMock) -> None:
        category_repo.delete_by_uid.return_value = True
        uid = uuid4()
        assert await service.delete_category(uid) is True
        category_repo.delete_by_uid.assert_called_once_with(uid)
