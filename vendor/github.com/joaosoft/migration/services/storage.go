package services

import (
	"database/sql"

	"fmt"

	errors "github.com/joaosoft/errors"
	"github.com/joaosoft/logger"
	manager "github.com/joaosoft/manager"
)

type StoragePostgres struct {
	logger logger.ILogger
	conn   manager.IDB
}

func NewStoragePostgres(logger logger.ILogger, connection manager.IDB) *StoragePostgres {
	return &StoragePostgres{
		logger: logger,
		conn:   connection,
	}
}

func (storage *StoragePostgres) GetMigration(idMigration string) (*Migration, error) {
	row := storage.conn.Get().QueryRow(`
	    SELECT
	    	mode,
		    "user",
			executed_at
		FROM migration
		WHERE id_migration = $1
		ORDER BY id_migration ASC
	`, idMigration)

	migration := &Migration{IdMigration: idMigration}
	if err := row.Scan(
		&migration.Mode,
		&migration.User,
		&migration.ExecutedAt); err != nil {

		if err != sql.ErrNoRows {
			return nil, errors.New(errors.ErrorLevel, 0, err)
		}

		return nil, nil
	}

	return migration, nil
}

func (storage *StoragePostgres) GetMigrations(values map[string][]string) (ListMigration, error) {
	query := `
	    SELECT
			id_migration,
			mode,
		    "user",
			executed_at
		FROM migration
	`

	index := 1
	params := make([]interface{}, 0)

	if values != nil {
		for key, value := range values {
			if index == 1 {
				query += ` WHERE `
			} else {
				query += ` AND `
			}
			query += fmt.Sprintf(`%s = $%d`, key, index)
			index = index + 1

			params = append(params, value[0])
		}
	}

	query += ` ORDER BY id_migration ASC`

	rows, err := storage.conn.Get().Query(query, params...)
	if err != nil {
		return nil, errors.New(errors.ErrorLevel, 0, err)
	}

	defer rows.Close()

	migrations := make(ListMigration, 0)
	for rows.Next() {
		migration := &Migration{}
		if err := rows.Scan(
			&migration.IdMigration,
			&migration.Mode,
			&migration.User,
			&migration.ExecutedAt); err != nil {

			if err != sql.ErrNoRows {
				return nil, errors.New(errors.ErrorLevel, 0, err)
			}
			return nil, nil
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func (storage *StoragePostgres) CreateMigration(newMigration *Migration) error {
	if _, err := storage.conn.Get().Exec(`
		INSERT INTO migration(
			id_migration,
			mode)
		VALUES($1, $2)
	`,
		newMigration.IdMigration, newMigration.Mode); err != nil {
		return errors.New(errors.ErrorLevel, 0, err)
	}

	return nil
}

func (storage *StoragePostgres) DeleteMigration(idMigration string) error {
	if _, err := storage.conn.Get().Exec(`
	    DELETE 
		FROM migration
		WHERE id_migration = $1
	`, idMigration); err != nil {
		return errors.New(errors.ErrorLevel, 0, err)
	}

	return nil
}

func (storage *StoragePostgres) DeleteMigrations() error {
	if _, err := storage.conn.Get().Exec(`
	    DELETE FROM migration`); err != nil {
		return errors.New(errors.ErrorLevel, 0, err)
	}

	return nil
}

func (storage *StoragePostgres) ExecuteMigration(migration string) error {
	if _, err := storage.conn.Get().Exec(migration); err != nil {
		return errors.New(errors.ErrorLevel, 0, err)
	}

	return nil
}
