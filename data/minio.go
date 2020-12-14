package data

type MinioRepository interface {
	MakeBucket(name string) (bool, error)
}
