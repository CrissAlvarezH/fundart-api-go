package ports

type ImageFile interface {
	Save() (string, error)
	Delete(path string) error
}
