#!/usr/bin/env python3

import json
import re
import sys
from pathlib import Path

TEMPLATE_VAR = re.compile(r"\{([^{}]+)\}")


def declared_names(item: dict, operation: dict) -> set[str]:
    names = set()
    for params in (item.get("parameters"), operation.get("parameters")):
        for param in params or []:
            if param.get("in") == "path":
                names.add(param.get("name"))
    return names


def main() -> int:
    path = Path(sys.argv[1])
    spec = json.loads(path.read_text())

    added: list[str] = []

    for route, item in spec.get("paths", {}).items():
        variables = TEMPLATE_VAR.findall(route)
        if not variables:
            continue
        operations = [op for op in item.values() if isinstance(op, dict) and "responses" in op]
        for name in variables:
            if all(name in declared_names(item, op) for op in operations):
                continue
            resource = route.strip("/").split("/", 1)[0]
            item.setdefault("parameters", []).append(
                {
                    "name": name,
                    "in": "path",
                    "description": f"The identifier for {resource}",
                    "required": True,
                    "schema": {"type": "string"},
                }
            )
            added.append(f"{route} -> {name}")

    if not added:
        print("add-missing-path-params: nothing to add (upstream spec declares its path parameters)")
        return 0

    path.write_text(json.dumps(spec, indent=2) + "\n")
    print(f"add-missing-path-params: declared {len(added)} missing path parameter(s):")
    for entry in added:
        print(f"  {entry}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
