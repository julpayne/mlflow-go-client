package features

import (
	"fmt"
	"regexp"
)

// Server health and version steps

func (tc *testContext) checkServerHealth() error {
	if tc.client == nil {
		return fmt.Errorf("client not initialized")
	}
	health, err := tc.client.GetHealth()
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.healthStatus = health
	return nil
}

func (tc *testContext) healthStatusShouldBe(expected string) error {
	if tc.healthStatus != expected {
		return fmt.Errorf("expected health status %s, got %s", expected, tc.healthStatus)
	}
	return nil
}

func (tc *testContext) checkServerVersion() error {
	if tc.client == nil {
		return fmt.Errorf("client not initialized")
	}
	version, err := tc.client.GetVersion()
	if err != nil {
		tc.lastError = err
		return err
	}
	tc.serverVersion = version
	return nil
}

func (tc *testContext) versionShouldMatch(pattern string) error {
	if err := tc.versionShouldNotBeEmpty(); err != nil {
		return err
	}
	matched, err := regexp.MatchString(pattern, tc.serverVersion)
	if err != nil {
		return fmt.Errorf("invalid regex pattern %s: %w", pattern, err)
	}
	if !matched {
		return fmt.Errorf("version %s does not match pattern %s", tc.serverVersion, pattern)
	}
	return nil
}

func (tc *testContext) versionShouldNotBeEmpty() error {
	if tc.serverVersion == "" {
		return fmt.Errorf("server version is empty")
	}
	return nil
}
