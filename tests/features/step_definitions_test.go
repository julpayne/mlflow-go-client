package features

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cucumber/godog"
	"github.com/julpayne/mlflow-go-client"
)

type testContext struct {
	client           *mlflow.Client
	experimentID     string
	experimentName   string
	runID            string
	modelName        string
	modelVersion     string
	lastError        error
	lastResponse     interface{}
	createdResources []resource
}

type resource struct {
	Type string
	ID   string
	Name string
}

func (ctx *testContext) reset() {
	ctx.experimentID = ""
	ctx.experimentName = ""
	ctx.runID = ""
	ctx.modelName = ""
	ctx.modelVersion = ""
	ctx.lastError = nil
	ctx.lastResponse = nil
}

func (ctx *testContext) cleanup() {
	// Clean up created resources in reverse order
	for i := len(ctx.createdResources) - 1; i >= 0; i-- {
		resource := ctx.createdResources[i]
		switch resource.Type {
		case "experiment":
			ctx.client.DeleteExperiment(resource.ID)
		case "run":
			ctx.client.DeleteRun(resource.ID)
		case "model":
			ctx.client.DeleteRegisteredModel(resource.Name)
		}
	}
	ctx.createdResources = nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := &testContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tc.reset()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		tc.cleanup()
		return ctx, nil
	})

	// Server setup steps
	ctx.Step(`^an MLflow server is running at "([^"]*)"$`, tc.serverIsRunning)
	ctx.Step(`^I have an MLflow client connected to the server$`, tc.clientConnected)

	// Experiment steps
	ctx.Step(`^I create an experiment named "([^"]*)"$`, tc.createExperiment)
	ctx.Step(`^the experiment should be created successfully$`, tc.experimentCreatedSuccessfully)
	ctx.Step(`^the experiment should have the name "([^"]*)"$`, tc.experimentHasName)
	ctx.Step(`^an experiment named "([^"]*)" exists$`, tc.experimentExists)
	ctx.Step(`^I get the experiment by ID$`, tc.getExperimentByID)
	ctx.Step(`^I get the experiment by name "([^"]*)"$`, tc.getExperimentByName)
	ctx.Step(`^the experiment should be returned$`, tc.experimentReturned)
	ctx.Step(`^multiple experiments exist$`, tc.multipleExperimentsExist)
	ctx.Step(`^I list all experiments$`, tc.listExperiments)
	ctx.Step(`^I should get a list of experiments$`, tc.getListOfExperiments)
	ctx.Step(`^the list should contain at least (\d+) experiment$`, tc.listContainsExperiments)
	ctx.Step(`^I search for experiments with filter "([^"]*)"$`, tc.searchExperiments)
	ctx.Step(`^I should find the experiment "([^"]*)"$`, tc.findExperiment)
	ctx.Step(`^I update the experiment name to "([^"]*)"$`, tc.updateExperimentName)
	ctx.Step(`^the experiment name should be "([^"]*)"$`, tc.experimentNameShouldBe)
	ctx.Step(`^I set tag "([^"]*)" with value "([^"]*)" on the experiment$`, tc.setExperimentTag)
	ctx.Step(`^the experiment should have tag "([^"]*)" with value "([^"]*)"$`, tc.experimentHasTag)
	ctx.Step(`^the experiment has tag "([^"]*)" with value "([^"]*)"$`, tc.experimentHasTagSet)
	ctx.Step(`^I delete tag "([^"]*)" from the experiment$`, tc.deleteExperimentTag)
	ctx.Step(`^the experiment should not have tag "([^"]*)"$`, tc.experimentDoesNotHaveTag)
	ctx.Step(`^I delete the experiment$`, tc.deleteExperiment)
	ctx.Step(`^the experiment should be deleted$`, tc.experimentDeleted)
	ctx.Step(`^a deleted experiment named "([^"]*)" exists$`, tc.deletedExperimentExists)
	ctx.Step(`^I restore the experiment$`, tc.restoreExperiment)
	ctx.Step(`^the experiment should be restored$`, tc.experimentRestored)

	// Run steps
	ctx.Step(`^an experiment named "([^"]*)" exists$`, tc.experimentExists)
	ctx.Step(`^I create a run in the experiment$`, tc.createRun)
	ctx.Step(`^the run should be created successfully$`, tc.runCreatedSuccessfully)
	ctx.Step(`^the run should have status "([^"]*)"$`, tc.runHasStatus)
	ctx.Step(`^a run exists in the experiment$`, tc.runExists)
	ctx.Step(`^I get the run by ID$`, tc.getRunByID)
	ctx.Step(`^the run should be returned$`, tc.runReturned)
	ctx.Step(`^the run should have valid metadata$`, tc.runHasValidMetadata)
	ctx.Step(`^I log metric "([^"]*)" with value ([\d.]+) to the run$`, tc.logMetric)
	ctx.Step(`^the metric should be logged successfully$`, tc.metricLoggedSuccessfully)
	ctx.Step(`^I log parameter "([^"]*)" with value "([^"]*)" to the run$`, tc.logParameter)
	ctx.Step(`^the parameter should be logged successfully$`, tc.parameterLoggedSuccessfully)
	ctx.Step(`^I set tag "([^"]*)" with value "([^"]*)" on the run$`, tc.setRunTag)
	ctx.Step(`^the run should have tag "([^"]*)" with value "([^"]*)"$`, tc.runHasTag)
	ctx.Step(`^the run has tag "([^"]*)" with value "([^"]*)"$`, tc.runHasTagSet)
	ctx.Step(`^I delete tag "([^"]*)" from the run$`, tc.deleteRunTag)
	ctx.Step(`^the run should not have tag "([^"]*)"$`, tc.runDoesNotHaveTag)
	ctx.Step(`^I log batch with (\d+) metrics and (\d+) parameters to the run$`, tc.logBatch)
	ctx.Step(`^the batch should be logged successfully$`, tc.batchLoggedSuccessfully)
	ctx.Step(`^the run should have (\d+) metrics$`, tc.runHasMetrics)
	ctx.Step(`^the run should have (\d+) parameters$`, tc.runHasParameters)
	ctx.Step(`^I update the run status to "([^"]*)"$`, tc.updateRunStatus)
	ctx.Step(`^the run status should be "([^"]*)"$`, tc.runStatusShouldBe)
	ctx.Step(`^multiple runs exist in the experiment$`, tc.multipleRunsExist)
	ctx.Step(`^I search for runs with filter "([^"]*)"$`, tc.searchRuns)
	ctx.Step(`^I should get a list of runs$`, tc.getListOfRuns)
	ctx.Step(`^I have logged metric "([^"]*)" multiple times to the run$`, tc.loggedMetricMultipleTimes)
	ctx.Step(`^I get the metric history for "([^"]*)"$`, tc.getMetricHistory)
	ctx.Step(`^I should get multiple metric values$`, tc.getMultipleMetricValues)
	ctx.Step(`^I list artifacts for the run$`, tc.listArtifacts)
	ctx.Step(`^I should get a list of artifacts$`, tc.getListOfArtifacts)
	ctx.Step(`^I delete the run$`, tc.deleteRun)
	ctx.Step(`^the run should be deleted$`, tc.runDeleted)
	ctx.Step(`^a deleted run exists in the experiment$`, tc.deletedRunExists)
	ctx.Step(`^I restore the run$`, tc.restoreRun)
	ctx.Step(`^the run should be restored$`, tc.runRestored)

	// Model steps
	ctx.Step(`^I create a registered model named "([^"]*)"$`, tc.createRegisteredModel)
	ctx.Step(`^the model should be created successfully$`, tc.modelCreatedSuccessfully)
	ctx.Step(`^the model should have the name "([^"]*)"$`, tc.modelHasName)
	ctx.Step(`^a registered model named "([^"]*)" exists$`, tc.registeredModelExists)
	ctx.Step(`^I get the registered model "([^"]*)"$`, tc.getRegisteredModel)
	ctx.Step(`^the model should be returned$`, tc.modelReturned)
	ctx.Step(`^the model name should be "([^"]*)"$`, tc.modelNameShouldBe)
	ctx.Step(`^multiple registered models exist$`, tc.multipleRegisteredModelsExist)
	ctx.Step(`^I list all registered models$`, tc.listRegisteredModels)
	ctx.Step(`^I should get a list of models$`, tc.getListOfModels)
	ctx.Step(`^the list should contain at least (\d+) model$`, tc.listContainsModels)
	ctx.Step(`^I search for models with filter "([^"]*)"$`, tc.searchRegisteredModels)
	ctx.Step(`^I should find the model "([^"]*)"$`, tc.findModel)
	ctx.Step(`^I update the model description to "([^"]*)"$`, tc.updateModelDescription)
	ctx.Step(`^the model description should be "([^"]*)"$`, tc.modelDescriptionShouldBe)
	ctx.Step(`^I rename the model to "([^"]*)"$`, tc.renameModel)
	ctx.Step(`^I set tag "([^"]*)" with value "([^"]*)" on the model$`, tc.setModelTag)
	ctx.Step(`^the model should have tag "([^"]*)" with value "([^"]*)"$`, tc.modelHasTag)
	ctx.Step(`^the model has tag "([^"]*)" with value "([^"]*)"$`, tc.modelHasTagSet)
	ctx.Step(`^I delete tag "([^"]*)" from the model$`, tc.deleteModelTag)
	ctx.Step(`^the model should not have tag "([^"]*)"$`, tc.modelDoesNotHaveTag)
	ctx.Step(`^I create a model version with source "([^"]*)"$`, tc.createModelVersion)
	ctx.Step(`^the model version should be created successfully$`, tc.modelVersionCreatedSuccessfully)
	ctx.Step(`^a model version exists for model "([^"]*)"$`, tc.modelVersionExists)
	ctx.Step(`^I get the model version$`, tc.getModelVersion)
	ctx.Step(`^the model version should be returned$`, tc.modelVersionReturned)
	ctx.Step(`^multiple model versions exist for model "([^"]*)"$`, tc.multipleModelVersionsExist)
	ctx.Step(`^I list model versions for "([^"]*)"$`, tc.listModelVersions)
	ctx.Step(`^I should get a list of versions$`, tc.getListOfVersions)
	ctx.Step(`^the list should contain at least (\d+) version$`, tc.listContainsVersions)
	ctx.Step(`^I search for model versions with filter "([^"]*)"$`, tc.searchModelVersions)
	ctx.Step(`^I should find at least one version$`, tc.findAtLeastOneVersion)
	ctx.Step(`^I get the latest model versions for "([^"]*)"$`, tc.getLatestModelVersions)
	ctx.Step(`^I should get at least one version$`, tc.getAtLeastOneVersion)
	ctx.Step(`^I update the model version description to "([^"]*)"$`, tc.updateModelVersionDescription)
	ctx.Step(`^the model version description should be "([^"]*)"$`, tc.modelVersionDescriptionShouldBe)
	ctx.Step(`^I transition the model version to stage "([^"]*)"$`, tc.transitionModelVersionStage)
	ctx.Step(`^the model version stage should be "([^"]*)"$`, tc.modelVersionStageShouldBe)
	ctx.Step(`^I set tag "([^"]*)" with value "([^"]*)" on the model version$`, tc.setModelVersionTag)
	ctx.Step(`^the model version should have tag "([^"]*)" with value "([^"]*)"$`, tc.modelVersionHasTag)
	ctx.Step(`^the model version has tag "([^"]*)" with value "([^"]*)"$`, tc.modelVersionHasTagSet)
	ctx.Step(`^I delete tag "([^"]*)" from the model version$`, tc.deleteModelVersionTag)
	ctx.Step(`^the model version should not have tag "([^"]*)"$`, tc.modelVersionDoesNotHaveTag)
	ctx.Step(`^I set alias "([^"]*)" pointing to the model version$`, tc.setModelAlias)
	ctx.Step(`^the model version should have alias "([^"]*)"$`, tc.modelVersionHasAlias)
	ctx.Step(`^a model version with alias "([^"]*)" exists for model "([^"]*)"$`, tc.modelVersionWithAliasExists)
	ctx.Step(`^I get the model version by alias "([^"]*)"$`, tc.getModelVersionByAlias)
	ctx.Step(`^I delete alias "([^"]*)" from the model$`, tc.deleteModelAlias)
	ctx.Step(`^the alias should be deleted$`, tc.aliasDeleted)
	ctx.Step(`^I delete the model version$`, tc.deleteModelVersion)
	ctx.Step(`^the model version should be deleted$`, tc.modelVersionDeleted)
	ctx.Step(`^I delete the registered model$`, tc.deleteRegisteredModel)
	ctx.Step(`^the model should be deleted$`, tc.modelDeleted)
}

