Feature: MLflow Runs API
  As a developer
  I want to manage MLflow runs
  So that I can track my machine learning experiments

  Background:
    Given an MLflow server is running at "http://localhost:5000"
    And I have an MLflow client connected to the server
    And an experiment named "runs-experiment" exists

  Scenario: Create a new run
    When I create a run in the experiment
    Then the run should be created successfully
    And the run should have status "RUNNING"

  Scenario: Get a run by ID
    Given a run exists in the experiment
    When I get the run by ID
    Then the run should be returned
    And the run should have valid metadata

  Scenario: Log a metric to a run
    Given a run exists in the experiment
    When I log metric "accuracy" with value 0.95 to the run
    Then the metric should be logged successfully

  Scenario: Log a parameter to a run
    Given a run exists in the experiment
    When I log parameter "learning_rate" with value "0.01" to the run
    Then the parameter should be logged successfully

  Scenario: Set a tag on a run
    Given a run exists in the experiment
    When I set tag "version" with value "1.0" on the run
    Then the run should have tag "version" with value "1.0"

  Scenario: Delete a tag from a run
    Given a run exists in the experiment
    And the run has tag "temp" with value "value"
    When I delete tag "temp" from the run
    Then the run should not have tag "temp"

  Scenario: Log batch metrics and parameters
    Given a run exists in the experiment
    When I log batch with 2 metrics and 2 parameters to the run
    Then the batch should be logged successfully
    And the run should have 2 metrics
    And the run should have 2 parameters

  Scenario: Update run status
    Given a run exists in the experiment
    When I update the run status to "FINISHED"
    Then the run status should be "FINISHED"

  Scenario: Search runs
    Given multiple runs exist in the experiment
    When I search for runs with filter "metrics.accuracy > 0.9"
    Then I should get a list of runs

  Scenario: Get metric history
    Given a run exists in the experiment
    And I have logged metric "loss" multiple times to the run
    When I get the metric history for "loss"
    Then I should get multiple metric values

  Scenario: List artifacts
    Given a run exists in the experiment
    When I list artifacts for the run
    Then I should get a list of artifacts

  Scenario: Delete a run
    Given a run exists in the experiment
    When I delete the run
    Then the run should be deleted

  Scenario: Restore a deleted run
    Given a deleted run exists in the experiment
    When I restore the run
    Then the run should be restored
