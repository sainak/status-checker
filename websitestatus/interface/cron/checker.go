package cron

import (
	"context"
	"fmt"
	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/logger"
	_websiteStatusService "github.com/sainak/status-checker/websitestatus/service"
	"time"
)

const (
	NumWorkers    = 4
	SleepDuration = 30 * time.Second
)

func SpawnWorkers(ctx context.Context, repo domain.WebsiteStatusStorer, broker <-chan domain.Website) {
	logger.Info("spawning workers")
	checker := _websiteStatusService.HttpChecker{Repo: repo}

	for i := 0; i < NumWorkers; i++ {
		go func(id int) {
			logger.Info(fmt.Sprintf("worker %d started", id))
			for {
				select {
				case <-ctx.Done():
					return
				case website, more := <-broker:
					if !more {
						return
					}
					_websiteStatusService.CheckWebsiteStatus(ctx, &checker, website)
				}
			}
		}(i)
	}
}

func RunChecker(repo domain.WebsiteStatusStorer) {
	logger.Info("starting checker cron job")
	broker := make(chan domain.Website, NumWorkers)
	defer close(broker)
	ctx := context.Background()

	SpawnWorkers(ctx, repo, broker)

mainLoop:
	for {
		logger.Info("creating jobs")
		var next = ""
	queryLoop:
		for {
			websites, next, err := repo.QueryWebsitesWithStatus(ctx, next, 100, nil)
			if err != nil {

			}
			logger.Info("creating jobs for ", len(websites), " websites")
			for _, website := range websites {
				broker <- website
			}
			if next == "" {
				break
			}
			select {
			case <-ctx.Done():
				break queryLoop
			default:
			}
		}
		logger.Info("finished creating jobs, sleeping for 1 minute")
		select {
		case <-ctx.Done():
			break mainLoop
		default:
			time.Sleep(SleepDuration)
		}
	}
}
