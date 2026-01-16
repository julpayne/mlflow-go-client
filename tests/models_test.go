package features

import (
	"fmt"
	"time"

	"github.com/julpayne/mlflow-go-client/pkg/mlflow"
)

// Model step implementations

func (tc *testContext) createRegisteredModel(name string) error {
	tc.modelName = name
	req := mlflow.CreateRegisteredModelRequest{
		Name: name,
	}
	resp, err := tc.client.CreateRegisteredModel(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.model = &resp.RegisteredModel
	tc.createdResources = append(tc.createdResources, resource{Type: "model", Name: name})
	return nil
}

func (tc *testContext) modelCreatedSuccessfully() error {
	if tc.modelName == "" {
		return fmt.Errorf("model name is empty")
	}
	if tc.model == nil {
		return fmt.Errorf("model is nil")
	}
	if tc.model.Name != tc.modelName {
		return fmt.Errorf("expected model name %s, got %s", tc.modelName, tc.model.Name)
	}
	return nil
}

func (tc *testContext) modelHasName(name string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	if model.RegisteredModel.Name != name {
		return fmt.Errorf("expected model name %s, got %s", name, model.RegisteredModel.Name)
	}
	return nil
}

func (tc *testContext) registeredModelExists(name string) error {
	return tc.getRegisteredModel(name)
}

func (tc *testContext) getRegisteredModel(name string) error {
	tc.modelName = name
	_, err := tc.client.GetRegisteredModel(name)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) modelReturned() error {
	if tc.model == nil {
		return fmt.Errorf("model is nil")
	}
	return nil
}

func (tc *testContext) modelNameShouldBe(name string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	if model.RegisteredModel.Name != name {
		return fmt.Errorf("expected name %s, got %s", name, model.RegisteredModel.Name)
	}
	return nil
}

func (tc *testContext) multipleRegisteredModelsExist() error {
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("multi-model-%d-%d", time.Now().Unix(), i)
		if err := tc.createRegisteredModel(name); err != nil {
			return err
		}
	}
	return nil
}

func (tc *testContext) searchRegisteredModels(filter string) error {
	req := mlflow.SearchRegisteredModelsRequest{
		Filter:     filter,
		MaxResults: 100,
	}
	resp, err := tc.client.SearchRegisteredModels(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) findModel(name string) error {
	resp, ok := tc.lastResponse.(*mlflow.SearchRegisteredModelsResponse)
	if !ok {
		return fmt.Errorf("expected SearchRegisteredModelsResponse")
	}
	for _, model := range resp.RegisteredModels {
		if model.Name == name {
			return nil
		}
	}
	return fmt.Errorf("model %s not found", name)
}

func (tc *testContext) updateModelDescription(description string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	return tc.client.UpdateRegisteredModel(tc.modelName, description)
}

func (tc *testContext) modelDescriptionShouldBe(description string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	if model.RegisteredModel.Description != description {
		return fmt.Errorf("expected description %s, got %s", description, model.RegisteredModel.Description)
	}
	return nil
}

func (tc *testContext) renameModel(newName string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.RenameRegisteredModelRequest{
		Name:    tc.modelName,
		NewName: newName,
	}
	_, err := tc.client.RenameRegisteredModel(req)
	if err != nil {
		return err
	}
	tc.modelName = newName
	return nil
}

func (tc *testContext) setModelTag(key, value string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.SetRegisteredModelTagRequest{
		Name:  tc.modelName,
		Key:   key,
		Value: value,
	}
	return tc.client.SetRegisteredModelTag(req)
}

func (tc *testContext) modelHasTag(key, value string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	for _, tag := range model.RegisteredModel.Tags {
		if tag.Key == key && tag.Value == value {
			return nil
		}
	}
	return fmt.Errorf("tag %s=%s not found", key, value)
}

func (tc *testContext) modelHasTagSet(key, value string) error {
	return tc.setModelTag(key, value)
}

func (tc *testContext) deleteModelTag(key string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.DeleteRegisteredModelTagRequest{
		Name: tc.modelName,
		Key:  key,
	}
	return tc.client.DeleteRegisteredModelTag(req)
}

func (tc *testContext) modelDoesNotHaveTag(key string) error {
	model, err := tc.client.GetRegisteredModel(tc.modelName)
	if err != nil {
		return err
	}
	for _, tag := range model.RegisteredModel.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag %s still exists", key)
		}
	}
	return nil
}

