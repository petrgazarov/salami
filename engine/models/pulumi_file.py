from pydantic import BaseModel, ConfigDict
from engine.models import PulumiResource


class Source(BaseModel):
    file_path: str


class PulumiFile(BaseModel):
    resources: list[PulumiResource] = []
    file_path: str
    source: Source

    model_config = ConfigDict(arbitrary_types_allowed=True)

    def pulumi_aws_imports(self):
        imports = []
        for resource in self.resources:
            resource_type = resource.resource_type.split(".")[0]
            if resource_type not in imports:
                imports.append(resource_type)
        return imports

    def other_packages(self):
        packages = []
        for resource in self.resources:
            for package in resource.other_packages:
                if package not in packages:
                    packages.append(package)
        return packages

    def code_for_imports(self):
        result = "import pulumi\n"
        other_packages_imports = ", ".join(self.other_packages())
        result += "import " + other_packages_imports + "\n"
        pulumi_aws_imports = ", ".join(self.pulumi_aws_imports())
        result += "from pulumi_aws import " + pulumi_aws_imports + "\n"
        return result

    def code_for_exports(self):
        result = ""
        for resource in self.resources:
            result += resource.code_for_exports()
        return result

    def code_for_resources(self):
        result = ""
        for resource in self.resources:
            result += resource.code_for_resource() + "\n"
        return result
