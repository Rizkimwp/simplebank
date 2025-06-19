package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
)

func TestCreateEntry(t *testing.T) {
	account1, err := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	arg := CreateEntryParams{
		AccountID: sql.NullInt64{Int64: account1.ID, Valid: true},
		Amount:    util.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	entry1, err := testQueries.GetEntry(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, 0)
	jsonData, err := json.MarshalIndent(entry1, "", "")
	require.NoError(t, err)
	fmt.Println(string(jsonData))
}

func TestDeleteEntri(t *testing.T) {
	entry, err := testQueries.DeleteEntry(context.Background(), 3)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("Tidak ada data dengan ID 3")
		return
	}
	require.NoError(t, err)
	require.NotEmpty(t, entry)


	_, err = testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	fmt.Println("Deleted entry:", entry.ID)
	
}
