package components

import (
	"testing"

	v1 "github.com/openshift-eng/ci-test-mapping/pkg/api/types/v1"
	"github.com/openshift-eng/ci-test-mapping/pkg/registry"
)

// TestJiraTagPriorityResolvesConflict verifies that the [Jira:X] tag
// match in FindMatch carries Priority 1, so it wins over generic
// substring matchers (Priority 0) from other components.
//
// Each synthetic test name contains both a [Jira:Hypershift] tag and a
// container name substring that another component would also claim.
// Before the fix, both claims had Priority 0, causing getHighestPriority
// to return "unable to resolve conflict". With the fix, the Jira-tag
// match at Priority 1 wins cleanly.
func TestJiraTagPriorityResolvesConflict(t *testing.T) {
	containers := []string{
		"cluster-network-operator",
		"cluster-etcd-operator",
		"cluster-monitoring-operator",
		"cluster-dns-operator",
		"insights-operator",
		"machine-config-operator",
		"oauth-proxy",
		"cloud-provider-aws-e2e",
		"cluster-version-operator",
	}

	componentRegistry := registry.NewComponentRegistry()
	ti := NewTestIdentifier(componentRegistry, nil)

	for _, container := range containers {
		t.Run(container, func(t *testing.T) {
			testName := "[sig-hypershift][Jira:Hypershift][Feature:ControlPlaneWorkloads] " +
				"Control Plane Workloads Container image pull policy " +
				container +
				" should have IfNotPresent pull policy for containers [control-plane-workloads]"
			testInfo := &v1.TestInfo{
				Name:  testName,
				Suite: "hypershift-e2e",
			}

			ownership, err := ti.Identify(testInfo)
			if err != nil {
				t.Fatalf("Identify() returned error: %v", err)
			}
			if ownership.Component != "HyperShift" {
				t.Errorf("Identify() assigned to %q, want %q", ownership.Component, "HyperShift")
			}
		})
	}
}
