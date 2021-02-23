package azurekeyvault_test

import (
	"testing"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/azurekeyvault"
	"github.com/stretchr/testify/assert"
)

func TestAzureProvider(t *testing.T) {
	provider := azurekeyvault.AzureKeyVaultProvider{
		AzureKeyVaultProvider: v1alpha1.AzureKeyVaultProvider{
			VaultName: "my-vault-name",
		},
	}
	assert.Equal(t, "my-vault-name", provider.Location())
}
