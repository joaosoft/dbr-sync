package main

import (
	"database/sql"
	"dbr-sync"
	"fmt"
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

	// start sync system
	dbrSync, err := dbr_sync.NewDbrSync()
	if err != nil {
		panic(err)
	}

	if err = dbrSync.Start(); err != nil {
		panic(err)
	}

	// test with dbr
	var db, _ = dbr.New(
		dbr.WithConfiguration(&dbr.DbrConfig{
			Db: &manager.DBConfig{
				DataSource: "postgres://user:password@localhost:7000/postgres?sslmode=disable&search_path=dbr-sync",
				Driver:     "postgres",
			},
		}),
		dbr.WithSuccessEventHandler(HandleSuccessEvent),
	)

	now := time.Now()
	example := &Example{
		IdExample:   "e4c15bfb-3aee-4477-b6f2-d5eb75c1f119",
		Name:        "joao",
		Description: "my first test",
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = db.
		Insert().
		Into("example").
		Record(example).
		Exec()

	if err != nil {
		panic(err)
	}

	<-time.After(10 * time.Second)
	dbrSync.Stop()
}

func HandleSuccessEvent(operation dbr.SqlOperation, table []string, query string, rows *sql.Rows, sqlResult sql.Result) error {
	fmt.Printf("\nSuccess event [operation: %s, tables: %s, query: %s]", operation, strings.Join(table, "; "), query)

	pm := manager.NewManager()

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s%s", "root", "password", "localhost", "5673", "/local")
	configRabbitmq := manager.NewRabbitmqConfig(uri, "dbr-sync-exchange", "direct")

	rabbitmqProducer, err := pm.NewSimpleRabbitmqProducer(configRabbitmq)
	if err != nil {
		return err
	}

	if err := rabbitmqProducer.Start(); err != nil {
		return err
	}

	oper := &dbr_sync.Operation{
		Operation: operation,
		Query:     &query,
	}

	opers := dbr_sync.OperationList{oper}

	message, err := json.Marshal(opers)
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
