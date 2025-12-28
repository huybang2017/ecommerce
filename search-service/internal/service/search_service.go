package service

import (
	"context"
	"fmt"
	"search-service/internal/domain"

	"go.uber.org/zap"
)

// SearchService contains the business logic for search operations
// This is the service layer - it orchestrates between repositories
// Following Clean Architecture: business logic is independent of infrastructure
type SearchService struct {
	searchRepo domain.SearchRepository
	logger     *zap.Logger
}

// NewSearchService creates a new search service with all dependencies
// Dependency injection: we inject all repositories and external services
func NewSearchService(
	searchRepo domain.SearchRepository,
	logger *zap.Logger,
) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
		logger:     logger,
	}
}

// SearchProducts performs a search with filters, sort, and pagination
func (s *SearchService) SearchProducts(ctx context.Context, req *domain.SearchRequest) (*domain.SearchResult, error) {
	// Validate request
	if req == nil {
		return nil, fmt.Errorf("search request cannot be nil")
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	// Perform search
	result, err := s.searchRepo.SearchProducts(req)
	if err != nil {
		s.logger.Error("failed to search products",
			zap.String("query", req.Query),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	s.logger.Info("search completed",
		zap.String("query", req.Query),
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page),
		zap.Int("limit", result.Limit),
	)

	return result, nil
}



