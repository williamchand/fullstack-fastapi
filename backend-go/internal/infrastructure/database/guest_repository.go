package database

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type guestRepository struct {
    queries *dbgen.Queries
    db      repositories.ConnectionPool
}

func NewGuestRepository(q *dbgen.Queries, db repositories.ConnectionPool) repositories.GuestRepository {
    return &guestRepository{queries: q, db: db}
}

func (r *guestRepository) WithTx(tx pgx.Tx) repositories.GuestRepository {
    return &guestRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *guestRepository) Add(ctx context.Context, g *entities.Guest) (*entities.Guest, error) {
    out, err := r.queries.AddGuest(ctx, dbgen.AddGuestParams{WeddingID: g.WeddingID, Name: g.Name, Contact: g.Contact, RsvpStatus: string(g.RSVPStatus), Message: toPgText(g.Message)})
    if err != nil { return nil, err }
    return r.toEntity(out), nil
}

func (r *guestRepository) Update(ctx context.Context, g *entities.Guest) (*entities.Guest, error) {
    out, err := r.queries.UpdateGuest(ctx, dbgen.UpdateGuestParams{ID: g.ID, Name: g.Name, Contact: g.Contact, RsvpStatus: string(g.RSVPStatus), Message: toPgText(g.Message)})
    if err != nil { return nil, err }
    return r.toEntity(out), nil
}

func (r *guestRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.queries.DeleteGuest(ctx, id)
}

func (r *guestRepository) ListByWedding(ctx context.Context, weddingID uuid.UUID) ([]*entities.Guest, error) {
    rows, err := r.queries.ListGuestsByWedding(ctx, weddingID)
    if err != nil { return nil, err }
    res := make([]*entities.Guest, 0, len(rows))
    for _, v := range rows { res = append(res, r.toEntity(v)) }
    return res, nil
}

func (r *guestRepository) toEntity(v dbgen.Guest) *entities.Guest {
    return &entities.Guest{
        ID:         v.ID,
        WeddingID:  v.WeddingID,
        Name:       v.Name,
        Contact:    v.Contact,
        RSVPStatus: entities.RSVPStatus(v.RsvpStatus),
        Message:    fromPgText(v.Message),
        CreatedAt:  v.CreatedAt.Time,
        DeletedAt:  fromPgTime(v.DeletedAt),
    }
}
