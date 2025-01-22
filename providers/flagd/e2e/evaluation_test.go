package e2e

import (
	"context"
	"flag"
	"github.com/open-feature/go-sdk-contrib/providers/flagd/e2e/containers"
	"testing"

	"github.com/cucumber/godog"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"
	"github.com/open-feature/go-sdk/openfeature"
)

func TestETestEvaluationFlagdInRPC(t *testing.T) {
	if testing.Short() {
		// skip e2e if testing -short
		t.Skip()
	}

	container, err := containers.NewFlagd(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	flag.Parse()

	name := "evaluation.feature"

	testSuite := godog.TestSuite{
		Name: name,
		TestSuiteInitializer: integration.InitializeTestSuite(func() openfeature.FeatureProvider {
			return flagd.NewProvider(flagd.WithPort(uint16(container.GetPort())))
		}),
		ScenarioInitializer: integration.InitializeEvaluationScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../spec/specification/assets/gherkin/evaluation.feature"},
			TestingT: t, // Testing instance that will run subtests.
			Strict:   true,
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run evaluation tests")
	}
}

func TestJsonEvaluatorFlagdInProcess(t *testing.T) {
	if testing.Short() {
		// skip e2e if testing -short
		t.Skip()
	}

	flag.Parse()

	name := "evaluation.feature"

	testSuite := godog.TestSuite{
		Name: name,
		TestSuiteInitializer: integration.InitializeTestSuite(func() openfeature.FeatureProvider {
			return flagd.NewProvider(flagd.WithInProcessResolver(), flagd.WithPort(9090))
		}),
		ScenarioInitializer: integration.InitializeEvaluationScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../spec/specification/assets/gherkin/evaluation.feature"},
			TestingT: t, // Testing instance that will run subtests.
			Strict:   true,
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run evaluation tests")
	}
}
