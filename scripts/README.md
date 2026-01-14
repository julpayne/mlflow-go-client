# MLflow Server Scripts

This directory contains scripts to download and run the MLflow server locally.

## Scripts

### `download_mlflow.sh`

Downloads and installs MLflow using pip.

**Usage:**
```bash
./scripts/download_mlflow.sh
```

**Requirements:**
- Python 3
- pip3

### `run_mlflow.sh`

Starts the MLflow server locally with default configuration.

**Usage:**
```bash
./scripts/run_mlflow.sh
```

**Default Configuration:**
- Host: `127.0.0.1`
- Port: `5000`
- Backend Store: `sqlite:///mlflow.db`
- Artifact Root: `./mlruns`

**Customization:**

You can customize the server settings using environment variables:

```bash
# Custom host and port
MLFLOW_HOST=0.0.0.0 MLFLOW_PORT=8080 ./scripts/run_mlflow.sh

# Use PostgreSQL backend
MLFLOW_BACKEND_STORE_URI=postgresql://user:password@localhost/mlflow \
MLFLOW_DEFAULT_ARTIFACT_ROOT=s3://my-bucket/mlflow-artifacts \
./scripts/run_mlflow.sh

# Use file system for artifacts
MLFLOW_DEFAULT_ARTIFACT_ROOT=/path/to/artifacts ./scripts/run_mlflow.sh
```

**Environment Variables:**
- `MLFLOW_HOST` - Server host (default: `127.0.0.1`)
- `MLFLOW_PORT` - Server port (default: `5000`)
- `MLFLOW_BACKEND_STORE_URI` - Backend store URI (default: `sqlite:///mlflow.db`)
- `MLFLOW_DEFAULT_ARTIFACT_ROOT` - Default artifact root (default: `./mlruns`)

## Quick Start

1. **Install MLflow:**
   ```bash
   ./scripts/download_mlflow.sh
   ```

2. **Start the server:**
   ```bash
   ./scripts/run_mlflow.sh
   ```

3. **Access the UI:**
   Open your browser and navigate to: http://localhost:5000

4. **Use with Go client:**
   ```go
   client := mlflow.NewClient("http://localhost:5000")
   ```

## Backend Store Options

### SQLite (Default)
```bash
MLFLOW_BACKEND_STORE_URI=sqlite:///mlflow.db ./scripts/run_mlflow.sh
```

### PostgreSQL
```bash
MLFLOW_BACKEND_STORE_URI=postgresql://user:password@localhost/mlflow \
./scripts/run_mlflow.sh
```

### MySQL
```bash
MLFLOW_BACKEND_STORE_URI=mysql+pymysql://user:password@localhost/mlflow \
./scripts/run_mlflow.sh
```

## Artifact Storage Options

### Local File System (Default)
```bash
MLFLOW_DEFAULT_ARTIFACT_ROOT=./mlruns ./scripts/run_mlflow.sh
```

### S3
```bash
MLFLOW_DEFAULT_ARTIFACT_ROOT=s3://my-bucket/mlflow-artifacts \
AWS_ACCESS_KEY_ID=your-key \
AWS_SECRET_ACCESS_KEY=your-secret \
./scripts/run_mlflow.sh
```

### Azure Blob Storage
```bash
MLFLOW_DEFAULT_ARTIFACT_ROOT=wasbs://container@account.blob.core.windows.net/mlflow \
AZURE_STORAGE_CONNECTION_STRING=your-connection-string \
./scripts/run_mlflow.sh
```

### Google Cloud Storage
```bash
MLFLOW_DEFAULT_ARTIFACT_ROOT=gs://my-bucket/mlflow-artifacts \
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json \
./scripts/run_mlflow.sh
```

## Troubleshooting

### Port Already in Use
If port 5000 is already in use, specify a different port:
```bash
MLFLOW_PORT=5001 ./scripts/run_mlflow.sh
```

### Permission Denied
Make sure the scripts are executable:
```bash
chmod +x scripts/*.sh
```

### MLflow Not Found
If MLflow is not found, make sure it's installed:
```bash
./scripts/download_mlflow.sh
```

Or install manually:
```bash
pip3 install mlflow
```
