# Godog BDD Tests

This directory contains BDD (Behavior-Driven Development) tests for the MLflow Go client using [godog](https://github.com/cucumber/godog).

## Overview

The tests are written in Gherkin format (`.feature` files) and cover:
- **Experiments API** - Creating, getting, listing, searching, updating, and deleting experiments
- **Runs API** - Creating runs, logging metrics/parameters, searching, and managing run lifecycle
- **Models API** - Managing registered models, model versions, tags, and aliases

## Running Tests

### Prerequisites

1. **Install MLflow server:**
   ```bash
   make install-mlflow
   ```

2. **Start MLflow server** (in a separate terminal):
   ```bash
   make run-mlflow
   ```
   
   Or use the test server automatically:
   ```bash
   make test-godog-server
   ```

### Running Tests

**Option 1: Use existing MLflow server**
```bash
# Set the server URL (optional, defaults to http://localhost:5000)
export MLFLOW_TEST_URL=http://localhost:5000

# Run tests
make test-godog
```

**Option 2: Automatic test server management**
```bash
# This will start a test server on port 5000 automatically
make test-godog-server
```

**Option 3: Run directly with go test**
```bash
cd tests
go test -v
```

## Test Configuration

You can configure the tests using environment variables:

- `MLFLOW_TEST_URL` - URL of existing MLflow server (e.g., `http://localhost:5000`)
- `MLFLOW_TEST_HOST` - Host for test server (default: `127.0.0.1`)
- `MLFLOW_TEST_PORT` - Port for test server (default: `5000`)
- `MLFLOW_TEST_BACKEND_STORE_URI` - Backend store URI for test server
- `MLFLOW_TEST_ARTIFACT_ROOT` - Artifact root for test server

## Test Structure

```
tests/features/
├── experiments.feature      # Experiment API tests
├── runs.feature            # Run API tests
├── models.feature          # Model API tests
├── step_definitions_test.go # Step definitions (Gherkin -> Go)
└── features_test.go        # Test runner
```

## Writing New Tests

1. **Add scenarios to feature files** (`.feature` files)
2. **Implement step definitions** in `step_definitions_test.go`
3. **Run tests** to verify

Example feature:
```gherkin
Scenario: Create a new experiment
  When I create an experiment named "my-experiment"
  Then the experiment should be created successfully
```

## Test Cleanup

Tests automatically clean up created resources (experiments, runs, models) after each scenario. The test server is also cleaned up automatically when using `test-godog-server`.

## Troubleshooting

### MLflow server not found
```bash
make install-mlflow
```

### Port already in use
Set a different port:
```bash
MLFLOW_TEST_PORT=5002 make test-godog-server
```

### Tests fail with connection errors
Make sure the MLflow server is running:
```bash
make run-mlflow
```

Or check the server URL:
```bash
export MLFLOW_TEST_URL=http://localhost:5000
make test-godog
```
