package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const (
	validUserID = "353218f0-515b-48cc-9eeb-3715c4cf7137"
)

func makeStringPtr(value string) *string {
	return &value
}

func makeInt32Ptr(value int32) *int32 {
	return &value
}

func GetDatabase(t *testing.T) Database {
	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		t.Fatalf("Could not find DATABASE_HOST environment variable")
	}
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		t.Fatalf("Could not find SERVICE_NAME environment variable")
	}
	log.Printf("Connecting to database at %s", host)
	url := fmt.Sprintf("postgres://postgres:password@%s/%s?sslmode=disable", host, serviceName)
	ret, err := NewDb("postgres", url)
	if err != nil {
		t.Fatalf("Could not get db connection: %v", err)
	}
	return ret
}

func TestIntegration_NewDB_BasePath(t *testing.T) {
	db := GetDatabase(t)
	if db == nil {
		t.Fatalf("DB object was nil")
	}
}

func TestIntegration_ClearTestData(t *testing.T) {
	db := GetDatabase(t)
	err := db.ClearTestData()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

}

func TestIntegration_CreateGoal_BasePath(t *testing.T) {
	db := GetDatabase(t)
	db.ClearTestData()
	defer db.ClearTestData()

	start, _ := time.Parse(time.RFC3339, "2020-07-01T00:00:00Z")
	expire, _ := time.Parse(time.RFC3339, "2020-07-07T00:00:00Z")

	err := db.CreateGoal(context.Background(), validUserID, "point", "weekly", 250, nil, nil, start, expire)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
