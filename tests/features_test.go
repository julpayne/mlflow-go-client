package features

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
			// Strict mode will fail immediately on undefined steps
			Strict: true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests. Check for undefined steps in feature files.")
	}
}

/*
// TestStepDefinitions validates that all steps in feature files have corresponding definitions
// This test will fail if any step in the feature files doesn't match a defined step pattern
func TestStepDefinitions(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
			Strict:   true, // Fail on undefined steps
		},
	}

	// Run the test suite
	_ = suite.Run()
}
*/

func TestMain(m *testing.M) {
	status := m.Run()
	os.Exit(status)
}
