package stores

type GocialStore interface {
	Save(state string, gocial []byte) error
	Get(state string) ([]byte, error)
	Delete(state string) error
}
