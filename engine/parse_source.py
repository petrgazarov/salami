import re
from engine.pulumi_resource import PulumiResource, ResourceType, LogicalName


def parse_source(text: str):
    lines = text.splitlines()

    properties = []
    uses = []
    exports = {}
    resource_type = ResourceType("")
    logical_name = LogicalName("")
    text_lines = []

    for line in lines:
        if line.startswith("@property"):
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
        elif line.startswith("Resource type: "):
            resource_type = ResourceType(line[len("Resource type: ") :])
        elif line.startswith("Logical name: "):
            logical_name = LogicalName(line[len("Logical name: ") :])
        else:
            text_lines.append(line)

    return PulumiResource(
        resource_type=resource_type,
        logical_name=logical_name,
        code="",
        other_packages=[],
        text="\n".join(text_lines),
        properties=properties,
        uses=uses,
        exports=exports,
    )
