package history

import (
	"context"
	"fmt"
	"strings"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/models"
)

type ViewRepo struct{}

func NewViewRepo() *ViewRepo {
	return &ViewRepo{}
}

func (v *ViewRepo) FetchOne(ctx context.Context, q []qm.QueryMod) (*models.AssetHistoryView, error) {
	if res, e := models.AssetHistoryViews(q...).One(ctx, boil.GetContextDB()); e != nil {
		return nil, e
	} else {
		return res, nil
	}
}

func (v *ViewRepo) FetchAll(ctx context.Context, q []qm.QueryMod) ([]*models.AssetHistoryView, error) {
	if res, e := models.AssetHistoryViews(q...).All(ctx, boil.GetContextDB()); e != nil {
		return nil, e
	} else {
		return res, nil
	}
}

func (v *ViewRepo) FetchWithKeys(ctx context.Context, q []qm.QueryMod, keys []string) ([]*models.AssetHistoryView, error) {
	if len(keys) > 0 {
		keysArr := "{" + strings.Join(keys, ",") + "}"
		selectClause := fmt.Sprintf("(SELECT jsonb_object_agg(key, value) FROM jsonb_each(data) WHERE key = ANY('%s')) as data", keysArr)
		q = append(q, qm.Select("asset_id", "asset_type_id", "ts", selectClause))
	} else {
		q = append(q, qm.Select("asset_id", "asset_type_id", "ts", "data"))
	}

	if res, e := models.AssetHistoryViews(q...).All(ctx, boil.GetContextDB()); e != nil {
		return nil, e
	} else {
		return res, nil
	}
}

func (v *ViewRepo) Count(ctx context.Context, q []qm.QueryMod) (int64, error) {
	if cnt, e := models.AssetHistoryViews(q...).CountG(ctx); e != nil {
		return 0, e
	} else {
		return cnt, nil
	}
}
