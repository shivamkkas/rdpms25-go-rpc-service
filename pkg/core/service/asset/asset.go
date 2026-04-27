package asset

import (
	"context"
	"encoding/hex"
	"rdpms25-go-rpc-service/pkg/core/domain"
	"rdpms25-go-rpc-service/pkg/models"
	"strings"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/ponder2000/rdpms25-template/pkg/util/sqlhelper"
)

type Service struct {
	*domain.SubscriptionHandler[*models.Asset]

	repo      repository.Asset
	tableName string
}

func (s *Service) queryBuilder(filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) []qm.QueryMod {
	q := make([]qm.QueryMod, 0)
	q = append(q, sqlhelper.FilterQueryBuilder(s.tableName, filters)...)
	q = append(q, complexFilters...)
	q = append(q, sqlhelper.LoaderQueryBuilder(loaders)...)
	return q
}

func (s *Service) GetOne(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (*models.Asset, error) {
	return s.repo.FetchOne(ctx, s.queryBuilder(filters, complexFilters, loaders))
}

func (s *Service) GetAll(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) ([]*models.Asset, error) {
	return s.repo.FetchAll(ctx, s.queryBuilder(filters, complexFilters, loaders))
}

func (s *Service) Edit(ctx context.Context, newObj *models.Asset, cols boil.Columns) (*models.Asset, error) {
	if cols.IsNone() {
		cols = boil.Infer()
	}
	res, e := s.repo.Update(ctx, newObj, cols)
	if e != nil {
		return nil, e
	}
	defer s.NotifyUpdate(res)
	return res, nil
}

func (s *Service) Save(ctx context.Context, newObj *models.Asset) (*models.Asset, error) {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := make([]qm.QueryMod, 0)
	q = append(q, models.AssetWhere.AssetTypeID.EQ(newObj.AssetTypeID))
	q = append(q, models.AssetWhere.OrganisationID.EQ(newObj.OrganisationID))
	q = append(q, qm.OrderBy(models.AssetColumns.HexCode+" DESC"))
	q = append(q, qm.Limit(1))
	obj, _ := s.repo.FetchAllWithTx(ctx, tx, q)
	if len(obj) > 0 {
		rawByte, _ := hex.DecodeString(obj[0].HexCode)
		nextHexCodeInt := int(rawByte[0]) + 1
		nextHexCode := strings.ToUpper(hex.EncodeToString([]byte{byte(nextHexCodeInt)}))
		newObj.HexCode = nextHexCode
	} else {
		newObj.HexCode = "0A"
	}
	res, e := s.repo.Insert(ctx, newObj, boil.Infer())
	if e != nil {
		return nil, e
	}
	defer s.NotifyInsert(res)
	return res, nil
}

func (s *Service) Delete(ctx context.Context, ids ...int) error {
	return s.repo.Delete(ctx, ids...)
}

func (s *Service) AssetDetailReport(ctx context.Context, filters *dto.AssetMakeFilters) ([]*dto.AssetMakeDistribution, error) {
	return s.repo.AssetDetailReport(ctx, filters)
}

func NewService(repo repository.Asset, tableName string) *Service {
	return &Service{SubscriptionHandler: domain.NewSubscriptionHandler[*models.Asset](), repo: repo, tableName: tableName}
}
