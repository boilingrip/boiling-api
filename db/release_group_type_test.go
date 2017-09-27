package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllReleaseGroupTypes(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	types, err := db.GetAllReleaseGroupTypes()
	require.Nil(t, err)
	require.NotNil(t, types)
	require.NotEmpty(t, types)
}
