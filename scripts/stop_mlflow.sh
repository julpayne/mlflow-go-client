#!/bin/bash

# Script to stop the MLflow server locally

echo "🛑 Stopping MLflow server..."
pkill -f "mlflow.server" || echo "No MLflow server process found"
