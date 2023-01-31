package sqlite

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/errors"
	"github.com/sainak/status-checker/core/logger"
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
		err = errors.ErrBadParamInput
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
		nextCursor = repo.EncodeCursor(websites[len(websites)-1].AddedAt)
	}
	return
}

func (s pgWebsiteStatusRepo) QueryWebsitesWithStatus(
	ctx context.Context,
	cursor string,
	num int64,
	filters map[string]string,
) (websites []domain.Website, nextCursor string, err error) {
	query := `SELECT websites.id, websites.url, websites.added_at, website_statuses.up, website_statuses.time
		FROM websites
		LEFT JOIN website_statuses ON websites.id = website_statuses.website_id
		AND website_statuses.time = (SELECT MAX(time) FROM website_statuses WHERE website_id = websites.id)
		WHERE added_at > $1 ORDER BY added_at LIMIT $2`

	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil {
		logger.Error(err)
		err = errors.ErrBadParamInput
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
		var up sql.NullBool
		var time sql.NullTime
		err = rows.Scan(&website.ID, &website.URL, &website.AddedAt, &up, &time)
		website.Status.Up = up.Bool
		website.Status.Time = time.Time
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
		nextCursor = repo.EncodeCursor(websites[len(websites)-1].AddedAt)
	}
	return
}

func (s pgWebsiteStatusRepo) InsertWebsite(
	ctx context.Context,
	website *domain.Website,
) (err error) {
	query := `INSERT INTO websites (id, url, added_at) VALUES (DEFAULT, $1, $2) RETURNING id`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
		return
	}
	err = stmt.QueryRowContext(ctx, website.URL, website.AddedAt).Scan(&website.ID)
	if err != nil {
		// check if the error is a duplicate key error
		//if _errors.Is(err, sql.) {
		//if err.Error() == "UNIQUE constraint failed: websites.url" {
		//	err = errors.ErrConflict
		//}
		panic(err)
		return
	}
	return err
}

func (s pgWebsiteStatusRepo) QueryWebsiteWithStatusByID(
	ctx context.Context,
	id int64,
) (website domain.Website, err error) {
	query := `SELECT websites.id, websites.url, websites.added_at, website_statuses.up, website_statuses.time
		FROM websites
		LEFT JOIN website_statuses ON websites.id = website_statuses.website_id
		AND website_statuses.time = (SELECT MAX(time) FROM website_statuses WHERE website_id = websites.id)
		WHERE websites.id = $1`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	var up sql.NullBool
	var time sql.NullTime
	err = stmt.QueryRowContext(ctx, id).Scan(&website.ID, &website.URL, &website.AddedAt, &up, &time)
	website.Status.Up = up.Bool
	website.Status.Time = time.Time
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

func (s pgWebsiteStatusRepo) InsertWebsiteStatus(
	ctx context.Context,
	status *domain.WebsiteStatus,
) (err error) {
	query := `INSERT INTO website_statuses (id, website_id, up, time) VALUES (DEFAULT, $1, $2, $3) RETURNING id`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, status.WebsiteID, status.Up, status.Time).Scan(&status.ID)
	return
}

func (s pgWebsiteStatusRepo) QueryStatusesForWebsite(
	ctx context.Context,
	websiteID int64,
	cursor string,
	num int64,
) (statuses []domain.WebsiteStatus, nextCursor string, err error) {
	query := `SELECT id, website_id, up, time 
		FROM website_statuses 
		WHERE website_id = $1 AND time > $2 ORDER BY time LIMIT $3`
	decodedCursor, err := repo.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		err = errors.ErrBadParamInput
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
		nextCursor = repo.EncodeCursor(statuses[len(statuses)-1].Time)
	}
	return
}
