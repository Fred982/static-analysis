package main

import (
	"log"

	pg "sa/src/database"
	"sa/src/models"
	"sa/src/workers"

	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := pg.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pg.CloseDatabase(db); err != nil {
			log.Fatal(err)
		}
	}()

	var repositories []models.GithubRepository
	if result := db.Limit(10).Find(&repositories); result.Error != nil {
		log.Fatal(result.Error)
	}

	repoCh := make(chan *models.GithubRepository, len(repositories))
	processedCh := make(chan int32)

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go workers.Worker(i, repoCh, processedCh)
	}

	for _, repo := range repositories {
		clonedRepo := repo
		repoCh <- &clonedRepo
	}
	close(repoCh)

	processed := make(map[int32]bool)
	for range repositories {
		repoID := <-processedCh
		if processed[repoID] {
			log.Fatalf("Duplicate processing detected for repo %d\n", repoID)
		}
		processed[repoID] = true
	}
}
