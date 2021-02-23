package health

import (
	"fmt"
	"time"

	"github.com/chrismellard/secret-operator/pkg/secretstores/api"
	"github.com/jenkins-x-plugins/secretfacade/pkg/secretstore"
)

const (
	secretName = "lastUpdated"
)

func CheckSecretStoreHealth(store secretstore.Interface, location api.SecretLocation) (bool, error) {
	now := time.Now().String()

	err := store.SetSecret(location.Location(), secretName, &secretstore.SecretValue{
		Value: now,
	})
	if err != nil {
		return false, fmt.Errorf("unable to set health secret: %w", err)
	}

	secret, err := store.GetSecret(location.Location(), secretName, "")
	if err != nil {
		return false, fmt.Errorf("unable to get health secret: %w", err)
	}

	if secret != now {
		return false, fmt.Errorf("retrieved secret differs from expected value, received %s but expectedd %s", secret, now)
	}

	return true, nil
}
