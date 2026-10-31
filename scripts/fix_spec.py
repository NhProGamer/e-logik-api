"""Clean NelmioApiDocBundle artefacts from openapi.json so oapi-codegen accepts it.

Usage: python3 scripts/fix_spec.py openapi.json > spec.fixed.json
"""
import json
import sys


def fix(o):
    if isinstance(o, list):
        return [fix(v) for v in o]
    if not isinstance(o, dict):
        return o
    o = {k: fix(v) for k, v in o.items()}
    if o.get("type") == "bool":
        o["type"] = "boolean"
    # Fields documented as date-time but with a date-only example expect
    # "2006-01-02"; keep them as plain strings so callers send that format.
    example = o.get("example")
    if o.get("format") == "date-time" and isinstance(example, str) and "T" not in example:
        del o["format"]
    # Nelmio sometimes emits the PHP enum class name instead of the values.
    if "enum" in o and not isinstance(o["enum"], list):
        del o["enum"]
    # oneOf [$ref] or [$ref, null] is just a nullable $ref.
    if "oneOf" in o:
        refs = [s for s in o["oneOf"] if s.get("enum") != [None]]
        if len(refs) == 1:
            del o["oneOf"]
            o["allOf"] = refs
            o["nullable"] = True
    return o


json.dump(fix(json.load(open(sys.argv[1]))), sys.stdout, ensure_ascii=False, indent=1)