// Server setup steps
func (tc *testContext) serverIsRunning(url string) error {
	// For testing, we'll use an existing server if MLFLOW_TEST_URL
	testURL := os.Getenv("MLFLOW_TEST_URL")
	if testURL != "" {
		tc.client = mlflow.NewClient(testURL)
		return nil
	}
	return fmt.Errorf("MLFLOW_TEST_URL is not set")
}

func (tc *testContext) clientConnected() error {
	if tc.client == nil {
		return fmt.Errorf("client not initialized")
	}
	return nil
}

// Experiment step implementations
func (tc *testContext) createExperiment(name string) error {
	tc.experimentName = name
	req := mlflow.CreateExperimentRequest{
		Name: name,
	}
	resp, err := tc.client.CreateExperiment(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.experimentID = resp.ExperimentID
	tc.createdResources = append(tc.createdResources, resource{Type: "experiment", ID: resp.ExperimentID, Name: name})
	return nil
}

func (tc *testContext) experimentCreatedSuccessfully() error {
	if tc.experimentID == "" {
		return fmt.Errorf("experiment ID is empty")
	}
	return nil
}

func (tc *testContext) experimentHasName(name string) error {
	exp, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		return err
	}
	if exp.Experiment.Name != name {
		return fmt.Errorf("expected experiment name %s, got %s", name, exp.Experiment.Name)
	}
	return nil
}

