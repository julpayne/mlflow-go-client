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

func (c *Client) CheckServer() error {
	minSupportedVersion := "3.8.0"

	// Check that the server is running and has a version that we can handle
	health, err := c.GetHealth()
	if err != nil {
		return fmt.Errorf("failed to get server health: %w", err)
	}
	if health != "OK" {
		return fmt.Errorf("server health is not OK: %s", health)
	}
	version, err := c.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to get server version: %w", err)
	}
	if version == "" {
		return fmt.Errorf("server version is empty")
	}
	if version < minSupportedVersion {
		return fmt.Errorf("server version %s is not supported, expected %s or higher", version, minSupportedVersion)
	}
	return nil
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

// unmarshalResponse unmarshals JSON response body into a struct of type T
func unmarshalResponse[T any](respBody []byte) (*T, error) {
	var response T
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &response, nil
}

// API endpoint constants
const (
	// Base API path
	apiBasePath = "/api/2.0/mlflow"

	// Base URLs for API sections
	experimentsBaseURL      = apiBasePath + "/experiments"
	runsBaseURL             = apiBasePath + "/runs"
	registeredModelsBaseURL = apiBasePath + "/registered-models"
	modelVersionsBaseURL    = apiBasePath + "/model-versions"

	// Server endpoints
	endpointHealth  = "/health"
	endpointVersion = "/version"

	// Experiments endpoints
	endpointExperimentsCreate        = experimentsBaseURL + "/create"
	endpointExperimentsGetBase       = experimentsBaseURL + "/get"
	endpointExperimentsGetByNameBase = experimentsBaseURL + "/get-by-name"
	endpointExperimentsDeleteBase    = experimentsBaseURL + "/delete"
	endpointExperimentsRestoreBase   = experimentsBaseURL + "/restore"
	endpointExperimentsUpdate        = experimentsBaseURL + "/update"
	endpointExperimentsSetTag        = experimentsBaseURL + "/set-experiment-tag"
	endpointExperimentsDeleteTag     = experimentsBaseURL + "/delete-experiment-tag"
	endpointExperimentsSearch        = experimentsBaseURL + "/search"

	// Runs endpoints
	endpointRunsCreate       = runsBaseURL + "/create"
	endpointRunsGet          = runsBaseURL + "/get"
	endpointRunsSearch       = runsBaseURL + "/search"
	endpointRunsUpdate       = runsBaseURL + "/update"
	endpointRunsDelete       = runsBaseURL + "/delete"
	endpointRunsRestore      = runsBaseURL + "/restore"
	endpointRunsLogMetric    = runsBaseURL + "/log-metric"
	endpointRunsLogParameter = runsBaseURL + "/log-parameter"
	endpointRunsSetTag       = runsBaseURL + "/set-tag"
	endpointRunsDeleteTag    = runsBaseURL + "/delete-tag"
	endpointRunsLogBatch     = runsBaseURL + "/log-batch"
	endpointRunsLogModel     = runsBaseURL + "/log-model"
	endpointRunsLogInputs    = runsBaseURL + "/log-inputs"

	// Metrics endpoints
	endpointMetricsGetHistoryBase = apiBasePath + "/metrics/get-history"

	// Artifacts endpoints
	endpointArtifactsListBase = apiBasePath + "/artifacts/list"

	// Registered Models endpoints
	endpointRegisteredModelsCreate                     = registeredModelsBaseURL + "/create"
	endpointRegisteredModelsGet                        = registeredModelsBaseURL + "/get"
	endpointRegisteredModelsList                       = registeredModelsBaseURL + "/list"
	endpointRegisteredModelsUpdate                     = registeredModelsBaseURL + "/update"
	endpointRegisteredModelsDelete                     = registeredModelsBaseURL + "/delete"
	endpointRegisteredModelsRename                     = registeredModelsBaseURL + "/rename"
	endpointRegisteredModelsGetLatestVersions          = registeredModelsBaseURL + "/get-latest-versions"
	endpointRegisteredModelsSearch                     = registeredModelsBaseURL + "/search"
	endpointRegisteredModelsSetTag                     = registeredModelsBaseURL + "/set-tag"
	endpointRegisteredModelsDeleteTagBase              = registeredModelsBaseURL + "/delete-tag"
	endpointRegisteredModelsAliasBase                  = registeredModelsBaseURL + "/alias"
	endpointRegisteredModelsGetModelVersionByAliasBase = registeredModelsBaseURL + "/get-model-version-by-alias"

	// Model Versions endpoints
	endpointModelVersionsCreate          = modelVersionsBaseURL + "/create"
	endpointModelVersionsGetBase         = modelVersionsBaseURL + "/get"
	endpointModelVersionsList            = modelVersionsBaseURL + "/list"
	endpointModelVersionsUpdate          = modelVersionsBaseURL + "/update"
	endpointModelVersionsDeleteBase      = modelVersionsBaseURL + "/delete"
	endpointModelVersionsTransitionStage = modelVersionsBaseURL + "/transition-stage"
	endpointModelVersionsSearch          = modelVersionsBaseURL + "/search"
	endpointModelVersionsGetDownloadURIs = modelVersionsBaseURL + "/get-download-uris"
	endpointModelVersionsSetTag          = modelVersionsBaseURL + "/set-tag"
	endpointModelVersionsDeleteTagBase   = modelVersionsBaseURL + "/delete-tag"
)

