package domain

import (
	"context"

	"gopkg.in/guregu/null.v4/zero"
)

type Website struct {
	ID      int64     `json:"id"`
	URL     string    `json:"url"`
	AddedAt zero.Time `json:"added_at"`
}

type Status struct {
	ID        int64     `json:"id"`
	Up        zero.Bool `json:"up"`
	CheckedAt zero.Time `json:"checked_at"`
	WebsiteID int64     `json:"-"`
}

type WebsiteStatus struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	AddedAt   zero.Time `json:"added_at"`
	StatusID  int64     `json:"status_id"`
	Up        zero.Bool `json:"up"`
	CheckedAt zero.Time `json:"checked_at"`
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
