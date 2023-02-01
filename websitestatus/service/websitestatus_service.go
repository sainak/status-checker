package service

import (
	"context"
	"time"

	"github.com/sainak/status-checker/core/domain"
)

func NewWebsiteStatusService(repo domain.WebsiteStatusStorer, timeout time.Duration) domain.WebsiteStatusService {
	return &websiteStatusService{repo, timeout}
}

type websiteStatusService struct {
	repo           domain.WebsiteStatusStorer
	contextTimeout time.Duration
}

func (w websiteStatusService) ListWebsitesStatus(ctx context.Context, cursor string, num int64, filters map[string]string) ([]domain.WebsiteStatus, string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.QueryWebsitesStatus(ctx, cursor, num, filters)
}

func (w websiteStatusService) CreateWebsite(ctx context.Context, website *domain.Website) error {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.InsertWebsite(ctx, website)
}

func (w websiteStatusService) GetWebsiteStatusByID(ctx context.Context, id int64) (res domain.WebsiteStatus, err error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.QueryWebsiteStatusByID(ctx, id)
}

func (w websiteStatusService) DeleteWebsite(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.DropWebsite(ctx, id)
}

func (w websiteStatusService) ListWebsiteStatuses(ctx context.Context, websiteID int64, cursor string, num int64) ([]domain.Status, string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.contextTimeout)
	defer cancel()

	return w.repo.QueryStatusesByWebsiteID(ctx, websiteID, cursor, num)
}
