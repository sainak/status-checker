package domain

import (
	"context"
	"time"

	"gopkg.in/guregu/null.v4"
)

type Website struct {
	ID      int64     `json:"id"`
	URL     string    `json:"url"`
	AddedAt time.Time `json:"added_at"`
}

type Status struct {
	ID        int64     `json:"id"`
	Up        bool      `json:"up"`
	CheckedAt time.Time `json:"checked_at"`
	WebsiteID int64     `json:"-"`
}

type WebsiteStatus struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	AddedAt   time.Time `json:"added_at"`
	StatusID  null.Int  `json:"status_id"`
	Up        null.Bool `json:"up"`
	CheckedAt null.Time `json:"checked_at"`
}

type WebsiteStatusStorer interface {
	QueryWebsites(ctx context.Context, cursor string, num int64, filters map[string]string) ([]Website, string, error)
	QueryWebsitesStatus(ctx context.Context, cursor string, num int64, filters map[string]string) ([]WebsiteStatus, string, error)
	QueryWebsiteStatusByID(ctx context.Context, id int64) (WebsiteStatus, error)
	InsertWebsite(ctx context.Context, website *Website) error
	DropWebsite(ctx context.Context, id int64) error
	InsertStatus(ctx context.Context, status *Status) error
	QueryStatusesByWebsiteID(ctx context.Context, websiteID int64, cursor string, num int64) ([]Status, string, error)
}

type WebsiteStatusService interface {
	ListWebsitesStatus(ctx context.Context, cursor string, num int64, filters map[string]string) ([]WebsiteStatus, string, error)
	GetWebsiteStatusByID(ctx context.Context, id int64) (WebsiteStatus, error)
	CreateWebsite(ctx context.Context, website *Website) error
	DeleteWebsite(ctx context.Context, id int64) error
	ListWebsiteStatuses(ctx context.Context, websiteID int64, cursor string, num int64) ([]Status, string, error)
}

type WebsiteStatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
	CreateStatus(ctx context.Context, status *Status) error
}
