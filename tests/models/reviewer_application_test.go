package modeltests

import (
	"log"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/models"
	"gopkg.in/go-playground/assert.v1"
	_ "gorm.io/driver/postgres"
	"syreclabs.com/go/faker"
)

func TestSaveApplication(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}
	newApplication := models.ReviewerApplication{
		Description: faker.Lorem().Sentence(3),
		Rating:      faker.RandomInt(1, 15000),
		State:       models.Recieved,
		AuthorID:    user.ID,
	}
	savedApplication, err := newApplication.SaveApplication(server.DB)
	if err != nil {
		t.Errorf("this is the error creating an application: %v\n", err)
		return
	}
	assert.Equal(t, newApplication.UUID, savedApplication.UUID)
	assert.Equal(t, newApplication.ID, savedApplication.ID)
	assert.Equal(t, newApplication.Description, savedApplication.Description)
	assert.Equal(t, newApplication.Rating, savedApplication.Rating)
	assert.Equal(t, newApplication.State, savedApplication.State)
	assert.Equal(t, newApplication.Author.ID, savedApplication.Author.ID)
}

func TestGetApplicationByUuid(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	application, user, err := seedOneUserAndOneApplication()
	if err != nil {
		log.Fatalf("cannot seed user and application: %v", err)
	}
	foundApplication, err := applicationInstance.FindApplicationByUUID(server.DB, application.UUID)
	if err != nil {
		t.Errorf("this is the error getting one application: %v\n", err)
		return
	}
	assert.Equal(t, foundApplication.ID, application.ID)
	assert.Equal(t, foundApplication.UUID, application.UUID)
	assert.Equal(t, foundApplication.Description, application.Description)
	assert.Equal(t, foundApplication.Rating, application.Rating)
	assert.Equal(t, foundApplication.State, application.State)
	assert.Equal(t, foundApplication.Author.ID, user.ID)
	assert.Equal(t, foundApplication.Author.UUID, user.UUID)
}

func TestUpdateApplication(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	application, user, err := seedOneUserAndOneApplication()
	if err != nil {
		log.Fatalf("cannot seed user, application: %v", err)
	}

	applicationUpdate := models.ReviewerApplication{
		State:    models.Recieved,
		AuthorID: user.ID,
	}
	updatedApplication, err := applicationUpdate.UpdateApplication(server.DB, application.UUID)
	if err != nil {
		t.Errorf("this is the error updating the application: %v\n", err)
		return
	}
	assert.Equal(t, updatedApplication.State, applicationUpdate.State)
}
