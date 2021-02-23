/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type SecretRef struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Key       string `json:"key"`
}

type ValueOrSecretKey struct {
	Value     *string    `json:"value,omitempty"`
	SecretRef *SecretRef `json:"secretRef,omitempty"`
}

type AzureKeyVaultProviderAuth struct {
	UseManagedIdentity bool              `json:"useManagedIdentity,omitempty"`
	SubscriptionId     ValueOrSecretKey  `json:"subscriptionId"`
	TenantId           ValueOrSecretKey  `json:"tenantId"`
	ClientId           *ValueOrSecretKey `json:"clientId,omitempty"`
	ClientSecret       *ValueOrSecretKey `json:"clientSecret,omitempty"`
}

type AzureKeyVaultProvider struct {
	VaultName string                    `json:"vaultName"`
	Auth      AzureKeyVaultProviderAuth `json:"auth"`
}

type GcpWorkloadIdentity struct {
	ServiceAccount    string `json:"serviceAccount"`
	Namespace         string `json:"namespace"`
	GcpServiceAccount string `json:"gcpServiceAccount"`
}

type GcpSecretsManagerAuth struct {
	WorkloadIdentity *GcpWorkloadIdentity `json:"workloadIdentity,omitempty"`
	CredentialsFile  *ValueOrSecretKey    `json:"credentialsFile,omitempty"`
}

type GcpSecretsManagerProvider struct {
	ProjectId string                `json:"projectId"`
	Auth      GcpSecretsManagerAuth `json:"auth"`
}

type Provider struct {
	AzureKeyVault     *AzureKeyVaultProvider     `json:"azureKeyVault,omitempty"`
	GcpSecretsManager *GcpSecretsManagerProvider `json:"gsm,omitempty"`
}

// SecretStoreSpec defines the desired state of SecretStore
type SecretStoreSpec struct {
	Provider Provider `json:"provider"`
}

// SecretStoreStatus defines the observed state of SecretStore
type SecretStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// SecretStore is the Schema for the secretstores API
type SecretStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretStoreSpec   `json:"spec,omitempty"`
	Status SecretStoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretStoreList contains a list of SecretStore
type SecretStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretStore{}, &SecretStoreList{})
}
