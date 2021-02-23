package azurekeyvault

import (
	"context"
	"fmt"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/helpers"
	"github.com/jenkins-x-plugins/secretfacade/pkg/iam/azureiam"
	"k8s.io/client-go/kubernetes"
)

type azureSecretStoreAuth struct {
	clientId           string
	clientSecret       string
	tenantId           string
	subscriptionId     string
	useManagedIdentity bool
}

func NewAzureSecretStoreAuth(ctx context.Context, client kubernetes.Interface, providerAuth v1alpha1.AzureKeyVaultProviderAuth) (azureiam.Credentials, error) {
	auth := azureSecretStoreAuth{}

	clientId, err := helpers.ValueOrSecretRef(ctx, providerAuth.ClientId, client)
	if err != nil {
		return nil, fmt.Errorf("expected to be able set clientId but failed: %w", err)
	}
	auth.clientId = clientId

	clientSecret, err := helpers.ValueOrSecretRef(ctx, providerAuth.ClientSecret, client)
	if err != nil {
		return nil, fmt.Errorf("expected to be able set clientSecret but failed: %w", err)
	}
	auth.clientSecret = clientSecret

	tenantId, err := helpers.ValueOrSecretRef(ctx, &providerAuth.TenantId, client)
	if err != nil {
		return nil, fmt.Errorf("expected to be able set tenantId but failed: %w", err)
	}
	auth.tenantId = tenantId

	subscriptionId, err := helpers.ValueOrSecretRef(ctx, &providerAuth.SubscriptionId, client)
	if err != nil {
		return nil, fmt.Errorf("expected to be able set tenantId but failed: %w", err)
	}
	auth.subscriptionId = subscriptionId

	auth.useManagedIdentity = providerAuth.UseManagedIdentity
	return auth, nil
}

func (p azureSecretStoreAuth) ClientID() string {
	return p.clientId
}

func (p azureSecretStoreAuth) ClientSecret() string {
	return p.clientSecret
}

func (p azureSecretStoreAuth) TenantID() string {
	return p.tenantId
}

func (p azureSecretStoreAuth) SubscriptionID() string {
	return p.subscriptionId
}

func (p azureSecretStoreAuth) UseManagedIdentity() bool {
	return p.useManagedIdentity
}
