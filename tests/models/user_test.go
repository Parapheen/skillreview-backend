package modeltests

import (
	"log"
	"strconv"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/models"
	"gopkg.in/go-playground/assert.v1"
	_ "gorm.io/driver/postgres"
	"syreclabs.com/go/faker"
)

func TestSaveUser(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		Nickname:  faker.Internet().UserName(),
		Email:     faker.Internet().Email(),
		Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Avatar:    faker.Internet().Url(),
		Rank:      "Ancient 2",
		Plan:      models.Basic,
	}
	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error creating a user: %v\n", err)
		return
	}
	assert.Equal(t, newUser.UUID, savedUser.UUID)
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
	assert.Equal(t, newUser.Avatar, savedUser.Avatar)
	assert.Equal(t, newUser.Rank, savedUser.Rank)
}

func TestSaveUserDashedNickname(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		Nickname:  "TAGANROK-MOSCOW",
		Email:     faker.Internet().Email(),
		Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Avatar:    faker.Internet().Url(),
		Rank:      "Ancient 4",
		Plan:      models.Basic,
	}
	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error creating a user: %v\n", err)
		return
	}
	assert.Equal(t, newUser.UUID, savedUser.UUID)
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
	assert.Equal(t, newUser.Avatar, savedUser.Avatar)
	assert.Equal(t, newUser.Rank, savedUser.Rank)
}

func TestGetUserByUuid(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	_, user, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatalf("cannot seed user and review request: %v", err)
	}
	foundUser, err := userInstance.FindUserByUIID(server.DB, user.UUID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.UUID, user.UUID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
	assert.Equal(t, len(foundUser.ReviewRequests), 1)
}

func TestUpdateUser(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}

	userUpdate := models.User{
		Nickname: faker.Internet().UserName(),
		Email:    faker.Internet().Email(),
	}
	updatedUser, err := userUpdate.UpdateUser(server.DB, user.UUID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}
