package modeltests

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"gorm.io/driver/postgres"
	"syreclabs.com/go/faker"

	"github.com/Parapheen/skillreview-backend/api/controllers"
	"github.com/Parapheen/skillreview-backend/api/models"
	"gorm.io/gorm"
)

var server = controllers.Server{}
var userInstance = models.User{}
var reviewRequestInstance = models.ReviewRequest{}
var reviewInstance = models.Review{}

var reviews = []models.Review{
	{
		Description:    faker.Hacker().SaySomethingSmart(),
		RateLaning:     faker.RandomInt(1, 5),
		RateTeamfights: faker.RandomInt(1, 5),
		RateOverall:    faker.RandomInt(1, 5),
		State:          models.Submitted,
	},
	{
		Description:    faker.Hacker().SaySomethingSmart(),
		RateLaning:     faker.RandomInt(1, 5),
		RateTeamfights: faker.RandomInt(1, 5),
		RateOverall:    faker.RandomInt(1, 5),
		State:          models.Submitted,
	},
	{
		Description:    faker.Hacker().SaySomethingSmart(),
		RateLaning:     faker.RandomInt(1, 5),
		RateTeamfights: faker.RandomInt(1, 5),
		RateOverall:    faker.RandomInt(1, 5),
		State:          models.Accepted,
	},
}

var reviewRequests = []models.ReviewRequest{
	{
		MatchId:    strconv.Itoa(faker.RandomInt(123, 12345678)),
		AuthorRank: "Ancient 2",
	},
	{
		MatchId:    strconv.Itoa(faker.RandomInt(123, 12345678)),
		AuthorRank: "Divine 2",
	},
	{
		MatchId:    strconv.Itoa(faker.RandomInt(123, 12345678)),
		AuthorRank: "Ancient 4",
	},
}

func TestMain(m *testing.M) {
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}

	err = server.DB.AutoMigrate(
		&models.Base{},
		&models.User{},
		// &models.Review{},
		// &models.ReviewRequest{},
	)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
}

func seedOneUser() (models.User, error) {

	user := models.User{
		Nickname:  faker.Internet().UserName(),
		Email:     faker.Internet().Email(),
		Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Avatar:    faker.Internet().Url(),
		Rank:      "Ancient 2",
		Plan:      models.Basic,
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed admins table: %v", err)
	}
	return user, nil
}

func refreshDatabase() error {
	err := server.DB.Migrator().DropTable(
		&models.Base{},
		&models.User{},
		// &models.ReviewRequest{},
		// &models.Review{},
	)
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(
		&models.Base{},
		&models.User{},
		// &models.ReviewRequest{},
		// &models.Review{},
	)
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneReviewRequest() (models.ReviewRequest, models.User, error) {

	user := models.User{
		Nickname:  faker.Internet().UserName(),
		Email:     faker.Internet().Email(),
		Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Avatar:    faker.Internet().Url(),
		Rank:      "Ancient 2",
		Plan:      models.Basic,
	}
	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, err
	}
	rr := models.ReviewRequest{
		Description:        faker.Lorem().Sentence(3),
		SelfRateLaning:     faker.RandomInt(1, 4),
		SelfRateTeamfights: faker.RandomInt(1, 4),
		SelfRateOverall:    faker.RandomInt(1, 4),
		State:              models.Open,
		AuthorID:           user.ID,
		AuthorUUID:         user.UUID,
	}
	err = server.DB.Model(&models.ReviewRequest{}).Create(&rr).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, err
	}
	return rr, user, nil
}

func seedOneUserAndOneReviewRequestAndOneReview() (models.ReviewRequest, models.User, models.Review, error) {

	user := models.User{
		Nickname:  faker.Internet().UserName(),
		Email:     faker.Internet().Email(),
		Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Avatar:    faker.Internet().Url(),
		Rank:      "Ancient 2",
		Plan:      models.Basic,
	}
	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, models.Review{}, err
	}
	rr := models.ReviewRequest{
		Description:        faker.Lorem().Sentence(3),
		SelfRateLaning:     faker.RandomInt(1, 4),
		SelfRateTeamfights: faker.RandomInt(1, 4),
		SelfRateOverall:    faker.RandomInt(1, 4),
		State:              models.Open,
		AuthorID:           user.ID,
		AuthorUUID:         user.UUID,
	}
	err = server.DB.Model(&models.ReviewRequest{}).Create(&rr).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, models.Review{}, err
	}
	review := models.Review{
		Description:       faker.Lorem().Sentence(3),
		ReviewRequestUUID: rr.UUID,
		AuthorID:          user.ID,
		ReviewRequestID:   rr.ID,
		State:             models.ReviewState(models.Submitted),
		RateLaning:        faker.RandomInt(1, 4),
		RateTeamfights:    faker.RandomInt(1, 4),
		RateOverall:       faker.RandomInt(1, 4),
	}
	err = server.DB.Model(&models.Review{}).Create(&review).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, models.Review{}, err
	}
	return rr, user, review, nil
}

