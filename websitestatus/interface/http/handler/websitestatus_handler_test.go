package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/domain/mocks"
	"github.com/sainak/status-checker/core/myerrors"
)

type WebsiteStatusHandlerTestSuite struct {
	suite.Suite
	handler *WebsiteStatusHandler
	service *mocks.WebsiteStatusService
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(WebsiteStatusHandlerTestSuite))
}

func (suite *WebsiteStatusHandlerTestSuite) SetupTest() {
	suite.service = &mocks.WebsiteStatusService{}
	suite.handler = NewWebsiteStatusHandler(suite.service)
}

func (suite *WebsiteStatusHandlerTestSuite) TearDownTest() {
	t := suite.T()
	suite.service.AssertExpectations(t)
}

func (suite *WebsiteStatusHandlerTestSuite) TestCreateWebsite() {
	t := suite.T()

	reqURL := "/sites"

	t.Run("create website with valid body", func(t *testing.T) {
		patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
			return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
		})
		if err != nil {
			log.Println(err)
		}
		defer func(patch *mpatch.Patch) {
			err := patch.Unpatch()
			if err != nil {
				log.Println(err)
			}
		}(patch)

		website := &domain.Website{
			URL:     "www.google.com",
			AddedAt: time.Now(),
		}

		suite.service.
			On("CreateWebsite", context.Background(), website).
			Return(nil)

		body, _ := json.Marshal(website)
		r := httptest.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.CreateWebsite(w, r)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(
			t,
			"{\"id\":0,\"url\":\"www.google.com\",\"added_at\":\"2020-11-01T00:00:00Z\"}\n",
			w.Body.String(),
		)
	})

	t.Run("create website with missing url", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, reqURL, bytes.NewReader([]byte("{}")))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.CreateWebsite(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "{\"message\":\"invalid request: url is required\"}\n", w.Body.String())
	})

	t.Run("create website with no body", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, reqURL, nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.CreateWebsite(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "{\"message\":\"invalid request: empty request body\"}\n", w.Body.String())
	})

	t.Run("create website fails with db error", func(t *testing.T) {
		patch, err := mpatch.PatchMethod(time.Now, func() time.Time {
			return time.Date(2021, 11, 01, 00, 00, 00, 0, time.UTC)
		})
		if err != nil {
			log.Println(err)
		}
		defer func(patch *mpatch.Patch) {
			err := patch.Unpatch()
			if err != nil {
				log.Println(err)
			}
		}(patch)

		website := &domain.Website{
			URL:     "www.google.com",
			AddedAt: time.Now(),
		}

		suite.service.
			On("CreateWebsite", context.Background(), website).
			Return(myerrors.ErrInternalServerError)

		body, _ := json.Marshal(website)
		r := httptest.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.CreateWebsite(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Equal(
			t,
			"{\"message\":\"internal server error occurred, we are checking...\"}\n",
			w.Body.String(),
		)
	})
}

func (suite *WebsiteStatusHandlerTestSuite) TestGetAllSites() {
	t := suite.T()

	reqURL := "/sites"

	t.Run("get all sites", func(t *testing.T) {
		suite.service.
			On("ListWebsitesStatus", context.Background(), "", int64(10), map[string]string(nil)).
			Return([]domain.WebsiteStatus{}, "test_cursor", nil)

		r := httptest.NewRequest(http.MethodGet, reqURL, nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.GetAllSites(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "[]\n", w.Body.String())
	})

	t.Run("get all sites fails with db error", func(t *testing.T) {
		suite.service.
			On("ListWebsitesStatus", context.Background(), "", int64(5), map[string]string(nil)).
			Return(nil, "", myerrors.ErrInternalServerError)

		r := httptest.NewRequest(http.MethodGet, reqURL+"?limit=5", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.GetAllSites(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Equal(
			t,
			"{\"message\":\"internal server error occurred, we are checking...\"}\n",
			w.Body.String(),
		)
	})
}
