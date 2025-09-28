package models

import "time"

type Metric struct {
	ID         int64             `db:"id" json:"id"`
	Service    string            `db:"service" json:"service"`
	Name       string            `db:"name" json:"name"`
	Value      float64           `db:"value" json:"value"`
	Labels     map[string]string `db:"labels" json:"labels"`
	ObservedAt time.Time         `db:"observed_at" json:"observed_at"`
	CreatedAt  time.Time         `db:"created_at" json:"created_at"`
}
