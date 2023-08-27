install:
	poetry install

test:
	SELECTED_MODULES=${SELECTED_MODULES} poetry run python -m pytest -n auto
