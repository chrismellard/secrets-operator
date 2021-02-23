package api

type SecretLocation interface {
	Location() string
}

type SecretStoreProvider interface {
	SecretLocation
}
