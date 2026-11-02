package forecast

import (
	"context"

	"github.com/shivamkkas/rdpms25-go-rpc-service/pkg/core/serv"
)

type forecaster struct {
	handler *serv.HistoryView
}

func NewForecaster(handler *serv.HistoryView) *forecaster {
	return &forecaster{
		handler: handler,
	}
}
func (f *forecaster) Preprocess(ctx context.Context, req any) (any, error) {

	return nil, nil

}
func (f *forecaster) Train(ctx context.Context, req any) (any, error) {
	return nil, nil
}
func (f *forecaster) Retrain(ctx context.Context, req any) (any, error) {
	return nil, nil
}
func (f *forecaster) Forecast(ctx context.Context, req any) (any, error) {
	return nil, nil
}
