test:
	poetry run python -m pytest -n 8 ${PATHS}

test_verbose:
	poetry run python -m pytest -s ${PATHS}