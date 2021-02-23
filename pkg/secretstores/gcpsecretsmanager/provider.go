package gcpsecretsmanager

import "github.com/chrismellard/secret-operator/api/v1alpha1"

type GcpSecretsManagerProvider struct {
	v1alpha1.GcpSecretsManagerProvider
}

func (a GcpSecretsManagerProvider) Location() string {
	return a.ProjectId
}
