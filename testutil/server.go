package testutil

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// MLflowServer manages an MLflow server process for testing
type MLflowServer struct {
	Host              string
	Port              string
	BackendStoreURI   string
	ArtifactRoot      string
	cmd               *exec.Cmd
	BaseURL           string
}

// NewMLflowServer creates a new MLflow server instance for testing
func NewMLflowServer() *MLflowServer {
	port := os.Getenv("MLFLOW_TEST_PORT")
	if port == "" {
		port = "5001" // Use different port than default to avoid conflicts
	}

	host := os.Getenv("MLFLOW_TEST_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	backendStoreURI := os.Getenv("MLFLOW_TEST_BACKEND_STORE_URI")
	if backendStoreURI == "" {
		backendStoreURI = fmt.Sprintf("sqlite:///test_mlflow_%s.db", port)
	}

	artifactRoot := os.Getenv("MLFLOW_TEST_ARTIFACT_ROOT")
	if artifactRoot == "" {
		artifactRoot = fmt.Sprintf("./test_mlruns_%s", port)
	}

	return &MLflowServer{
		Host:            host,
		Port:            port,
		BackendStoreURI: backendStoreURI,
		ArtifactRoot:    artifactRoot,
		BaseURL:         fmt.Sprintf("http://%s:%s", host, port),
	}
}

// Start starts the MLflow server
func (s *MLflowServer) Start(ctx context.Context) error {
	// Check if mlflow command exists
	if _, err := exec.LookPath("mlflow"); err != nil {
		return fmt.Errorf("mlflow command not found: %w", err)
	}

	// Create artifact root directory
	if err := os.MkdirAll(s.ArtifactRoot, 0755); err != nil {
		return fmt.Errorf("failed to create artifact root: %w", err)
	}

	// Start MLflow server
	s.cmd = exec.CommandContext(ctx, "mlflow", "server",
		"--host", s.Host,
		"--port", s.Port,
		"--backend-store-uri", s.BackendStoreURI,
		"--default-artifact-root", s.ArtifactRoot,
	)

	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MLflow server: %w", err)
	}

	// Wait for server to be ready
	if err := s.waitForServer(ctx, 30*time.Second); err != nil {
		s.Stop()
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

// waitForServer waits for the server to be ready
func (s *MLflowServer) waitForServer(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

		client := &http.Client{
			Timeout: 1 * time.Second,
		}

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if time.Now().After(deadline) {
					return fmt.Errorf("timeout waiting for server to start")
				}

				// Try to connect to the server
				resp, err := client.Get(s.BaseURL + "/health")
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode == 200 {
						return nil
					}
				}
			}
		}
}

// Stop stops the MLflow server
func (s *MLflowServer) Stop() error {
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}

// Cleanup cleans up test artifacts
func (s *MLflowServer) Cleanup() error {
	// Remove test database
	if s.BackendStoreURI != "" {
		dbPath := s.BackendStoreURI
		if len(dbPath) > 10 && dbPath[:10] == "sqlite:///" {
			dbPath = dbPath[10:]
			os.Remove(dbPath)
		}
	}

	// Remove artifact root
	if s.ArtifactRoot != "" {
		os.RemoveAll(s.ArtifactRoot)
	}

	return nil
}
