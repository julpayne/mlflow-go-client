Feature: MLflow Models API
  As a developer
  I want to manage MLflow registered models
  So that I can version and deploy my models

  Background:
    Given an MLflow server is running at "http://localhost:5000"
    And I have an MLflow client connected to the server

  Scenario: Create a registered model
    When I create a registered model named "test-model"
    Then the model should be created successfully
    And the model should have the name "test-model"

  Scenario: Get a registered model
    Given a registered model named "get-model" exists
    When I get the registered model "get-model"
    Then the model should be returned
    And the model name should be "get-model"

  Scenario: List registered models
    Given multiple registered models exist
    When I list all registered models
    Then I should get a list of models
    And the list should contain at least 1 model

  Scenario: Search registered models
    Given a registered model named "search-model" exists
    When I search for models with filter "name='search-model'"
    Then I should find the model "search-model"

  Scenario: Update a registered model
    Given a registered model named "update-model" exists
    When I update the model description to "Updated description"
    Then the model description should be "Updated description"

  Scenario: Rename a registered model
    Given a registered model named "old-name" exists
    When I rename the model to "new-name"
    Then the model name should be "new-name"

  Scenario: Set registered model tag
    Given a registered model named "tag-model" exists
    When I set tag "type" with value "classification" on the model
    Then the model should have tag "type" with value "classification"

  Scenario: Delete registered model tag
    Given a registered model named "delete-tag-model" exists
    And the model has tag "temp" with value "value"
    When I delete tag "temp" from the model
    Then the model should not have tag "temp"

  Scenario: Create a model version
    Given a registered model named "version-model" exists
    When I create a model version with source "runs:/test-run/model"
    Then the model version should be created successfully

  Scenario: Get a model version
    Given a model version exists for model "get-version-model"
    When I get the model version
    Then the model version should be returned

  Scenario: List model versions
    Given multiple model versions exist for model "list-versions-model"
    When I list model versions for "list-versions-model"
    Then I should get a list of versions
    And the list should contain at least 1 version

  Scenario: Search model versions
    Given a model version exists for model "search-version-model"
    When I search for model versions with filter "name='search-version-model'"
    Then I should find at least one version

  Scenario: Get latest model versions
    Given a model version exists for model "latest-version-model"
    When I get the latest model versions for "latest-version-model"
    Then I should get at least one version

  Scenario: Update a model version
    Given a model version exists for model "update-version-model"
    When I update the model version description to "Updated version"
    Then the model version description should be "Updated version"

  Scenario: Transition model version stage
    Given a model version exists for model "stage-model"
    When I transition the model version to stage "Production"
    Then the model version stage should be "Production"

  Scenario: Set model version tag
    Given a model version exists for model "tag-version-model"
    When I set tag "deployed" with value "true" on the model version
    Then the model version should have tag "deployed" with value "true"

  Scenario: Delete model version tag
    Given a model version exists for model "delete-tag-version-model"
    And the model version has tag "temp" with value "value"
    When I delete tag "temp" from the model version
    Then the model version should not have tag "temp"

  Scenario: Set model alias
    Given a model version exists for model "alias-model"
    When I set alias "production" pointing to the model version
    Then the model version should have alias "production"

  Scenario: Get model version by alias
    Given a model version with alias "prod-alias" exists for model "alias-get-model"
    When I get the model version by alias "prod-alias"
    Then the model version should be returned

  Scenario: Delete model alias
    Given a model version with alias "delete-alias" exists for model "alias-delete-model"
    When I delete alias "delete-alias" from the model
    Then the alias should be deleted

  Scenario: Delete a model version
    Given a model version exists for model "delete-version-model"
    When I delete the model version
    Then the model version should be deleted

  Scenario: Delete a registered model
    Given a registered model named "delete-model" exists
    When I delete the registered model
    Then the model should be deleted
