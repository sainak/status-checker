package pg

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/logger"
	"github.com/sainak/status-checker/core/myerrors"
	"github.com/sainak/status-checker/helpers/repo"
)

const (
	queryWebsitesQuery          = `SELECT websites.id, websites.url, websites.added_at FROM websites WHERE added_at > $1 ORDER BY added_at LIMIT $2`
	queryWebsiteStatusQuery     = `SELECT websites.id, websites.added_at, websites.url, website_statuses.id AS status_id, website_statuses.up, website_statuses.checked_at FROM websites LEFT JOIN website_statuses ON websites.id = website_statuses.website_id AND website_statuses.checked_at = (SELECT MAX(checked_at) FROM website_statuses WHERE website_id = websites.id) WHERE added_at > $1 ORDER BY added_at LIMIT $2`
	insertWebsiteQuery          = `INSERT INTO websites (url, added_at) VALUES ($1, $2) RETURNING id`
	queryWebsiteStatusByIDQuery = `SELECT websites.id, websites.added_at, websites.url,  website_statuses.id AS status_id, website_statuses.up, website_statuses.checked_at FROM websites LEFT JOIN website_statuses ON websites.id = website_statuses.website_id AND website_statuses.time = (SELECT MAX(checked_at) FROM website_statuses WHERE website_id = websites.id) WHERE websites.id = $1`
	dropWebsiteQuery            = `DELETE FROM websites WHERE id = $1`
	insertStatusQuery           = `INSERT INTO website_statuses (id, website_id, up, checked_at) VALUES (DEFAULT, $1, $2, $3) RETURNING id`
	queryStatusQuery            = `SELECT id, website_id, up, checked_at FROM website_statuses  WHERE website_id = $1 AND checked_at > $2 ORDER BY checked_at LIMIT $3`
)

func NewWebsiteStatusRepo(db *sqlx.DB) domain.WebsiteStatusStorer {
	return &pgWebsiteStatusRepo{db}
}

type pgWebsiteStatusRepo struct {
	DB *sqlx.DB
}

func (r pgWebsiteStatusRepo) QueryWebsites(
	ctx context.Context,
	cursor string,
	num int64,
	filters map[string]string,
) ([]domain.Website, string, error) {

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		logger.Error(err)
		return []domain.Website{}, "", myerrors.ErrBadCursor
	}

	stmt, err := r.DB.PreparexContext(ctx, queryWebsitesQuery)
	if err != nil {
		logger.Error(err)
		return []domain.Website{}, "", myerrors.ErrInternalServerError
	}

	rows, err := stmt.QueryxContext(ctx, decodedCursor, num)
	if err != nil {
		logger.Error(err)
		return []domain.Website{}, "", myerrors.ParseDBError(err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)

	var websites []domain.Website
	for rows.Next() {
		website := domain.Website{}
		err = rows.Scan(&website.ID, &website.URL, &website.AddedAt)
		if err != nil {
			logger.Error(err)
			return []domain.Website{}, "", myerrors.ParseDBError(err)
		}
		websites = append(websites, website)
	}
	var nextCursor string
	if len(websites) == int(num) {
		nextCursor = repo.EncodeCursor(websites[len(websites)-1].AddedAt)
	}
	return websites, nextCursor, nil
}

func (r pgWebsiteStatusRepo) QueryWebsitesStatus(
	ctx context.Context,
	cursor string,
	num int64,
	filters map[string]string,
) ([]domain.WebsiteStatus, string, error) {

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		logger.Error(err)
		return []domain.WebsiteStatus{}, "", myerrors.ErrBadCursor
	}

	stmt, err := r.DB.PreparexContext(ctx, queryWebsiteStatusQuery)
	if err != nil {
		logger.Error(err)
		return []domain.WebsiteStatus{}, "", myerrors.ErrInternalServerError
	}

	rows, err := stmt.QueryxContext(ctx, decodedCursor, num)
	if err != nil {
		logger.Error(err)
		return []domain.WebsiteStatus{}, "", myerrors.ParseDBError(err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)

	var websites []domain.WebsiteStatus
	for rows.Next() {
		website := domain.WebsiteStatus{}
		err = rows.StructScan(&website)
		if err != nil {
			logger.Error(err)
			return []domain.WebsiteStatus{}, "", myerrors.ParseDBError(err)
		}
		websites = append(websites, website)
	}
	var nextCursor string
	if len(websites) == int(num) {
		nextCursor = repo.EncodeCursor(websites[len(websites)-1].AddedAt)
	}
	return websites, nextCursor, nil
}

