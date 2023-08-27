# Salami

<img src="salami-icon.svg" alt="Salami Icon" width="150" height="150">

## Development Setup

See Makefile for available commands. A few things to note:

1. Running tests requires MongoDB to be running on localhost:27017. You can then see the details of all tests by quering `test_results` collection in Mongo.

2. Running tests costs $! This is because it calls OpenAI API. To run individual tests, pass modules via `SELECTED_MODULES` variable:

```
make test SELECTED_MODULES=tests/ial/ecs/service,tests/ial/s3
```
