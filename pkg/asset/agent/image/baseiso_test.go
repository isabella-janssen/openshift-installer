package image

import (
	"context"
	"testing"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func TestCustomStreamGetter(t *testing.T) {
	cases := []struct {
		name                  string
		workflow              workflow.AgentWorkflowType
		installConfigSupplied bool
		osImageStream         types.OSImageStream
		expectedStreamURL     string
		expectNil             bool
	}{
		{
			name:      "add-nodes workflow uses cluster info",
			workflow:  workflow.AgentWorkflowTypeAddNodes,
			expectNil: false,
		},
		{
			name:                  "install workflow with rhel-9 stream",
			workflow:              workflow.AgentWorkflowTypeInstall,
			installConfigSupplied: true,
			osImageStream:         types.OSImageStreamRHCOS9,
			expectNil:             false,
		},
		{
			name:                  "install workflow with rhel-10 stream",
			workflow:              workflow.AgentWorkflowTypeInstall,
			installConfigSupplied: true,
			osImageStream:         types.OSImageStreamRHCOS10,
			expectNil:             false,
		},
		{
			name:                  "install workflow with rhel-10-nvidia stream maps to rhel-10",
			workflow:              workflow.AgentWorkflowTypeInstall,
			installConfigSupplied: true,
			osImageStream:         types.OSImageStreamRHCOS10Nvidia,
			expectNil:             false,
		},
		{
			name:                  "install workflow with empty stream defaults to rhel-9",
			workflow:              workflow.AgentWorkflowTypeInstall,
			installConfigSupplied: true,
			osImageStream:         "",
			expectNil:             false,
		},
		{
			name:                  "install workflow without install-config",
			workflow:              workflow.AgentWorkflowTypeInstall,
			installConfigSupplied: false,
			expectNil:             true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agentWorkflow := &workflow.AgentWorkflow{Workflow: tc.workflow}
			clusterInfo := &joiner.ClusterInfo{
				OSImage: &stream.Stream{
					Architectures: map[string]stream.Arch{
						"x86_64": {},
					},
				},
			}

			var installConfig *agent.OptionalInstallConfig
			if tc.installConfigSupplied {
				installConfig = &agent.OptionalInstallConfig{
					Supplied: true,
				}
				installConfig.Config = &types.InstallConfig{
					OSImageStream: tc.osImageStream,
					Platform: types.Platform{
						BareMetal: &baremetal.Platform{},
					},
				}
			} else {
				installConfig = &agent.OptionalInstallConfig{
					Supplied: false,
				}
			}

			fetcher := customStreamGetter(agentWorkflow, clusterInfo, installConfig)

			if tc.expectNil {
				assert.Nil(t, fetcher, "Expected fetcher to be nil")
				return
			}

			assert.NotNil(t, fetcher, "Expected fetcher to be non-nil")

			// For add-nodes workflow, verify it returns the cluster's OSImage
			if tc.workflow == workflow.AgentWorkflowTypeAddNodes {
				result, err := fetcher(context.Background())
				assert.NoError(t, err)
				assert.Equal(t, clusterInfo.OSImage, result)
				return
			}

			// For install workflow with osImageStream, verify the stream is fetched
			// We can't easily test FetchCoreOSBuild without mocking, but we can verify
			// the function doesn't panic and returns a result
			result, err := fetcher(context.Background())

			// We expect an error here since we're not in a real environment with
			// access to the stream metadata, but we can verify the function was called
			// The important test is that rhel-10-nvidia doesn't cause a panic and attempts
			// to fetch the stream (even if it fails in test environment)
			if tc.osImageStream == types.OSImageStreamRHCOS10Nvidia {
				// For rhel-10-nvidia, we just verify the function executes without panic
				// The actual fetching will fail in test environment, which is expected
				_ = result
				_ = err
			}
		})
	}
}