func (tc *testContext) experimentExists(name string) error {
	return tc.createExperiment(name)
}

func (tc *testContext) getExperimentByID() error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	_, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) getExperimentByName(name string) error {
	_, err := tc.client.GetExperimentByName(name)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) experimentReturned() error {
	return nil // Already checked in getExperimentByID/getExperimentByName
}

func (tc *testContext) multipleExperimentsExist() error {
	// Create a few experiments
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("multi-exp-%d-%d", time.Now().Unix(), i)
		if err := tc.createExperiment(name); err != nil {
			return err
		}
	}
	return nil
}

func (tc *testContext) listExperiments() error {
	resp, err := tc.client.ListExperiments(100, "")
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getListOfExperiments() error {
	resp, ok := tc.lastResponse.(*mlflow.ListExperimentsResponse)
	if !ok {
		return fmt.Errorf("expected ListExperimentsResponse")
	}
	if resp.Experiments == nil {
		return fmt.Errorf("experiments list is nil")
	}
	return nil
}

func (tc *testContext) listContainsExperiments(count int) error {
	resp, ok := tc.lastResponse.(*mlflow.ListExperimentsResponse)
	if !ok {
		return fmt.Errorf("expected ListExperimentsResponse")
	}
	if len(resp.Experiments) < count {
		return fmt.Errorf("expected at least %d experiments, got %d", count, len(resp.Experiments))
	}
	return nil
}

func (tc *testContext) searchExperiments(filter string) error {
	req := mlflow.SearchExperimentsRequest{
		Filter: filter,
	}
	resp, err := tc.client.SearchExperiments(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) findExperiment(name string) error {
	resp, ok := tc.lastResponse.(*mlflow.SearchExperimentsResponse)
	if !ok {
		return fmt.Errorf("expected SearchExperimentsResponse")
	}
	for _, exp := range resp.Experiments {
		if exp.Name == name {
			return nil
		}
	}
	return fmt.Errorf("experiment %s not found", name)
}

func (tc *testContext) updateExperimentName(newName string) error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	return tc.client.UpdateExperiment(tc.experimentID, newName)
}

func (tc *testContext) experimentNameShouldBe(name string) error {
	exp, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		return err
	}
	if exp.Experiment.Name != name {
		return fmt.Errorf("expected name %s, got %s", name, exp.Experiment.Name)
	}
	return nil
}

