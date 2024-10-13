package db

import (
	"Project/utils"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := TestQueries.CreateAccount(context.Background(), arg)

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
	newaccount := CreateRandomAccount(t)
	account1, err := TestQueries.GetAccount(context.Background(), newaccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, newaccount.ID, account1.ID)
	require.Equal(t, newaccount.Owner, account1.Owner)
	require.Equal(t, newaccount.Balance, account1.Balance)
	require.Equal(t, newaccount.Currency, account1.Currency)
	require.WithinDuration(t, newaccount.CreatedAt, account1.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	newAccount := CreateRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: utils.RandomMoney(),
	}
	account1, err := TestQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, newAccount.ID, account1.ID)
	require.Equal(t, newAccount.Owner, account1.Owner)
	require.Equal(t, arg.Balance, account1.Balance)
	require.Equal(t, newAccount.Currency, account1.Currency)
	require.WithinDuration(t, newAccount.CreatedAt, account1.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account

	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := TestQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
