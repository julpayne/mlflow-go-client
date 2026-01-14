# MLflow Go Client

A comprehensive Go client library for interacting with the MLflow REST API.

## Installation

```bash
go get github.com/julpayne/mlflow-go-client/pkg/mlflow
```

## Usage

### Basic Setup

```go
package main

import (
    "fmt"
    "github.com/julpayne/mlflow-go-client/pkg/mlflow"
)

func main() {
    // Create a new client
    client := mlflow.NewClient("http://localhost:5000")
    
    // Optional: Set authentication token
    client.SetAuthToken("your-token-here")
    
    // Optional: Set custom timeout
    client.SetTimeout(60 * time.Second)
}
```

### Experiments

#### Create an Experiment

```go
req := mlflow.CreateExperimentRequest{
    Name: "my-experiment",
    Tags: []mlflow.ExperimentTag{
        {Key: "team", Value: "ml-team"},
    },
}

resp, err := client.CreateExperiment(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created experiment with ID: %s\n", resp.ExperimentID)
```

#### Get an Experiment

```go
// By ID
experiment, err := client.GetExperiment("experiment-id")
if err != nil {
    log.Fatal(err)
}

// By name
experiment, err := client.GetExperimentByName("my-experiment")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Experiment: %s\n", experiment.Experiment.Name)
```

#### List Experiments

```go
experiments, err := client.ListExperiments(100, "")
if err != nil {
    log.Fatal(err)
}

for _, exp := range experiments.Experiments {
    fmt.Printf("Experiment: %s (ID: %s)\n", exp.Name, exp.ExperimentID)
}
```

#### Search Experiments

```go
// Search experiments with filters
searchReq := mlflow.SearchExperimentsRequest{
    ViewType:   "ACTIVE_ONLY", // ACTIVE_ONLY, DELETED_ONLY, or ALL
    MaxResults: 100,
    Filter:     "name LIKE '%test%'",
    OrderBy:    []string{"name ASC"},
}

results, err := client.SearchExperiments(searchReq)
if err != nil {
    log.Fatal(err)
}

for _, exp := range results.Experiments {
    fmt.Printf("Experiment: %s\n", exp.Name)
}
```

#### Update and Delete Experiments

```go
// Update experiment name
err := client.UpdateExperiment("experiment-id", "new-name")

// Set experiment tag
err := client.SetExperimentTag("experiment-id", "key", "value")

// Delete experiment tag
err := client.DeleteExperimentTag("experiment-id", "key")

// Delete experiment
err := client.DeleteExperiment("experiment-id")

// Restore deleted experiment
err := client.RestoreExperiment("experiment-id")
```

### Runs

#### Create a Run

```go
req := mlflow.CreateRunRequest{
    ExperimentID: "experiment-id",
    RunName:      "my-run",
    Tags: []mlflow.RunTag{
        {Key: "version", Value: "1.0"},
    },
}

run, err := client.CreateRun(req)
if err != nil {
    log.Fatal(err)
}

runID := run.Run.Info.RunID
fmt.Printf("Created run with ID: %s\n", runID)
```

#### Log Metrics and Parameters

```go
// Log a metric
err := client.LogMetric(mlflow.LogMetricRequest{
    RunID: runID,
    Key:   "accuracy",
    Value: 0.95,
    Step:  1,
})

// Log a parameter
err := client.LogParam(mlflow.LogParamRequest{
    RunID: runID,
    Key:   "learning_rate",
    Value: "0.01",
})

// Set a tag
err := client.SetTag(mlflow.SetTagRequest{
    RunID: runID,
    Key:   "model_type",
    Value: "neural_network",
})
```

#### Log Multiple Metrics/Params at Once

```go
metrics := []mlflow.Metric{
    {Key: "loss", Value: 0.5, Step: 1, Timestamp: time.Now().UnixMilli()},
    {Key: "accuracy", Value: 0.95, Step: 1, Timestamp: time.Now().UnixMilli()},
}

params := []mlflow.Param{
    {Key: "epochs", Value: "100"},
    {Key: "batch_size", Value: "32"},
}

tags := []mlflow.RunTag{
    {Key: "framework", Value: "pytorch"},
}

err := client.LogBatch(runID, metrics, params, tags)
```

#### Get and Search Runs

```go
// Get a run by ID
run, err := client.GetRun(runID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Run status: %s\n", run.Run.Info.Status)

// Search for runs
searchReq := mlflow.SearchRunsRequest{
    ExperimentIDs: []string{"experiment-id"},
    Filter:        "metrics.accuracy > 0.9",
    MaxResults:    10,
    OrderBy:       []string{"metrics.accuracy DESC"},
}

results, err := client.SearchRuns(searchReq)
if err != nil {
    log.Fatal(err)
}

for _, run := range results.Runs {
    fmt.Printf("Run: %s, Accuracy: %f\n", run.Info.RunID, 
        getMetricValue(run.Data.Metrics, "accuracy"))
}
```

#### Log Model and Inputs

