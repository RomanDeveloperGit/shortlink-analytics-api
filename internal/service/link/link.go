package link

import (
	"context"
	"fmt"
	"time"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/model"
)

type LinkStats struct {
	ShortCode   string
	FullURL     string
	TotalVisits int
	// UniqueVisits int -- потом, когда всякие ip, user agent будем засыпать в кликхаус (а по-хорошему это надо делать)
	FirstVisitAt *time.Time
	LastVisitAt  *time.Time
	CreatedAt    time.Time
}

type LinkRepository interface {
	GetLinkCreatedEventsWithUnexpiredStatusAndWithAuthorEmail(ctx context.Context, from, to time.Time) ([]model.LinkCreatedEvent, error)
	GetLinkVisitedEventsByLinkIDs(ctx context.Context, from, to time.Time, linkIDs []int) ([]model.LinkVisitedEvent, error)
}

type service struct {
	linkRepository LinkRepository
}

func NewService(
	linkRepository LinkRepository,
) *service {
	return &service{
		linkRepository: linkRepository,
	}
}

func (s *service) convertLinkVisitedEventsToMapByLinkIDCount(
	linkVisitedEvents []model.LinkVisitedEvent,
) map[int][]model.LinkVisitedEvent {
	eventsMap := make(map[int][]model.LinkVisitedEvent, len(linkVisitedEvents))

	for _, event := range linkVisitedEvents {
		eventsMap[int(event.LinkID)] = append(eventsMap[int(event.LinkID)], event)
	}

	return eventsMap
}

func (s *service) convertLinkCreatedEventsToMapByAuthorEmail(
	linkCreatedEvents []model.LinkCreatedEvent,
) map[string][]model.LinkCreatedEvent {
	eventsMap := make(map[string][]model.LinkCreatedEvent, len(linkCreatedEvents))

	for _, event := range linkCreatedEvents {
		eventsMap[*event.AuthorEmail] = append(eventsMap[*event.AuthorEmail], event)
	}

	return eventsMap
}

// TODO: лучше этот метод выносить внутрь cron-job?
func (s *service) GetStatsByAuthorsInPeriod(
	ctx context.Context,
	from, to time.Time,
) (map[string][]LinkStats, error) {
	linkCreatedEvents, err := s.linkRepository.GetLinkCreatedEventsWithUnexpiredStatusAndWithAuthorEmail(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get link created events: %v", err)
	}

	uniqueLinkIDs := make([]int, 0, len(linkCreatedEvents))

	for _, event := range linkCreatedEvents {
		uniqueLinkIDs = append(uniqueLinkIDs, int(event.LinkID))
	}

	linkVisitedEvents, err := s.linkRepository.GetLinkVisitedEventsByLinkIDs(ctx, from, to, uniqueLinkIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get link visited events: %v", err)
	}

	linkVisitedEventsMapByLinkID := s.convertLinkVisitedEventsToMapByLinkIDCount(linkVisitedEvents)
	linkCreatedEventsMapByAuthorEmail := s.convertLinkCreatedEventsToMapByAuthorEmail(linkCreatedEvents)

	result := make(map[string][]LinkStats, len(linkCreatedEventsMapByAuthorEmail))

	for authorEmail, linkCreatedEvents := range linkCreatedEventsMapByAuthorEmail {
		for _, linkCreatedEvent := range linkCreatedEvents {
			linkVisitedEvents, ok := linkVisitedEventsMapByLinkID[int(linkCreatedEvent.LinkID)]

			var firstVisitedAt, lastVisitedAt *time.Time

			if ok {
				firstVisitedAt = &linkVisitedEvents[0].Timestamp
				lastVisitedAt = &linkVisitedEvents[len(linkVisitedEvents)-1].Timestamp
			}

			linkStats := LinkStats{
				ShortCode:    linkCreatedEvent.ShortCode,
				FullURL:      linkCreatedEvent.FullURL,
				TotalVisits:  len(linkVisitedEvents),
				FirstVisitAt: firstVisitedAt,
				LastVisitAt:  lastVisitedAt,
				CreatedAt:    linkCreatedEvent.CreatedAt,
			}

			result[authorEmail] = append(result[authorEmail], linkStats)
		}
	}

	return result, nil
}
