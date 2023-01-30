package domain

import "context"

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
	CreateStatus(ctx context.Context, status *WebsiteStatus) error
}
