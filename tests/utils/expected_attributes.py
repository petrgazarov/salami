import os
import inspect
import importlib.util


def expected_attributes():
    caller_dir = os.path.dirname(inspect.stack()[1].filename)
    expected_attrs_path = os.path.join(caller_dir, "expected_attributes.py")

    spec = importlib.util.spec_from_file_location("expected_attrs", expected_attrs_path)
    if spec is None or spec.loader is None:
        raise ImportError(f"No module found at {expected_attrs_path}")
    expected_attrs = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(expected_attrs)
    return [
        (attr, getattr(expected_attrs, attr))
        for attr in dir(expected_attrs)
        if not attr.startswith("__")
    ]
