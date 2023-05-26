package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/I1Asyl/ginBerliner/models"
	"github.com/I1Asyl/ginBerliner/pkg/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var db *sql.DB
var services *Services
var repo *repository.Repository

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
			fmt.Println(err)
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
	services.AddUser(models.User{
		Username:  "asyl",
		Password:  "Qqwerty1!.",
		Email:     "email@som.com",
		LastName:  "Yerassyl",
		FirstName: "Altay",
	})
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
