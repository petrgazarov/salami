from pydantic import BaseModel
from engine.types import Resource


class FileManager(BaseModel):
    source_mapping: dict[str, list[Resource]] = {}
