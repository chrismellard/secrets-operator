package gcpsecretsmanager_test

import (
	"testing"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/gcpsecretsmanager"
	"github.com/stretchr/testify/assert"
)

func TestGcpSecretsManagerProvider(t *testing.T) {
	provider := gcpsecretsmanager.GcpSecretsManagerProvider{
		GcpSecretsManagerProvider: v1alpha1.GcpSecretsManagerProvider{
			ProjectId: "exciting-new-project",
			Auth: v1alpha1.GcpSecretsManagerAuth{
				WorkloadIdentity: &v1alpha1.GcpWorkloadIdentity{
					ServiceAccount:    "",
					Namespace:         "",
					GcpServiceAccount: "",
				},
				CredentialsFile: &v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "",
						Name:      "",
						Key:       "",
					},
				},
			},
		},
	}
	assert.Equal(t, "exciting-new-project", provider.Location())
}
