package password

import (
	"github.com/pkg/errors"
	"github.com/sethvargo/go-password/password"
)

// DefaultGenerateSecret generates a secret using sensible defaults
func GeneratePassword(length int, allowedSymbols string, numDigits int, numSymbols int, allowRepeat bool, noUpper bool) (string, error) {
	input := password.GeneratorInput{
		Symbols: allowedSymbols,
	}

	generator, err := password.NewGenerator(&input)
	if err != nil {
		return "", errors.Wrap(err, "unable to create password generator")
	}

	secret, err := generator.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)

	if err != nil {
		return "", errors.Wrap(err, "unable to generate secret")
	}
	return secret, nil
}
