package history

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/ponder2000/rdpms25-template/pkg/core/dto"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/core/repo"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/models"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/util/sqlhelper"
)

type ViewService struct {
	repo      repo.HistoryView
	tableName string
}

func (s *ViewService) queryBuilder(filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) []qm.QueryMod {
	q := make([]qm.QueryMod, 0)
	q = append(q, sqlhelper.FilterQueryBuilder(s.tableName, filters)...)
	q = append(q, complexFilters...)
	q = append(q, sqlhelper.LoaderQueryBuilder(loaders)...)
	return q
}

func (s *ViewService) GetOne(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (*models.AssetHistoryView, error) {
	return s.repo.FetchOne(ctx, s.queryBuilder(filters, complexFilters, loaders))
}

func (s *ViewService) GetAll(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) ([]*models.AssetHistoryView, error) {
	return s.repo.FetchAll(ctx, s.queryBuilder(filters, complexFilters, loaders))
}

func (s *ViewService) GetAllWithKeys(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice, keys []string) ([]*models.AssetHistoryView, error) {
	return s.repo.FetchWithKeys(ctx, s.queryBuilder(filters, complexFilters, loaders), keys)
}

func (s *ViewService) GetPaginated(ctx context.Context, pageRequest *dto.PageRequest, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (*dto.PageResponse[*models.AssetHistoryView], error) {
	cnt, e := s.repo.Count(ctx, s.queryBuilder(filters, complexFilters, nil))
	if e != nil {
		return nil, e
	}

	complexFilters = append(complexFilters, []qm.QueryMod{qm.OrderBy(pageRequest.OrderByClause()), qm.Offset(pageRequest.Offset), qm.Limit(pageRequest.Limit)}...)
	items, e := s.repo.FetchAll(ctx, s.queryBuilder(filters, complexFilters, loaders))
	if e != nil {
		return nil, e
	}

	return &dto.PageResponse[*models.AssetHistoryView]{
		TotalRecords: cnt,
		Offset:       pageRequest.Offset,
		Limit:        pageRequest.Limit,
		Items:        items,
	}, nil
}

func NewViewService(repo repo.HistoryView, tableName string) *ViewService {
	return &ViewService{repo: repo, tableName: tableName}
}
