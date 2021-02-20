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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type PasswordGenerator struct {
	// +kubebuilder:default=12
	Length int `json:"length,omitempty"`
	// +kubebuilder:default=2
	NumSymbols int `json:"numSymbols,omitempty"`
	// +kubebuilder:default=2
	NumDigits int `json:"numDigits,omitempty"`
	// +kubebuilder:default=false
	NoUpper bool `json:"noUpper,omitempty"`
	// +kubebuilder:default=true
	AllowRepeat bool `json:"allowRepeat,omitempty"`
	// +kubebuilder:default="~!#%^_+-=?,."
	AllowedSymbols string `json:"allowedSymbols,omitempty"`
}

type PropertyGenerator struct {
	Hmac     bool               `json:"hmac,omitempty"`
	Password *PasswordGenerator `json:"password,omitempty"`
}
type PropertySource struct {
	PropertyGenerator *PropertyGenerator `json:"generator,omitempty"`
}

type SecretClaimProperty struct {
	Name           string         `json:"name,omitempty"`
	PropertySource PropertySource `json:"source,omitempty"`
}

type KubernetesClaim struct {
	Name        string                `json:"name,omitempty"`
	Namespace   string                `json:"namespace,omitempty"`
	SecretType  v1.SecretType         `json:"secretType,omitempty"`
	Labels      map[string]string     `json:"labels,omitempty"`
	Annotations map[string]string     `json:"annotations,omitempty"`
	Properties  []SecretClaimProperty `json:"properties,omitempty"`
}

// SecretClaimSpec defines the desired state of SecretClaim
type SecretClaimSpec struct {
	KubernetesClaim *KubernetesClaim `json:"kubernetes,omitempty"`
}

// SecretClaimStatus defines the observed state of SecretClaim
type SecretClaimStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// SecretClaim is the Schema for the secretclaims API
type SecretClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretClaimSpec   `json:"spec,omitempty"`
	Status SecretClaimStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecretClaimList contains a list of SecretClaim
type SecretClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretClaim{}, &SecretClaimList{})
}
