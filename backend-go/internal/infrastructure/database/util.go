package database

import (
    "encoding/json"
    "time"

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

func toPgJSON(m map[string]interface{}) pgtype.JSONB {
    if m == nil {
        return pgtype.JSONB{Valid: false}
    }
    b, _ := json.Marshal(m)
    return pgtype.JSONB{Bytes: b, Valid: true}
}

func fromPgJSON(j pgtype.JSONB) map[string]interface{} {
    if !j.Valid {
        return nil
    }
    var m map[string]interface{}
    _ = json.Unmarshal(j.Bytes, &m)
    return m
}