```go
// Log a model to a run
err := client.LogModel(mlflow.LogModelRequest{
    RunID:    runID,
    ModelJSON: `{"model": "sklearn", "version": "1.0"}`,
})

// Log inputs (datasets and model inputs)
err := client.LogInputs(mlflow.LogInputsRequest{
    RunID: runID,
    Datasets: []mlflow.Dataset{
        {
            Name:       "training_data",
            Digest:     "abc123",
            SourceType: "LOCAL",
            Source:     "/path/to/data",
        },
    },
})
```

#### Get Metric History and List Artifacts

```go
// Get metric history for a run
history, err := client.GetMetricHistory(mlflow.GetMetricHistoryRequest{
    RunUUID:   runID,
    MetricKey: "accuracy",
    MaxResults: 100,
})
if err != nil {
    log.Fatal(err)
}

for _, metric := range history.Metrics {
    fmt.Printf("Step %d: %f\n", metric.Step, metric.Value)
}

// List artifacts for a run
artifacts, err := client.ListArtifacts(runID, "", "")
if err != nil {
    log.Fatal(err)
}

for _, file := range artifacts.Files {
    fmt.Printf("Artifact: %s (size: %d)\n", file.Path, file.FileSize)
}
```

#### Update and End a Run

```go
// Update run status
err := client.UpdateRun(mlflow.UpdateRunRequest{
    RunID:  runID,
    Status: "FINISHED",
    EndTime: time.Now().UnixMilli(),
})

// Delete a run
err := client.DeleteRun(runID)

// Restore a deleted run
err := client.RestoreRun(runID)
```

### Models

#### Create a Registered Model

```go
req := mlflow.CreateRegisteredModelRequest{
    Name:        "my-model",
    Description: "A machine learning model",
    Tags: []mlflow.RegisteredModelTag{
        {Key: "type", Value: "classification"},
    },
}

model, err := client.CreateRegisteredModel(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created model: %s\n", model.RegisteredModel.Name)
```

#### Create a Model Version

```go
req := mlflow.CreateModelVersionRequest{
    Name:   "my-model",
    Source: "runs:/run-id/model",
    RunID:  "run-id",
    Description: "Version 1.0 of the model",
}

version, err := client.CreateModelVersion(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created version: %s\n", version.ModelVersion.Version)
```

#### Get Models and Versions

```go
// Get registered model
model, err := client.GetRegisteredModel("my-model")

// List all registered models
models, err := client.ListRegisteredModels(100, "")

// Get model version
version, err := client.GetModelVersion("my-model", "1")

// List model versions
versions, err := client.ListModelVersions("my-model", 100, "")
```

#### Search Models and Versions

```go
// Search registered models
searchReq := mlflow.SearchRegisteredModelsRequest{
    Filter:     "name LIKE '%classifier%'",
    MaxResults: 50,
    OrderBy:    []string{"name ASC"},
}

models, err := client.SearchRegisteredModels(searchReq)
if err != nil {
    log.Fatal(err)
}

// Search model versions
versionSearchReq := mlflow.SearchModelVersionsRequest{
    Filter:     "name='my-model' AND version='1'",
    MaxResults: 100,
}

versions, err := client.SearchModelVersions(versionSearchReq)
if err != nil {
    log.Fatal(err)
}
```

#### Update and Transition Model Versions

```go
// Rename registered model
renamed, err := client.RenameRegisteredModel(mlflow.RenameRegisteredModelRequest{
    Name:    "old-name",
    NewName: "new-name",
})

// Get latest model versions
latest, err := client.GetLatestModelVersions(mlflow.GetLatestModelVersionsRequest{
    Name:   "my-model",
    Stages: []string{"Production", "Staging"},
})

// Update model version
err := client.UpdateModelVersion("my-model", "1", "Updated description", "Production")

// Transition model version stage
version, err := client.TransitionModelVersionStage(
    "my-model", 
    "1", 
    "Production",
    "true", // archive existing versions
)

// Get download URIs for model artifacts
uris, err := client.GetDownloadURIs(mlflow.GetDownloadURIsRequest{
    Name:    "my-model",
    Version: "1",
    Paths:   []string{"model.pkl", "requirements.txt"},
})
```

#### Model Tags and Aliases

```go
// Set registered model tag
err := client.SetRegisteredModelTag(mlflow.SetRegisteredModelTagRequest{
    Name:  "my-model",
    Key:   "team",
    Value: "ml-team",
})

// Set model version tag
err := client.SetModelVersionTag(mlflow.SetModelVersionTagRequest{
    Name:    "my-model",
    Version: "1",
    Key:     "deployed",
    Value:   "true",
})

// Delete tags
err := client.DeleteRegisteredModelTag(mlflow.DeleteRegisteredModelTagRequest{
    Name: "my-model",
    Key:  "team",
})

err := client.DeleteModelVersionTag(mlflow.DeleteModelVersionTagRequest{
    Name:    "my-model",
    Version: "1",
    Key:     "deployed",
})

// Set model alias
err := client.SetRegisteredModelAlias(mlflow.SetRegisteredModelAliasRequest{
    Name:    "my-model",
    Alias:   "production",
    Version: "1",
})

// Get model version by alias
version, err := client.GetModelVersionByAlias(mlflow.GetModelVersionByAliasRequest{
    Name:  "my-model",
    Alias: "production",
})

// Delete alias
err := client.DeleteRegisteredModelAlias(mlflow.DeleteRegisteredModelAliasRequest{
    Name:  "my-model",
    Alias: "production",
})

// Delete model version
err := client.DeleteModelVersion("my-model", "1")

// Delete registered model
err := client.DeleteRegisteredModel("my-model")
```

