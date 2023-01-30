package service

import (
	"context"
	"github.com/sainak/status-checker/core/domain"
	"time"
)

func NewWebsiteStatusService(repo domain.WebsiteStatusRepo, timeout time.Duration) domain.WebsiteStatusService {
	return &websiteStatusService{repo, timeout}
}

type websiteStatusService struct {
	repo           domain.WebsiteStatusRepo
	contextTimeout time.Duration
}

func (w websiteStatusService) ListWebsites(ctx context.Context, cursor string, num int64, filters map[string]string) ([]domain.Website, string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.FetchWebsites(ctx, cursor, num, filters)
}

func (w websiteStatusService) CreateWebsite(ctx context.Context, website *domain.Website) error {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.InsertWebsite(ctx, website)
}

func (w websiteStatusService) GetWebsiteByID(ctx context.Context, id int64) (res domain.Website, err error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.FetchWebsiteByID(ctx, id)
}

func (w websiteStatusService) UpdateWebsite(ctx context.Context, id int64, website *domain.Website) error {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.UpdateIntoWebsite(ctx, id, website)
}

func (w websiteStatusService) DeleteWebsite(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.DropWebsite(ctx, id)
}

func (w websiteStatusService) CreateWebsiteStatus(ctx context.Context, status *domain.WebsiteStatus) error {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.InsertWebsiteStatus(ctx, status)
}

func (w websiteStatusService) ListWebsiteStatuses(ctx context.Context, websiteID int64, cursor string, num int64) ([]domain.WebsiteStatus, string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.FetchWebsiteStatuses(ctx, websiteID, cursor, num)
}
