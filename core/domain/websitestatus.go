package domain

import (
	"context"
	"time"
)

type Status struct {
	Up   bool
	Time time.Time
}

type Website struct {
	ID      int64     `json:"id"`
	URL     string    `json:"url"`
	AddedAt time.Time `json:"added_at"`
	Status  Status
}

type WebsiteStatus struct {
	ID        int64     `json:"id"`
	Up        bool      `json:"up"`
	Time      time.Time `json:"time"`
	WebsiteID int64     `json:"website_id"`
}

type WebsiteStatusStorer interface {
	FetchWebsites(ctx context.Context, cursor string, num int64, filters map[string]string) ([]Website, string, error)
	InsertWebsite(ctx context.Context, website *Website) error
	FetchWebsiteByID(ctx context.Context, id int64) (Website, error)
	UpdateIntoWebsite(ctx context.Context, id int64, website *Website) error
	DropWebsite(ctx context.Context, id int64) error
	InsertWebsiteStatus(ctx context.Context, status *WebsiteStatus) error
	FetchWebsiteStatuses(ctx context.Context, websiteID int64, cursor string, num int64) ([]WebsiteStatus, string, error)
}

type WebsiteStatusService interface {
	ListWebsites(ctx context.Context, cursor string, num int64, filters map[string]string) ([]Website, string, error)
	CreateWebsite(ctx context.Context, website *Website) error
	GetWebsiteByID(ctx context.Context, id int64) (Website, error)
	UpdateWebsite(ctx context.Context, id int64, website *Website) error
	DeleteWebsite(ctx context.Context, id int64) error
	CreateWebsiteStatus(ctx context.Context, status *WebsiteStatus) error
	ListWebsiteStatuses(ctx context.Context, websiteID int64, cursor string, num int64) ([]WebsiteStatus, string, error)
}
