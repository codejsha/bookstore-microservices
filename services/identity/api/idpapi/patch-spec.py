#!/usr/bin/env python3

import json
import sys
from pathlib import Path

SPEC = Path(__file__).parent / "schema" / "admin-openapi.json"
OLD = "/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}/composites/clients/{client-uuid}"
NEW = "/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}/composites/clients/{composite-client-uuid}"
PARAM = {
    "name": "composite-client-uuid",
    "in": "path",
    "description": "client whose client-level roles are filtered from the composite",
    "required": True,
    "schema": {"type": "string"},
}


def main() -> int:
    spec = json.loads(SPEC.read_text())
    paths = spec.get("paths", {})
    if OLD not in paths:
        return 0

    item = paths.pop(OLD)
    for op in item.values():
        if isinstance(op, dict) and "parameters" in op:
            names = [p.get("name") for p in op["parameters"]]
            if PARAM["name"] not in names:
                op["parameters"].insert(1, dict(PARAM))
    paths[NEW] = item

    SPEC.write_text(json.dumps(spec, indent=2) + "\n")
    print(f"patched duplicate path template: {OLD}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
