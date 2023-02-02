package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v4"

	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/domain/mocks"
	"github.com/sainak/status-checker/core/myerrors"
)

var (
	defTime = null.NewTime(time.Now(), true)
	defBool = null.NewBool(true, true)
	defInt  = null.NewInt(1, true)
)

type WebsiteStatusServiceTestSuit struct {
	suite.Suite
	service domain.WebsiteStatusService
	repo    *mocks.WebsiteStatusStorer
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WebsiteStatusServiceTestSuit))
}

func (suite *WebsiteStatusServiceTestSuit) SetupTest() {
	suite.repo = &mocks.WebsiteStatusStorer{}
	suite.service = NewWebsiteStatusService(suite.repo, 30*time.Second)
}

func (suite *WebsiteStatusServiceTestSuit) TearDownSuite() {
	suite.repo.AssertExpectations(suite.T())
}

func (suite *WebsiteStatusServiceTestSuit) TestListWebsitesStatus() {
	t := suite.T()

	list := []domain.WebsiteStatus{
		{1, "https://google.com", defTime.ValueOrZero(), defInt, defBool, defTime},
		{2, "https://facebook.com", defTime.ValueOrZero(), defInt, defBool, defTime},
		{3, "https://twitter.com", defTime.ValueOrZero(), defInt, defBool, defTime},
	}

	t.Run("when list website is successful", func(t *testing.T) {
		suite.repo.
			On("QueryWebsitesStatus", mock.Anything, "", int64(2), mock.Anything).
			Return(list, "", nil)
		res, nextCursor, err := suite.service.ListWebsitesStatus(context.Background(), "", int64(2), nil)
		require.Equal(t, list, res)
		require.Equal(t, "", nextCursor)
		require.Nil(t, err)
	})

	t.Run("when list website is failed", func(t *testing.T) {
		//either call `suite.SetupTest()` or register mock call with different arguments
		suite.repo.
			On("QueryWebsitesStatus", mock.Anything, "invalid cusror", int64(10), mock.Anything).
			Return([]domain.WebsiteStatus{}, "", myerrors.ErrInternalServerError)
		res, nextCursor, err := suite.service.ListWebsitesStatus(context.Background(), "invalid cusror", int64(10), nil)
		require.Empty(t, res)
		require.Equal(t, "", nextCursor)
		require.ErrorAs(t, err, &myerrors.ErrInternalServerError)
	})

}

func (suite *WebsiteStatusServiceTestSuit) TestCreateWebsite() {
	t := suite.T()

	website := domain.Website{ID: 0, URL: "https://google.com", AddedAt: time.Now()}

	t.Run("when create website is successful", func(t *testing.T) {
		suite.repo.
			On("InsertWebsite", mock.Anything, &website).
			Return(nil)
		err := suite.service.CreateWebsite(context.Background(), &website)
		require.Nil(t, err)
	})

	t.Run("when create website is failed", func(t *testing.T) {
		suite.repo.
			On("InsertWebsite", mock.Anything, &domain.Website{}).
			Return(myerrors.ErrInternalServerError)
		err := suite.service.CreateWebsite(context.Background(), &domain.Website{})
		require.ErrorAs(t, err, &myerrors.ErrInternalServerError)
	})
}

