import os
import inspect


def open_relative_file(filename):
    caller_file = inspect.stack()[1].filename
    caller_dir = os.path.dirname(os.path.abspath(caller_file))
    file_path = os.path.join(caller_dir, filename)
    with open(file_path, "r") as file:
        return file.read()


def get_full_path(filename):
    caller_file = inspect.stack()[1].filename
    caller_dir = os.path.dirname(os.path.abspath(caller_file))
    return os.path.join(caller_dir, filename)