// Endpoint helper functions for parameterized endpoints

// endpointMetricsGetHistory returns the endpoint for getting metric history with query parameters
func endpointMetricsGetHistory(runUUID, metricKey string, maxResults int, pageToken string) string {
	endpoint := fmt.Sprintf("%s?run_uuid=%s&metric_key=%s",
		endpointMetricsGetHistoryBase, url.QueryEscape(runUUID), url.QueryEscape(metricKey))
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
	return endpoint
}

// endpointArtifactsList returns the endpoint for listing artifacts with query parameters
func endpointArtifactsList(runID, path, pageToken string) string {
	endpoint := fmt.Sprintf("%s?run_id=%s", endpointArtifactsListBase, url.QueryEscape(runID))
	if path != "" {
		endpoint += "&path=" + url.QueryEscape(path)
	}
	if pageToken != "" {
		endpoint += "&page_token=" + url.QueryEscape(pageToken)
	}
	return endpoint
}

// endpointRegisteredModelsListWithParams returns the endpoint for listing registered models with query parameters
func endpointRegisteredModelsListWithParams(maxResults int, pageToken string) string {
	endpoint := endpointRegisteredModelsList
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
	return endpoint
}

// endpointModelVersionsListWithParams returns the endpoint for listing model versions with query parameters
func endpointModelVersionsListWithParams(name string, maxResults int, pageToken string) string {
	endpoint := fmt.Sprintf("%s?name=%s", endpointModelVersionsList, url.QueryEscape(name))
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
	return endpoint
}

// endpointRegisteredModelsGetLatestVersionsWithParams returns the endpoint for getting latest model versions with query parameters
func endpointRegisteredModelsGetLatestVersionsWithParams(name string, stages []string) string {
	endpoint := fmt.Sprintf("%s?name=%s", endpointRegisteredModelsGetLatestVersions, url.QueryEscape(name))
	if len(stages) > 0 {
		params := url.Values{}
		for _, stage := range stages {
			params.Add("stages", stage)
		}
		endpoint += "&" + params.Encode()
	}
	return endpoint
}

