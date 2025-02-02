package coffy

import (
	"errors"
	"testing"
)

// Working configuration, must contain all config parameters
var validConfig = `
server:
    port: 8080
database:
    path: ./coffy_path/coffy_machine.db
`

var missingServer = `
database:
    path: ./coffy_path/coffy_machine.db
`

var missingServerPort = `
server: 
    unknown: any
database:
    path: ./coffy_path/coffy_machine.db
`

var missingDatabase = `
server:
    port: 8080
`

var missingDatabasePath = `
server:
    port: 8080
database:
    unknown: any
`

var coffyDBpath = "./coffy_path/coffy_machine.db"

func TestParse(t *testing.T) {
	config, err := Parse(validConfig)
	if err != nil {
		t.Errorf("couldn't parse config: %v", err)
		return
	}
	if config.Server.Port != 8080 {
		t.Errorf("invalid server port: %v", config.Server.Port)
	}
	if config.Database.Path != coffyDBpath {
		t.Errorf("invalid database path: %v", config.Database.Path)
	}

}

func TestParseMissingServer(t *testing.T) {
	_, err := Parse(missingServer)
	if err == nil {
		t.Errorf("Expected error for missing server port")
		return
	}
	var expectedErr = &MissingPropertyError{}
	if !errors.As(err, expectedErr) {
		t.Errorf("Expected missing property error, got: %v", err)
		return
	}
}

func TestParseMissingServerPort(t *testing.T) {
	_, err := Parse(missingServerPort)
	if err == nil {
		t.Errorf("Expected error for missing server port")
		return
	}
	var expectedErr = &MissingPropertyError{}
	if !errors.As(err, expectedErr) {
		t.Errorf("Expected missing property error, got: %v", err)
		return
	}
	if err.Error() != "missing property: 'port'" {
		t.Errorf("Expected message: missing property: 'port', got: %v", err)
	}
}

func TestParseMissingDatabase(t *testing.T) {
	_, err := Parse(missingDatabase)
	if err == nil {
		t.Errorf("Expected error for missing database config entry")
	}
	var expectedErr = &MissingPropertyError{}
	if !errors.As(err, expectedErr) {
		t.Errorf("Expected missing property error, got: %v", err)
		return
	}
}

func TestParseMissingDatabasePath(t *testing.T) {
	_, err := Parse(missingDatabasePath)
	if err == nil {
		t.Errorf("Expected error for missing database config entry")
	}
	var expectedErr = &MissingPropertyError{}
	if !errors.As(err, expectedErr) {
		t.Errorf("Expected missing property error, got: %v", err)
		return
	}
	if err.Error() != "missing property: 'path'" {
		t.Errorf("Expected message: 'missing property: 'path'', got: %v", err)
	}
}
