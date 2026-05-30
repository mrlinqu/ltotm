package webserver

type storage interface {
	Put(id string, text []byte) error
	Get(id string) ([]byte, error)
	Del(id string) error
}
