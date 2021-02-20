package claimhandlers

type ClaimHandler interface {
	Handle() error
}
