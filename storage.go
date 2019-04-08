package profile

import (
	"github.com/joaosoft/dbr"
)

type StoragePostgres struct {
	config *ProfileConfig
	db     *dbr.Dbr
}

func NewStoragePostgres(config *ProfileConfig) (*StoragePostgres, error) {
	dbr, err := dbr.New(dbr.WithConfiguration(config.Dbr))
	if err != nil {
		return nil, err
	}

	return &StoragePostgres{
		config: config,
		db:     dbr,
	}, nil
}

func (storage *StoragePostgres) GetSections() (Sections, error) {
	sections := make(Sections, 0)

	_, err := storage.db.
		Select("*").
		From(dbr.As(profileTableSection, "s")).
		Where("s.active").
		Load(&sections)

	if err != nil {
		return nil, err
	}

	return sections, nil
}

func (storage *StoragePostgres) GetSection(sectionKey string) (*Section, error) {
	section := Section{}

	count, err := storage.db.
		Select("*").
		From(dbr.As(profileTableSection, "s")).
		Where("s.key = ?", sectionKey).
		Where("s.active").
		Load(&section)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	return &section, nil
}

func (storage *StoragePostgres) GetSectionContents(sectionKey string) (Contents, error) {
	contents := make(Contents, 0)

	_, err := storage.db.
		Select("*").
		From(dbr.As(profileTableSection, "s")).
		Join(dbr.As(profileTableContent, "c"), "c.fk_section = s.id_section").
		Where("s.active").
		Where("c.active").
		Where("s.key = ?", sectionKey).
		Load(&contents)

	if err != nil {
		return nil, err
	}

	return contents, nil
}