func (tc *testContext) createModelVersion(source string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.CreateModelVersionRequest{
		Name:   tc.modelName,
		Source: source,
	}
	resp, err := tc.client.CreateModelVersion(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.modelVersion = resp.ModelVersion.Version
	return nil
}

func (tc *testContext) modelVersionCreatedSuccessfully() error {
	if tc.modelVersion == "" {
		return fmt.Errorf("model version is empty")
	}
	return nil
}

func (tc *testContext) modelVersionExists(modelName string) error {
	if err := tc.getRegisteredModel(modelName); err != nil {
		return err
	}
	return tc.createModelVersion("runs:/test-run/model")
}

func (tc *testContext) getModelVersion() error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	_, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		tc.lastError = err
		return err
	}
	return nil
}

func (tc *testContext) modelVersionReturned() error {
	return nil
}

func (tc *testContext) multipleModelVersionsExist(modelName string) error {
	if err := tc.createRegisteredModel(modelName); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if err := tc.createModelVersion(fmt.Sprintf("runs:/test-run-%d/model", i)); err != nil {
			return err
		}
	}
	return nil
}

func (tc *testContext) searchModelVersions(filter string) error {
	req := mlflow.SearchModelVersionsRequest{
		Filter: filter,
	}
	resp, err := tc.client.SearchModelVersions(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) findAtLeastOneVersion() error {
	resp, ok := tc.lastResponse.(*mlflow.SearchModelVersionsResponse)
	if !ok {
		return fmt.Errorf("expected SearchModelVersionsResponse")
	}
	if len(resp.ModelVersions) == 0 {
		return fmt.Errorf("no versions found")
	}
	return nil
}

func (tc *testContext) getLatestModelVersions(modelName string) error {
	tc.modelName = modelName
	req := mlflow.GetLatestModelVersionsRequest{
		Name: modelName,
	}
	resp, err := tc.client.GetLatestModelVersions(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.lastResponse = resp
	return nil
}

func (tc *testContext) getAtLeastOneVersion() error {
	resp, ok := tc.lastResponse.(*mlflow.GetLatestModelVersionsResponse)
	if !ok {
		return fmt.Errorf("expected GetLatestModelVersionsResponse")
	}
	if len(resp.ModelVersions) == 0 {
		return fmt.Errorf("no versions found")
	}
	return nil
}

func (tc *testContext) updateModelVersionDescription(description string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	return tc.client.UpdateModelVersion(tc.modelName, tc.modelVersion, description, "")
}

func (tc *testContext) modelVersionDescriptionShouldBe(description string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	if version.ModelVersion.Description != description {
		return fmt.Errorf("expected description %s, got %s", description, version.ModelVersion.Description)
	}
	return nil
}

func (tc *testContext) transitionModelVersionStage(stage string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	_, err := tc.client.TransitionModelVersionStage(tc.modelName, tc.modelVersion, stage, "")
	if err != nil {
		return err
	}
	return nil
}

func (tc *testContext) modelVersionStageShouldBe(stage string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	if version.ModelVersion.CurrentStage != stage {
		return fmt.Errorf("expected stage %s, got %s", stage, version.ModelVersion.CurrentStage)
	}
	return nil
}

func (tc *testContext) setModelVersionTag(key, value string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	req := mlflow.SetModelVersionTagRequest{
		Name:    tc.modelName,
		Version: tc.modelVersion,
		Key:     key,
		Value:   value,
	}
	return tc.client.SetModelVersionTag(req)
}

func (tc *testContext) modelVersionHasTag(key, value string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	for _, tag := range version.ModelVersion.Tags {
		if tag.Key == key && tag.Value == value {
			return nil
		}
	}
	return fmt.Errorf("tag %s=%s not found", key, value)
}

func (tc *testContext) modelVersionHasTagSet(key, value string) error {
	return tc.setModelVersionTag(key, value)
}

func (tc *testContext) deleteModelVersionTag(key string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	req := mlflow.DeleteModelVersionTagRequest{
		Name:    tc.modelName,
		Version: tc.modelVersion,
		Key:     key,
	}
	return tc.client.DeleteModelVersionTag(req)
}

func (tc *testContext) modelVersionDoesNotHaveTag(key string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	for _, tag := range version.ModelVersion.Tags {
		if tag.Key == key {
			return fmt.Errorf("tag %s still exists", key)
		}
	}
	return nil
}

func (tc *testContext) setModelAlias(alias string) error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	req := mlflow.SetRegisteredModelAliasRequest{
		Name:    tc.modelName,
		Alias:   alias,
		Version: tc.modelVersion,
	}
	return tc.client.SetRegisteredModelAlias(req)
}

func (tc *testContext) modelVersionHasAlias(alias string) error {
	version, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err != nil {
		return err
	}
	for _, a := range version.ModelVersion.Aliases {
		if a == alias {
			return nil
		}
	}
	return fmt.Errorf("alias %s not found", alias)
}

func (tc *testContext) modelVersionWithAliasExists(alias, modelName string) error {
	if err := tc.getRegisteredModel(modelName); err != nil {
		return err
	}
	if err := tc.createModelVersion("runs:/test-run/model"); err != nil {
		return err
	}
	return tc.setModelAlias(alias)
}

func (tc *testContext) getModelVersionByAlias(alias string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	req := mlflow.GetModelVersionByAliasRequest{
		Name:  tc.modelName,
		Alias: alias,
	}
	resp, err := tc.client.GetModelVersionByAlias(req)
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.modelVersion = resp.ModelVersion.Version
	return nil
}

func (tc *testContext) deleteModelAlias(alias string) error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	if tc.modelVersion == "" {
		return fmt.Errorf("no model version set")
	}
	req := mlflow.DeleteRegisteredModelAliasRequest{
		Name:    tc.modelName,
		Alias:   alias,
		Version: tc.modelVersion,
	}
	return tc.client.DeleteRegisteredModelAlias(req)
}

