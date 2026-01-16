package features

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cucumber/godog"
	"github.com/julpayne/mlflow-go-client/pkg/mlflow"
)

type testContext struct {
	client           *mlflow.Client
	experimentID     string
	experimentName   string
	runID            string
	modelName        string
	modelVersion     string
	model            *mlflow.RegisteredModel
	healthStatus     string
	serverVersion    string
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
	ctx.healthStatus = ""
	ctx.serverVersion = ""
	ctx.lastError = nil
	ctx.lastResponse = nil
}

func (ctx *testContext) cleanup() {
	// Clean up created resources in reverse order
	for i := len(ctx.createdResources) - 1; i >= 0; i-- {
		resource := ctx.createdResources[i]
		switch resource.Type {
		case "experiment":
			_ = ctx.client.DeleteExperiment(resource.ID)
		case "run":
			_ = ctx.client.DeleteRun(resource.ID)
		case "model":
			_ = ctx.client.DeleteRegisteredModel(resource.Name)
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

	// Server health and version steps
	ctx.Step(`^I check the server health$`, tc.checkServerHealth)
	ctx.Step(`^the health status should be "([^"]*)"$`, tc.healthStatusShouldBe)
	ctx.Step(`^I check the server version$`, tc.checkServerVersion)
	ctx.Step(`^the version should match "([^"]*)"$`, tc.versionShouldMatch)
	ctx.Step(`^the version should not be empty$`, tc.versionShouldNotBeEmpty)

	// Experiment steps
	ctx.Step(`^I create an experiment named "([^"]*)"$`, tc.createExperiment)
	ctx.Step(`^the experiment should be created successfully$`, tc.experimentCreatedSuccessfully)
	ctx.Step(`^the experiment should have the name "([^"]*)"$`, tc.experimentHasName)
	ctx.Step(`^an experiment named "([^"]*)" exists$`, tc.experimentExists)
	ctx.Step(`^an experiment with a unique name exists$`, tc.experimentUniqueNameExists)
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
	ctx.Step(`^multiple runs exist in the experiment with tag "([^"]*)" equals "([^"]*)"$`, tc.multipleRunsExistWithTag)
	ctx.Step(`^I search for runs with filter "([^"]*)"$`, tc.searchRuns)
	ctx.Step(`^I should get a non-empty list of runs$`, tc.getListOfRuns)
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

	// Other steps
	ctx.Step(`^fix this step$`, tc.fixThisStep)
}

func debugLog(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.Println(msg)
}

// Server setup steps
func (tc *testContext) serverIsRunning(url string) error {
	// If the cliebnt is already connected then no need to check again
	if tc.clientConnected() == nil {
		return nil
	}
	// For testing, we'll use an existing server if MLFLOW_TEST_URL
	testURL := os.Getenv("MLFLOW_TEST_URL")
	if testURL != "" {
		client := mlflow.NewClient(testURL)
		// now chekc the server is running and has a version that we can handle
		if err := client.CheckServer(); err != nil {
			return err
		}
		tc.client = client
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

func (tc *testContext) fixThisStep() error {
	debugLog("TODO: fix this step")
	return godog.ErrSkip
}
