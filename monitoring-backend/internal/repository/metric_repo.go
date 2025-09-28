package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Metric struct {
	ID         int64     `db:"id"`
	Service    string    `db:"service"`
	Name       string    `db:"name"`
	Value      float64   `db:"value"`
	LabelsRaw  []byte    `db:"labels"`
	ObservedAt time.Time `db:"observed_at"`
	CreatedAt  time.Time `db:"created_at"`
}

type Repo struct{ DB *sqlx.DB }

func New(dataSource string) (*Repo, error) {
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &Repo{DB: db}, nil
}

func migrate(db *sqlx.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS metrics (
  id BIGSERIAL PRIMARY KEY,
  service TEXT NOT NULL,
  name TEXT NOT NULL,
  value DOUBLE PRECISION NOT NULL,
  labels JSONB NOT NULL DEFAULT '{}'::jsonb,
  observed_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_metrics_service ON metrics(service);
CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics(name);
CREATE INDEX IF NOT EXISTS idx_metrics_observed_at ON metrics(observed_at);
CREATE INDEX IF NOT EXISTS idx_metrics_labels_gin ON metrics USING GIN(labels);
`
	_, err := db.Exec(schema)
	return err
}

func (r *Repo) InsertMetric(ctx context.Context, m map[string]any) error {
	labels := map[string]string{}
	if v, ok := m["labels"].(map[string]any); ok {
		for k, vv := range v {
			labels[k] = toString(vv)
		}
	}
	labelsJSON, _ := json.Marshal(labels)
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO metrics(service,name,value,labels,observed_at) VALUES($1,$2,$3,$4,$5)`,
		m["service"], m["name"], m["value"], labelsJSON, m["observed_at"],
	)
	return err
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%g", t)
	case int64:
		return fmt.Sprintf("%d", t)
	default:
		return ""
	}
}

func (r *Repo) QueryMetrics(ctx context.Context, service, name string, from, to time.Time, limit int) ([]Metric, error) {
	q := `SELECT id,service,name,value,labels,observed_at,created_at
	      FROM metrics
	      WHERE observed_at BETWEEN $1 AND $2
	      AND ($3 = '' OR service = $3)
	      AND ($4 = '' OR name = $4)
	      ORDER BY observed_at DESC
	      LIMIT $5`
	rows, err := r.DB.QueryxContext(ctx, q, from, to, service, name, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Metric
	for rows.Next() {
		var m Metric
		if err := rows.StructScan(&m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (m *Metric) Labels() (map[string]string, error) {
	if len(m.LabelsRaw) == 0 {
		return map[string]string{}, nil
	}
	var res map[string]string
	return res, json.Unmarshal(m.LabelsRaw, &res)
}

func (r *Repo) Ping(ctx context.Context) error { return r.DB.DB.PingContext(ctx) }
