# MLflow Go Client

A comprehensive Go client library for interacting with the MLflow REST API.

## Installation

```bash
go get github.com/julpayne/mlflow-go-client
```

## Usage

### Basic Setup

```go
package main

import (
    "fmt"
    "github.com/julpayne/mlflow-go-client"
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

#### Update and Delete Experiments

```go
// Update experiment name
err := client.UpdateExperiment("experiment-id", "new-name")

// Set experiment tag
err := client.SetExperimentTag("experiment-id", "key", "value")

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

#### Update and Transition Model Versions

```go
// Update model version
err := client.UpdateModelVersion("my-model", "1", "Updated description", "Production")

// Transition model version stage
version, err := client.TransitionModelVersionStage(
    "my-model", 
    "1", 
    "Production",
    "true", // archive existing versions
)

// Delete model version
err := client.DeleteModelVersion("my-model", "1")

// Delete registered model
err := client.DeleteRegisteredModel("my-model")
```

## API Coverage

This client supports the following MLflow API endpoints:

### Experiments
- ✅ Create experiment
- ✅ Get experiment (by ID or name)
- ✅ List experiments
- ✅ Update experiment
- ✅ Delete experiment
- ✅ Restore experiment
- ✅ Set experiment tag

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

### Models
- ✅ Create registered model
- ✅ Get registered model
- ✅ List registered models
- ✅ Update registered model
- ✅ Delete registered model
- ✅ Create model version
- ✅ Get model version
- ✅ List model versions
- ✅ Update model version
- ✅ Delete model version
- ✅ Transition model version stage

## Error Handling

All methods return errors that can be checked:

```go
experiment, err := client.GetExperiment("id")
if err != nil {
    // Handle error
    fmt.Printf("Error: %v\n", err)
    return
}
```

## Authentication

The client supports Bearer token authentication:

```go
client.SetAuthToken("your-api-token")
```

## License

MIT
