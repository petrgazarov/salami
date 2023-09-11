from typing import Any, Dict
import copy


def remove_titles(schema: Dict[str, Any]) -> Dict[str, Any]:
    schema.pop("title", None)

    for value in schema.values():
        if isinstance(value, dict):
            remove_titles(value)
        elif isinstance(value, list):
            for item in value:
                if isinstance(item, dict):
                    remove_titles(item)

    return schema


def inline_enums(schema):
    definitions = schema.pop("definitions", {})

    def replace_ref(obj):
        if isinstance(obj, dict):
            if "$ref" in obj:
                ref_name = obj["$ref"].split("/")[-1]
                # preserve original description
                description = obj.get("description")
                updated_def = definitions[ref_name].copy()
                if description:
                    updated_def["description"] = description
                return updated_def
            if "allOf" in obj:
                # keep the original description
                description = obj.get("description")
                updated_obj = replace_ref(obj["allOf"][0])
                # restore the original description
                if description:
                    updated_obj["description"] = description
                return updated_obj
            return {k: replace_ref(v) for k, v in obj.items()}
        if isinstance(obj, list):
            return [replace_ref(item) for item in obj]
        return obj

    return replace_ref(schema)


def pydantic_to_openai(pydantic_class):
    parameters_schema = pydantic_class.input_schema.model_json_schema()
    parameters_schema = copy.deepcopy(parameters_schema)
    normalized_parameters = inline_enums(remove_titles(parameters_schema))

    return {
        "name": pydantic_class.name,
        "description": pydantic_class.description,
        "parameters": normalized_parameters,
    }
