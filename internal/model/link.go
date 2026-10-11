package model

import "time"

type LinkCreatedEvent struct {
	LinkID      int32     `ch:"link_id"`
	ShortCode   string    `ch:"short_code"`
	FullURL     string    `ch:"full_url"`
	AuthorEmail *string   `ch:"author_email"`
	ExpiresAt   time.Time `ch:"expires_at"`
	CreatedAt   time.Time `ch:"created_at"`
}

type LinkVisitedEvent struct {
	LinkID    int32     `ch:"link_id"`
	Timestamp time.Time `ch:"timestamp"`
}