func seedOneUserAndOneReviewRequestAndReviews() (models.ReviewRequest, models.User, error) {

	user := models.User{
		Nickname:  faker.Internet().UserName(),
		Email:     faker.Internet().Email(),
		Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
		Avatar:    faker.Internet().Url(),
		Rank:      "Ancient 2",
		Plan:      models.Basic,
	}
	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, err
	}
	rr := models.ReviewRequest{
		Description:        faker.Lorem().Sentence(3),
		SelfRateLaning:     faker.RandomInt(1, 4),
		SelfRateTeamfights: faker.RandomInt(1, 4),
		SelfRateOverall:    faker.RandomInt(1, 4),
		State:              models.Open,
		AuthorID:           user.ID,
		AuthorUUID:         user.UUID,
	}
	err = server.DB.Model(&models.ReviewRequest{}).Create(&rr).Error
	if err != nil {
		return models.ReviewRequest{}, models.User{}, err
	}
	for i := range reviews {
		reviews[i].ReviewRequestID = rr.ID
		reviews[i].AuthorID = user.ID
		err = server.DB.Model(&models.Review{}).Create(&reviews[i]).Error
		if err != nil {
			log.Fatalf("cannot generate reviews: %v", err)
		}
	}
	return rr, user, nil
}

func seedUserAndRequestReviews() ([]models.User, []models.ReviewRequest, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.ReviewRequest{}, err
	}
	var users = []models.User{
		{
			Nickname:  faker.Internet().UserName(),
			Email:     faker.Internet().Email(),
			Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
			Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
			Avatar:    faker.Internet().Url(),
			Rank:      "Ancient 2",
			Plan:      models.Basic,
		},
		{
			Nickname:  faker.Internet().UserName(),
			Email:     faker.Internet().Email(),
			Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
			Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
			Avatar:    faker.Internet().Url(),
			Rank:      "Ancient 3",
			Plan:      models.Basic,
		},
		{
			Nickname:  faker.Internet().UserName(),
			Email:     faker.Internet().Email(),
			Steam64ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
			Steam32ID: strconv.Itoa(faker.RandomInt(123, 12345678)),
			Avatar:    faker.Internet().Url(),
			Rank:      "Ancient 4",
			Plan:      models.Pro,
		},
	}
	var reviewRequests = []models.ReviewRequest{
		{
			Description:        faker.Lorem().Sentence(3),
			SelfRateLaning:     faker.RandomInt(1, 4),
			SelfRateTeamfights: faker.RandomInt(1, 4),
			SelfRateOverall:    faker.RandomInt(1, 4),
			State:              models.Open,
		},
		{
			Description:        faker.Lorem().Sentence(3),
			SelfRateLaning:     faker.RandomInt(1, 4),
			SelfRateTeamfights: faker.RandomInt(1, 4),
			SelfRateOverall:    faker.RandomInt(1, 4),
			State:              models.Closed,
		},
		{
			Description:        faker.Lorem().Sentence(3),
			SelfRateLaning:     faker.RandomInt(1, 4),
			SelfRateTeamfights: faker.RandomInt(1, 4),
			SelfRateOverall:    faker.RandomInt(1, 4),
			State:              models.Open,
		},
		{
			Description:        faker.Lorem().Sentence(3),
			SelfRateLaning:     faker.RandomInt(1, 4),
			SelfRateTeamfights: faker.RandomInt(1, 4),
			SelfRateOverall:    faker.RandomInt(1, 4),
			State:              models.Closed,
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		reviewRequests[i].AuthorID = users[i].ID
		reviewRequests[i].AuthorUUID = users[i].UUID

		err = server.DB.Model(&models.ReviewRequest{}).Create(&reviewRequests[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, reviewRequests, nil
}