// endpointModelVersionsSearchWithParams returns the endpoint for searching model versions with query parameters
func endpointModelVersionsSearchWithParams(filter string, maxResults int, orderBy []string, pageToken string) string {
	endpoint := endpointModelVersionsSearch
	params := url.Values{}
	if filter != "" {
		params.Add("filter", filter)
	}
	if maxResults > 0 {
		params.Add("max_results", fmt.Sprintf("%d", maxResults))
	}
	if len(orderBy) > 0 {
		for _, order := range orderBy {
			params.Add("order_by", order)
		}
	}
	if pageToken != "" {
		params.Add("page_token", pageToken)
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	return endpoint
}

// Experiments API

// CreateExperiment creates a new experiment
func (c *Client) CreateExperiment(req CreateExperimentRequest) (*CreateExperimentResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointExperimentsCreate, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[CreateExperimentResponse](respBody)
}

// GetExperiment gets an experiment by ID
func (c *Client) GetExperiment(experimentID string) (*GetExperimentResponse, error) {
	req := GetExperimentRequest{
		ExperimentID: experimentID,
	}
	respBody, err := c.doRequest(http.MethodGet, endpointExperimentsGetBase, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetExperimentResponse](respBody)
}

// GetExperimentByName gets an experiment by name
func (c *Client) GetExperimentByName(experimentName string) (*GetExperimentResponse, error) {
	req := GetExperimentByNameRequest{
		ExperimentName: experimentName,
	}
	respBody, err := c.doRequest(http.MethodGet, endpointExperimentsGetByNameBase, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetExperimentResponse](respBody)
}

// DeleteExperiment deletes an experiment
func (c *Client) DeleteExperiment(experimentID string) error {
	req := map[string]string{
		"experiment_id": experimentID,
	}
	_, err := c.doRequest(http.MethodPost, endpointExperimentsDeleteBase, req)
	return err
}

// RestoreExperiment restores a deleted experiment
func (c *Client) RestoreExperiment(experimentID string) error {
	req := map[string]string{
		"experiment_id": experimentID,
	}
	_, err := c.doRequest(http.MethodPost, endpointExperimentsRestoreBase, req)
	return err
}

// UpdateExperiment updates an experiment
func (c *Client) UpdateExperiment(experimentID, newName string) error {
	req := map[string]interface{}{
		"experiment_id": experimentID,
		"new_name":      newName,
	}
	_, err := c.doRequest(http.MethodPost, endpointExperimentsUpdate, req)
	return err
}

// SetExperimentTag sets a tag on an experiment
func (c *Client) SetExperimentTag(experimentID, key, value string) error {
	req := map[string]string{
		"experiment_id": experimentID,
		"key":           key,
		"value":         value,
	}
	_, err := c.doRequest(http.MethodPost, endpointExperimentsSetTag, req)
	return err
}

// DeleteExperimentTag deletes a tag from an experiment
func (c *Client) DeleteExperimentTag(experimentID, key string) error {
	req := map[string]string{
		"experiment_id": experimentID,
		"key":           key,
	}
	_, err := c.doRequest(http.MethodPost, endpointExperimentsDeleteTag, req)
	return err
}

// SearchExperiments searches for experiments
func (c *Client) SearchExperiments(req SearchExperimentsRequest) (*SearchExperimentsResponse, error) {
	if req.MaxResults <= 0 {
		return nil, fmt.Errorf("max_results must be greater than zero when provided")
	}
	respBody, err := c.doRequest(http.MethodPost, endpointExperimentsSearch, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[SearchExperimentsResponse](respBody)
}

// Runs API

// CreateRun creates a new run
func (c *Client) CreateRun(req CreateRunRequest) (*CreateRunResponse, error) {
	if req.StartTime == 0 {
		req.StartTime = time.Now().UnixMilli()
	}

	respBody, err := c.doRequest(http.MethodPost, endpointRunsCreate, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[CreateRunResponse](respBody)
}

// GetRun gets a run by ID
func (c *Client) GetRun(runID string) (*GetRunResponse, error) {
	req := GetRunRequest{
		RunID: runID,
	}
	respBody, err := c.doRequest(http.MethodPost, endpointRunsGet, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetRunResponse](respBody)
}

// SearchRuns searches for runs
func (c *Client) SearchRuns(req SearchRunsRequest) (*SearchRunsResponse, error) {
	if req.MaxResults <= 0 {
		return nil, fmt.Errorf("max_results must be greater than zero when provided")
	}
	respBody, err := c.doRequest(http.MethodPost, endpointRunsSearch, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[SearchRunsResponse](respBody)
}

// UpdateRun updates a run
func (c *Client) UpdateRun(req UpdateRunRequest) (*UpdateRunResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointRunsUpdate, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[UpdateRunResponse](respBody)
}

// DeleteRun deletes a run
func (c *Client) DeleteRun(runID string) error {
	req := map[string]string{
		"run_id": runID,
	}
	_, err := c.doRequest(http.MethodPost, endpointRunsDelete, req)
	return err
}

// RestoreRun restores a deleted run
func (c *Client) RestoreRun(runID string) error {
	req := map[string]string{
		"run_id": runID,
	}
	_, err := c.doRequest(http.MethodPost, endpointRunsRestore, req)
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

	_, err := c.doRequest(http.MethodPost, endpointRunsLogMetric, req)
	return err
}

// LogParam logs a parameter to a run
func (c *Client) LogParam(req LogParamRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointRunsLogParameter, req)
	return err
}

// SetTag sets a tag on a run
func (c *Client) SetTag(req SetTagRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointRunsSetTag, req)
	return err
}

// DeleteTag deletes a tag from a run
func (c *Client) DeleteTag(runID, key string) error {
	req := map[string]string{
		"run_id": runID,
		"key":    key,
	}
	_, err := c.doRequest(http.MethodPost, endpointRunsDeleteTag, req)
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
	_, err := c.doRequest(http.MethodPost, endpointRunsLogBatch, req)
	return err
}

// LogModel logs a model to a run
func (c *Client) LogModel(req LogModelRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointRunsLogModel, req)
	return err
}

// LogInputs logs inputs (datasets and/or model inputs) to a run
func (c *Client) LogInputs(req LogInputsRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointRunsLogInputs, req)
	return err
}