## API Coverage

This client supports the following MLflow API endpoints:

### Experiments
- ✅ Create experiment
- ✅ Get experiment (by ID)
- ✅ Get experiment (by name)
- ✅ List experiments
- ✅ Search experiments
- ✅ Update experiment
- ✅ Delete experiment
- ✅ Restore experiment
- ✅ Set experiment tag
- ✅ Delete experiment tag

### Runs
- ✅ Create run
- ✅ Get run
- ✅ Search runs
- ✅ Update run
- ✅ Delete run
- ✅ Restore run
- ✅ Log metric
- ✅ Log parameter
- ✅ Set tag
- ✅ Delete tag
- ✅ Log batch (multiple metrics/params/tags)
- ✅ Log model
- ✅ Log inputs (datasets and model inputs)
- ✅ Get metric history
- ✅ List artifacts

### Models
- ✅ Create registered model
- ✅ Get registered model
- ✅ List registered models
- ✅ Search registered models
- ✅ Update registered model
- ✅ Rename registered model
- ✅ Delete registered model
- ✅ Create model version
- ✅ Get model version
- ✅ List model versions
- ✅ Search model versions
- ✅ Get latest model versions
- ✅ Update model version
- ✅ Delete model version
- ✅ Transition model version stage
- ✅ Get download URIs for model version artifacts
- ✅ Set registered model tag
- ✅ Set model version tag
- ✅ Delete registered model tag
- ✅ Delete model version tag
- ✅ Set registered model alias
- ✅ Delete registered model alias
- ✅ Get model version by alias

## Error Handling

All methods return errors that can be checked. The client uses a custom `APIError` type that provides detailed information about API errors:

```go
experiment, err := client.GetExperiment("id")
if err != nil {
    // Check if it's an APIError to access detailed information
    if apiErr, ok := mlflow.IsAPIError(err); ok {
        fmt.Printf("Status Code: %d\n", apiErr.GetStatusCode())
        fmt.Printf("Error Code: %s\n", apiErr.GetErrorCode())
        fmt.Printf("Message: %s\n", apiErr.GetMessage())
        fmt.Printf("Response Body: %s\n", apiErr.GetResponseBodyString())
    } else {
        // Handle other types of errors (network, etc.)
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

The `APIError` type provides the following methods:
- `GetStatusCode()` - Returns the HTTP status code
- `GetErrorCode()` - Returns the MLflow error code (if available)
- `GetMessage()` - Returns the error message
- `GetResponseBody()` - Returns the raw response body as bytes
- `GetResponseBodyString()` - Returns the response body as a string

## Authentication

The client supports Bearer token authentication:

```go
client.SetAuthToken("your-api-token")
```

## Running MLflow Server Locally

This repository includes scripts and Makefile targets to easily download and run the MLflow server locally for testing and development.

### Quick Start with Makefile

1. **Install MLflow:**
   ```bash
   make install-mlflow
   ```

2. **Start the server:**
   ```bash
   make run-mlflow
   ```

3. **Use with the Go client:**
   ```go
   client := mlflow.NewClient("http://localhost:5000")
   ```

The server will be available at `http://localhost:5000` by default.

### Custom Configuration with Makefile

You can customize the server settings using environment variables:

```bash
# Custom host and port
MLFLOW_PORT=8080 make run-mlflow

# Use PostgreSQL backend
MLFLOW_BACKEND_STORE_URI=postgresql://user:password@localhost/mlflow \
make run-mlflow

# Stop the server
make stop-mlflow
```

### Using Scripts Directly

Alternatively, you can use the scripts directly:

```bash
# Install MLflow
./scripts/download_mlflow.sh

# Start the server
./scripts/run_mlflow.sh

# With custom configuration
MLFLOW_HOST=0.0.0.0 MLFLOW_PORT=8080 ./scripts/run_mlflow.sh
```

### Available Makefile Targets

- `make install-mlflow` - Download and install MLflow server
- `make run-mlflow` - Run MLflow server locally
- `make stop-mlflow` - Stop MLflow server (if running)
- `make build` - Build the Go client library
- `make test` - Run tests
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make lint` - Run golangci-lint (if installed)
- `make example` - Build and run the example
- `make clean` - Clean build artifacts
- `make help` - Show all available targets

For more details, see [scripts/README.md](scripts/README.md).

## License

MIT
