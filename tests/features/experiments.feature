Feature: MLflow Experiments API
  As a developer
  I want to manage MLflow experiments
  So that I can organize my machine learning runs

  Background:
    Given an MLflow server is running at "http://localhost:5000"
    And I have an MLflow client connected to the server

  Scenario: Create a new experiment
    When I create an experiment named "test-experiment"
    Then the experiment should be created successfully
    And the experiment should have the name "test-experiment"

  Scenario: Get an experiment by ID
    Given an experiment named "get-experiment" exists
    When I get the experiment by ID
    Then the experiment should be returned
    And the experiment name should be "get-experiment"

  Scenario: Get an experiment by name
    Given an experiment named "get-by-name-experiment" exists
    When I get the experiment by name "get-by-name-experiment"
    Then the experiment should be returned
    And the experiment name should be "get-by-name-experiment"

  Scenario: List all experiments
    Given multiple experiments exist
    When I list all experiments
    Then I should get a list of experiments
    And the list should contain at least 1 experiment

  Scenario: Search experiments
    Given an experiment named "search-test" exists
    When I search for experiments with filter "name='search-test'"
    Then I should find the experiment "search-test"

  Scenario: Update an experiment
    Given an experiment named "update-test" exists
    When I update the experiment name to "updated-name"
    Then the experiment name should be "updated-name"

  Scenario: Set experiment tag
    Given an experiment named "tag-test" exists
    When I set tag "team" with value "ml-team" on the experiment
    Then the experiment should have tag "team" with value "ml-team"

  Scenario: Delete experiment tag
    Given an experiment named "delete-tag-test" exists
    And the experiment has tag "temp" with value "value"
    When I delete tag "temp" from the experiment
    Then the experiment should not have tag "temp"

  Scenario: Delete an experiment
    Given an experiment named "delete-test" exists
    When I delete the experiment
    Then the experiment should be deleted

  Scenario: Restore a deleted experiment
    Given a deleted experiment named "restore-test" exists
    When I restore the experiment
    Then the experiment should be restored
