#!/usr/bin/env python3

import json
import sys
from pathlib import Path


def main() -> int:
    path = Path(sys.argv[1])
    spec = json.loads(path.read_text())

    declared = set(spec.get("components", {}).get("securitySchemes", {}))
    removed: list[str] = []

    for route, item in spec.get("paths", {}).items():
        for method, operation in item.items():
            if not isinstance(operation, dict) or "security" not in operation:
                continue
            kept = []
            for requirement in operation["security"] or []:
                undeclared = set(requirement) - declared
                if undeclared:
                    removed.append(f"{method.upper()} {route} -> {', '.join(sorted(undeclared))}")
                else:
                    kept.append(requirement)
            operation["security"] = kept

    if not removed:
        print("strip-undeclared-security: nothing to strip (upstream spec is valid)")
        return 0

    path.write_text(json.dumps(spec, indent=2) + "\n")
    print(f"strip-undeclared-security: dropped {len(removed)} undeclared requirement(s):")
    for entry in removed:
        print(f"  {entry}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
