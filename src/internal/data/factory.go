package data

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/divideandconquer/go-consul-client/src/balancer"
)

// Factory can get a DB
type Factory interface {
	Get() (Database, error)
}

type factory struct {
	lb          balancer.DNS
	dbUser      string
	dbPass      string
	serviceName string
	log         *logrus.Entry
}

// Get fetches a DB
func (f factory) Get() (Database, error) {
	//setup db
	dbLoc, err := f.lb.FindService(fmt.Sprintf("%s-db", f.serviceName))
	if err != nil {
		f.log.Printf("Error finding to db: %v", err)
		return nil, err
	}

	db, err := NewDb("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", f.dbUser, f.dbPass, dbLoc.URL, dbLoc.Port, f.serviceName))
	if err != nil {
		f.log.Printf("Error connecting to db: %v", err)
		return nil, err
	}
	return db, nil
}

// GetDBFactory sets up the db factory
func GetDBFactory(lb balancer.DNS, dbUser, dbPass, serviceName string, log *logrus.Entry) Factory {
	return factory{lb: lb, dbUser: dbUser, dbPass: dbPass, serviceName: serviceName, log: log}
}