func (r pgWebsiteStatusRepo) InsertWebsite(
	ctx context.Context,
	website *domain.Website,
) error {

	stmt, err := r.DB.PreparexContext(ctx, insertWebsiteQuery)
	if err != nil {
		logger.Error(err)
		return myerrors.ErrInternalServerError
	}

	err = stmt.QueryRowxContext(ctx, website.URL, website.AddedAt).Scan(&website.ID)
	if err != nil {
		logger.Error(err)
		return myerrors.ParseDBError(err)
	}
	return nil
}

func (r pgWebsiteStatusRepo) QueryWebsiteStatusByID(
	ctx context.Context,
	id int64,
) (domain.WebsiteStatus, error) {

	stmt, err := r.DB.PreparexContext(ctx, queryWebsiteStatusByIDQuery)
	if err != nil {
		logger.Error(err)
		return domain.WebsiteStatus{}, myerrors.ErrInternalServerError
	}

	var website domain.WebsiteStatus
	err = stmt.QueryRowxContext(ctx, id).StructScan(&website)
	if err != nil {
		logger.Error(err)
		return domain.WebsiteStatus{}, myerrors.ParseDBError(err)
	}
	return website, nil
}

func (r pgWebsiteStatusRepo) DropWebsite(
	ctx context.Context,
	id int64,
) error {

	stmt, err := r.DB.PrepareContext(ctx, dropWebsiteQuery)
	if err != nil {
		logger.Error(err)
		return myerrors.ErrInternalServerError
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		logger.Error(err)
		return myerrors.ParseDBError(err)
	}
	return nil
}

func (r pgWebsiteStatusRepo) InsertStatus(
	ctx context.Context,
	status *domain.Status,
) error {

	stmt, err := r.DB.PrepareContext(ctx, insertStatusQuery)
	if err != nil {
		logger.Error(err)
		return myerrors.ErrInternalServerError
	}

	err = stmt.QueryRowContext(ctx, status.WebsiteID, status.Up, status.CheckedAt).Scan(&status.ID)
	if err != nil {
		logger.Error(err)
		return myerrors.ParseDBError(err)
	}
	return nil
}

func (r pgWebsiteStatusRepo) QueryStatusesByWebsiteID(
	ctx context.Context,
	websiteID int64,
	cursor string,
	num int64,
) ([]domain.Status, string, error) {

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		logger.Error(err)
		return []domain.Status{}, "", myerrors.ErrBadCursor
	}

	stmt, err := r.DB.PreparexContext(ctx, queryStatusQuery)
	if err != nil {
		logger.Error(err)
		return []domain.Status{}, "", myerrors.ErrInternalServerError
	}

	rows, err := stmt.QueryxContext(ctx, websiteID, decodedCursor, num)
	if err != nil {
		logger.Error(err)
		return []domain.Status{}, "", myerrors.ParseDBError(err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)

	var statuses []domain.Status
	err = rows.StructScan(&statuses)
	if err != nil {
		logger.Error(err)
		return []domain.Status{}, "", myerrors.ParseDBError(err)
	}
	var nextCursor string
	if len(statuses) == int(num) {
		nextCursor = repo.EncodeCursor(statuses[len(statuses)-1].CheckedAt)
	}
	return statuses, nextCursor, nil
}
