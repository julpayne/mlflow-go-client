package mlflow

// Experiment represents an MLflow experiment
type Experiment struct {
	ExperimentID   string            `json:"experiment_id"`
	Name           string            `json:"name"`
	ArtifactLocation string          `json:"artifact_location"`
	LifecycleStage string            `json:"lifecycle_stage"`
	LastUpdateTime int64             `json:"last_update_time"`
	CreationTime   int64             `json:"creation_time"`
	Tags           []ExperimentTag   `json:"tags"`
}

// ExperimentTag represents a tag on an experiment
type ExperimentTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateExperimentRequest represents a request to create an experiment
type CreateExperimentRequest struct {
	Name             string            `json:"name"`
	ArtifactLocation string            `json:"artifact_location,omitempty"`
	Tags             []ExperimentTag   `json:"tags,omitempty"`
}

// CreateExperimentResponse represents the response from creating an experiment
type CreateExperimentResponse struct {
	ExperimentID string `json:"experiment_id"`
}

// Run represents an MLflow run
type Run struct {
	Info     RunInfo     `json:"info"`
	Data     RunData     `json:"data"`
	Inputs   RunInputs   `json:"inputs,omitempty"`
}

// RunInfo contains metadata about a run
type RunInfo struct {
	RunID           string    `json:"run_id"`
	RunUUID         string    `json:"run_uuid"`
	RunName         string    `json:"run_name,omitempty"`
	ExperimentID    string    `json:"experiment_id"`
	UserID          string    `json:"user_id,omitempty"`
	Status          string    `json:"status"`
	StartTime       int64     `json:"start_time"`
	EndTime         int64     `json:"end_time,omitempty"`
	ArtifactURI     string    `json:"artifact_uri"`
	LifecycleStage  string    `json:"lifecycle_stage"`
}

// RunData contains metrics, parameters, and tags for a run
type RunData struct {
	Metrics   []Metric   `json:"metrics"`
	Params    []Param    `json:"params"`
	Tags      []RunTag   `json:"tags"`
}

// RunInputs contains input datasets and tags
type RunInputs struct {
	Datasets []Dataset `json:"datasets,omitempty"`
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
	Name     string            `json:"name"`
	Digest   string            `json:"digest"`
	SourceType string          `json:"source_type"`
	Source   string            `json:"source"`
	Schema   DatasetSchema     `json:"schema,omitempty"`
	Profile  string            `json:"profile,omitempty"`
	Tags     []DatasetTag      `json:"tags,omitempty"`
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
	ExperimentID string            `json:"experiment_id"`
	UserID       string            `json:"user_id,omitempty"`
	RunName      string            `json:"run_name,omitempty"`
	StartTime    int64             `json:"start_time,omitempty"`
	Tags         []RunTag          `json:"tags,omitempty"`
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
	ExperimentIDs []string          `json:"experiment_ids,omitempty"`
	Filter        string            `json:"filter,omitempty"`
	RunViewType   string            `json:"run_view_type,omitempty"`
	MaxResults    int               `json:"max_results,omitempty"`
	OrderBy       []string          `json:"order_by,omitempty"`
	PageToken     string            `json:"page_token,omitempty"`
}

// SearchRunsResponse represents the response from searching runs
type SearchRunsResponse struct {
	Runs         []Run  `json:"runs"`
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
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	CreationTimestamp int64        `json:"creation_timestamp"`
	LastUpdatedTimestamp int64     `json:"last_updated_timestamp"`
	UserID       string            `json:"user_id,omitempty"`
	CurrentStage string            `json:"current_stage"`
	Description  string            `json:"description,omitempty"`
	Source       string            `json:"source"`
	RunID        string            `json:"run_id,omitempty"`
	Status       string            `json:"status"`
	StatusMessage string           `json:"status_message,omitempty"`
	Tags         []ModelVersionTag `json:"tags,omitempty"`
	Aliases      []string          `json:"aliases,omitempty"`
}

// ModelVersionTag represents a tag on a model version
type ModelVersionTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RegisteredModel represents a registered model in MLflow
type RegisteredModel struct {
	Name         string            `json:"name"`
	CreationTimestamp int64        `json:"creation_timestamp"`
	LastUpdatedTimestamp int64     `json:"last_updated_timestamp"`
	UserID       string            `json:"user_id,omitempty"`
	Description  string            `json:"description,omitempty"`
	LatestVersions []ModelVersion  `json:"latest_versions,omitempty"`
	Tags         []RegisteredModelTag `json:"tags,omitempty"`
	Aliases      []string          `json:"aliases,omitempty"`
}

// RegisteredModelTag represents a tag on a registered model
type RegisteredModelTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateRegisteredModelRequest represents a request to create a registered model
type CreateRegisteredModelRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Tags        []RegisteredModelTag  `json:"tags,omitempty"`
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

// ListExperimentsResponse represents the response from listing experiments
type ListExperimentsResponse struct {
	Experiments      []Experiment `json:"experiments"`
	NextPageToken    string       `json:"next_page_token,omitempty"`
}

// GetExperimentResponse represents the response from getting an experiment
type GetExperimentResponse struct {
	Experiment Experiment `json:"experiment"`
}

// GetRunResponse represents the response from getting a run
type GetRunResponse struct {
	Run Run `json:"run"`
}

// GetRegisteredModelResponse represents the response from getting a registered model
type GetRegisteredModelResponse struct {
	RegisteredModel RegisteredModel `json:"registered_model"`
}

// GetModelVersionResponse represents the response from getting a model version
type GetModelVersionResponse struct {
	ModelVersion ModelVersion `json:"model_version"`
}

// ListRegisteredModelsResponse represents the response from listing registered models
type ListRegisteredModelsResponse struct {
	RegisteredModels []RegisteredModel `json:"registered_models"`
	NextPageToken    string             `json:"next_page_token,omitempty"`
}

// ListModelVersionsResponse represents the response from listing model versions
type ListModelVersionsResponse struct {
	ModelVersions []ModelVersion `json:"model_versions"`
	NextPageToken string         `json:"next_page_token,omitempty"`
}

// ErrorResponse represents an error response from the MLflow API
type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	Message      string `json:"message"`
}