func (tc *testContext) aliasDeleted() error {
	// Verify alias is deleted by trying to get it (should fail)
	req := mlflow.GetModelVersionByAliasRequest{
		Name:  tc.modelName,
		Alias: "deleted-alias",
	}
	_, err := tc.client.GetModelVersionByAlias(req)
	if err == nil {
		return fmt.Errorf("alias still exists")
	}
	return nil
}

func (tc *testContext) deleteModelVersion() error {
	if tc.modelName == "" || tc.modelVersion == "" {
		return fmt.Errorf("model name or version not set")
	}
	return tc.client.DeleteModelVersion(tc.modelName, tc.modelVersion)
}

func (tc *testContext) modelVersionDeleted() error {
	_, err := tc.client.GetModelVersion(tc.modelName, tc.modelVersion)
	if err == nil {
		return fmt.Errorf("model version still exists")
	}
	return nil
}

func (tc *testContext) deleteRegisteredModel() error {
	if tc.modelName == "" {
		return fmt.Errorf("no model name set")
	}
	return tc.client.DeleteRegisteredModel(tc.modelName)
}

func (tc *testContext) modelDeleted() error {
	_, err := tc.client.GetRegisteredModel(tc.modelName)
	if err == nil {
		return fmt.Errorf("model still exists")
	}
	return nil
}
