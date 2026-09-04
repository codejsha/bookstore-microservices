from typing import Annotated

from pydantic import AfterValidator, StringConstraints


def _require_non_blank(value: str) -> str:
    if not value.strip():
        raise ValueError("must not be blank")
    return value


def non_blank_str(max_length: int) -> type[str]:
    return Annotated[str, StringConstraints(max_length=max_length), AfterValidator(_require_non_blank)]


NonBlankStr255 = non_blank_str(255)
NonBlankStr10000 = non_blank_str(10000)
