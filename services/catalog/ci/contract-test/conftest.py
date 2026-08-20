import json
import os
from pathlib import Path

import schemathesis

SCHEMA_PATH = Path(os.environ["CONTRACT_SCHEMA_PATH"])
BASE_URL = os.environ.get("CONTRACT_BASE_URL", "http://localhost:8080")
BEARER_TOKEN = os.environ.get("CONTRACT_BEARER_TOKEN", "")


def _load_schema_dict() -> dict:
    raw = json.loads(SCHEMA_PATH.read_text())
    # Schemathesis 3.x supports up to OpenAPI 3.0
    version = raw.get("openapi", "")
    if version.startswith(("3.1", "3.2")):
        raw["openapi"] = "3.0.3"
    return raw


schema = schemathesis.from_dict(
    _load_schema_dict(),
    base_url=BASE_URL,
    force_schema_version="30",
)


@schemathesis.hook
def before_call(context, case):
    if BEARER_TOKEN:
        case.headers = {**(case.headers or {}), "Authorization": f"Bearer {BEARER_TOKEN}"}
