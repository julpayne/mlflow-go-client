package main

import (
	"fmt"
	"log"
	"time"

	"github.com/julpayne/mlflow-go-client"
)

func main() {
	// Create a new MLflow client
	client := mlflow.NewClient("http://localhost:5000")

	// Example 1: Create an experiment
	fmt.Println("Creating experiment...")
	expReq := mlflow.CreateExperimentRequest{
		Name: "go-client-example",
		Tags: []mlflow.ExperimentTag{
			{Key: "language", Value: "go"},
			{Key: "example", Value: "true"},
		},
	}

	expResp, err := client.CreateExperiment(expReq)
	if err != nil {
		log.Fatalf("Failed to create experiment: %v", err)
	}
	fmt.Printf("Created experiment with ID: %s\n", expResp.ExperimentID)

	// Example 2: Create a run
	fmt.Println("\nCreating run...")
	runReq := mlflow.CreateRunRequest{
		ExperimentID: expResp.ExperimentID,
		RunName:      "example-run",
		Tags: []mlflow.RunTag{
			{Key: "version", Value: "1.0"},
		},
	}

	runResp, err := client.CreateRun(runReq)
	if err != nil {
		log.Fatalf("Failed to create run: %v", err)
	}
	runID := runResp.Run.Info.RunID
	fmt.Printf("Created run with ID: %s\n", runID)

	// Example 3: Log metrics
	fmt.Println("\nLogging metrics...")
	for i := 0; i < 5; i++ {
		err := client.LogMetric(mlflow.LogMetricRequest{
			RunID:  runID,
			Key:    "accuracy",
			Value:  0.9 + float64(i)*0.01,
			Step:   int64(i),
			Timestamp: time.Now().UnixMilli(),
		})
		if err != nil {
			log.Printf("Failed to log metric: %v", err)
		}

		err = client.LogMetric(mlflow.LogMetricRequest{
			RunID:  runID,
			Key:    "loss",
			Value:  0.1 - float64(i)*0.01,
			Step:   int64(i),
			Timestamp: time.Now().UnixMilli(),
		})
		if err != nil {
			log.Printf("Failed to log metric: %v", err)
		}
	}
	fmt.Println("Logged metrics successfully")

	// Example 4: Log parameters
	fmt.Println("\nLogging parameters...")
	params := []mlflow.LogParamRequest{
		{RunID: runID, Key: "learning_rate", Value: "0.01"},
		{RunID: runID, Key: "epochs", Value: "100"},
		{RunID: runID, Key: "batch_size", Value: "32"},
	}

	for _, param := range params {
		err := client.LogParam(param)
		if err != nil {
			log.Printf("Failed to log parameter: %v", err)
		}
	}
	fmt.Println("Logged parameters successfully")

	// Example 5: Get the run
	fmt.Println("\nRetrieving run...")
	run, err := client.GetRun(runID)
	if err != nil {
		log.Fatalf("Failed to get run: %v", err)
	}
	fmt.Printf("Run Status: %s\n", run.Run.Info.Status)
	fmt.Printf("Metrics: %d\n", len(run.Run.Data.Metrics))
	fmt.Printf("Params: %d\n", len(run.Run.Data.Params))

	// Example 6: Search runs
	fmt.Println("\nSearching runs...")
	searchReq := mlflow.SearchRunsRequest{
		ExperimentIDs: []string{expResp.ExperimentID},
		MaxResults:    10,
		OrderBy:       []string{"metrics.accuracy DESC"},
	}

	searchResp, err := client.SearchRuns(searchReq)
	if err != nil {
		log.Fatalf("Failed to search runs: %v", err)
	}
	fmt.Printf("Found %d runs\n", len(searchResp.Runs))

	// Example 7: Update run status
	fmt.Println("\nUpdating run status...")
	updateResp, err := client.UpdateRun(mlflow.UpdateRunRequest{
		RunID:   runID,
		Status:  "FINISHED",
		EndTime: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Fatalf("Failed to update run: %v", err)
	}
	fmt.Printf("Run updated. Status: %s\n", updateResp.RunInfo.Status)

	// Example 8: List experiments
	fmt.Println("\nListing experiments...")
	experiments, err := client.ListExperiments(10, "")
	if err != nil {
		log.Fatalf("Failed to list experiments: %v", err)
	}
	fmt.Printf("Found %d experiments\n", len(experiments.Experiments))
	for _, exp := range experiments.Experiments {
		fmt.Printf("  - %s (ID: %s)\n", exp.Name, exp.ExperimentID)
	}

	fmt.Println("\nExample completed successfully!")
}