// GetMetricHistory gets the history of a metric for a run
func (c *Client) GetMetricHistory(req GetMetricHistoryRequest) (*GetMetricHistoryResponse, error) {
	respBody, err := c.doRequest(http.MethodGet, endpointMetricsGetHistory(req.RunUUID, req.MetricKey, req.MaxResults, req.PageToken), nil)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetMetricHistoryResponse](respBody)
}

// ListArtifacts lists artifacts for a run
func (c *Client) ListArtifacts(runID, path string, pageToken string) (*ListArtifactsResponse, error) {
	respBody, err := c.doRequest(http.MethodGet, endpointArtifactsList(runID, path, pageToken), nil)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[ListArtifactsResponse](respBody)
}

// Models API

// CreateRegisteredModel creates a new registered model
func (c *Client) CreateRegisteredModel(req CreateRegisteredModelRequest) (*CreateRegisteredModelResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointRegisteredModelsCreate, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[CreateRegisteredModelResponse](respBody)
}

// GetRegisteredModel gets a registered model by name
func (c *Client) GetRegisteredModel(name string) (*GetRegisteredModelResponse, error) {
	req := GetRegisteredModelRequest{
		Name: name,
	}
	respBody, err := c.doRequest(http.MethodPost, endpointRegisteredModelsGet, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetRegisteredModelResponse](respBody)
}

// ListRegisteredModels lists all registered models
func (c *Client) ListRegisteredModels(maxResults int, pageToken string) (*ListRegisteredModelsResponse, error) {
	respBody, err := c.doRequest(http.MethodGet, endpointRegisteredModelsListWithParams(maxResults, pageToken), nil)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[ListRegisteredModelsResponse](respBody)
}

// UpdateRegisteredModel updates a registered model
func (c *Client) UpdateRegisteredModel(name, description string) error {
	req := map[string]string{
		"name": name,
	}
	if description != "" {
		req["description"] = description
	}
	_, err := c.doRequest(http.MethodPatch, endpointRegisteredModelsUpdate, req)
	return err
}

// DeleteRegisteredModel deletes a registered model
func (c *Client) DeleteRegisteredModel(name string) error {
	req := map[string]string{
		"name": name,
	}
	_, err := c.doRequest(http.MethodPost, endpointRegisteredModelsDelete, req)
	return err
}

// CreateModelVersion creates a new model version
func (c *Client) CreateModelVersion(req CreateModelVersionRequest) (*CreateModelVersionResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointModelVersionsCreate, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[CreateModelVersionResponse](respBody)
}

// GetModelVersion gets a model version
func (c *Client) GetModelVersion(name, version string) (*GetModelVersionResponse, error) {
	req := GetModelVersionRequest{
		Name:    name,
		Version: version,
	}
	respBody, err := c.doRequest(http.MethodPost, endpointModelVersionsGetBase, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetModelVersionResponse](respBody)
}

// ListModelVersions lists model versions for a registered model
func (c *Client) ListModelVersions(name string, maxResults int, pageToken string) (*ListModelVersionsResponse, error) {
	respBody, err := c.doRequest(http.MethodGet, endpointModelVersionsListWithParams(name, maxResults, pageToken), nil)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[ListModelVersionsResponse](respBody)
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
	_, err := c.doRequest(http.MethodPatch, endpointModelVersionsUpdate, req)
	return err
}

// DeleteModelVersion deletes a model version
func (c *Client) DeleteModelVersion(name, version string) error {
	req := map[string]string{
		"name":    name,
		"version": version,
	}
	_, err := c.doRequest(http.MethodPost, endpointModelVersionsDeleteBase, req)
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
	respBody, err := c.doRequest(http.MethodPost, endpointModelVersionsTransitionStage, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetModelVersionResponse](respBody)
}

