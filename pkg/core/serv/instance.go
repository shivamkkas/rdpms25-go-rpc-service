package serv

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/core/repo"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/core/serv/history"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/models"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/util/sqlhelper"
)

type HistoryView interface {
	ReadableView[*models.AssetHistoryView]
	GetAllWithKeys(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice, keys []string) ([]*models.AssetHistoryView, error)
}

func NewHistoryView(repo repo.HistoryView) HistoryView {
	return history.NewViewService(repo, models.ViewNames.AssetHistoryView)
}
