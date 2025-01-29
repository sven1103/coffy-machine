package coffy

import (
	"errors"
	"testing"
)

var validConfig = `
server:
    port: 8080
`

var missingServer = `
`

func TestParse(t *testing.T) {
	config, err := Parse(validConfig)
	if err != nil {
		t.Errorf("couldn't parse config: %v", err)
		return
	}
	if config.Server.Port != 8080 {
		t.Errorf("invalid server port: %v", config.Server.Port)
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
