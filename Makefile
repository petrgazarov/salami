install:
	poetry install

test:
	poetry run python -m pytest -n auto ${PATHS}

test_verbose:
	poetry run python -m pytest -s ${PATHS}