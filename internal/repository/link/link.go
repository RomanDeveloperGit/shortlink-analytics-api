package link

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/model"
)

type repository struct {
	db clickhouse.Conn
}

func NewRepository(db clickhouse.Conn) *repository {
	return &repository{db}
}

func (r *repository) GetLinkCreatedEventsWithUnexpiredStatusAndWithAuthorEmail(
	ctx context.Context,
	from, to time.Time,
) ([]model.LinkCreatedEvent, error) {
	var events []model.LinkCreatedEvent

	q := `
		SELECT
			link_id,
			short_code,
			full_url,
			author_email,
			expires_at,
			created_at
		FROM link_created_event
		WHERE expires_at >= ?
		AND expires_at <= ?
		AND created_at >= ?
		AND created_at <= ?
		AND author_email IS NOT NULL
		ORDER BY created_at
	`

	err := r.db.Select(ctx, &events, q, from, to, from, to)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *repository) GetLinkVisitedEventsByLinkIDs(
	ctx context.Context,
	from, to time.Time,
	linkIDs []int,
) ([]model.LinkVisitedEvent, error) {
	var events []model.LinkVisitedEvent

	q := `
		SELECT
			link_id,
			timestamp
		FROM link_visited_event
		WHERE link_id IN ?
		AND timestamp >= ?
		AND timestamp <= ?
		ORDER BY timestamp
	`

	err := r.db.Select(ctx, &events, q, linkIDs, from, to)
	if err != nil {
		return nil, err
	}

	return events, nil
}
