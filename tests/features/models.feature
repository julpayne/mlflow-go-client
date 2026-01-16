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
    When I create a registered model named "test-model"
    And I get the registered model "test-model"
    Then the model should be returned
    And the model name should be "test-model"

  Scenario: Search registered models
    When I create a registered model named "test-model"
    And a registered model named "test-model" exists
    When I search for models with filter "name='test-model'"
    Then I should find the model "test-model"

  Scenario: Update a registered model
    When I create a registered model named "test-model"
    And a registered model named "test-model" exists
    When I update the model description to "Updated description"
    Then the model description should be "Updated description"

  Scenario: Rename a registered model
    When I create a registered model named "test-model"
    And a registered model named "test-model" exists
    When I rename the model to "new-test-model"
    Then the model name should be "new-test-model"

  Scenario: Set registered model tag
    When I create a registered model named "tag-model"
    And a registered model named "tag-model" exists
    When I set tag "type" with value "classification" on the model
    Then the model should have tag "type" with value "classification"

  Scenario: Delete registered model tag
    When I create a registered model named "delete-tag-model"
    And a registered model named "delete-tag-model" exists
    And I set tag "temp" with value "value" on the model
    And the model has tag "temp" with value "value"
    When I delete tag "temp" from the model
    Then the model should not have tag "temp"

  Scenario: Create a model version
    When I create a registered model named "version-model"
    And a registered model named "version-model" exists
    When I create a model version with source "runs:/test-run/model"
    Then the model version should be created successfully

  Scenario: Get a model version
    When I create a registered model named "get-version-model"
    And a model version exists for model "get-version-model"
    When I get the model version
    Then the model version should be returned

  Scenario: Search model versions
    When I create a registered model named "search-version-model"
    And a model version exists for model "search-version-model"
    When I search for model versions with filter "name='search-version-model'"
    Then I should find at least one version

  Scenario: Get latest model versions
    When I create a registered model named "latest-version-model"
    And a model version exists for model "latest-version-model"
    When I get the latest model versions for "latest-version-model"
    Then I should get at least one version

  Scenario: Update a model version
    When I create a registered model named "update-version-model"
    And a model version exists for model "update-version-model"
    When I update the model version description to "Updated version"
    Then the model version description should be "Updated version"

  Scenario: Transition model version stage
    When I create a registered model named "stage-model"
    And a model version exists for model "stage-model"
    When I transition the model version to stage "Production"
    Then the model version stage should be "Production"

  Scenario: Set model version tag
    When I create a registered model named "tag-version-model"
    And a model version exists for model "tag-version-model"
    When I set tag "deployed" with value "true" on the model version
    Then the model version should have tag "deployed" with value "true"

  Scenario: Delete model version tag
    When I create a registered model named "tag-version-model"
    And a model version exists for model "tag-version-model"
    And the model version has tag "temp" with value "value"
    When I delete tag "temp" from the model version
    Then the model version should not have tag "temp"

  Scenario: Set model alias
    When I create a registered model named "alias-model"
    And a model version exists for model "alias-model"
    When I set alias "production" pointing to the model version
    Then the model version should have alias "production"

  Scenario: Get model version by alias
    When I create a registered model named "alias-get-model"
    And a model version with alias "prod-alias" exists for model "alias-get-model"
    When I get the model version by alias "prod-alias"
    Then the model version should be returned

  Scenario: Delete model alias
    When I create a registered model named "alias-delete-model"
    And a model version with alias "delete-alias" exists for model "alias-delete-model"
    When I delete alias "delete-alias" from the model
    Then the alias should be deleted

  Scenario: Delete a model version
    When I create a registered model named "delete-version-model"
    And a model version exists for model "delete-version-model"
    When I delete the model version
    Then the model version should be deleted

  Scenario: Delete a registered model
    When I create a registered model named "delete-model"
    And a registered model named "delete-model" exists
    When I delete the registered model
    Then the model should be deleted
