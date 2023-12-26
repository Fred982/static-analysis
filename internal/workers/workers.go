package workers

import (
	"fmt"
	"log"
	"sa/pkg/database/models"
	"sa/pkg/git"
	"time"
)

func Worker(workerID int, repoCh <-chan *models.GithubRepository, processedCh chan<- int32) {
	for repo := range repoCh {
		fmt.Printf("Worker %d processing repo %d\n", workerID, repo.ID)
		repoURL := "https://github.com/" + repo.FullName
		localPath := "../repo"

		if err := git.CloneRepository(repoURL, localPath); err != nil {
			log.Println(repoURL + " " + err.Error())
		}

		time.Sleep(1 * time.Second)

		processedCh <- repo.ID
	}
}
