package features

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/julpayne/mlflow-go-client/pkg/mlflow"
)

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
	resp, err := tc.client.GetExperimentByName(name)
	if err != nil {
		debugLog("Error whilst getting experiment by name %s: %s", name, err.Error())
		return tc.createExperiment(name)
	}
	if resp.Experiment.LifecycleStage != "active" {
		return fmt.Errorf("experiment %s is not active", name)
	}
	if resp.Experiment.Name != name {
		return fmt.Errorf("expected experiment name %s, got %s", name, resp.Experiment.Name)
	}
	return nil
}

func (tc *testContext) experimentUniqueNameExists() error {
	return tc.experimentExists(fmt.Sprintf("test-experiment-%s", uuid.New().String()))
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
	req := mlflow.SearchExperimentsRequest{
		MaxResults: 200,
	}
	resp, err := tc.client.SearchExperiments(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getListOfExperiments() error {
	resp, ok := tc.lastResponse.(*mlflow.SearchExperimentsResponse)
	if !ok {
		return fmt.Errorf("expected SearchExperimentsResponse")
	}
	if resp.Experiments == nil {
		return fmt.Errorf("experiments list is nil")
	}
	return nil
}

func (tc *testContext) listContainsExperiments(count int) error {
	resp, ok := tc.lastResponse.(*mlflow.SearchExperimentsResponse)
	if !ok {
		return fmt.Errorf("expected SearchExperimentsResponse")
	}
	if len(resp.Experiments) < count {
		return fmt.Errorf("expected at least %d experiments, got %d", count, len(resp.Experiments))
	}
	return nil
}

func (tc *testContext) searchExperiments(filter string) error {
	req := mlflow.SearchExperimentsRequest{
		MaxResults: 200,
		Filter:     filter,
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
	err := tc.client.DeleteExperiment(tc.experimentID)
	if err != nil {
		return err
	}
	debugLog("Experiment %s deleted", tc.experimentID)
	return nil
}

func (tc *testContext) experimentDeleted() error {
	_, err := tc.client.GetExperiment(tc.experimentID)
	if err != nil {
		// Check if it's an API error with appropriate status
		if apiError, ok := mlflow.IsAPIError(err); ok && apiError.GetStatusCode() == 404 {
			// Experiment should not be found or should be deleted
			return nil
		}
		debugLog("Experiment %s still exists: %v", tc.experimentID, err)
		return fmt.Errorf("experiment %s still exists: %v", tc.experimentID, err)
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
