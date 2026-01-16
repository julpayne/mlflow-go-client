package mlflow

import "fmt"

// Experiment represents an MLflow experiment
type Experiment struct {
	ExperimentID     string          `json:"experiment_id"`
	Name             string          `json:"name"`
	ArtifactLocation string          `json:"artifact_location"`
	LifecycleStage   string          `json:"lifecycle_stage"`
	LastUpdateTime   int64           `json:"last_update_time"`
	CreationTime     int64           `json:"creation_time"`
	Tags             []ExperimentTag `json:"tags"`
}

// ExperimentTag represents a tag on an experiment
type ExperimentTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateExperimentRequest represents a request to create an experiment
type CreateExperimentRequest struct {
	Name             string          `json:"name"`
	ArtifactLocation string          `json:"artifact_location,omitempty"`
	Tags             []ExperimentTag `json:"tags,omitempty"`
}

// CreateExperimentResponse represents the response from creating an experiment
type CreateExperimentResponse struct {
	ExperimentID string `json:"experiment_id"`
}

// Run represents an MLflow run
type Run struct {
	Info    RunInfo    `json:"info"`
	Data    RunData    `json:"data"`
	Inputs  RunInputs  `json:"inputs,omitempty"`
	Outputs RunOutputs `json:"outputs,omitempty"`
}

// RunInfo contains metadata about a run
type RunInfo struct {
	RunID          string `json:"run_id"`
	RunName        string `json:"run_name,omitempty"`
	ExperimentID   string `json:"experiment_id"`
	UserID         string `json:"user_id,omitempty"`
	Status         string `json:"status"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time,omitempty"`
	ArtifactURI    string `json:"artifact_uri"`
	LifecycleStage string `json:"lifecycle_stage"`
}

// RunData contains metrics, parameters, and tags for a run
type RunData struct {
	Metrics []Metric `json:"metrics"`
	Params  []Param  `json:"params"`
	Tags    []RunTag `json:"tags"`
}

// RunInputs contains input datasets and tags
type RunInputs struct {
	Datasets    []Dataset    `json:"datasets,omitempty"`
	ModelInputs []ModelInput `json:"model_inputs,omitempty"`
}

// RunOutputs contains outputs of a Run
type RunOutputs struct {
	ModelOutputs []ModelOutput `json:"model_outputs,omitempty"`
}

// ModelOutput represents a model output
type ModelOutput struct {
	ModelName    string          `json:"model_name"`
	ModelVersion string          `json:"model_version,omitempty"`
	ModelStage   string          `json:"model_stage,omitempty"`
	Alias        string          `json:"alias,omitempty"`
	Tags         []ModelInputTag `json:"tags,omitempty"`
}

// Metric represents a metric logged to a run
type Metric struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Step      int64   `json:"step"`
}

// Param represents a parameter for a run
type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RunTag represents a tag on a run
type RunTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Dataset represents an input dataset
type Dataset struct {
	Name       string        `json:"name"`
	Digest     string        `json:"digest"`
	SourceType string        `json:"source_type"`
	Source     string        `json:"source"`
	Schema     DatasetSchema `json:"schema,omitempty"`
	Profile    string        `json:"profile,omitempty"`
	Tags       []DatasetTag  `json:"tags,omitempty"`
}

// DatasetSchema represents the schema of a dataset
type DatasetSchema struct {
	Columns []DatasetSchemaColumn `json:"columns,omitempty"`
}

// DatasetSchemaColumn represents a column in a dataset schema
type DatasetSchemaColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DatasetTag represents a tag on a dataset
type DatasetTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateRunRequest represents a request to create a run
type CreateRunRequest struct {
	ExperimentID string   `json:"experiment_id"`
	UserID       string   `json:"user_id,omitempty"`
	RunName      string   `json:"run_name,omitempty"`
	StartTime    int64    `json:"start_time,omitempty"`
	Tags         []RunTag `json:"tags,omitempty"`
}

// CreateRunResponse represents the response from creating a run
type CreateRunResponse struct {
	Run Run `json:"run"`
}

