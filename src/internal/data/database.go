package data

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/healthimation/go-glitch/glitch"
	"github.com/healthimation/go-glitch/postgres"
)

// Database is an interface to the data store
type Database interface {
	// ClearTestData will clear out test data
	ClearTestData() glitch.DataError
}

// keeps track of the db pool for each connection string
// if the db moves and new connection string comes in it will create a new pool, otherwise we keep the existing one
var connMap map[string]*sql.DB
var mapLock sync.RWMutex

func init() {
	mapLock.Lock()
	defer mapLock.Unlock()
	connMap = make(map[string]*sql.DB)
}

// getDbFromMap is a thread safe way to get a connection from the map or create a new one
func getDbFromMap(driver, conn string) (result *sql.DB, err error) {
	key := fmt.Sprintf("%s|%s", driver, conn)
	var ok bool
	mapLock.RLock()
	result, ok = connMap[key]
	mapLock.RUnlock()
	if ok {
		return
	}
	// need to create a connection
	mapLock.Lock()
	defer mapLock.Unlock()
	result, ok = connMap[key] // check to make sure something didnt just create it
	if ok {
		return
	}
	result, err = sql.Open(driver, conn)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to DB: %v", err)
	}
	result.SetMaxOpenConns(30)
	result.SetMaxIdleConns(30)
	connMap[key] = result
	return
}

type db struct {
	conn *sql.DB
}

// NewDb will return an initialized Database object
func NewDb(driver, conn string) (Database, error) {
	d, err := getDbFromMap(driver, conn)
	if err != nil {
		return nil, err
	}
	return &db{conn: d}, nil
}

// ClearTestData will clear out test data
func (d *db) ClearTestData() glitch.DataError {
	query := "SELECT * FROM clear_test_data()"
	_, err := d.conn.Exec(query)
	return postgres.ToDataError(err, "Error executing query")
}
