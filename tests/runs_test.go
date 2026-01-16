package features

import (
	"fmt"
	"time"

	"github.com/julpayne/mlflow-go-client/pkg/mlflow"
)

// Run step implementations

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