// LogMetricRequest represents a request to log a metric
type LogMetricRequest struct {
	RunID     string  `json:"run_id"`
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp,omitempty"`
	Step      int64   `json:"step,omitempty"`
}

// LogParamRequest represents a request to log a parameter
type LogParamRequest struct {
	RunID string `json:"run_id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SetTagRequest represents a request to set a tag
type SetTagRequest struct {
	RunID string `json:"run_id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SearchRunsRequest represents a request to search for runs
type SearchRunsRequest struct {
	ExperimentIDs []string `json:"experiment_ids,omitempty"`
	Filter        string   `json:"filter,omitempty"`
	RunViewType   string   `json:"run_view_type,omitempty"`
	MaxResults    int      `json:"max_results,omitempty"`
	OrderBy       []string `json:"order_by,omitempty"`
	PageToken     string   `json:"page_token,omitempty"`
}

// SearchRunsResponse represents the response from searching runs
type SearchRunsResponse struct {
	Runs          []Run  `json:"runs"`
	NextPageToken string `json:"next_page_token,omitempty"`
}

// UpdateRunRequest represents a request to update a run
type UpdateRunRequest struct {
	RunID   string `json:"run_id"`
	Status  string `json:"status,omitempty"`
	EndTime int64  `json:"end_time,omitempty"`
}

// UpdateRunResponse represents the response from updating a run
type UpdateRunResponse struct {
	RunInfo RunInfo `json:"run_info"`
}

// ModelVersion represents a model version in MLflow
type ModelVersion struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	CreationTimestamp    int64             `json:"creation_timestamp"`
	LastUpdatedTimestamp int64             `json:"last_updated_timestamp"`
	UserID               string            `json:"user_id,omitempty"`
	CurrentStage         string            `json:"current_stage"`
	Description          string            `json:"description,omitempty"`
	Source               string            `json:"source"`
	RunID                string            `json:"run_id,omitempty"`
	Status               string            `json:"status"`
	StatusMessage        string            `json:"status_message,omitempty"`
	Tags                 []ModelVersionTag `json:"tags,omitempty"`
	Aliases              []string          `json:"aliases,omitempty"`
}

// ModelVersionTag represents a tag on a model version
type ModelVersionTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RegisteredModel represents a registered model in MLflow
type RegisteredModel struct {
	Name                 string               `json:"name"`
	CreationTimestamp    int64                `json:"creation_timestamp"`
	LastUpdatedTimestamp int64                `json:"last_updated_timestamp"`
	UserID               string               `json:"user_id,omitempty"`
	Description          string               `json:"description,omitempty"`
	LatestVersions       []ModelVersion       `json:"latest_versions,omitempty"`
	Tags                 []RegisteredModelTag `json:"tags,omitempty"`
	Aliases              []string             `json:"aliases,omitempty"`
}

// RegisteredModelTag represents a tag on a registered model
type RegisteredModelTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateRegisteredModelRequest represents a request to create a registered model
type CreateRegisteredModelRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Tags        []RegisteredModelTag `json:"tags,omitempty"`
}

// CreateRegisteredModelResponse represents the response from creating a registered model
type CreateRegisteredModelResponse struct {
	RegisteredModel RegisteredModel `json:"registered_model"`
}

