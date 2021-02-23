package gcpsecretsmanager

import (
	"context"
	"fmt"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/helpers"
	"golang.org/x/oauth2/google"
	"k8s.io/client-go/kubernetes"
)

func NewGcpSecretStoreAuth(ctx context.Context, client kubernetes.Interface, providerAuth v1alpha1.GcpSecretsManagerAuth) (google.Credentials, error) {
	if providerAuth.CredentialsFile != nil {
		credentialJson, err := helpers.ValueOrSecretRef(ctx, providerAuth.CredentialsFile, client)
		if err != nil {
			return google.Credentials{}, fmt.Errorf("expected to be able set credentialsFile but failed: %w", err)
		}
		creds, err := google.CredentialsFromJSON(ctx, []byte(credentialJson))
		if err != nil {
			return google.Credentials{}, fmt.Errorf("unable to configure credentials object from credentials json: %w", err)
		}
		return *creds, nil
	}
	return google.Credentials{}, fmt.Errorf("unable to create google credentials from provider config")
}