func (tc *testContext) setExperimentTag(key, value string) error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	return tc.client.SetExperimentTag(tc.experimentID, key, value)
}

func (tc *testContext) experimentHasTag(key, value string) error {
	exp, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		return err
	}
	for _, tag := range exp.Experiment.Tags {
		if tag.Key == key && tag.Value == value {
			return nil
		}
	}
	return fmt.Errorf("tag %s=%s not found", key, value)
}

func (tc *testContext) experimentHasTagSet(key, value string) error {
	return tc.setExperimentTag(key, value)
}

func (tc *testContext) deleteExperimentTag(key string) error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	return tc.client.DeleteExperimentTag(tc.experimentID, key)
}

func (tc *testContext) experimentDoesNotHaveTag(key string) error {
	exp, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		return err
	}
	for _, tag := range exp.Experiment.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag %s still exists", key)
		}
	}
	return nil
}

func (tc *testContext) deleteExperiment() error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	return tc.client.DeleteExperiment(tc.experimentID)
}

func (tc *testContext) experimentDeleted() error {
	_, err := tc.client.GetExperiment(tc.experimentID)
	if err == nil {
		return fmt.Errorf("experiment still exists")
	}
	// Check if it's an API error with appropriate status
	if _, ok := mlflow.IsAPIError(err); ok {
		// Experiment should not be found or should be deleted
		return nil
	}
	return err
}

func (tc *testContext) deletedExperimentExists(name string) error {
	if err := tc.createExperiment(name); err != nil {
		return err
	}
	return tc.deleteExperiment()
}

func (tc *testContext) restoreExperiment() error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	return tc.client.RestoreExperiment(tc.experimentID)
}

func (tc *testContext) experimentRestored() error {
	exp, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		return err
	}
	if exp.Experiment.LifecycleStage != "active" {
		return fmt.Errorf("experiment not restored, lifecycle stage: %s", exp.Experiment.LifecycleStage)
	}
	return nil
}

// Run step implementations (continuing in next part due to length)
func (tc *testContext) createRun() error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	req := mlflow.CreateRunRequest{
		ExperimentID: tc.experimentID,
		RunName:      fmt.Sprintf("test-run-%d", time.Now().Unix()),
	}
	resp, err := tc.client.CreateRun(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.runID = resp.Run.Info.RunID
	tc.createdResources = append(tc.createdResources, resource{Type: "run", ID: tc.runID})
	return nil
}

