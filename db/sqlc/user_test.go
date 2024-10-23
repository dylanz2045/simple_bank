package db

import (
	"Project/utils"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {

	hashedPassword, err := utils.HashedPassword("secret")
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpodateUserOnlyFullName(t *testing.T) {
	olduser := CreateRandomUser(t)

	newFullName := utils.RandomOwner()

	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: olduser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, olduser.FullName, newUser.FullName)
	require.Equal(t, olduser.Username, newUser.Username)
	require.Equal(t, newUser.FullName, newFullName)
	require.Equal(t, olduser.Email, newUser.Email)
	require.Equal(t, olduser.HashedPassword, newUser.HashedPassword)
}

func TestUpodateUserOnlyEmail(t *testing.T) {
	olduser := CreateRandomUser(t)

	newEmail := utils.RandomEmail()
	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: olduser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, olduser.Email, newUser.Email)
	require.Equal(t, olduser.Username, newUser.Username)
	require.Equal(t, newUser.Email, newEmail)
	require.Equal(t, olduser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, olduser.FullName, newUser.FullName)
}

func TestUpodateUserOnlyPassword(t *testing.T) {
	olduser := CreateRandomUser(t)

	newPassword := utils.RandomString(6)
	newHashedPassword, err := utils.HashedPassword(newPassword)
	require.NoError(t, err)
	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: olduser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, olduser.PasswordChangedAt, newUser.PasswordChangedAt)
	require.NotEqual(t, olduser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, olduser.Username, newUser.Username)
	require.Equal(t, newUser.Email, olduser.Email)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.Equal(t, olduser.FullName, newUser.FullName)
}

func TestUpdateAllFields(t *testing.T) {
	olduser := CreateRandomUser(t)

	newEmail := utils.RandomEmail()

	newPassword := utils.RandomString(6)
	newHashedPassword, err := utils.HashedPassword(newPassword)
	require.NoError(t, err)

	newFullName := utils.RandomOwner()

	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: olduser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, olduser.PasswordChangedAt, newUser.PasswordChangedAt)
	require.Equal(t, newEmail, newUser.Email)
	require.NotEqual(t, olduser.FullName, newUser.FullName)
	require.Equal(t, newFullName, newUser.FullName)
	require.NotEqual(t, olduser.Email, newUser.Email)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.NotEqual(t, olduser.HashedPassword, newUser.HashedPassword)

}
