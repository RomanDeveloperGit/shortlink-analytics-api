package register

import (
	"context"
	"fmt"
	"time"

	analyticsv1 "github.com/RomanDeveloperGit/shortlink-protos/gen/go/shortlink/analytics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/RomanDeveloperGit/shortlink-analytics-api/internal/service/link"
)

type AnalyticsService interface {
	GetStatsByAuthorsInPeriod(
		ctx context.Context,
		from, to time.Time,
	) (map[string][]link.LinkStats, error)
}

type analyticsRegister struct {
	analyticsv1.UnimplementedAnalyticsServiceServer

	analyticsService AnalyticsService
}

func newAnalyticsHandler(analyticsService AnalyticsService) *analyticsRegister {
	return &analyticsRegister{
		analyticsService: analyticsService,
	}
}

func (r *analyticsRegister) GetStatsByAuthorsInPeriod(
	ctx context.Context,
	req *analyticsv1.GetStatsByAuthorsInPeriodRequest,
) (*analyticsv1.GetStatsByAuthorsInPeriodResponse, error) {
	var from, to time.Time

	if req.GetFrom() == nil {
		return nil, status.Error(codes.InvalidArgument, "field 'from' is required")
	}
	if req.GetTo() == nil {
		return nil, status.Error(codes.InvalidArgument, "field 'to' is required")
	}

	if req.GetFrom().AsTime().After(req.GetTo().AsTime()) {
		return nil, status.Error(codes.InvalidArgument, "'from' must be before 'to'")
	}

	if req.GetFrom() != nil {
		from = req.GetFrom().AsTime()
	}
	if req.GetTo() != nil {
		to = req.GetTo().AsTime()
	}

	statsByAuthor, err := r.analyticsService.GetStatsByAuthorsInPeriod(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("get stats by authors in period: %w", err)
	}

	resp := &analyticsv1.GetStatsByAuthorsInPeriodResponse{
		StatsByAuthor: make(map[string]*analyticsv1.LinkStatsList, len(statsByAuthor)),
	}

	for authorEmail, stats := range statsByAuthor {
		list := &analyticsv1.LinkStatsList{
			Items: make([]*analyticsv1.LinkStats, 0, len(stats)),
		}

		for _, s := range stats {
			linkStats := &analyticsv1.LinkStats{
				ShortLink:   s.ShortCode,
				FullUrl:     s.FullURL,
				TotalVisits: int32(s.TotalVisits),
				CreatedAt:   timestamppb.New(s.CreatedAt),
			}

			if s.FirstVisitAt != nil {
				linkStats.FirstVisitAt = timestamppb.New(*s.FirstVisitAt)
			}
			if s.LastVisitAt != nil {
				linkStats.LastVisitAt = timestamppb.New(*s.LastVisitAt)
			}

			list.Items = append(list.Items, linkStats)
		}

		resp.StatsByAuthor[authorEmail] = list
	}

	return resp, nil
}

func NewRegisterAnalyticsService(analyticsService AnalyticsService) func(*grpc.Server) {
	return func(server *grpc.Server) {
		analyticsv1.RegisterAnalyticsServiceServer(server, newAnalyticsHandler(analyticsService))
	}
}
