#!/bin/bash

# Script to download and install MLflow server
# This script installs MLflow using pip

set -e

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 is not installed. Please install Python 3 first."
    exit 1
fi

# Check if pip is installed
if ! command -v pip3 &> /dev/null; then
    echo "❌ Error: pip3 is not installed. Please install pip3 first."
    exit 1
fi

# Get Python version
PYTHON_VERSION=$(python3 --version 2>&1 | awk '{print $2}')
echo "📦 Python version: $PYTHON_VERSION"

# Install MLflow
echo "📥 Installing MLflow..."
pip3 install mlflow

# Verify installation
if command -v mlflow &> /dev/null; then
    MLFLOW_VERSION=$(mlflow --version 2>&1 | head -n 1)
    echo "✅ MLflow installed successfully!"
    echo "   Version: $MLFLOW_VERSION"
    echo ""
    echo "🎉 MLflow is ready to use!"
    echo "   Run 'make run-mlflow' to start the server"
else
    echo "❌ Error: MLflow installation failed or not found in PATH"
    exit 1
fi
