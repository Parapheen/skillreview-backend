package controllertests

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"syreclabs.com/go/faker"

	"github.com/Parapheen/skillreview-backend/api/controllers"
	"github.com/Parapheen/skillreview-backend/api/models"
	"gorm.io/gorm"
)

var server = controllers.Server{}

func TestMain(m *testing.M) {
	Database()

	os.Exit(m.Run())
}

func CreateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
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
		&models.ReviewRequest{},
		&models.Review{},
	)

	if err != nil {
		log.Fatalf("Error automigrating %v", err)
	}
}

func refreshUserTable() error {
	err := server.DB.Migrator().DropTable(&models.User{})
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table users")
	return nil
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

func refreshUserAndReviewRequestTables() error {

	err := server.DB.Migrator().DropTable(&models.User{}, &models.ReviewRequest{})
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.ReviewRequest{})
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables users and review_requests")
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
