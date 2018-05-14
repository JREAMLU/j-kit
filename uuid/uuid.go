package uuid

import (
	uuid "github.com/satori/go.uuid"
)

// Generate generate
func Generate() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