func (tc *testContext) runCreatedSuccessfully() error {
	if tc.runID == "" {
		return fmt.Errorf("run ID is empty")
	}
	return nil
}

func (tc *testContext) runHasStatus(status string) error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	if run.Run.Info.Status != status {
		return fmt.Errorf("expected status %s, got %s", status, run.Run.Info.Status)
	}
	return nil
}

func (tc *testContext) runExists() error {
	return tc.createRun()
}

func (tc *testContext) getRunByID() error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	_, err := tc.client.GetRun(tc.runID)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) runReturned() error {
	return nil
}

func (tc *testContext) runHasValidMetadata() error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	if run.Run.Info.RunID == "" {
		return fmt.Errorf("run ID is empty")
	}
	if run.Run.Info.ExperimentID == "" {
		return fmt.Errorf("experiment ID is empty")
	}
	return nil
}

func (tc *testContext) logMetric(key string, value float64) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	req := mlflow.LogMetricRequest{
		RunID:     tc.runID,
		Key:       key,
		Value:     value,
		Step:      1,
		Timestamp: time.Now().UnixMilli(),
	}
	return tc.client.LogMetric(req)
}

func (tc *testContext) metricLoggedSuccessfully() error {
	return nil
}

func (tc *testContext) logParameter(key, value string) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	req := mlflow.LogParamRequest{
		RunID: tc.runID,
		Key:   key,
		Value: value,
	}
	return tc.client.LogParam(req)
}

func (tc *testContext) parameterLoggedSuccessfully() error {
	return nil
}

func (tc *testContext) setRunTag(key, value string) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	req := mlflow.SetTagRequest{
		RunID: tc.runID,
		Key:   key,
		Value: value,
	}
	return tc.client.SetTag(req)
}

func (tc *testContext) runHasTag(key, value string) error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	for _, tag := range run.Run.Data.Tags {
		if tag.Key == key && tag.Value == value {
			return nil
		}
	}
	return fmt.Errorf("tag %s=%s not found", key, value)
}

func (tc *testContext) runHasTagSet(key, value string) error {
	return tc.setRunTag(key, value)
}

func (tc *testContext) deleteRunTag(key string) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	return tc.client.DeleteTag(tc.runID, key)
}

func (tc *testContext) runDoesNotHaveTag(key string) error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	for _, tag := range run.Run.Data.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag %s still exists", key)
		}
	}
	return nil
}

func (tc *testContext) logBatch(metricCount, paramCount int) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	metrics := make([]mlflow.Metric, metricCount)
	for i := 0; i < metricCount; i++ {
		metrics[i] = mlflow.Metric{
			Key:       fmt.Sprintf("metric_%d", i),
			Value:     float64(i) * 0.1,
			Step:      int64(i),
			Timestamp: time.Now().UnixMilli(),
		}
	}
	params := make([]mlflow.Param, paramCount)
	for i := 0; i < paramCount; i++ {
		params[i] = mlflow.Param{
			Key:   fmt.Sprintf("param_%d", i),
			Value: fmt.Sprintf("value_%d", i),
		}
	}
	return tc.client.LogBatch(tc.runID, metrics, params, nil)
}

func (tc *testContext) batchLoggedSuccessfully() error {
	return nil
}

func (tc *testContext) runHasMetrics(count int) error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	if len(run.Run.Data.Metrics) < count {
		return fmt.Errorf("expected at least %d metrics, got %d", count, len(run.Run.Data.Metrics))
	}
	return nil
}

func (tc *testContext) runHasParameters(count int) error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	if len(run.Run.Data.Params) < count {
		return fmt.Errorf("expected at least %d parameters, got %d", count, len(run.Run.Data.Params))
	}
	return nil
}

func (tc *testContext) updateRunStatus(status string) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	req := mlflow.UpdateRunRequest{
		RunID:   tc.runID,
		Status:  status,
		EndTime: time.Now().UnixMilli(),
	}
	_, err := tc.client.UpdateRun(req)
	return err
}

func (tc *testContext) runStatusShouldBe(status string) error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	if run.Run.Info.Status != status {
		return fmt.Errorf("expected status %s, got %s", status, run.Run.Info.Status)
	}
	return nil
}

func (tc *testContext) multipleRunsExist() error {
	for i := 0; i < 3; i++ {
		if err := tc.createRun(); err != nil {
			return err
		}
	}
	return nil
}

