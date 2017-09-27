package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddGetAllReleaseProperties(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	p, err := db.GetAllReleaseProperties()
	require.Nil(t, err)
	require.NotNil(t, p)
	require.NotEmpty(t, p)

	l := len(p)

	err = db.AddReleaseProperty("SquareVinylApproved")
	require.Nil(t, err)

	p, err = db.GetAllReleaseProperties()
	require.Nil(t, err)
	require.Equal(t, l+1, len(p))
}
