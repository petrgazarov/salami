import os
from rich.logging import RichHandler
from rich.text import Text
import logging


class RichHandlerColor(RichHandler):
    def render_message(self, record, message):
        if "[" in message and "]" in message:  # if it looks like rich's color tag
            message = Text.from_markup(message)
        else:  # let RichHandler handle the message as usual
            message = super().render_message(record, message)
        return message


logging.basicConfig(
    level=os.getenv("LOG_LEVEL", "INFO").upper(),
    format="%(message)s",
    datefmt="[%X]",
    handlers=[RichHandlerColor(rich_tracebacks=True, markup=True)],
)

logger = logging.getLogger("rich")


async def log_tool_run(tool, func):
    tool_name = type(tool).name
    logger.info(f"Executing tool [green]{tool_name}[/green] with args:\n{tool.input}")
    result = await func()
    logger.info(f"Exiting tool [green]{tool_name}[/green] with result:\n{result}")
    return result
