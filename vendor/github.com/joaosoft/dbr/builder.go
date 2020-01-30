package dbr

type Builder interface {
	Build() (string, error)
}

type functionBuilder interface {
	Build(db *db) (string, error)
}
