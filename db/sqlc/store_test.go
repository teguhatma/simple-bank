package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	id := createRandomAccount(t)
	id2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: id,
				ToAccountID:   id2,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	// existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, id, transfer.FromAccountID)
		require.Equal(t, id2, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, id, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Check entries
		toEntry := result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, id2, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, id, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, id2, toAccount.ID)

		// // Check account balance
		// fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		// diff1 := randAccount.Balance - fromAccount.Balance
		// diff2 := toAccount.Balance - randAccount1.Balance
		// require.Equal(t, diff1, diff2)
		// require.True(t, diff1 > 0)
		// require.True(t, diff1%amount == 0)

		// k := int(diff1 / amount)
		// require.True(t, k >= 1 && k <= n)
		// require.NotContains(t, existed, k)
		// existed[k] = true
	}

	// Check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), id2)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)
	// require.Equal(t, randAccount.Balance-int64(n)*amount, updateAccount1.Balance)
	// require.Equal(t, randAccount1.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	id := createRandomAccount(t)
	id2 := createRandomAccount(t)
	// fmt.Println(">> before:", randAccount.Balance, randAccount1.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := id
		toAccountID := id2

		if i%2 == 1 {
			fromAccountID = id2
			toAccountID = id
		}

		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final updated balances
	updateAccount1, err := testQueries.GetAccount(context.Background(), id)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), id2)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)
	// require.Equal(t, randAccount.Balance, updateAccount1.Balance)
	// require.Equal(t, randAccount1.Balance, updateAccount2.Balance)
}
