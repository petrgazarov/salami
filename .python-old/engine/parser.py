import re
from engine.types import LogicalName, ResourceType, ParsedResource


def split_text_into_resources(text: str):
    text_blocks = text.split("\n\n\n")
    text_blocks = map(str.strip, text_blocks)
    return text_blocks


def parse_resource_text(text: str):
    lines = text.splitlines()
    
    resource_type = ResourceType("")
    logical_name = LogicalName("")
    text_lines = []
    uses = []
    properties = []
    exports = {}
    referenced_variables = []

    for line in lines:
        if line.startswith("Resource type: "):
            resource_type = ResourceType(line[len("Resource type: ") :])
            break
        elif line.startswith("Logical name: "):
            logical_name = LogicalName(line[len("Logical name: ") :])
            break
        elif line.startswith("@property"):
            properties.append(line[len("@property(") : -1])
        elif line.startswith("@uses"):
            uses += [
                LogicalName(item.strip())
                for item in line[len("@uses(") : -1].split(", ")
            ]
        elif line.startswith("@exports"):
            match = re.search(r"@exports\(([^:]+): ([^\)]+)\)", line)
            if match:
                exports[match.group(1).strip()] = match.group(2).strip()
        else:
            text_lines.append(line)

        variable_matches = re.findall(r"\{(\w+)\}", line)
        if variable_matches:
            referenced_variables.extend(variable_matches)
    referenced_variables = list(set(referenced_variables))

    return ParsedResource(
        resource_type=resource_type,
        logical_name=logical_name,
        free_text="\n".join(text_lines),
        properties=properties,
        uses=uses,
        exports=exports,
        referenced_variables=referenced_variables,
    )
