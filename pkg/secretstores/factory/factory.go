package factory

import (
	"context"
	"fmt"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/api"
	"github.com/chrismellard/secret-operator/pkg/secretstores/azurekeyvault"
	"github.com/chrismellard/secret-operator/pkg/secretstores/gcpsecretsmanager"
	"github.com/jenkins-x-plugins/secretfacade/pkg/secretstore"
	"github.com/jenkins-x-plugins/secretfacade/pkg/secretstore/azuresecrets"
	gcpsecrets "github.com/jenkins-x-plugins/secretfacade/pkg/secretstore/gcpsecretsmanager"
	"k8s.io/client-go/kubernetes"
)

func CreateSecretStoreProvider(v v1alpha1.Provider) (api.SecretStoreProvider, error) {
	if v.AzureKeyVault != nil {
		return azurekeyvault.AzureKeyVaultProvider{AzureKeyVaultProvider: *v.AzureKeyVault}, nil
	} else if v.GcpSecretsManager != nil {
		return gcpsecretsmanager.GcpSecretsManagerProvider{GcpSecretsManagerProvider: *v.GcpSecretsManager}, nil
	}
	return nil, fmt.Errorf("unable to determine which SecretStoreProvider to instantiate")
}

func CreateSecretStoreManager(ctx context.Context, provider v1alpha1.Provider, kubeClient kubernetes.Interface) (secretstore.Interface, error) {

	if provider.AzureKeyVault != nil {
		auth, err := azurekeyvault.NewAzureSecretStoreAuth(ctx, kubeClient, provider.AzureKeyVault.Auth)
		if err != nil {
			return nil, fmt.Errorf("unable to create Azure Key Vault credentials: %w", err)
		}
		return azuresecrets.NewAzureKeyVaultSecretManager(auth), nil
	} else if provider.GcpSecretsManager != nil {
		creds, err := gcpsecretsmanager.NewGcpSecretStoreAuth(ctx, kubeClient, provider.GcpSecretsManager.Auth)
		if err != nil {
			return nil, fmt.Errorf("unable to create Gcp secret manager credentials: %w", err)
		}
		return gcpsecrets.NewGcpSecretsManager(creds), nil
	}
	return nil, fmt.Errorf("unable to create secret store manager")
}
