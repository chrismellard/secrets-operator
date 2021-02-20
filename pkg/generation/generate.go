package generation

import (
	"fmt"

	"github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/generation/generators/hmac"
	"github.com/chrismellard/secret-operator/pkg/generation/generators/password"
)

func Generate(propertyGenerator v1alpha1.PropertyGenerator) (string, error) {
	if propertyGenerator.Hmac {
		return hmac.Hmac()
	} else if propertyGenerator.Password != nil {
		return password.GeneratePassword(
			propertyGenerator.Password.Length,
			propertyGenerator.Password.AllowedSymbols,
			propertyGenerator.Password.NumDigits,
			propertyGenerator.Password.NumSymbols,
			propertyGenerator.Password.AllowRepeat,
			propertyGenerator.Password.NoUpper)
	}
	return "", fmt.Errorf("unable to determine property generator")
}
