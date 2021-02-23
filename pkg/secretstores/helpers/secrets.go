package helpers

import (
	"context"
	"fmt"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ValueOrSecretRef(ctx context.Context, r *v1alpha1.ValueOrSecretKey, client kubernetes.Interface) (string, error) {
	if r == nil {
		return "", nil
	}
	if r.Value != nil {
		return *r.Value, nil
	}
	if r.SecretRef != nil {
		secret, err := client.CoreV1().Secrets(r.SecretRef.Namespace).Get(ctx, r.SecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("unable to retrieve secret %s from namespace %s: %w", r.SecretRef.Name, r.SecretRef.Namespace, err)
		}
		if value, ok := secret.Data[r.SecretRef.Key]; !ok {
			return "", fmt.Errorf("unable to find key %s in secret %s in namespace %s", r.SecretRef.Key, r.SecretRef.Name, r.SecretRef.Namespace)
		} else {
			return string(value), nil
		}
	}
	return "", nil
}
