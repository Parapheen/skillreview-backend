package modeltests

import (
	"log"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/models"
	"gopkg.in/go-playground/assert.v1"
	_ "gorm.io/driver/postgres"
	"syreclabs.com/go/faker"
)

func TestSaveReviewRequest(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}
	newReviewRequest := models.ReviewRequest{
		Description:        faker.Lorem().Sentence(3),
		SelfRateLaning:     faker.RandomInt(1, 4),
		SelfRateTeamfights: faker.RandomInt(1, 4),
		SelfRateOverall:    faker.RandomInt(1, 4),
		State:              models.Open,
		AuthorID:           user.ID,
	}
	savedRR, err := newReviewRequest.SaveReviewRequest(server.DB)
	if err != nil {
		t.Errorf("this is the error creating a review request: %v\n", err)
		return
	}
	assert.Equal(t, newReviewRequest.UUID, savedRR.UUID)
	assert.Equal(t, newReviewRequest.ID, savedRR.ID)
	assert.Equal(t, newReviewRequest.Description, savedRR.Description)
	assert.Equal(t, newReviewRequest.SelfRateLaning, savedRR.SelfRateLaning)
	assert.Equal(t, newReviewRequest.SelfRateTeamfights, savedRR.SelfRateTeamfights)
	assert.Equal(t, newReviewRequest.SelfRateOverall, savedRR.SelfRateOverall)
	assert.Equal(t, newReviewRequest.State, savedRR.State)
}

func TestGetReviewRequestByUuid(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	rr, user, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatalf("cannot seed user and review request: %v", err)
	}
	foundRR, err := reviewRequestInstance.FindReviewRequestByUUID(server.DB, rr.UUID)
	if err != nil {
		t.Errorf("this is the error getting one review request: %v\n", err)
		return
	}
	assert.Equal(t, foundRR.ID, rr.ID)
	assert.Equal(t, foundRR.UUID, rr.UUID)
	assert.Equal(t, foundRR.Description, rr.Description)
	assert.Equal(t, foundRR.SelfRateLaning, rr.SelfRateLaning)
	assert.Equal(t, foundRR.SelfRateOverall, rr.SelfRateOverall)
	assert.Equal(t, foundRR.SelfRateTeamfights, rr.SelfRateTeamfights)
	assert.Equal(t, foundRR.State, rr.State)
	assert.Equal(t, foundRR.Author.ID, user.ID)
	assert.Equal(t, foundRR.Author.UUID, user.UUID)
}

func TestUpdateReviewRequest(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	rr, user, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatalf("cannot seed user and review request: %v", err)
	}

	rrUpdate := models.ReviewRequest{
		State:       models.Closed,
		Description: faker.Lorem().Sentence(3),
		AuthorID:    user.ID,
	}
	updatedRR, err := rrUpdate.UpdateReviewRequest(server.DB, rr.UUID)
	if err != nil {
		t.Errorf("this is the error updating the revuew request: %v\n", err)
		return
	}
	assert.Equal(t, updatedRR.State, rrUpdate.State)
	assert.Equal(t, updatedRR.Description, rrUpdate.Description)
}
