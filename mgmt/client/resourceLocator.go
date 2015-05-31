package client

type ResourceLocator interface {
	Get(string) ([]byte, error)
}
