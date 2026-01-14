#!/bin/bash

# Script to run MLflow server locally
# Default configuration: runs on http://localhost:5000

set -e

# Default values
HOST=${MLFLOW_HOST:-"127.0.0.1"}
PORT=${MLFLOW_PORT:-"5000"}
BACKEND_URI=${MLFLOW_BACKEND_STORE_URI:-"sqlite:///mlflow.db"}
DEFAULT_ARTIFACT_ROOT=${MLFLOW_DEFAULT_ARTIFACT_ROOT:-"./mlruns"}

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting MLflow server...${NC}"
echo ""
echo -e "Configuration:"
echo -e "  ${YELLOW}Host:${NC} $HOST"
echo -e "  ${YELLOW}Port:${NC} $PORT"
echo -e "  ${YELLOW}Backend Store URI:${NC} $BACKEND_URI"
echo -e "  ${YELLOW}Default Artifact Root:${NC} $DEFAULT_ARTIFACT_ROOT"
echo ""

# Check if MLflow is installed
if ! command -v mlflow &> /dev/null; then
    echo -e "${YELLOW}⚠️  MLflow not found. Installing...${NC}"
    ./scripts/download_mlflow.sh
    echo ""
fi

# Create artifact root directory if it doesn't exist
if [ ! -d "$DEFAULT_ARTIFACT_ROOT" ]; then
    echo -e "${YELLOW}📁 Creating artifact root directory: $DEFAULT_ARTIFACT_ROOT${NC}"
    mkdir -p "$DEFAULT_ARTIFACT_ROOT"
fi

# Create backend database directory if using SQLite
if [[ "$BACKEND_URI" == sqlite://* ]]; then
    DB_PATH=$(echo "$BACKEND_URI" | sed 's|sqlite:///||')
    DB_DIR=$(dirname "$DB_PATH")
    if [ "$DB_DIR" != "." ] && [ ! -d "$DB_DIR" ]; then
        echo -e "${YELLOW}📁 Creating database directory: $DB_DIR${NC}"
        mkdir -p "$DB_DIR"
    fi
fi

echo -e "${GREEN}✅ Starting MLflow server...${NC}"
echo -e "${BLUE}📍 Server will be available at: http://$HOST:$PORT${NC}"
echo -e "${YELLOW}💡 Press Ctrl+C to stop the server${NC}"
echo ""

# Start MLflow server
mlflow server \
    --host "$HOST" \
    --port "$PORT" \
    --backend-store-uri "$BACKEND_URI" \
    --default-artifact-root "$DEFAULT_ARTIFACT_ROOT"
