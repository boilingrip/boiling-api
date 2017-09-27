package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllPrivileges(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	p, err := db.GetAllPrivileges()
	require.Nil(t, err)
	require.NotNil(t, p)
	require.NotEmpty(t, p)
}
