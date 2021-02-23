package gcpsecretsmanager_test

import (
	"context"
	"testing"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/gcpsecretsmanager"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes/fake"
)

type credentialsFile struct {
	Type string `json:"type"` // serviceAccountKey or userCredentialsKey

	// Service Account fields
	ClientEmail  string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	TokenURL     string `json:"token_uri"`
	ProjectID    string `json:"project_id"`
}

func TestGcpSecretStoreAuthFromValue(t *testing.T) {
	credObject := &credentialsFile{
		Type:      "service_account",
		ProjectID: "projectabcd",
	}
	credBytes, err := json.Marshal(credObject)
	creds, err := gcpsecretsmanager.NewGcpSecretStoreAuth(context.TODO(), nil, v1alpha1.GcpSecretsManagerAuth{
		CredentialsFile: &v1alpha1.ValueOrSecretKey{
			Value: func() *string { s := string(credBytes); return &s }(),
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, creds)
	assert.Equal(t, credObject.ProjectID, creds.ProjectID)
}

func TestGcpSecretStoreAuthFromSecret(t *testing.T) {
	credObject := &credentialsFile{
		Type:      "service_account",
		ProjectID: "project1234",
	}
	credBytes, err := json.Marshal(credObject)
	kubeClient := fake.NewSimpleClientset([]runtime.Object{
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gcp-secret",
				Namespace: "rarity",
			},
			Data: map[string][]byte{
				"serviceAccountKey": credBytes,
			},
		},
	}...)
	creds, err := gcpsecretsmanager.NewGcpSecretStoreAuth(context.TODO(), kubeClient, v1alpha1.GcpSecretsManagerAuth{
		CredentialsFile: &v1alpha1.ValueOrSecretKey{
			SecretRef: &v1alpha1.SecretRef{
				Namespace: "rarity",
				Name:      "gcp-secret",
				Key:       "serviceAccountKey",
			},
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, creds)
	assert.Equal(t, credObject.ProjectID, creds.ProjectID)
}
