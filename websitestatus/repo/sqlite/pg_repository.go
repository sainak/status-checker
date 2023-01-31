package sqlite

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/logger"
	"github.com/sainak/status-checker/core/myerrors"
	"github.com/sainak/status-checker/helpers/repo"
)

func NewWebsiteStatusRepo(db *sqlx.DB) domain.WebsiteStatusStorer {
	return &pgWebsiteStatusRepo{db}
}

type pgWebsiteStatusRepo struct {
	DB *sqlx.DB
}

func (s pgWebsiteStatusRepo) QueryWebsites(
	ctx context.Context,
	cursor string,
	num int64,
	filters map[string]string,
) (websites []domain.Website, nextCursor string, err error) {
	query := `SELECT websites.id, websites.url, websites.added_at
		FROM websites
		WHERE added_at > $1 ORDER BY added_at LIMIT $2`

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		logger.Error(err)
		err = myerrors.ErrBadParamInput
		return
	}
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	rows, err := stmt.QueryContext(ctx, decodedCursor, num)
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)

	for rows.Next() {
		website := domain.Website{}
		err = rows.Scan(&website.ID, &website.URL, &website.AddedAt)
		if err != nil {
			logger.Error(err)
		}
		websites = append(websites, website)
	}

	if err != nil {
		logger.Error(err)
		return
	}
	if len(websites) == int(num) {
		nextCursor = repo.EncodeCursor(websites[len(websites)-1].AddedAt.ValueOrZero())
	}
	return
}

func (s pgWebsiteStatusRepo) QueryWebsitesStatus(
	ctx context.Context,
	cursor string,
	num int64,
	filters map[string]string,
) (websites []domain.WebsiteStatus, nextCursor string, err error) {
	query := `SELECT websites.id, websites.added_at, websites.url, 
		website_statuses.id AS status_id, website_statuses.up, website_statuses.checked_at
		FROM websites
		LEFT JOIN website_statuses ON websites.id = website_statuses.website_id
		AND website_statuses.checked_at = (SELECT MAX(checked_at) FROM website_statuses WHERE website_id = websites.id)
		WHERE added_at > $1 ORDER BY added_at LIMIT $2`

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		logger.Error(err)
		err = myerrors.ErrBadParamInput
		return
	}
	stmt, err := s.DB.PreparexContext(ctx, query)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.QueryxContext(ctx, decodedCursor, num)
	if err != nil {
		panic(err)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error(err)
		}
	}(rows)

	for rows.Next() {
		website := domain.WebsiteStatus{}
		err = rows.StructScan(&website)
		if err != nil {
			logger.Error(err)
		}
		websites = append(websites, website)
	}

	if err != nil {
		logger.Error(err)
		return
	}
	if len(websites) == int(num) {
		nextCursor = repo.EncodeCursor(websites[len(websites)-1].AddedAt.ValueOrZero())
	}
	return
}

func (s pgWebsiteStatusRepo) InsertWebsite(
	ctx context.Context,
	website *domain.Website,
) (err error) {
	query := `INSERT INTO websites (id, url, added_at) VALUES (DEFAULT, $1, $2) RETURNING id`
	stmt, err := s.DB.PreparexContext(ctx, query)
	if err != nil {
		panic(err)
		return
	}
	err = stmt.QueryRowxContext(ctx, website.URL, website.AddedAt).Scan(&website.ID)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (s pgWebsiteStatusRepo) QueryWebsiteStatusByID(
	ctx context.Context,
	id int64,
) (website domain.WebsiteStatus, err error) {
	query := `SELECT websites.id, websites.added_at, websites.url, 
		website_statuses.id AS status_id, website_statuses.up, website_statuses.checked_at
		FROM websites
		LEFT JOIN website_statuses ON websites.id = website_statuses.website_id
		AND website_statuses.time = (SELECT MAX(checked_at) FROM website_statuses WHERE website_id = websites.id)
		WHERE websites.id = $1`
	stmt, err := s.DB.PreparexContext(ctx, query)
	if err != nil {
		panic(err)
	}
	err = stmt.QueryRowxContext(ctx, id).StructScan(&website)
	if err != nil {
		return
	}
	return
}

func (s pgWebsiteStatusRepo) UpdateIntoWebsite(
	ctx context.Context,
	id int64,
	website *domain.Website,
) (err error) {
	query := `UPDATE websites SET url = $2 WHERE id = $1`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, id, website.URL)
	return
}

func (s pgWebsiteStatusRepo) DropWebsite(
	ctx context.Context,
	id int64,
) (err error) {
	query := `DELETE FROM websites WHERE id = $1`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, id)
	return
}

func (s pgWebsiteStatusRepo) InsertStatus(
	ctx context.Context,
	status *domain.Status,
) (err error) {
	query := `INSERT INTO website_statuses (id, website_id, up, checked_at) VALUES (DEFAULT, $1, $2, $3) RETURNING id`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, status.WebsiteID, status.Up, status.CheckedAt).Scan(&status.ID)
	return
}

func (s pgWebsiteStatusRepo) QueryStatusesByWebsiteID(
	ctx context.Context,
	websiteID int64,
	cursor string,
	num int64,
) (statuses []domain.Status, nextCursor string, err error) {
	query := `SELECT id, website_id, up, checked_at 
		FROM website_statuses 
		WHERE website_id = $1 AND checked_at > $2 ORDER BY checked_at LIMIT $3`
	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		err = myerrors.ErrBadParamInput
		return
	}
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	rows, err := stmt.QueryContext(ctx, websiteID, decodedCursor, num)
	if err != nil {
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	err = rows.Scan(&statuses)
	if err != nil {
		return
	}
	if len(statuses) == int(num) {
		nextCursor = repo.EncodeCursor(statuses[len(statuses)-1].CheckedAt.ValueOrZero())
	}
	return
}