// RenameRegisteredModel renames a registered model
func (c *Client) RenameRegisteredModel(req RenameRegisteredModelRequest) (*RenameRegisteredModelResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointRegisteredModelsRename, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[RenameRegisteredModelResponse](respBody)
}

// GetLatestModelVersions gets the latest model versions for a registered model
func (c *Client) GetLatestModelVersions(req GetLatestModelVersionsRequest) (*GetLatestModelVersionsResponse, error) {
	respBody, err := c.doRequest(http.MethodGet, endpointRegisteredModelsGetLatestVersionsWithParams(req.Name, req.Stages), nil)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetLatestModelVersionsResponse](respBody)
}

// SearchModelVersions searches for model versions
func (c *Client) SearchModelVersions(req SearchModelVersionsRequest) (*SearchModelVersionsResponse, error) {
	if req.MaxResults <= 0 {
		return nil, fmt.Errorf("max_results must be greater than zero when provided")
	}
	respBody, err := c.doRequest(http.MethodGet, endpointModelVersionsSearchWithParams(req.Filter, req.MaxResults, req.OrderBy, req.PageToken), nil)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[SearchModelVersionsResponse](respBody)
}

// GetDownloadURIs gets download URIs for model version artifacts
func (c *Client) GetDownloadURIs(req GetDownloadURIsRequest) (*GetDownloadURIsResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointModelVersionsGetDownloadURIs, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetDownloadURIsResponse](respBody)
}

// SearchRegisteredModels searches for registered models
func (c *Client) SearchRegisteredModels(req SearchRegisteredModelsRequest) (*SearchRegisteredModelsResponse, error) {
	if req.MaxResults <= 0 {
		return nil, fmt.Errorf("max_results must be greater than zero when provided")
	}
	respBody, err := c.doRequest(http.MethodPost, endpointRegisteredModelsSearch, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[SearchRegisteredModelsResponse](respBody)
}

// SetRegisteredModelTag sets a tag on a registered model
func (c *Client) SetRegisteredModelTag(req SetRegisteredModelTagRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointRegisteredModelsSetTag, req)
	return err
}

// SetModelVersionTag sets a tag on a model version
func (c *Client) SetModelVersionTag(req SetModelVersionTagRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointModelVersionsSetTag, req)
	return err
}

// DeleteRegisteredModelTag deletes a tag from a registered model
func (c *Client) DeleteRegisteredModelTag(req DeleteRegisteredModelTagRequest) error {
	reqBody := map[string]string{
		"name": req.Name,
		"key":  req.Key,
	}
	_, err := c.doRequest(http.MethodPost, endpointRegisteredModelsDeleteTagBase, reqBody)
	return err
}

// DeleteModelVersionTag deletes a tag from a model version
func (c *Client) DeleteModelVersionTag(req DeleteModelVersionTagRequest) error {
	reqBody := map[string]string{
		"name":    req.Name,
		"version": req.Version,
		"key":     req.Key,
	}
	_, err := c.doRequest(http.MethodPost, endpointModelVersionsDeleteTagBase, reqBody)
	return err
}

// SetRegisteredModelAlias sets an alias for a registered model
func (c *Client) SetRegisteredModelAlias(req SetRegisteredModelAliasRequest) error {
	_, err := c.doRequest(http.MethodPost, endpointRegisteredModelsAliasBase, req)
	return err
}

// DeleteRegisteredModelAlias deletes an alias from a registered model
func (c *Client) DeleteRegisteredModelAlias(req DeleteRegisteredModelAliasRequest) error {
	reqBody := map[string]string{
		"name":  req.Name,
		"alias": req.Alias,
	}
	_, err := c.doRequest(http.MethodPost, endpointRegisteredModelsAliasBase, reqBody)
	return err
}

// GetModelVersionByAlias gets a model version by alias
func (c *Client) GetModelVersionByAlias(req GetModelVersionByAliasRequest) (*GetModelVersionByAliasResponse, error) {
	respBody, err := c.doRequest(http.MethodPost, endpointRegisteredModelsGetModelVersionByAliasBase, req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[GetModelVersionByAliasResponse](respBody)
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
	return c.doTextRequest(http.MethodGet, endpointHealth)
}

// GetVersion gets the version of the MLflow server
func (c *Client) GetVersion() (string, error) {
	return c.doTextRequest(http.MethodGet, endpointVersion)
}
