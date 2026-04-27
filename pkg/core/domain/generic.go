package domain

type Distinct struct {
	Id    int64 `json:"id" db:"id" boil:"id"`
	Value any   `json:"value" db:"value" boil:"value"`
}

type DistinctSlice []*Distinct
