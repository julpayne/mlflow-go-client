package mlflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents an MLflow API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewClient creates a new MLflow client
func NewClient(baseURL string) *Client {
	// Ensure baseURL doesn't end with a slash
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAuthToken sets the authentication token for the client
func (c *Client) SetAuthToken(token string) {
	c.AuthToken = token
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

// doRequest performs an HTTP request to the MLflow API
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp ErrorResponse
		apiErr := &APIError{
			StatusCode:   resp.StatusCode,
			ResponseBody: respBody,
		}

		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			apiErr.ErrorCode = errorResp.ErrorCode
			apiErr.Message = errorResp.Message
		} else {
			// If we can't parse the error response, use the raw body as message
			apiErr.Message = string(respBody)
		}

		return nil, apiErr
	}

	return respBody, nil
}

// Experiments API

// CreateExperiment creates a new experiment
func (c *Client) CreateExperiment(req CreateExperimentRequest) (*CreateExperimentResponse, error) {
	endpoint := "/api/2.0/mlflow/experiments/create"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response CreateExperimentResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetExperiment gets an experiment by ID
func (c *Client) GetExperiment(experimentID string) (*GetExperimentResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/experiments/get?experiment_id=%s", url.QueryEscape(experimentID))
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetExperimentResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetExperimentByName gets an experiment by name
func (c *Client) GetExperimentByName(experimentName string) (*GetExperimentResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/experiments/get-by-name?experiment_name=%s", url.QueryEscape(experimentName))
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetExperimentResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListExperiments lists all experiments
func (c *Client) ListExperiments(maxResults int, pageToken string) (*ListExperimentsResponse, error) {
	endpoint := "/api/2.0/mlflow/experiments/list"
	if maxResults > 0 || pageToken != "" {
		params := url.Values{}
		if maxResults > 0 {
			params.Add("max_results", fmt.Sprintf("%d", maxResults))
		}
		if pageToken != "" {
			params.Add("page_token", pageToken)
		}
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ListExperimentsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// DeleteExperiment deletes an experiment
func (c *Client) DeleteExperiment(experimentID string) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/experiments/delete?experiment_id=%s", url.QueryEscape(experimentID))
	_, err := c.doRequest("POST", endpoint, nil)
	return err
}

// RestoreExperiment restores a deleted experiment
func (c *Client) RestoreExperiment(experimentID string) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/experiments/restore?experiment_id=%s", url.QueryEscape(experimentID))
	_, err := c.doRequest("POST", endpoint, nil)
	return err
}

// UpdateExperiment updates an experiment
func (c *Client) UpdateExperiment(experimentID, newName string) error {
	req := map[string]interface{}{
		"experiment_id": experimentID,
		"new_name":      newName,
	}
	endpoint := "/api/2.0/mlflow/experiments/update"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// SetExperimentTag sets a tag on an experiment
func (c *Client) SetExperimentTag(experimentID, key, value string) error {
	req := map[string]string{
		"experiment_id": experimentID,
		"key":           key,
		"value":         value,
	}
	endpoint := "/api/2.0/mlflow/experiments/set-experiment-tag"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// DeleteExperimentTag deletes a tag from an experiment
func (c *Client) DeleteExperimentTag(experimentID, key string) error {
	req := map[string]string{
		"experiment_id": experimentID,
		"key":           key,
	}
	endpoint := "/api/2.0/mlflow/experiments/delete-experiment-tag"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// SearchExperiments searches for experiments
func (c *Client) SearchExperiments(req SearchExperimentsRequest) (*SearchExperimentsResponse, error) {
	endpoint := "/api/2.0/mlflow/experiments/search"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response SearchExperimentsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// Runs API

// CreateRun creates a new run
func (c *Client) CreateRun(req CreateRunRequest) (*CreateRunResponse, error) {
	if req.StartTime == 0 {
		req.StartTime = time.Now().UnixMilli()
	}

	endpoint := "/api/2.0/mlflow/runs/create"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response CreateRunResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetRun gets a run by ID
func (c *Client) GetRun(runID string) (*GetRunResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/runs/get?run_id=%s", url.QueryEscape(runID))
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetRunResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// SearchRuns searches for runs
func (c *Client) SearchRuns(req SearchRunsRequest) (*SearchRunsResponse, error) {
	endpoint := "/api/2.0/mlflow/runs/search"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response SearchRunsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// UpdateRun updates a run
func (c *Client) UpdateRun(req UpdateRunRequest) (*UpdateRunResponse, error) {
	endpoint := "/api/2.0/mlflow/runs/update"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response UpdateRunResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// DeleteRun deletes a run
func (c *Client) DeleteRun(runID string) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/runs/delete?run_id=%s", url.QueryEscape(runID))
	_, err := c.doRequest("POST", endpoint, nil)
	return err
}

// RestoreRun restores a deleted run
func (c *Client) RestoreRun(runID string) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/runs/restore?run_id=%s", url.QueryEscape(runID))
	_, err := c.doRequest("POST", endpoint, nil)
	return err
}

// LogMetric logs a metric to a run
func (c *Client) LogMetric(req LogMetricRequest) error {
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().UnixMilli()
	}
	if req.Step == 0 {
		req.Step = 0
	}

	endpoint := "/api/2.0/mlflow/runs/log-metric"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// LogParam logs a parameter to a run
func (c *Client) LogParam(req LogParamRequest) error {
	endpoint := "/api/2.0/mlflow/runs/log-parameter"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// SetTag sets a tag on a run
func (c *Client) SetTag(req SetTagRequest) error {
	endpoint := "/api/2.0/mlflow/runs/set-tag"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// DeleteTag deletes a tag from a run
func (c *Client) DeleteTag(runID, key string) error {
	req := map[string]string{
		"run_id": runID,
		"key":    key,
	}
	endpoint := "/api/2.0/mlflow/runs/delete-tag"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// LogBatch logs multiple metrics, parameters, and tags in a single request
func (c *Client) LogBatch(runID string, metrics []Metric, params []Param, tags []RunTag) error {
	req := map[string]interface{}{
		"run_id":  runID,
		"metrics": metrics,
		"params":  params,
		"tags":    tags,
	}
	endpoint := "/api/2.0/mlflow/runs/log-batch"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// LogModel logs a model to a run
func (c *Client) LogModel(req LogModelRequest) error {
	endpoint := "/api/2.0/mlflow/runs/log-model"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// LogInputs logs inputs (datasets and/or model inputs) to a run
func (c *Client) LogInputs(req LogInputsRequest) error {
	endpoint := "/api/2.0/mlflow/runs/log-inputs"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// GetMetricHistory gets the history of a metric for a run
func (c *Client) GetMetricHistory(req GetMetricHistoryRequest) (*GetMetricHistoryResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/metrics/get-history?run_uuid=%s&metric_key=%s",
		url.QueryEscape(req.RunUUID), url.QueryEscape(req.MetricKey))
	if req.MaxResults > 0 || req.PageToken != "" {
		params := url.Values{}
		if req.MaxResults > 0 {
			params.Add("max_results", fmt.Sprintf("%d", req.MaxResults))
		}
		if req.PageToken != "" {
			params.Add("page_token", req.PageToken)
		}
		endpoint += "&" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetMetricHistoryResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListArtifacts lists artifacts for a run
func (c *Client) ListArtifacts(runID, path string, pageToken string) (*ListArtifactsResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/artifacts/list?run_id=%s", url.QueryEscape(runID))
	if path != "" {
		endpoint += "&path=" + url.QueryEscape(path)
	}
	if pageToken != "" {
		endpoint += "&page_token=" + url.QueryEscape(pageToken)
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ListArtifactsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// Models API

// CreateRegisteredModel creates a new registered model
func (c *Client) CreateRegisteredModel(req CreateRegisteredModelRequest) (*CreateRegisteredModelResponse, error) {
	endpoint := "/api/2.0/mlflow/registered-models/create"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response CreateRegisteredModelResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetRegisteredModel gets a registered model by name
func (c *Client) GetRegisteredModel(name string) (*GetRegisteredModelResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/registered-models/get?name=%s", url.QueryEscape(name))
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetRegisteredModelResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListRegisteredModels lists all registered models
func (c *Client) ListRegisteredModels(maxResults int, pageToken string) (*ListRegisteredModelsResponse, error) {
	endpoint := "/api/2.0/mlflow/registered-models/list"
	if maxResults > 0 || pageToken != "" {
		params := url.Values{}
		if maxResults > 0 {
			params.Add("max_results", fmt.Sprintf("%d", maxResults))
		}
		if pageToken != "" {
			params.Add("page_token", pageToken)
		}
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ListRegisteredModelsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// UpdateRegisteredModel updates a registered model
func (c *Client) UpdateRegisteredModel(name, description string) error {
	req := map[string]string{
		"name": name,
	}
	if description != "" {
		req["description"] = description
	}
	endpoint := "/api/2.0/mlflow/registered-models/update"
	_, err := c.doRequest("PATCH", endpoint, req)
	return err
}

// DeleteRegisteredModel deletes a registered model
func (c *Client) DeleteRegisteredModel(name string) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/registered-models/delete?name=%s", url.QueryEscape(name))
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// CreateModelVersion creates a new model version
func (c *Client) CreateModelVersion(req CreateModelVersionRequest) (*CreateModelVersionResponse, error) {
	endpoint := "/api/2.0/mlflow/model-versions/create"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response CreateModelVersionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetModelVersion gets a model version
func (c *Client) GetModelVersion(name, version string) (*GetModelVersionResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/model-versions/get?name=%s&version=%s",
		url.QueryEscape(name), url.QueryEscape(version))
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetModelVersionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListModelVersions lists model versions for a registered model
func (c *Client) ListModelVersions(name string, maxResults int, pageToken string) (*ListModelVersionsResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/model-versions/list?name=%s", url.QueryEscape(name))
	if maxResults > 0 || pageToken != "" {
		params := url.Values{}
		if maxResults > 0 {
			params.Add("max_results", fmt.Sprintf("%d", maxResults))
		}
		if pageToken != "" {
			params.Add("page_token", pageToken)
		}
		endpoint += "&" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ListModelVersionsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// UpdateModelVersion updates a model version
func (c *Client) UpdateModelVersion(name, version, description, stage string) error {
	req := map[string]string{
		"name":    name,
		"version": version,
	}
	if description != "" {
		req["description"] = description
	}
	if stage != "" {
		req["stage"] = stage
	}
	endpoint := "/api/2.0/mlflow/model-versions/update"
	_, err := c.doRequest("PATCH", endpoint, req)
	return err
}

// DeleteModelVersion deletes a model version
func (c *Client) DeleteModelVersion(name, version string) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/model-versions/delete?name=%s&version=%s",
		url.QueryEscape(name), url.QueryEscape(version))
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// TransitionModelVersionStage transitions a model version to a new stage
func (c *Client) TransitionModelVersionStage(name, version, stage, archiveExistingVersions string) (*GetModelVersionResponse, error) {
	req := map[string]string{
		"name":    name,
		"version": version,
		"stage":   stage,
	}
	if archiveExistingVersions != "" {
		req["archive_existing_versions"] = archiveExistingVersions
	}
	endpoint := "/api/2.0/mlflow/model-versions/transition-stage"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response GetModelVersionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// RenameRegisteredModel renames a registered model
func (c *Client) RenameRegisteredModel(req RenameRegisteredModelRequest) (*RenameRegisteredModelResponse, error) {
	endpoint := "/api/2.0/mlflow/registered-models/rename"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response RenameRegisteredModelResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetLatestModelVersions gets the latest model versions for a registered model
func (c *Client) GetLatestModelVersions(req GetLatestModelVersionsRequest) (*GetLatestModelVersionsResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/registered-models/get-latest-versions?name=%s", url.QueryEscape(req.Name))
	if len(req.Stages) > 0 {
		params := url.Values{}
		for _, stage := range req.Stages {
			params.Add("stages", stage)
		}
		endpoint += "&" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetLatestModelVersionsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// SearchModelVersions searches for model versions
func (c *Client) SearchModelVersions(req SearchModelVersionsRequest) (*SearchModelVersionsResponse, error) {
	endpoint := "/api/2.0/mlflow/model-versions/search"
	// Build query parameters for GET request
	params := url.Values{}
	if req.Filter != "" {
		params.Add("filter", req.Filter)
	}
	if req.MaxResults > 0 {
		params.Add("max_results", fmt.Sprintf("%d", req.MaxResults))
	}
	if len(req.OrderBy) > 0 {
		for _, order := range req.OrderBy {
			params.Add("order_by", order)
		}
	}
	if req.PageToken != "" {
		params.Add("page_token", req.PageToken)
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response SearchModelVersionsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetDownloadURIs gets download URIs for model version artifacts
func (c *Client) GetDownloadURIs(req GetDownloadURIsRequest) (*GetDownloadURIsResponse, error) {
	endpoint := "/api/2.0/mlflow/model-versions/get-download-uris"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response GetDownloadURIsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// SearchRegisteredModels searches for registered models
func (c *Client) SearchRegisteredModels(req SearchRegisteredModelsRequest) (*SearchRegisteredModelsResponse, error) {
	endpoint := "/api/2.0/mlflow/registered-models/search"
	respBody, err := c.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var response SearchRegisteredModelsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// SetRegisteredModelTag sets a tag on a registered model
func (c *Client) SetRegisteredModelTag(req SetRegisteredModelTagRequest) error {
	endpoint := "/api/2.0/mlflow/registered-models/set-tag"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// SetModelVersionTag sets a tag on a model version
func (c *Client) SetModelVersionTag(req SetModelVersionTagRequest) error {
	endpoint := "/api/2.0/mlflow/model-versions/set-tag"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// DeleteRegisteredModelTag deletes a tag from a registered model
func (c *Client) DeleteRegisteredModelTag(req DeleteRegisteredModelTagRequest) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/registered-models/delete-tag?name=%s&key=%s",
		url.QueryEscape(req.Name), url.QueryEscape(req.Key))
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// DeleteModelVersionTag deletes a tag from a model version
func (c *Client) DeleteModelVersionTag(req DeleteModelVersionTagRequest) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/model-versions/delete-tag?name=%s&version=%s&key=%s",
		url.QueryEscape(req.Name), url.QueryEscape(req.Version), url.QueryEscape(req.Key))
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// SetRegisteredModelAlias sets an alias for a registered model
func (c *Client) SetRegisteredModelAlias(req SetRegisteredModelAliasRequest) error {
	endpoint := "/api/2.0/mlflow/registered-models/alias"
	_, err := c.doRequest("POST", endpoint, req)
	return err
}

// DeleteRegisteredModelAlias deletes an alias from a registered model
func (c *Client) DeleteRegisteredModelAlias(req DeleteRegisteredModelAliasRequest) error {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/registered-models/alias?name=%s&alias=%s",
		url.QueryEscape(req.Name), url.QueryEscape(req.Alias))
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// GetModelVersionByAlias gets a model version by alias
func (c *Client) GetModelVersionByAlias(req GetModelVersionByAliasRequest) (*GetModelVersionByAliasResponse, error) {
	endpoint := fmt.Sprintf("/api/2.0/mlflow/registered-models/get-model-version-by-alias?name=%s&alias=%s",
		url.QueryEscape(req.Name), url.QueryEscape(req.Alias))
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response GetModelVersionByAliasResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// doTextRequest performs an HTTP request and returns the response as a string
func (c *Client) doTextRequest(method, endpoint string) (string, error) {
	req, err := http.NewRequest(method, c.BaseURL+endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode:   resp.StatusCode,
			ResponseBody: respBody,
			Message:      string(respBody),
		}
		return "", apiErr
	}

	return string(respBody), nil
}

// GetHealth gets the health status of the MLflow server
func (c *Client) GetHealth() (string, error) {
	return c.doTextRequest("GET", "/health")
}

// GetVersion gets the version of the MLflow server
func (c *Client) GetVersion() (string, error) {
	return c.doTextRequest("GET", "/version")
}
