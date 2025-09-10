#!/usr/bin/env python3
"""Drop security requirements that name a scheme the spec never declares.

Hyperswitch's published spec is invalid: two subscription operations require a
`client_secret` scheme that is absent from components.securitySchemes. That is a
hard error for oapi-codegen, so `make payapi` could not regenerate the client at
all -- against the pinned tag or any other.

The offending requirements are dropped rather than the missing scheme invented:
we do not call the subscriptions API, and guessing at how `client_secret` is
transported would be fabricating an auth contract we have not verified.

Runs against the *bundled* spec, so the mirrored upstream artifact stays pristine.
Prints what it removed -- if a future Hyperswitch version fixes the bug, or breaks
a scheme we actually use, that shows up in the build log rather than passing
silently.
"""

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
