Feature: MLflow Server Health and Version
  As a developer
  I want to check the MLflow server health and version
  So that I can verify the server is running and check its version

  Background:
    Given an MLflow server is running at "http://localhost:5000"
    And I have an MLflow client connected to the server

  Scenario: Check server health
    When I check the server health
    Then the health status should be "OK"

  Scenario: Check server version matches specific pattern
    When I check the server version
    # Only matches 3.8.X versions
    Then the version should match "^3.8.[0-9]+$"
