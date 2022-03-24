package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teguhatma/simple-bank/util"
)

func createRandomAccount(t *testing.T) int64 {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	id, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	return id
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	id := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, id, account.ID)
}

func TestUpdateAccount(t *testing.T) {
	id := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      id,
		Balance: util.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, id, account.ID)
}

func TestDeleteAccount(t *testing.T) {
	id := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), id)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), id)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
