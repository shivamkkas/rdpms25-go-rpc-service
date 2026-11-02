package repo

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/core/repo/history"
	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/models"
)

type HistoryView interface {
	ReadableView[*models.AssetHistoryView]
	FetchWithKeys(ctx context.Context, q []qm.QueryMod, keys []string) ([]*models.AssetHistoryView, error)
}

func NewHistoryView() HistoryView { return history.NewViewRepo() }
