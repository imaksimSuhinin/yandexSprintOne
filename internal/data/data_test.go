package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
	db := NewDataBase()

	testValue := "Max"
	db.Write("name", testValue)
	usernameReal, err := db.ReadValue("name")

	require.NoError(t, err)
	require.Equal(t, testValue, usernameReal)
}
