package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/teguhatma/simple-bank/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	randAccount := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), randAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randAccount.ID, account.ID)
	require.Equal(t, randAccount.Owner, account.Owner)
	require.Equal(t, randAccount.Balance, account.Balance)
	require.Equal(t, randAccount.Currency, account.Currency)
	require.WithinDuration(t, randAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      randAccount.ID,
		Balance: util.RandomMoney(),
	}

	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, randAccount.ID, account.ID)
	require.Equal(t, randAccount.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, randAccount.Currency, account.Currency)
	require.WithinDuration(t, randAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), randAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), randAccount.ID)
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
