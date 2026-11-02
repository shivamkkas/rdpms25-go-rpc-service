package stub

import (
	"context"
)

type Trainer interface {
	Train(ctx context.Context, req any) (any, error)
	Retrain(ctx context.Context, req any) (any, error)
}

// Predictor handles prediction (classification / anomaly)
type Predictor interface {
	Predict(ctx context.Context, req any) (any, error)
}

// Forecaster handles time-series forecasting
type Forecaster interface {
	Forecast(ctx context.Context, req any) (any, error)
}

// Preprocessor handles preprocessing (mainly forecasting)
type Preprocessor interface {
	Preprocess(ctx context.Context, req any) (any, error)
}

// ModelManager handles shared model operations
type ModelManager interface {
	GetModelDetails(ctx context.Context, req any) (any, error)
	ListModels(ctx context.Context, req any) (any, error)
	DeleteModel(ctx context.Context, req any) (any, error)
	HealthCheck(ctx context.Context, req any) (any, error)
}

// ForecastingModule supports full forecasting pipeline
type ForecastingModule interface {
	Trainer
	Forecaster
	Preprocessor
	ModelManager
}

// ClassificationModule supports classification pipeline
type ClassificationModule interface {
	Trainer
	Predictor
	ModelManager
}

// AnomalyModule supports anomaly detection pipeline
type AnomalyModule interface {
	Trainer
	Predictor
	ModelManager
}

// ModuleRegistry gives access to all modules
type ModuleRegistry interface {
	GetForecasting() ForecastingModule
	GetClassification() ClassificationModule
	GetAnomaly() AnomalyModule
}