// CreateModelVersionRequest represents a request to create a model version
type CreateModelVersionRequest struct {
	Name        string            `json:"name"`
	Source      string            `json:"source"`
	RunID       string            `json:"run_id,omitempty"`
	Tags        []ModelVersionTag `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
}

// CreateModelVersionResponse represents the response from creating a model version
type CreateModelVersionResponse struct {
	ModelVersion ModelVersion `json:"model_version"`
}

// GetExperimentRequest represents a request to get an experiment
type GetExperimentRequest struct {
	ExperimentID string `json:"experiment_id"`
}

// GetExperimentByNameRequest represents a request to get an experiment by name
type GetExperimentByNameRequest struct {
	ExperimentName string `json:"experiment_name"`
}

// GetExperimentResponse represents the response from getting an experiment
type GetExperimentResponse struct {
	Experiment Experiment `json:"experiment"`
}

// GetRunRequest represents a request to get a run
type GetRunRequest struct {
	RunID string `json:"run_id"`
}

// GetRunResponse represents the response from getting a run
type GetRunResponse struct {
	Run Run `json:"run"`
}

// GetRegisteredModelRequest represents a request to get a registered model
type GetRegisteredModelRequest struct {
	Name string `json:"name"`
}

// GetRegisteredModelResponse represents the response from getting a registered model
type GetRegisteredModelResponse struct {
	RegisteredModel RegisteredModel `json:"registered_model"`
}

// GetModelVersionRequest represents a request to get a model version
type GetModelVersionRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// GetModelVersionResponse represents the response from getting a model version
type GetModelVersionResponse struct {
	ModelVersion ModelVersion `json:"model_version"`
}

// ErrorResponse represents an error response from the MLflow API
type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// APIError represents an error from the MLflow API
type APIError struct {
	StatusCode   int
	Message      string
	ResponseBody []byte
	ErrorCode    string
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("MLflow API error [%d]: %s - %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("MLflow API error [%d]: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("MLflow API error [%d]: %s", e.StatusCode, string(e.ResponseBody))
}

// GetStatusCode returns the HTTP status code
func (e *APIError) GetStatusCode() int {
	return e.StatusCode
}

// GetResponseBody returns the raw response body
func (e *APIError) GetResponseBody() []byte {
	return e.ResponseBody
}

// GetResponseBodyString returns the response body as a string
func (e *APIError) GetResponseBodyString() string {
	return string(e.ResponseBody)
}

// GetErrorCode returns the MLflow error code if available
func (e *APIError) GetErrorCode() string {
	return e.ErrorCode
}

// GetMessage returns the error message
func (e *APIError) GetMessage() string {
	return e.Message
}

// IsAPIError checks if an error is an APIError and returns it
func IsAPIError(err error) (*APIError, bool) {
	if err == nil {
		return nil, false
	}
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}

// SearchExperimentsRequest represents a request to search experiments
type SearchExperimentsRequest struct {
	ViewType   string   `json:"view_type,omitempty"` // ACTIVE_ONLY, DELETED_ONLY, or ALL
	MaxResults int      `json:"max_results,omitempty"`
	PageToken  string   `json:"page_token,omitempty"`
	Filter     string   `json:"filter,omitempty"`
	OrderBy    []string `json:"order_by,omitempty"`
}

// SearchExperimentsResponse represents the response from searching experiments
type SearchExperimentsResponse struct {
	Experiments   []Experiment `json:"experiments"`
	NextPageToken string       `json:"next_page_token,omitempty"`
}

// FileInfo represents information about a file in artifacts
type FileInfo struct {
	Path     string `json:"path"`
	IsDir    bool   `json:"is_dir"`
	FileSize int64  `json:"file_size,omitempty"`
}

// ListArtifactsRequest represents a request to list artifacts
type ListArtifactsRequest struct {
	RunID     string `json:"run_id"`
	Path      string `json:"path,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

// ListArtifactsResponse represents the response from listing artifacts
type ListArtifactsResponse struct {
	RootURI       string     `json:"root_uri"`
	Files         []FileInfo `json:"files"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

// GetMetricHistoryRequest represents a request to get metric history
type GetMetricHistoryRequest struct {
	RunID      string `json:"run_id"`
	MetricKey  string `json:"metric_key"`
	MaxResults int    `json:"max_results,omitempty"`
	PageToken  string `json:"page_token,omitempty"`
}

// GetMetricHistoryResponse represents the response from getting metric history
type GetMetricHistoryResponse struct {
	Metrics       []Metric `json:"metrics"`
	NextPageToken string   `json:"next_page_token,omitempty"`
}

// LogBatchRequest represents a request to log a batch of metrics
type LogBatchRequest struct {
	RunID   string   `json:"run_id"`
	Metrics []Metric `json:"metrics,omitempty"`
	Params  []Param  `json:"params,omitempty"`
	Tags    []RunTag `json:"tags,omitempty"`
}

// LogModelRequest represents a request to log a model
type LogModelRequest struct {
	RunID     string `json:"run_id"`
	ModelJSON string `json:"model_json"`
}

// LogInputsRequest represents a request to log inputs
type LogInputsRequest struct {
	RunID       string       `json:"run_id"`
	Datasets    []Dataset    `json:"datasets,omitempty"`
	ModelInputs []ModelInput `json:"model_inputs,omitempty"`
}

// ModelInput represents a model input
type ModelInput struct {
	ModelName    string          `json:"model_name"`
	ModelVersion string          `json:"model_version,omitempty"`
	ModelStage   string          `json:"model_stage,omitempty"`
	Alias        string          `json:"alias,omitempty"`
	Tags         []ModelInputTag `json:"tags,omitempty"`
}

// ModelInputTag represents a tag on a model input
type ModelInputTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RenameRegisteredModelRequest represents a request to rename a registered model
type RenameRegisteredModelRequest struct {
	Name    string `json:"name"`
	NewName string `json:"new_name"`
}

// RenameRegisteredModelResponse represents the response from renaming a registered model
type RenameRegisteredModelResponse struct {
	RegisteredModel RegisteredModel `json:"registered_model"`
}

// GetLatestModelVersionsRequest represents a request to get latest model versions
type GetLatestModelVersionsRequest struct {
	Name   string   `json:"name"`
	Stages []string `json:"stages,omitempty"`
}

// GetLatestModelVersionsResponse represents the response from getting latest model versions
type GetLatestModelVersionsResponse struct {
	ModelVersions []ModelVersion `json:"model_versions"`
}

// SearchModelVersionsRequest represents a request to search model versions
type SearchModelVersionsRequest struct {
	Filter     string   `json:"filter,omitempty"`
	MaxResults int      `json:"max_results,omitempty"`
	OrderBy    []string `json:"order_by,omitempty"`
	PageToken  string   `json:"page_token,omitempty"`
}

// SearchModelVersionsResponse represents the response from searching model versions
type SearchModelVersionsResponse struct {
	ModelVersions []ModelVersion `json:"model_versions"`
	NextPageToken string         `json:"next_page_token,omitempty"`
}

// GetDownloadURIsRequest represents a request to get download URIs
type GetDownloadURIsRequest struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Paths   []string `json:"paths,omitempty"`
}

// GetDownloadURIsResponse represents the response from getting download URIs
type GetDownloadURIsResponse struct {
	Files []DownloadURIInfo `json:"files"`
}

// DownloadURIInfo represents download URI information
type DownloadURIInfo struct {
	Path        string `json:"path"`
	ArtifactURI string `json:"artifact_uri"`
}

// SearchRegisteredModelsRequest represents a request to search registered models
type SearchRegisteredModelsRequest struct {
	Filter     string   `json:"filter,omitempty"`
	MaxResults int      `json:"max_results,omitempty"`
	OrderBy    []string `json:"order_by,omitempty"`
	PageToken  string   `json:"page_token,omitempty"`
}

// SearchRegisteredModelsResponse represents the response from searching registered models
type SearchRegisteredModelsResponse struct {
	RegisteredModels []RegisteredModel `json:"registered_models"`
	NextPageToken    string            `json:"next_page_token,omitempty"`
}

// SetRegisteredModelTagRequest represents a request to set a registered model tag
type SetRegisteredModelTagRequest struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SetModelVersionTagRequest represents a request to set a model version tag
type SetModelVersionTagRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

// DeleteRegisteredModelTagRequest represents a request to delete a registered model tag
type DeleteRegisteredModelTagRequest struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// DeleteModelVersionTagRequest represents a request to delete a model version tag
type DeleteModelVersionTagRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Key     string `json:"key"`
}

// SetRegisteredModelAliasRequest represents a request to set a registered model alias
type SetRegisteredModelAliasRequest struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Version string `json:"version"`
}

// DeleteRegisteredModelAliasRequest represents a request to delete a registered model alias
type DeleteRegisteredModelAliasRequest struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Version string `json:"version"`
}

// GetModelVersionByAliasRequest represents a request to get model version by alias
type GetModelVersionByAliasRequest struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// GetModelVersionByAliasResponse represents the response from getting model version by alias
type GetModelVersionByAliasResponse struct {
	ModelVersion ModelVersion `json:"model_version"`
}