func (suite *WebsiteStatusServiceTestSuit) TestGetWebsiteStatusByID() {
	t := suite.T()

	defTime := null.NewTime(time.Now(), true)
	defBool := null.NewBool(true, true)

	website := domain.WebsiteStatus{
		ID:        1,
		URL:       "https://google.com",
		AddedAt:   defTime.ValueOrZero(),
		StatusID:  null.NewInt(1, true),
		Up:        defBool,
		CheckedAt: defTime,
	}

	t.Run("when get website status by id is successful", func(t *testing.T) {
		suite.repo.
			On("QueryWebsiteStatusByID", mock.Anything, int64(1)).
			Return(website, nil)
		res, err := suite.service.GetWebsiteStatusByID(context.Background(), int64(1))
		require.Equal(t, website, res)
		require.Nil(t, err)
	})

	t.Run("when get website status by id is failed", func(t *testing.T) {
		suite.repo.
			On("QueryWebsiteStatusByID", mock.Anything, int64(-1)).
			Return(domain.WebsiteStatus{}, myerrors.ErrInternalServerError)
		res, err := suite.service.GetWebsiteStatusByID(context.Background(), int64(-1))
		require.Empty(t, res)
		require.ErrorAs(t, err, &myerrors.ErrInternalServerError)
	})

	t.Run("when get website status by id is not found", func(t *testing.T) {
		suite.repo.
			On("QueryWebsiteStatusByID", mock.Anything, int64(100)).
			Return(domain.WebsiteStatus{}, myerrors.ErrNotFound)
		res, err := suite.service.GetWebsiteStatusByID(context.Background(), int64(100))
		require.Empty(t, res)
		require.ErrorAs(t, err, &myerrors.ErrNotFound)
	})
}

func (suite *WebsiteStatusServiceTestSuit) TestDeleteWebsite() {
	t := suite.T()

	t.Run("when delete website is successful", func(t *testing.T) {
		suite.repo.
			On("DropWebsite", mock.Anything, int64(1)).
			Return(nil)
		err := suite.service.DeleteWebsite(context.Background(), int64(1))
		require.Nil(t, err)
	})

	t.Run("when delete website is failed", func(t *testing.T) {
		suite.repo.
			On("DropWebsite", mock.Anything, int64(-1)).
			Return(myerrors.ErrInternalServerError)
		err := suite.service.DeleteWebsite(context.Background(), int64(-1))
		require.ErrorAs(t, err, &myerrors.ErrInternalServerError)
	})

	t.Run("when delete website, website is not found", func(t *testing.T) {
		suite.repo.
			On("DropWebsite", mock.Anything, int64(100)).
			Return(myerrors.ErrNotFound)
		err := suite.service.DeleteWebsite(context.Background(), int64(100))
		require.ErrorAs(t, err, &myerrors.ErrNotFound)
	})
}

func (suite *WebsiteStatusServiceTestSuit) TestListWebsiteStatuses() {
	t := suite.T()

	list := []domain.Status{
		{1, true, time.Now(), 1},
		{12, false, time.Now(), 2},
		{13, true, time.Now(), 3},
	}

	t.Run("when list website status is successful", func(t *testing.T) {
		suite.repo.
			On("QueryStatusesByWebsiteID", mock.Anything, int64(1), "", int64(10)).
			Return(list, "", nil)
		res, cursor, err := suite.service.ListWebsiteStatuses(context.Background(), int64(1), "", int64(10))
		require.Equal(t, list, res)
		require.Equal(t, "", cursor)
		require.Nil(t, err)
	})

	t.Run("when list website status is failed", func(t *testing.T) {
		suite.repo.
			On("QueryStatusesByWebsiteID", mock.Anything, int64(-1), "", int64(10)).
			Return([]domain.Status{}, "", myerrors.ErrInternalServerError)
		res, cursor, err := suite.service.ListWebsiteStatuses(context.Background(), int64(-1), "", int64(10))
		require.Empty(t, res)
		require.Equal(t, "", cursor)
		require.ErrorAs(t, err, &myerrors.ErrInternalServerError)
	})

	t.Run("when list website status is empty", func(t *testing.T) {
		suite.repo.
			On("QueryStatusesByWebsiteID", mock.Anything, int64(100), "", int64(10)).
			Return([]domain.Status{}, "", nil)
		res, cursor, err := suite.service.ListWebsiteStatuses(context.Background(), int64(100), "", int64(10))
		require.Empty(t, res)
		require.Equal(t, "", cursor)
		require.Nil(t, err)
	})
}