func (tc *testContext) searchRuns(filter string) error {
	if tc.experimentID == "" {
		return fmt.Errorf("no experiment ID set")
	}
	req := mlflow.SearchRunsRequest{
		ExperimentIDs: []string{tc.experimentID},
		Filter:        filter,
	}
	resp, err := tc.client.SearchRuns(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getListOfRuns() error {
	resp, ok := tc.lastResponse.(*mlflow.SearchRunsResponse)
	if !ok {
		return fmt.Errorf("expected SearchRunsResponse")
	}
	if resp.Runs == nil {
		return fmt.Errorf("runs list is nil")
	}
	return nil
}

func (tc *testContext) loggedMetricMultipleTimes(metricKey string) error {
	for i := 0; i < 5; i++ {
		if err := tc.logMetric(metricKey, float64(i)*0.1); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	}
	return nil
}

func (tc *testContext) getMetricHistory(metricKey string) error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	req := mlflow.GetMetricHistoryRequest{
		RunUUID:   tc.runID,
		MetricKey: metricKey,
	}
	resp, err := tc.client.GetMetricHistory(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getMultipleMetricValues() error {
	resp, ok := tc.lastResponse.(*mlflow.GetMetricHistoryResponse)
	if !ok {
		return fmt.Errorf("expected GetMetricHistoryResponse")
	}
	if len(resp.Metrics) < 2 {
		return fmt.Errorf("expected multiple metric values, got %d", len(resp.Metrics))
	}
	return nil
}

func (tc *testContext) listArtifacts() error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	resp, err := tc.client.ListArtifacts(tc.runID, "", "")
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getListOfArtifacts() error {
	resp, ok := tc.lastResponse.(*mlflow.ListArtifactsResponse)
	if !ok {
		return fmt.Errorf("expected ListArtifactsResponse")
	}
	if resp.Files == nil {
		return fmt.Errorf("files list is nil")
	}
	return nil
}

func (tc *testContext) deleteRun() error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	return tc.client.DeleteRun(tc.runID)
}

func (tc *testContext) runDeleted() error {
	_, err := tc.client.GetRun(tc.runID)
	if err == nil {
		return fmt.Errorf("run still exists")
	}
	return nil
}

func (tc *testContext) deletedRunExists() error {
	if err := tc.createRun(); err != nil {
		return err
	}
	return tc.deleteRun()
}

func (tc *testContext) restoreRun() error {
	if tc.runID == "" {
		return fmt.Errorf("no run ID set")
	}
	return tc.client.RestoreRun(tc.runID)
}

func (tc *testContext) runRestored() error {
	run, err := tc.client.GetRun(tc.runID)
	if err != nil {
		return err
	}
	if run.Run.Info.LifecycleStage != "active" {
		return fmt.Errorf("run not restored, lifecycle stage: %s", run.Run.Info.LifecycleStage)
	}
	return nil
}

// Model step implementations
func (tc *testContext) createRegisteredModel(name string) error {
	tc.modelName = name
	req := mlflow.CreateRegisteredModelRequest{
		Name: name,
	}
	_, err := tc.client.CreateRegisteredModel(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.createdResources = append(tc.createdResources, resource{Type: "model", Name: name})
	return nil
}

func (tc *testContext) modelCreatedSuccessfully() error {
	if tc.modelName == "" {
		return fmt.Errorf("model name is empty")
	}
	return nil
}

func (tc *testContext) modelHasName(name string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	if model.RegisteredModel.Name != name {
		return fmt.Errorf("expected model name %s, got %s", name, model.RegisteredModel.Name)
	}
	return nil
}

func (tc *testContext) registeredModelExists(name string) error {
	return tc.createRegisteredModel(name)
}

func (tc *testContext) getRegisteredModel(name string) error {
	tc.modelName = name
	_, err := tc.client.GetRegisteredModel(name)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) modelReturned() error {
	return nil
}

func (tc *testContext) modelNameShouldBe(name string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	if model.RegisteredModel.Name != name {
		return fmt.Errorf("expected name %s, got %s", name, model.RegisteredModel.Name)
	}
	return nil
}

func (tc *testContext) multipleRegisteredModelsExist() error {
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("multi-model-%d-%d", time.Now().Unix(), i)
		if err := tc.createRegisteredModel(name); err != nil {
			return err
		}
	}
	return nil
}

func (tc *testContext) listRegisteredModels() error {
	resp, err := tc.client.ListRegisteredModels(100, "")
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getListOfModels() error {
	resp, ok := tc.lastResponse.(*mlflow.ListRegisteredModelsResponse)
	if !ok {
		return fmt.Errorf("expected ListRegisteredModelsResponse")
	}
	if resp.RegisteredModels == nil {
		return fmt.Errorf("models list is nil")
	}
	return nil
}

func (tc *testContext) listContainsModels(count int) error {
	resp, ok := tc.lastResponse.(*mlflow.ListRegisteredModelsResponse)
	if !ok {
		return fmt.Errorf("expected ListRegisteredModelsResponse")
	}
	if len(resp.RegisteredModels) < count {
		return fmt.Errorf("expected at least %d models, got %d", count, len(resp.RegisteredModels))
	}
	return nil
}

func (tc *testContext) searchRegisteredModels(filter string) error {
	req := mlflow.SearchRegisteredModelsRequest{
		Filter: filter,
	}
	resp, err := tc.client.SearchRegisteredModels(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) findModel(name string) error {
	resp, ok := tc.lastResponse.(*mlflow.SearchRegisteredModelsResponse)
	if !ok {
		return fmt.Errorf("expected SearchRegisteredModelsResponse")
	}
	for _, model := range resp.RegisteredModels {
		if model.Name == name {
			return nil
		}
	}
	return fmt.Errorf("model %s not found", name)
}

func (tc *testContext) updateModelDescription(description string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	return tc.client.UpdateRegisteredModel(tc.modelName, description)
}

func (tc *testContext) modelDescriptionShouldBe(description string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	if model.RegisteredModel.Description != description {
		return fmt.Errorf("expected description %s, got %s", description, model.RegisteredModel.Description)
	}
	return nil
}

func (tc *testContext) renameModel(newName string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.RenameRegisteredModelRequest{
		Name:    tc.modelName,
		NewName: newName,
	}
	_, err := tc.client.RenameRegisteredModel(req)
	if err != nil {
		return err
	}
	tc.modelName = newName
	return nil
}

func (tc *testContext) setModelTag(key, value string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.SetRegisteredModelTagRequest{
		Name:  tc.modelName,
		Key:   key,
		Value: value,
	}
	return tc.client.SetRegisteredModelTag(req)
}

func (tc *testContext) modelHasTag(key, value string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	for _, tag := range model.RegisteredModel.Tags {
		if tag.Key == key && tag.Value == value {
			return nil
		}
	}
	return fmt.Errorf("tag %s=%s not found", key, value)
}

func (tc *testContext) modelHasTagSet(key, value string) error {
	return tc.setModelTag(key, value)
}

func (tc *testContext) deleteModelTag(key string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.DeleteRegisteredModelTagRequest{
		Name: tc.modelName,
		Key:  key,
	}
	return tc.client.DeleteRegisteredModelTag(req)
}

func (tc *testContext) modelDoesNotHaveTag(key string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	for _, tag := range model.RegisteredModel.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag %s still exists", key)
		}
	}
	return nil
}

func (tc *testContext) createModelVersion(source string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.CreateModelVersionRequest{
		Name:   tc.modelName,
		Source: source,
	}
	resp, err := tc.client.CreateModelVersion(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.modelVersion = resp.ModelVersion.Version
	return nil
}

func (tc *testContext) modelVersionCreatedSuccessfully() error {
	if tc.modelVersion == "" {
		return fmt.Errorf("model version is empty")
	}
	return nil
}

func (tc *testContext) modelVersionExists(modelName string) error {
	if err := tc.createRegisteredModel(modelName); err != nil {
		return err
	}
	return tc.createModelVersion("runs:/test-run/model")
}

func (tc *testContext) getModelVersion() error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	_, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) modelVersionReturned() error {
	return nil
}

func (tc *testContext) multipleModelVersionsExist(modelName string) error {
	if err := tc.createRegisteredModel(modelName); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if err := tc.createModelVersion(fmt.Sprintf("runs:/test-run-%d/model", i)); err != nil {
			return err
		}
	}
	return nil
}

func (tc *testContext) listModelVersions(modelName string) error {
	tc.modelName = modelName
	resp, err := tc.client.ListModelVersions(modelName, 100, "")
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getListOfVersions() error {
	resp, ok := tc.lastResponse.(*mlflow.ListModelVersionsResponse)
	if !ok {
		return fmt.Errorf("expected ListModelVersionsResponse")
	}
	if resp.ModelVersions == nil {
		return fmt.Errorf("versions list is nil")
	}
	return nil
}

func (tc *testContext) listContainsVersions(count int) error {
	resp, ok := tc.lastResponse.(*mlflow.ListModelVersionsResponse)
	if !ok {
		return fmt.Errorf("expected ListModelVersionsResponse")
	}
	if len(resp.ModelVersions) < count {
		return fmt.Errorf("expected at least %d versions, got %d", count, len(resp.ModelVersions))
	}
	return nil
}

func (tc *testContext) searchModelVersions(filter string) error {
	req := mlflow.SearchModelVersionsRequest{
		Filter: filter,
	}
	resp, err := tc.client.SearchModelVersions(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) findAtLeastOneVersion() error {
	resp, ok := tc.lastResponse.(*mlflow.SearchModelVersionsResponse)
	if !ok {
		return fmt.Errorf("expected SearchModelVersionsResponse")
	}
	if len(resp.ModelVersions) == 0 {
		return fmt.Errorf("no versions found")
	}
	return nil
}

func (tc *testContext) getLatestModelVersions(modelName string) error {
	tc.modelName = modelName
	req := mlflow.GetLatestModelVersionsRequest{
		Name: modelName,
	}
	resp, err := tc.client.GetLatestModelVersions(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getAtLeastOneVersion() error {
	resp, ok := tc.lastResponse.(*mlflow.GetLatestModelVersionsResponse)
	if !ok {
		return fmt.Errorf("expected GetLatestModelVersionsResponse")
	}
	if len(resp.ModelVersions) == 0 {
		return fmt.Errorf("no versions found")
	}
	return nil
}

func (tc *testContext) updateModelVersionDescription(description string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	return tc.client.UpdateModelVersion(tc.modelName, tc.modelVersion, description, "")
}

func (tc *testContext) modelVersionDescriptionShouldBe(description string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	if version.ModelVersion.Description != description {
		return fmt.Errorf("expected description %s, got %s", description, version.ModelVersion.Description)
	}
	return nil
}

func (tc *testContext) transitionModelVersionStage(stage string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	_, err := tc.client.TransitionModelVersionStage(tc.modelName, tc.modelVersion, stage, "")
	if err != nil {
		return err
	}
	return nil
}

func (tc *testContext) modelVersionStageShouldBe(stage string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	if version.ModelVersion.CurrentStage != stage {
		return fmt.Errorf("expected stage %s, got %s", stage, version.ModelVersion.CurrentStage)
	}
	return nil
}

func (tc *testContext) setModelVersionTag(key, value string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	req := mlflow.SetModelVersionTagRequest{
		Name:    tc.modelName,
		Version: tc.modelVersion,
		Key:     key,
		Value:   value,
	}
	return tc.client.SetModelVersionTag(req)
}

func (tc *testContext) modelVersionHasTag(key, value string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	for _, tag := range version.ModelVersion.Tags {
		if tag.Key == key && tag.Value == value {
			return nil
		}
	}
	return fmt.Errorf("tag %s=%s not found", key, value)
}

func (tc *testContext) modelVersionHasTagSet(key, value string) error {
	return tc.setModelVersionTag(key, value)
}

func (tc *testContext) deleteModelVersionTag(key string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	req := mlflow.DeleteModelVersionTagRequest{
		Name:    tc.modelName,
		Version: tc.modelVersion,
		Key:     key,
	}
	return tc.client.DeleteModelVersionTag(req)
}

func (tc *testContext) modelVersionDoesNotHaveTag(key string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	for _, tag := range version.ModelVersion.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag %s still exists", key)
		}
	}
	return nil
}

func (tc *testContext) setModelAlias(alias string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	req := mlflow.SetRegisteredModelAliasRequest{
		Name:    tc.modelName,
		Alias:   alias,
		Version: tc.modelVersion,
	}
	return tc.client.SetRegisteredModelAlias(req)
}

func (tc *testContext) modelVersionHasAlias(alias string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	for _, a := range version.ModelVersion.Aliases {
		if a == alias {
			return nil
		}
	}
	return fmt.Errorf("alias %s not found", alias)
}

func (tc *testContext) modelVersionWithAliasExists(alias, modelName string) error {
	if err := tc.createRegisteredModel(modelName); err != nil {
		return err
	}
	if err := tc.createModelVersion("runs:/test-run/model"); err != nil {
		return err
	}
	return tc.setModelAlias(alias)
}

func (tc *testContext) getModelVersionByAlias(alias string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.GetModelVersionByAliasRequest{
		Name:  tc.modelName,
		Alias: alias,
	}
	resp, err := tc.client.GetModelVersionByAlias(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.modelVersion = resp.ModelVersion.Version
	return nil
}

func (tc *testContext) deleteModelAlias(alias string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.DeleteRegisteredModelAliasRequest{
		Name:  tc.modelName,
		Alias: alias,
	}
	return tc.client.DeleteRegisteredModelAlias(req)
}

func (tc *testContext) aliasDeleted() error {
	// Verify alias is deleted by trying to get it (should fail)
	req := mlflow.GetModelVersionByAliasRequest{
		Name:  tc.modelName,
		Alias: "deleted-alias",
	}
	_, err := tc.client.GetModelVersionByAlias(req)
	if err == nil {
		return fmt.Errorf("alias still exists")
	}
	return nil
}

func (tc *testContext) deleteModelVersion() error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	return tc.client.DeleteModelVersion(tc.modelName, tc.modelVersion)
}

func (tc *testContext) modelVersionDeleted() error {
	_, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err == nil {
		return fmt.Errorf("model version still exists")
	}
	return nil
}

func (tc *testContext) deleteRegisteredModel() error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	return tc.client.DeleteRegisteredModel(tc.modelName)
}

func (tc *testContext) modelDeleted() error {
	_, err := tc.client.GetRegisteredModel(tc.modelName)
	if err == nil {
		return fmt.Errorf("model still exists")
	}
	return nil
}
