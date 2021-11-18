package modeltests

import (
	"log"
	"testing"

	"github.com/Parapheen/skillreview-backend/api/models"
	"gopkg.in/go-playground/assert.v1"
	_ "gorm.io/driver/postgres"
	"syreclabs.com/go/faker"
)

func TestSaveReview(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	rr, user, err := seedOneUserAndOneReviewRequest()
	if err != nil {
		log.Fatalf("Cannot seed user and review request: %v\n", err)
	}
	newReview := models.Review{
		Description:     faker.Lorem().Sentence(3),
		RateLaning:      faker.RandomInt(1, 4),
		RateTeamfights:  faker.RandomInt(1, 4),
		RateOverall:     faker.RandomInt(1, 4),
		State:           models.Submitted,
		AuthorID:        user.ID,
		ReviewRequestID: rr.ID,
	}
	savedRR, err := newReview.SaveReview(server.DB)
	if err != nil {
		t.Errorf("this is the error creating a review: %v\n", err)
		return
	}
	assert.Equal(t, newReview.UUID, savedRR.UUID)
	assert.Equal(t, newReview.ID, savedRR.ID)
	assert.Equal(t, newReview.Description, savedRR.Description)
	assert.Equal(t, newReview.RateLaning, savedRR.RateLaning)
	assert.Equal(t, newReview.RateTeamfights, savedRR.RateTeamfights)
	assert.Equal(t, newReview.RateOverall, savedRR.RateOverall)
	assert.Equal(t, newReview.State, savedRR.State)
	assert.Equal(t, newReview.Author.ID, savedRR.Author.ID)
	assert.Equal(t, newReview.ReviewRequestID, savedRR.ReviewRequestID)
}

func TestGetReviewByUuid(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	rr, user, review, err := seedOneUserAndOneReviewRequestAndOneReview()
	if err != nil {
		log.Fatalf("cannot seed user, review and review request: %v", err)
	}
	foundReview, err := reviewInstance.FindReviewByUUID(server.DB, review.UUID)
	if err != nil {
		t.Errorf("this is the error getting one review: %v\n", err)
		return
	}
	assert.Equal(t, foundReview.ID, review.ID)
	assert.Equal(t, foundReview.UUID, review.UUID)
	assert.Equal(t, foundReview.Description, review.Description)
	assert.Equal(t, foundReview.RateLaning, review.RateLaning)
	assert.Equal(t, foundReview.RateOverall, review.RateOverall)
	assert.Equal(t, foundReview.RateTeamfights, review.RateTeamfights)
	assert.Equal(t, foundReview.State, review.State)
	assert.Equal(t, foundReview.Author.ID, user.ID)
	assert.Equal(t, foundReview.Author.UUID, user.UUID)
	assert.Equal(t, foundReview.ReviewRequestID, rr.ID)
}

func TestUpdateReview(t *testing.T) {

	err := refreshDatabase()
	if err != nil {
		log.Fatal(err)
	}

	_, user, review, err := seedOneUserAndOneReviewRequestAndOneReview()
	if err != nil {
		log.Fatalf("cannot seed user, review and review request: %v", err)
	}

	reviewUpdate := models.Review{
		State:    models.Reviewed,
		AuthorID: user.ID,
	}
	updatedReview, err := reviewUpdate.UpdateReview(server.DB, review.UUID)
	if err != nil {
		t.Errorf("this is the error updating the revuew request: %v\n", err)
		return
	}
	assert.Equal(t, updatedReview.State, reviewUpdate.State)
}
