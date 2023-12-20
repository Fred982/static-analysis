package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/lpernett/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}
type GithubRepository struct {
	ID          int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FullName    string
	Description string
}

func (GithubRepository) TableName() string {
	return "github_repository"
}

func connectToDatabase() (*gorm.DB, error) {
	dbName := os.Getenv("DATABASE_NAME")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func closeDatabase(db *gorm.DB) error {
	postgres, err := db.DB()
	if err != nil {
		return err
	}
	return postgres.Close()
}

func cloneRepository(repoURL, localPath string) error {
	parts := strings.Split(repoURL, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	repoPath := filepath.Join(localPath, repoName)
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error cloning repository: %v", err)
	}
	return nil
}

func cloneWorker(repoCh <-chan GithubRepository, done chan<- bool) {
	for repo := range repoCh {
		repoURL := "https://github.com/" + repo.FullName
		localPath := "../repo"

		if err := cloneRepository(repoURL, localPath); err != nil {
			fmt.Println(err)
		}
	}
	done <- true
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := connectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := closeDatabase(db); err != nil {
			log.Fatal(err)
		}
	}()
	var repositories []GithubRepository
	if result := db.Limit(100).Find(&repositories); result.Error != nil {
		log.Fatal(result.Error)
	}

	repoCh := make(chan GithubRepository, len(repositories))
	done := make(chan bool)

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go cloneWorker(repoCh, done)
	}

	for _, repo := range repositories {
		repoCh <- repo
	}
	close(repoCh)

	for i := 0; i < numWorkers; i++ {
		<-done
	}

}
