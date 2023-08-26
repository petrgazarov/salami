from engine.utils.file import get_full_path, open_relative_file
from .pulumi_resources import pulumi_resources
from . import expected_attributes as _expected_attributes

expected_attributes = {
    name: getattr(_expected_attributes, name)
    for name in dir(_expected_attributes)
    if not name.startswith("__")
}

source_filepath = get_full_path("source.txt")
expected_code = open_relative_file("expected_code.txt")

__all__ = [
    "pulumi_resources",
    "expected_attributes",
    "expected_code",
    "source_filepath",
]
