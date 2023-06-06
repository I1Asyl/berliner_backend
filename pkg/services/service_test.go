package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var services *Services
var repo *repository.Repository
var testUser models.User

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "", []string{"MYSQL_ROOT_PASSWORD=secret", "MYSQL_DATABASE=berliner"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	testUser = models.User{
		Id:        1,
		Username:  "asyl",
		FirstName: "Yerassyl",
		LastName:  "Altay",
		Email:     "altayerasyl@gmail.com",
		Password:  "Qqwerty1!.",
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dsn := fmt.Sprintf("root:secret@tcp(localhost:%s)/berliner", resource.GetPort("3306/tcp"))
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	repo = repository.NewRepository(dsn, "file://../../migrations")
	services = NewService(repo)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	testTable := []struct {
		name      string
		inputUser models.User
		expected  map[string]string
	}{
		{
			name: "success",
			inputUser: models.User{
				Username:  "test",
				Email:     "email@som.com",
				Password:  "Qqwerty1!.",
				LastName:  "Yerassyl",
				FirstName: "Altay",
			},
			expected: map[string]string{},
		},
		{
			name: "error username",
			inputUser: models.User{
				Username:  "t",
				Email:     "email@som.com",
				Password:  "Qqwerty1!.",
				LastName:  "Yerassyl",
				FirstName: "Altay",
			},
			expected: map[string]string{
				"username": "Invalid username",
			},
		},
		{
			name: "error email",
			inputUser: models.User{
				Username:  "test",
				Email:     "emailsom.com",
				Password:  "Qqwerty1!.",
				LastName:  "Yerassyl",
				FirstName: "Altay",
			},
			expected: map[string]string{
				"email": "Invalid email",
			},
		},
	}
	err := repo.Migration.Up()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := services.AddUser(testCase.inputUser)
			ans := reflect.DeepEqual(err, testCase.expected)
			if !ans {
				t.Errorf("Expected %v, got %v", testCase.expected, err)
			}
		})
	}
	err = repo.Migration.Down()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
}

func TestCheckUserAndPassword(t *testing.T) {
	testTable := []struct {
		name      string
		inputUser models.AuthorizationForm
		expected  bool
	}{
		{
			name: "success",
			inputUser: models.AuthorizationForm{
				Username: "asyl",
				Password: "Qqwerty1!.",
			},
			expected: true,
		},
		{
			name: "error",
			inputUser: models.AuthorizationForm{
				Username: "asyl",
				Password: "Qqwerty1!",
			},
			expected: false,
		},
	}

	err := repo.Migration.Up()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
	services.AddUser(testUser)
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ans, _ := services.CheckUserAndPassword(testCase.inputUser)
			if ans != testCase.expected {
				t.Errorf("Expected %v, got %v", testCase.expected, ans)
			}
		})
	}
	err = repo.Migration.Down()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
}

func TestGenarateToken(t *testing.T) {
	testTable := []struct {
		name       string
		issueTime  time.Time
		expireTime time.Time
		inputUser  models.AuthorizationForm
		expected   string
		jwt_secret string
	}{
		{
			name:       "success",
			issueTime:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			expireTime: time.Date(2009, time.November, 10, 23, 30, 0, 0, time.UTC),
			inputUser: models.AuthorizationForm{
				Username: "asyl",
				Password: "Qqwerty1!.",
			},
			expected:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImFzeWwiLCJpc3MiOiJ0ZXN0Iiwic3ViIjoic29tZWJvZHkiLCJleHAiOjEyNTc4OTU4MDAsImlhdCI6MTI1Nzg5NDAwMH0.cWHFSBmmpznRvLw56mokDKpa1Olv4Wy7Pf5YGp3gKFw",
			jwt_secret: "randomJWTSecret",
		},
		{
			name:       "success2",
			issueTime:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			expireTime: time.Date(3000, time.November, 10, 23, 30, 0, 0, time.UTC),
			inputUser: models.AuthorizationForm{
				Username: "asyl",
				Password: "Qqwerty1!.",
			},
			expected:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImFzeWwiLCJpc3MiOiJ0ZXN0Iiwic3ViIjoic29tZWJvZHkiLCJleHAiOjMyNTMwODA3ODAwLCJpYXQiOjEyNTc4OTQwMDB9.yGe-6MApCd8jvvsuwZH4O9tc3AB-ISBDMYx3xSP_Ork",
			jwt_secret: "randomJWTSecret",
		},
	}
	err := repo.Migration.Up()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			os.Setenv("JWT_SECRET", testCase.jwt_secret)
			ans, err := services.GenerateToken(testCase.inputUser, testCase.issueTime, testCase.expireTime)
			if ans != testCase.expected {
				t.Errorf("Expected %v, got %v, error: %s", testCase.expected, ans, err)
			}
		})
	}
	err = repo.Migration.Down()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
}

