from uuid import uuid4

import pytest
from pydantic import ValidationError

from internal.domain.constant.author_role import AuthorRole
from internal.domain.model.command.category_command import CreateCategoryCommand, UpdateCategoryCommand
from internal.domain.model.command.faq_command import CreateFaqCommand, UpdateFaqCommand
from internal.domain.model.command.ticket_command import (
    AddCommentCommand,
    CreateTicketCommand,
    UpdateTicketCommand,
)


class TestTicketCommands:
    def test_create_ticket_command_every_field_valid_passes(self):
        CreateTicketCommand(customer_uid=uuid4(), subject="Broken book", description="Pages missing")

    @pytest.mark.parametrize("field", ["subject", "description"])
    def test_create_ticket_command_blank_field_raises_validation_error(self, field):
        kwargs = {"customer_uid": uuid4(), "subject": "s", "description": "d"}
        kwargs[field] = "  "
        with pytest.raises(ValidationError):
            CreateTicketCommand(**kwargs)

    def test_update_ticket_command_every_field_none_passes(self):
        UpdateTicketCommand()

    def test_update_ticket_command_blank_subject_raises_validation_error(self):
        with pytest.raises(ValidationError):
            UpdateTicketCommand(subject=" ")

    def test_add_comment_command_blank_body_raises_validation_error(self):
        with pytest.raises(ValidationError):
            AddCommentCommand(ticket_uid=uuid4(), author_uid=uuid4(), author_role=AuthorRole.CUSTOMER, body=" ")

    @pytest.mark.parametrize(("field", "limit"), [("subject", 255), ("description", 16000)])
    def test_ticket_commands_over_length_field_raise_validation_error(self, field, limit):
        kwargs = {"customer_uid": uuid4(), "subject": "s", "description": "d"}
        kwargs[field] = "x" * (limit + 1)
        with pytest.raises(ValidationError):
            CreateTicketCommand(**kwargs)
        with pytest.raises(ValidationError):
            UpdateTicketCommand(**{field: "x" * (limit + 1)})

    def test_create_ticket_command_subject_at_limit_passes(self):
        command = CreateTicketCommand(customer_uid=uuid4(), subject="x" * 255, description="y" * 16000)
        assert len(command.subject) == 255

    def test_add_comment_command_over_length_body_raises_validation_error(self):
        with pytest.raises(ValidationError):
            AddCommentCommand(
                ticket_uid=uuid4(),
                author_uid=uuid4(),
                author_role=AuthorRole.CUSTOMER,
                body="x" * 16001,
            )


class TestFaqCommands:
    def test_create_faq_command_every_field_valid_passes(self):
        CreateFaqCommand(question="How to refund?", answer="Go to orders.")

    def test_create_faq_command_blank_answer_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateFaqCommand(question="How to refund?", answer="")

    def test_update_faq_command_blank_question_raises_validation_error(self):
        with pytest.raises(ValidationError):
            UpdateFaqCommand(question=" ")

    @pytest.mark.parametrize(("field", "limit"), [("question", 500), ("answer", 16000)])
    def test_faq_commands_over_length_field_raise_validation_error(self, field, limit):
        kwargs = {"question": "q", "answer": "a"}
        kwargs[field] = "x" * (limit + 1)
        with pytest.raises(ValidationError):
            CreateFaqCommand(**kwargs)
        with pytest.raises(ValidationError):
            UpdateFaqCommand(**{field: "x" * (limit + 1)})


class TestCategoryCommands:
    def test_create_category_command_every_field_valid_passes(self):
        CreateCategoryCommand(name="Orders")

    def test_create_category_command_blank_name_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCategoryCommand(name=" ")

    def test_update_category_command_empty_description_keeps_it_empty(self):
        assert UpdateCategoryCommand(description="").description == ""

    def test_category_commands_over_length_name_raise_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCategoryCommand(name="x" * 101)
        with pytest.raises(ValidationError):
            UpdateCategoryCommand(name="x" * 101)

    def test_category_commands_over_length_description_raise_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCategoryCommand(name="Orders", description="x" * 501)
        with pytest.raises(ValidationError):
            UpdateCategoryCommand(description="x" * 501)

    def test_create_faq_command_answer_has_surrounding_whitespace_preserves_it(self):
        body = "\n    def f():\n        pass\n"
        cmd = CreateFaqCommand(question="How?", answer=body)
        assert cmd.answer == body
