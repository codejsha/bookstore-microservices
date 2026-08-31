from typing import Annotated

from pydantic import AfterValidator, StringConstraints


def _require_non_blank(value: str) -> str:
    if not value.strip():
        raise ValueError("must not be blank")
    return value


NonBlankStr100 = Annotated[str, StringConstraints(max_length=100), AfterValidator(_require_non_blank)]
NonBlankStr255 = Annotated[str, StringConstraints(max_length=255), AfterValidator(_require_non_blank)]
NonBlankStr500 = Annotated[str, StringConstraints(max_length=500), AfterValidator(_require_non_blank)]
NonBlankStr16000 = Annotated[str, StringConstraints(max_length=16000), AfterValidator(_require_non_blank)]
Str500 = Annotated[str, StringConstraints(max_length=500)]
