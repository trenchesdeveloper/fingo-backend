package db_test

import (
	"context"
	"github.com/trenchesdeveloper/fingo-backend/utils"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	db "github.com/trenchesdeveloper/fingo-backend/db/sqlc"
)

var wg sync.WaitGroup

func clean_up() {
	err := testQueries.DeleteAllUsers(context.Background())

	if err != nil {
		log.Fatalln(err)
	}
}

func createRandomUser(t *testing.T) db.User {
	hashedPassword, err := utils.GenerateHashedPassword(utils.RandomString(6))

	assert.NoError(t, err)

	arg := db.CreateUserParams{
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, arg.Email, user.Email)

	return user

}

func TestCreateUser(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	newPassword, err := utils.GenerateHashedPassword(utils.RandomString(6))

	assert.NoError(t, err)

	arg := db.UpdateUserPasswordParams{
		ID:             user.ID,
		HashedPassword: newPassword,
		UpdatedAt:      time.Now(),
	}

	user2, err := testQueries.UpdateUserPassword(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, newPassword, user2.HashedPassword)
	assert.WithinDuration(t, arg.UpdatedAt, user2.UpdatedAt, 2*time.Second)

}

func TestGetUserByEmail(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, user.Email, user2.Email)
	assert.Equal(t, user.HashedPassword, user2.HashedPassword)
}

func TestGetUserByID(t *testing.T) {
	defer clean_up()
	user := createRandomUser(t)

	user2, err := testQueries.GetUserByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, user2)
	assert.Equal(t, user.ID, user2.ID)
	assert.Equal(t, user.Email, user2.Email)
	assert.Equal(t, user.HashedPassword, user2.HashedPassword)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)

	assert.NoError(t, err)

	user2, err := testQueries.GetUserByID(context.Background(), user.ID)

	assert.Error(t, err)
	assert.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	defer clean_up()

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()
			createRandomUser(t)
		}()
	}

	wg.Wait()

	arg := db.ListUsersParams{
		Limit:  10,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	assert.NoError(t, err)
	assert.Len(t, users, 10)

	for _, user := range users {
		assert.NotEmpty(t, user)
	}
}
