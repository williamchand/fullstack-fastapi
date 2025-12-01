package database

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper functions for pgtype conversion
func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func fromPgText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func fromPgTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// Helper functions for pgtype conversion
func toPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func toPgJSON(m map[string]any) []byte {
	if m == nil {
		return []byte{}
	}
	b, _ := json.Marshal(m)

	return b
}

func fromPgJSON(j []byte) map[string]any {
	var m map[string]any
	_ = json.Unmarshal(j, &m)
	return m
}

func toPgNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	n := &pgtype.Numeric{Valid: true}
	n.Int = new(big.Int)
	n.Exp = -2
	n.Int.SetInt64(int64(*f * 100))
	return *n
}

func toPgUUIDPtr(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: [16]byte(*u), // convert value to fixed array
		Valid: true,
	}
}

func toPgNumericFromFloat64(f float64) pgtype.Numeric {
	var num pgtype.Numeric
	if err := num.Scan(f); err != nil {
		return pgtype.Numeric{}
	}
	return num
}

func fromPgNumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	amount, _ := n.Int.Float64()
	return amount
}
