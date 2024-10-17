package db

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func TestTransferTx(t *testing.T) {

// 	account1 := CreateRandomAccount(t)
// 	account2 := CreateRandomAccount(t)
// 	fmt.Println(">>before trade", account1.Balance, account2.Balance)

// 	n := 2
// 	amount := int64(10)

// 	errs := make(chan error)
// 	results := make(chan TransferTxResult)

// 	for i := 0; i < n; i++ {
// 		txname := fmt.Sprintf("tx %d", i+1)
// 		go func() {
// 			ctx := context.WithValue(context.Background(), txKey, txname)
// 			result, err := testStore.TransferTx(ctx, TransferTxParams{
// 				FromAccountID: account1.ID,
// 				ToAccountID:   account2.ID,
// 				Amount:        amount,
// 			})

// 			errs <- err
// 			results <- result
// 		}()
// 	}

// 	existed := make(map[int]bool)

// 	for i := 0; i < n; i++ {
// 		err := <-errs
// 		require.NoError(t, err)

// 		result := <-results
// 		require.NotEmpty(t, result)

// 		//检查这个转账记录
// 		Transfer := result.Transfer
// 		require.NotEmpty(t, Transfer)
// 		require.Equal(t, account1.ID, Transfer.FromAccountID)
// 		require.Equal(t, account2.ID, Transfer.ToAccountID)
// 		require.Equal(t, amount, Transfer.Amount)
// 		require.NotZero(t, Transfer.ID)
// 		require.NotZero(t, Transfer.CreatedAt)

// 		_, err = testStore.GetTransfer(context.Background(), Transfer.ID)
// 		require.NoError(t, err)

// 		//接下来是检查这个账单目录是否有存在
// 		formEntry := result.FromEntry
// 		require.NotEmpty(t, formEntry)
// 		require.Equal(t, account1.ID, formEntry.AccountID)
// 		require.Equal(t, -amount, formEntry.Amount)
// 		require.NotZero(t, formEntry.ID)
// 		require.NotZero(t, formEntry.CreatedAt)

// 		_, err = testStore.GetEntry(context.Background(), formEntry.ID)
// 		require.NoError(t, err)

// 		ToEntry := result.ToEntry
// 		require.NotEmpty(t, ToEntry)
// 		require.Equal(t, account2.ID, ToEntry.AccountID)
// 		require.Equal(t, amount, ToEntry.Amount)
// 		require.NotZero(t, ToEntry.ID)
// 		require.NotZero(t, ToEntry.CreatedAt)

// 		_, err = testStore.GetEntry(context.Background(), ToEntry.ID)
// 		require.NoError(t, err)

// 		formAccount := result.FromAccount
// 		require.NotEmpty(t, formAccount)
// 		require.Equal(t, account1.ID, formAccount.ID)

// 		toAccount := result.ToAccount
// 		require.NotEmpty(t, toAccount)
// 		require.Equal(t, account2.ID, toAccount.ID)

// 		// //最终还需要更新余额的操作

// 		fmt.Println(">>tx :", formAccount.Balance, toAccount.Balance)

// 		diff1 := account1.Balance - formAccount.Balance
// 		diff2 := toAccount.Balance - account2.Balance
// 		require.Equal(t, diff1, diff2)
// 		require.True(t, diff1 > 0)
// 		require.True(t, diff1%amount == 0)

// 		k := int(diff1 / amount)
// 		require.True(t, k >= 1 && k <= n)
// 		require.NotContains(t, existed, k)
// 		existed[k] = true
// 	}

// 	uupdateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	uupdateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
// 	require.NoError(t, err)

// 	fmt.Println(">>After trade", uupdateAccount1.Balance, uupdateAccount2.Balance)

// 	require.Equal(t, account1.Balance-int64(n)*amount, uupdateAccount1.Balance)
// 	require.Equal(t, account2.Balance+int64(n)*amount, uupdateAccount2.Balance)

// }
