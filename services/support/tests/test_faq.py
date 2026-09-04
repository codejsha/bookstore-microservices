from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.error import UnknownReferenceError
from internal.domain.model.command.faq_command import CreateFaqCommand, UpdateFaqCommand
from internal.domain.model.option.faq_option import FaqSearchOption
from internal.domain.service.support_service import SupportService
from tests.conftest import make_faq


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


class TestCreateFaq:
    async def test_create_faq_without_published_persists_unpublished_faq(
        self, service: SupportService, faq_repo: MagicMock
    ) -> None:
        result = await service.create_faq(CreateFaqCommand(question="Q?", answer="A."))
        assert result.published is False
        assert result.view_count == 0
        assert result.category_id is None
        faq_repo.save.assert_called_once()

    async def test_create_faq_unknown_category_uid_raises_unknown_reference_error(
        self,
        service: SupportService,
        faq_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        category_repo.find_id_by_uid.return_value = None
        with pytest.raises(UnknownReferenceError):
            await service.create_faq(CreateFaqCommand(question="Q?", answer="A.", category_uid=uuid4()))
        faq_repo.save.assert_not_called()

    async def test_create_faq_with_category_uid_resolves_category_id(
        self,
        service: SupportService,
        category_repo: MagicMock,
    ) -> None:
        category_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 3
        result = await service.create_faq(
            CreateFaqCommand(question="Q?", answer="A.", category_uid=category_uid, published=True)
        )
        assert result.category_id == 3
        assert result.published is True


class TestGetFaq:
    async def test_get_faq_missing_is_none(self, service: SupportService, faq_repo: MagicMock) -> None:
        faq_repo.find_by_uid.return_value = None
        assert await service.get_faq(uuid4()) is None

    async def test_get_faq_faq_published_increments_view_count(
        self, service: SupportService, faq_repo: MagicMock
    ) -> None:
        faq = make_faq(published=True, view_count=5)
        faq_repo.find_by_uid.return_value = faq
        result = await service.get_faq(faq.uid)
        assert result is not None
        assert result.view_count == 6
        faq_repo.increment_view_count.assert_called_once_with(faq.uid)

    async def test_get_faq_unpublished_leaves_view_count(self, service: SupportService, faq_repo: MagicMock) -> None:
        faq = make_faq(published=False, view_count=5)
        faq_repo.find_by_uid.return_value = faq
        result = await service.get_faq(faq.uid)
        assert result is not None
        assert result.view_count == 5
        faq_repo.increment_view_count.assert_not_called()

    async def test_get_faq_increment_view_false_leaves_view_count(
        self, service: SupportService, faq_repo: MagicMock
    ) -> None:
        faq = make_faq(published=True, view_count=5)
        faq_repo.find_by_uid.return_value = faq
        result = await service.get_faq(faq.uid, increment_view=False)
        assert result is not None
        assert result.view_count == 5
        faq_repo.increment_view_count.assert_not_called()


class TestSearchFaqs:
    async def test_search_faqs_with_option_delegates_to_repository(
        self, service: SupportService, faq_repo: MagicMock
    ) -> None:
        option = FaqSearchOption(query="hello")
        faqs = [make_faq()]
        faq_repo.search.return_value = (faqs, 1)
        result, total = await service.search_faqs(option)
        assert result == faqs
        assert total == 1
        faq_repo.search.assert_called_once_with(option)


class TestUpdateFaq:
    async def test_update_faq_missing_is_none(self, service: SupportService, faq_repo: MagicMock) -> None:
        faq_repo.find_by_uid.return_value = None
        assert await service.update_faq(uuid4(), UpdateFaqCommand(question="X")) is None
        faq_repo.save.assert_not_called()

    async def test_update_faq_partial_fields_updates_only_those(
        self, service: SupportService, faq_repo: MagicMock
    ) -> None:
        faq = make_faq(published=False)
        original_answer = faq.answer
        faq_repo.find_by_uid.return_value = faq
        result = await service.update_faq(faq.uid, UpdateFaqCommand(question="New Q?", published=True))
        assert result is not None
        assert result.question == "New Q?"
        assert result.published is True
        assert result.answer == original_answer

    async def test_update_faq_with_category_uid_resolves_category_id(
        self,
        service: SupportService,
        faq_repo: MagicMock,
        category_repo: MagicMock,
    ) -> None:
        faq = make_faq()
        faq_repo.find_by_uid.return_value = faq
        category_uid = uuid4()
        category_repo.find_id_by_uid.return_value = 8
        result = await service.update_faq(faq.uid, UpdateFaqCommand(category_uid=category_uid))
        assert result is not None
        assert result.category_id == 8


class TestDeleteFaq:
    async def test_delete_faq_delegates_to_repository(self, service: SupportService, faq_repo: MagicMock) -> None:
        faq_repo.delete_by_uid.return_value = True
        uid = uuid4()
        assert await service.delete_faq(uid) is True
        faq_repo.delete_by_uid.assert_called_once_with(uid)
