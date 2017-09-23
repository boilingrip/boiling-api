package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllMedia(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	media, err := db.GetAllMedia()
	require.Nil(t, err)
	require.NotNil(t, media)
	require.NotEmpty(t, media)
}
