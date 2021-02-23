package azurekeyvault_test

import (
	"context"
	"testing"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/secretstores/azurekeyvault"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestAzureAuth(t *testing.T) {

	type TestCase struct {
		Secrets      []*v1.Secret
		ProviderAuth v1alpha1.AzureKeyVaultProviderAuth
		Result       testAzureSecretStoreAuth
		ShouldError  bool
	}

	testCases := []TestCase{
		{
			ProviderAuth: v1alpha1.AzureKeyVaultProviderAuth{
				UseManagedIdentity: false,
				SubscriptionId:     v1alpha1.ValueOrSecretKey{Value: func() *string { s := "4096edc0-2b72-4577-9023-9813e4ac44f6"; return &s }()},
				TenantId:           v1alpha1.ValueOrSecretKey{Value: func() *string { s := "79888d8a-d066-45b0-be6a-da686a6640b9"; return &s }()},
				ClientId:           &v1alpha1.ValueOrSecretKey{Value: func() *string { s := "2c9e56ef-e36e-4d56-b4e1-cfe4e5db51f4"; return &s }()},
				ClientSecret:       &v1alpha1.ValueOrSecretKey{Value: func() *string { s := "sfk247j-sdfk29-fsdfm2kf0-*&&^"; return &s }()},
			},
			Result: testAzureSecretStoreAuth{
				clientId:           "2c9e56ef-e36e-4d56-b4e1-cfe4e5db51f4",
				clientSecret:       "sfk247j-sdfk29-fsdfm2kf0-*&&^",
				tenantId:           "79888d8a-d066-45b0-be6a-da686a6640b9",
				subscriptionId:     "4096edc0-2b72-4577-9023-9813e4ac44f6",
				useManagedIdentity: false,
			},
			ShouldError: false,
		},
		{
			ProviderAuth: v1alpha1.AzureKeyVaultProviderAuth{
				UseManagedIdentity: false,
				SubscriptionId: v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "twilight",
						Name:      "azure-creds",
						Key:       "subscriptionId",
					},
				},
				TenantId: v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "twilight",
						Name:      "azure-creds",
						Key:       "tenantId",
					},
				},
				ClientId: &v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "twilight",
						Name:      "azure-creds",
						Key:       "clientId",
					},
				},
				ClientSecret: &v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "twilight",
						Name:      "azure-creds",
						Key:       "clientSecret",
					},
				},
			},
			Secrets: []*v1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "azure-creds",
						Namespace: "twilight",
					},
					Data: map[string][]byte{
						"clientSecret":   []byte("08daec26874b495cbfb199b9b93cbae4"),
						"clientId":       []byte("b646a450-22db-4787-a241-9624d89397ca"),
						"tenantId":       []byte("9ff84e09-46e4-45a4-a914-cf390e91d66f"),
						"subscriptionId": []byte("bbdf2688-8d1e-43bc-8c18-862a01582c80"),
					},
				},
			},
			Result: testAzureSecretStoreAuth{
				clientId:           "b646a450-22db-4787-a241-9624d89397ca",
				clientSecret:       "08daec26874b495cbfb199b9b93cbae4",
				tenantId:           "9ff84e09-46e4-45a4-a914-cf390e91d66f",
				subscriptionId:     "bbdf2688-8d1e-43bc-8c18-862a01582c80",
				useManagedIdentity: false,
			},
			ShouldError: false,
		},
		{
			ProviderAuth: v1alpha1.AzureKeyVaultProviderAuth{
				UseManagedIdentity: false,
				SubscriptionId: v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "pinkiepie",
						Name:      "azure-creds",
						Key:       "subscriptionId",
					},
				},
				TenantId: v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "pinkiepie",
						Name:      "azure-creds",
						Key:       "tenantId",
					},
				},
				ClientId: &v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "pinkiepie",
						Name:      "azure-creds",
						Key:       "clientId",
					},
				},
				ClientSecret: &v1alpha1.ValueOrSecretKey{
					Value: nil,
					SecretRef: &v1alpha1.SecretRef{
						Namespace: "pinkiepie",
						Name:      "azure-creds",
						Key:       "clientSecret",
					},
				},
			},
			Secrets: []*v1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "azure-creds",
						Namespace: "twilight",
					},
					Data: map[string][]byte{
						"tenantId":       []byte("9ff84e09-46e4-45a4-a914-cf390e91d66f"),
						"subscriptionId": []byte("bbdf2688-8d1e-43bc-8c18-862a01582c80"),
					},
				},
			},
			ShouldError: true,
		},
	}

	toRuntimeObjects := func(secrets ...*v1.Secret) []runtime.Object {
		var out []runtime.Object
		for _, s := range secrets {
			out = append(out, s)
		}
		return out
	}

	for _, testCase := range testCases {

		runtimeObjects := toRuntimeObjects(testCase.Secrets...)

		intArr := []int{1, 2}
		intAppend := []int{3, 4}
		intArr = append(intArr, intAppend...)

		kubeClient := fake.NewSimpleClientset(runtimeObjects...)
		result, err := azurekeyvault.NewAzureSecretStoreAuth(context.TODO(), kubeClient, testCase.ProviderAuth)
		if testCase.ShouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, testCase.Result.UseManagedIdentity(), result.UseManagedIdentity())
			assert.Equal(t, testCase.Result.ClientSecret(), result.ClientSecret())
			assert.Equal(t, testCase.Result.ClientID(), result.ClientID())
			assert.Equal(t, testCase.Result.SubscriptionID(), result.SubscriptionID())
			assert.Equal(t, testCase.Result.TenantID(), result.TenantID())
		}

	}
}

type testAzureSecretStoreAuth struct {
	clientId           string
	clientSecret       string
	tenantId           string
	subscriptionId     string
	useManagedIdentity bool
}

func (t testAzureSecretStoreAuth) ClientID() string {
	return t.clientId
}

func (t testAzureSecretStoreAuth) ClientSecret() string {
	return t.clientSecret
}

func (t testAzureSecretStoreAuth) TenantID() string {
	return t.tenantId
}

func (t testAzureSecretStoreAuth) SubscriptionID() string {
	return t.subscriptionId
}

func (t testAzureSecretStoreAuth) UseManagedIdentity() bool {
	return t.useManagedIdentity
}
