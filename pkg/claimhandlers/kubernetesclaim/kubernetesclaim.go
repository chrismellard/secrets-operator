package kubernetesclaim

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/claimhandlers"
	"github.com/chrismellard/secret-operator/pkg/clients/kube"
	"github.com/chrismellard/secret-operator/pkg/source"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type handler struct {
	ctx   context.Context
	claim v1alpha1.SecretClaim
}

func (h handler) Handle() error {
	kubernetesClaim := h.claim.Spec.KubernetesClaim
	secretProperties := map[string][]byte{}
	for _, property := range kubernetesClaim.Properties {
		sourcedProperty, err := source.HandleProperty(property.PropertySource)
		if err != nil {
			return fmt.Errorf("error sourcing property %s", property.Name)
		}
		sourcePropertyBytes := []byte(sourcedProperty)
		encodedSecret := make([]byte, base64.StdEncoding.EncodedLen(len(sourcePropertyBytes)))
		base64.StdEncoding.Encode(encodedSecret, sourcePropertyBytes)
		secretProperties[property.Name] = encodedSecret
	}

	secret := createSecret(h.claim, secretProperties)
	err := applySecret(h.ctx, secret, h.claim)
	if err != nil {
		return fmt.Errorf("error when applying secret %w", err)
	}

	return nil
}

func applySecret(ctx context.Context, secret v1.Secret, claim v1alpha1.SecretClaim) error {

	client, err := kube.CreateClientSet()
	if err != nil {
		return err
	}

	secretClient := client.CoreV1().Secrets(secret.Namespace)
	existingSecret, err := secretClient.Get(ctx, secret.Name, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		_, err := secretClient.Create(ctx, &secret, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("error creating secret %s: %w", secret.Name, err)
		}
	} else {
		if !checkOwnership(claim, existingSecret.OwnerReferences) {
			return fmt.Errorf("existing secret %s is not owned by this claim %s", secret.Name, claim.Name)
		}
		_, err := secretClient.Update(ctx, &secret, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("error updating secret %s: %w", secret.Name, err)
		}
	}
	return nil
}

func createSecret(claim v1alpha1.SecretClaim, secretProperties map[string][]byte) v1.Secret {
	kubernetesClaim := claim.Spec.KubernetesClaim
	ownerRef := createOwnerReference(claim)
	return v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            kubernetesClaim.Name,
			Namespace:       kubernetesClaim.Namespace,
			Labels:          kubernetesClaim.Labels,
			Annotations:     kubernetesClaim.Annotations,
			OwnerReferences: []metav1.OwnerReference{ownerRef},
		},
		Data: secretProperties,
		Type: kubernetesClaim.SecretType,
	}
}

func NewHandler(claim v1alpha1.SecretClaim, ctx context.Context) claimhandlers.ClaimHandler {
	return &handler{ctx: ctx, claim: claim}
}

func createOwnerReference(claim v1alpha1.SecretClaim) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: claim.TypeMeta.APIVersion,
		Kind:       claim.TypeMeta.Kind,
		Name:       claim.ObjectMeta.Name,
		UID:        claim.UID,
	}
}

func checkOwnership(claim v1alpha1.SecretClaim, ownerRefs []metav1.OwnerReference) bool {
	for _, ownerRef := range ownerRefs {
		if ownerRef.Kind != claim.Kind {
			return false
		}
		if ownerRef.APIVersion != claim.APIVersion {
			return false
		}
		if ownerRef.Name != claim.Name {
			return false
		}
		if ownerRef.UID != claim.UID {
			return false
		}
		return true
	}
	return false
}
