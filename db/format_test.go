package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllFormats(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	f, err := db.GetAllFormats()
	require.Nil(t, err)
	require.NotNil(t, f)
	require.NotEmpty(t, f)
}
