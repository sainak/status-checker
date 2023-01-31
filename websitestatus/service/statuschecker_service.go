package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/logger"
	"gopkg.in/guregu/null.v4/zero"
)

type httpChecker struct {
	http.Client
	repo domain.WebsiteStatusStorer
}

func NewHttpChecker(r domain.WebsiteStatusStorer) domain.StatusChecker {
	return &httpChecker{
		repo: r,
	}

}

func (h *httpChecker) Check(ctx context.Context, name string) (status bool, err error) {
	logger.Info("checking website: ", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+name, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	resp, err := h.Do(req)
	if err != nil {
		// handle no such host error
		//logger.Error(err)
		//if err, ok := err.(*url.Error); ok {
		//	if err, ok := err.Err.(*net.OpError); ok {
		//		if _, ok := err.Err.(*net.DNSError); !ok {
		//			return false, err
		//		}
		//	}
		//}
		status, err = false, nil
	} else {
		status = resp.StatusCode == http.StatusOK
		_ = resp.Body.Close() // nolint:errcheck // ignore error on close because we only care about the status code
	}
	logger.Info(fmt.Sprint(name, " is up: ", status))
	return
}

func (h *httpChecker) CreateStatus(ctx context.Context, status *domain.Status) error {
	return h.repo.InsertStatus(ctx, status)
}

func CheckWebsiteStatus(ctx context.Context, checker domain.StatusChecker, website domain.Website) {
	status, err := checker.Check(ctx, website.URL)
	if err != nil {
		logger.Error(err)
		return
	}
	websiteStatus := domain.Status{
		WebsiteID: website.ID,
		Up:        zero.NewBool(status, true),
		CheckedAt: zero.NewTime(time.Now(), true),
	}
	err = checker.CreateStatus(ctx, &websiteStatus)
	if err != nil {
		logger.Error(err)
		return
	}
}
