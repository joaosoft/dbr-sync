package main

import (
	"database/sql"
	"dbr-sync"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"

	"github.com/joaosoft/dbr"
	"github.com/joaosoft/json"
	"github.com/joaosoft/manager"
)

type Example struct {
	IdExample   string    `json:"id_example" db:"id_example"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func main() {

	// start consumer
	dbrSync, err := dbr_sync.NewDbrSync()
	if err != nil {
		panic(err)
	}

	if err = dbrSync.Start(); err != nil {
		panic(err)
	}

	// start producer
	var db, _ = dbr.New(
		dbr.WithConfiguration(&dbr.DbrConfig{
			Db: &manager.DBConfig{
				DataSource: "postgres://user:password@localhost:7000/postgres?sslmode=disable&search_path=dbr-sync-origin",
				Driver:     "postgres",
			},
		}),
		dbr.WithSuccessEventHandler(HandleSuccessEventProducer),
	)

	if err = insert("first", db); err != nil {
		panic(err)
	}

	if err = insert("second", db); err != nil {
		panic(err)
	}

	if err = update("first", db); err != nil {
		panic(err)
	}

	if err = delete("second", db); err != nil {
		panic(err)
	}

	<-time.After(10 * time.Second)
	dbrSync.Stop()
}

func insert(name string, db *dbr.Dbr) error {
	id, _ := uuid.NewV4()
	now := time.Now()
	example := &Example{
		IdExample:   id.String(),
		Name:        name,
		Description: fmt.Sprintf("my %s test", name),
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := db.
		Insert().
		Into("example").
		Record(example).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func update(name string, db *dbr.Dbr) error {
	_, err := db.
		Update("example").
		Set("description", fmt.Sprintf("my %s test updated", name)).
		Where("name = ?", name).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func delete(name string, db *dbr.Dbr) error {
	_, err := db.
		Delete().
		From("example").
		Where("name = ?", name).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func HandleSuccessEventProducer(operation dbr.SqlOperation, table []string, query string, rows *sql.Rows, sqlResult sql.Result) error {
	fmt.Printf("\nSuccess event [operation: %s, tables: %s, query: %s]", operation, strings.Join(table, "; "), query)

	pm := manager.NewManager()

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s%s", "root", "password", "localhost", "5672", "/local")
	configRabbitmq := manager.NewRabbitmqConfig(uri, "dbr-sync-exchange", "direct")

	rabbitmqProducer, err := pm.NewSimpleRabbitmqProducer(configRabbitmq)
	if err != nil {
		return err
	}

	if err := rabbitmqProducer.Start(); err != nil {
		return err
	}

	var operList dbr_sync.OperationList
	mode := 2

	switch mode {
	case 1:
		operList = append(operList, &dbr_sync.Operation{
			Operation: operation,
			Query:     &query,
		})
	case 2:
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		for rows.Next() {
			if err = rows.Scan(values...); err != nil {
				return err
			}

			columnsWithValues := make(map[string]interface{})
			for i, column := range columns {
				columnsWithValues[column] = values[i]
			}

			operList = append(operList, &dbr_sync.Operation{
				Operation: operation,
				Details: &dbr_sync.Details{
					Table:  table[0],
					Values: columnsWithValues,
				},
			})
		}
	}

	message, err := json.Marshal(operList)
	fmt.Printf("HERE: %s", string(message))
	if err != nil {
		return err
	}

	err = rabbitmqProducer.Publish("new.sync", message, true)
	if err != nil {
		return err
	}

	if err := rabbitmqProducer.Stop(); err != nil {
		return err
	}

	return nil
}
