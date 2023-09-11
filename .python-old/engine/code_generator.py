import asyncio
from pydantic import BaseModel
from engine.types import LogicalName, Target
from engine.symbol_table import SymbolTable
from engine.get_pulumi_code import get_pulumi_code


def code_for_exports(exports: dict[str, str], variable_name: str):
    result = ""
    for property, name in exports.items():
        result += f'pulumi.export("{name}", {variable_name}.{property})\n'
    return result


def code_for_resource(code: str, variable_name: str):
    return f"{variable_name} = {code}\n"


def variable_name_from_logical_name(logical_name: LogicalName):
    return "".join(
        ["_" + i.lower() if i.isupper() else i for i in logical_name]
    ).lstrip("_")


class CodeGenerator(BaseModel):
    target: Target = Target.PULUMI_PYTHON
    generated_code: dict[LogicalName, str] = {}
    symbol_table: SymbolTable

    async def run(self):
        logical_names = list(self.symbol_table.resources.keys())
        tasks = [
            get_pulumi_code(
                resource=symbol_table.resources[logical_name],
                symbol_table=symbol_table,
            )
            for logical_name in logical_names
        ]
        results = await asyncio.gather(*tasks)
        for logical_name, result in zip(logical_names, results):
            self.generated_code[logical_name] = result