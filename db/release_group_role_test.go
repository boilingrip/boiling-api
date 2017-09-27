package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllReleaseGroupRoles(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	roles, err := db.GetAllReleaseGroupRoles()
	require.Nil(t, err)
	require.NotNil(t, roles)
	require.NotEmpty(t, roles)
}
