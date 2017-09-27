package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateGetPasskey(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	pk, err := db.GenerateNewPasskeyForUser(1)
	require.Nil(t, err)
	require.NotEmpty(t, pk)

	got, err := db.GetPasskeyForUser(1)
	require.Nil(t, err)
	require.NotNil(t, got)
	require.Equal(t, 1, got.Uid)
	require.Equal(t, pk, got.Passkey)
	require.Equal(t, true, got.Valid)
}

func TestGetAllPasskeysForUser(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	pks1, err := db.GetAllPasskeysForUser(1)
	require.Nil(t, err)

	_, err = db.GenerateNewPasskeyForUser(1)
	require.Nil(t, err)

	pks2, err := db.GetAllPasskeysForUser(1)
	require.Nil(t, err)
	require.Equal(t, len(pks1)+1, len(pks2))
}