func TestParseToken(t *testing.T) {
	testTable := []struct {
		name             string
		token            string
		jwt_secret       string
		expectedUsername string
	}{
		{
			name:             "success",
			token:            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImFzeWwiLCJpc3MiOiJ0ZXN0Iiwic3ViIjoic29tZWJvZHkiLCJleHAiOjMyNTMwODA3ODAwLCJpYXQiOjEyNTc4OTQwMDB9.yGe-6MApCd8jvvsuwZH4O9tc3AB-ISBDMYx3xSP_Ork",
			jwt_secret:       "randomJWTSecret",
			expectedUsername: "asyl",
		},
	}
	err := repo.Migration.Up()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			os.Setenv("JWT_SECRET", testCase.jwt_secret)
			ans, err := services.ParseToken(testCase.token)
			if ans != testCase.expectedUsername {
				t.Errorf("Expected %v, got %v, error: %s", testCase.expectedUsername, ans, err)
			}
		})
	}
	err = repo.Migration.Down()
	if err != nil {
		t.Errorf("Migration problems %s ", err)
	}
}

func TestHashPassword(t *testing.T) {
	testTable := []struct {
		name     string
		password string
		hash     string
	}{
		{
			name:     "success",
			password: "Qqwerty1!.",
			hash:     "$2a$11$IZph6sLg28fsOA2qD6xhsO2pWvnL9ihKkalZgpAeG.Nl6I8QN.Y4m",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := bcrypt.CompareHashAndPassword([]byte(testCase.hash), []byte(testCase.password))
			if err != nil {
				t.Errorf("Expected %v, is not a hash of %v", testCase.hash, testCase.password)
			}
		})
	}
}

func TestCreateTeam(t *testing.T) {
	testTable := []struct {
		name       string
		team       models.Team
		teamLeader models.User
		expected   map[string]string
	}{
		{
			name: "success",
			team: models.Team{
				TeamName:        "Team",
				TeamLeaderId:    1,
				TeamDescription: "hoho",
			},
			teamLeader: testUser,
			expected:   map[string]string{},
		},
		{
			name: "error",
			team: models.Team{
				TeamName:        "Team",
				TeamLeaderId:    1,
				TeamDescription: "",
			},
			teamLeader: testUser,
			expected: map[string]string{
				"teamDescription": "Team description can not be empty",
			},
		},
	}

	for _, testCase := range testTable {

		t.Run(testCase.name, func(t *testing.T) {
			if err := repo.Migration.Up(); err != nil {
				t.Errorf("Migration problems %s ", err)
			}
			services.AddUser(testUser)
			err := services.CreateTeam(testCase.team, testUser)
			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %v, got %v", testCase.expected, err)
			}
			if err := repo.Migration.Down(); err != nil {
				t.Errorf("Migration problems %s ", err)
			}
		})

	}

}

// func TestGetFollowing(t *testing.T) {
// 	testTable := []struct {
// 		name     string
// 		user     models.User
// 		expected []string
// 	}{
// 		{
// 			name: "success",
// 			user: models.User{
// 				Username:  "asyl",
// 				FirstName: "XXXXX",
// 				LastName:  "XXXXX",
// 				Email:     "XXXXXXXXXXXXX",
// 				Password:  "Qqwerty1!.",
// 			},
// 		},
// 	}
// }
