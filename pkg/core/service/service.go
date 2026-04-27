package service

import (
	"context"
	"rdpms25-go-rpc-service/pkg/core/domain"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/ponder2000/rdpms25-template/pkg/core/dto"
	"github.com/ponder2000/rdpms25-template/pkg/util/sqlhelper"
)

type Deletable interface {
	Delete(ctx context.Context, ids ...int) error
}

type Creatable[T any] interface {
	Save(ctx context.Context, newObj T) (T, error)
}

type Updatable[T any] interface {
	Edit(ctx context.Context, newObj T, cols boil.Columns) (T, error)
}

type Readable[T any] interface {
	GetOne(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (T, error)
	GetAll(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) ([]T, error)
}

type PaginatedReadable[T any] interface {
	GetPaginated(ctx context.Context, pageRequest *dto.PageRequest, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (*dto.PageResponse[T], error)
}

// type Countable[T any] interface {
// 	Count(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod) (int64, error)
// }

type CRUD[T any] interface {
	Creatable[T]
	Updatable[T]
	Readable[T]
	Deletable
}

type PaginatedCRUD[T any] interface {
	CRUD[T]
	PaginatedReadable[T]
}

type ReadableView[T any] interface {
	Readable[T]
	PaginatedReadable[T]
	// Countable[T]
}

type Subscription[T any] interface {
	Subscribe(id string) (<-chan *domain.OperationEvent[T], error)
	UnSubscribe(id string) error
}

type ReportView[T any, R any] interface {
	Live[T, R]
	Detail[T, R]
}

type Live[T any, R any] interface {
	Live(ctx context.Context, filters T) ([]R, error)
}

type Detail[T any, R any] interface {
	Detail(ctx context.Context, filters T) ([]R, error)
}

type Upsertable[T any] interface {
	Upsert(ctx context.Context, newObj T) (T, error)
}

type ReadableWithRelations[T any, R any] interface {
	GetOneWithRelation(ctx context.Context, filters sqlhelper.FilterSlice, complexFilters []qm.QueryMod, loaders sqlhelper.LoaderSlice) (R, error)
}

type Countable[T any, R any] interface {
	Countable(ctx context.Context, filters T) ([]R, error)
}
