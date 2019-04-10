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

func (storage *StoragePostgres) GetSections() (SectionList, error) {
	sections := make(SectionList, 0)

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

func (storage *StoragePostgres) GetSectionsContents() (SectionsContentsList, error) {
	sectionsContents := make(SectionsContentsList, 0)

	_, err := storage.db.
		Select(
			"*",
			dbr.OnNull(
				storage.db.Select(dbr.ArrayToJson(dbr.ArrayAgg(dbr.RowToJson("t")))).
					From(
						dbr.As(storage.db.Select("*").
							From(dbr.As(profileTableContent, "c")).
							Where("c.fk_section = s.id_section").
							Where("c.active").
							OrderAsc("c.position"), "t")),
				"[]", "contents")).
		From(dbr.As(profileTableSection, "s")).
		Where("s.active").
		OrderAsc("s.position").
		Load(&sectionsContents)

	if err != nil {
		return nil, err
	}

	return sectionsContents, nil
}

func (storage *StoragePostgres) GetSection(sectionKey string) (*Section, error) {
	section := Section{}

	count, err := storage.db.
		Select("*").
		From(dbr.As(profileTableSection, "s")).
		Where("s.key = ?", sectionKey).
		Where("s.active").
		OrderAsc("s.position").
		Load(&section)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	return &section, nil
}

func (storage *StoragePostgres) GetSectionContents(sectionKey string) (ContentList, error) {
	contents := make(ContentList, 0)

	_, err := storage.db.
		Select("*").
		From(dbr.As(profileTableSection, "s")).
		Join(dbr.As(profileTableContent, "c"), "c.fk_section = s.id_section").
		Where("s.active").
		Where("c.active").
		Where("s.key = ?", sectionKey).
		OrderAsc("c.position").
		Load(&contents)

	if err != nil {
		return nil, err
	}

	return contents, nil
}
