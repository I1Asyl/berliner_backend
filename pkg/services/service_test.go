package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/I1Asyl/berliner_backend/models"
	"github.com/I1Asyl/berliner_backend/pkg/repository"
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

func TestCreateChannel(t *testing.T) {
	testTable := []struct {
		name          string
		channel       models.Channel
		channelLeader models.User
		expected      map[string]string
	}{
		{
			name: "success",
			channel: models.Channel{
				Name:        "Channel",
				LeaderId:    1,
				Description: "hoho",
			},
			channelLeader: testUser,
			expected:      map[string]string{},
		},
		{
			name: "error",
			channel: models.Channel{
				Name:        "Channel",
				LeaderId:    1,
				Description: "",
			},
			channelLeader: testUser,
			expected: map[string]string{
				"description": "Channel description can not be empty",
			},
		},
	}

	for _, testCase := range testTable {

		t.Run(testCase.name, func(t *testing.T) {
			if err := repo.Migration.Up(); err != nil {
				t.Errorf("Migration problems %s ", err)
			}
			services.AddUser(testUser)
			err := services.CreateChannel(testCase.channel, testUser)
			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %v, got %v", testCase.expected, err)
			}
			if err := repo.Migration.Down(); err != nil {
				t.Errorf("Migration problems %s ", err)
			}
		})

	}

}

// gets User model by username in the transaction
func TestGetUserByUsername(t *testing.T) {
	testTable := []struct {
		name     string
		username string
		expected models.User
	}{
		{
			name:     "success",
			username: "asyl",
			expected: testUser,
		},
		{
			name:     "error",
			username: "x",
			expected: models.User{},
		},
	}
	if err := repo.Migration.Up(); err != nil {
		t.Errorf("Migration problems %s ", err)
	}
	services.AddUser(testUser)
	for _, testCase := range testTable {
		user, _ := services.GetUserByUsername(testCase.username)

		// user's password is hashed, so we don't need to compare it
		user.Password = testCase.expected.Password

		if !reflect.DeepEqual(user, testCase.expected) {
			t.Errorf("Expected %v, got %v", testCase.expected, user)
		}
	}
	if err := repo.Migration.Down(); err != nil {
		t.Errorf("Migration problems %s ", err)
	}
}

// // create a new channel in the database for the given user
// func (a ApiService) CreateChannel(channel models.Channel, user models.User) map[string]string {

// 	invalid := channel.IsValid()
// 	channel.LeaderId = user.Id
// 	tx := a.repo.SqlQueries.StartTransaction()

// 	if len(invalid) == 0 {
// 		if err := tx.AddChannel(channel); err != nil {
// 			invalid["common"] = err.Error()
// 		} else {
// 			channel, _ = tx.GetChannelByName(channel.Name)
// 			membership := models.Membership{UserId: channel.LeaderId, ChannelId: channel.Id, IsEditor: true}
// 			tx.AddMembership(membership)

// 		}
// 	}
// 	err := tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 	}

// 	return invalid
// }

// // create a new post in the database for the given user or channel
// func (a ApiService) CreatePost(post models.Post) map[string]string {

// 	invalid := post.IsValid()

// 	if len(invalid) == 0 {
// 		err := a.repo.SqlQueries.AddPost(post)
// 		if err != nil {
// 			invalid["common"] = err.Error()
// 		}
// 	}

// 	return invalid
// }

// // get all user's channel posts from the database
// func (a ApiService) GetPostsFromChannels(user models.User) ([]models.Post, error) {
// 	posts, err := a.repo.SqlQueries.GetChannelPosts(user)
// 	return posts, err
// }

// // get all user's following's posts from the database
// func (a ApiService) GetPostsFromUsers(user models.User) ([]models.Post, error) {
// 	posts, err := a.repo.SqlQueries.GetUserPosts(user)
// 	return posts, err
// }

// // get all posts available for the given user from the database
// func (a ApiService) GetAllPosts(user models.User) ([]models.Post, error) {
// 	channelPosts, _ := a.GetPostsFromChannels(user)
// 	userPosts, _ := a.GetPostsFromUsers(user)
// 	posts := channelPosts
// 	posts = append(posts, userPosts...)
// 	return posts, nil
// }

// func (a ApiService) GetFollowing(user models.User) ([]models.User, error) {
// 	users, err := a.repo.SqlQueries.GetFollowing(user)
// 	return users, err
// }

// func (a ApiService) DeleteChannel(channel models.Channel) error {
// 	err := a.repo.SqlQueries.DeleteChannel(channel)
// 	return err
// }

// func (a ApiService) UpdateChannel(channel models.Channel) error {
// 	err := a.repo.SqlQueries.UpdateChannel(channel)
// 	return err
// }
