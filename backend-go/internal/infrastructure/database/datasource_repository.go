package database

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/database/dbgen"
)

type dataSourceRepository struct {
    queries *dbgen.Queries
    db      repositories.ConnectionPool
}

func NewDataSourceRepository(queries *dbgen.Queries, db repositories.ConnectionPool) repositories.DataSourceRepository {
    return &dataSourceRepository{queries: queries, db: db}
}

func (r *dataSourceRepository) WithTx(tx pgx.Tx) repositories.DataSourceRepository {
    return &dataSourceRepository{queries: r.queries.WithTx(tx), db: r.db}
}

func (r *dataSourceRepository) Create(ctx context.Context, ds *entities.DataSource) (*entities.DataSource, error) {
    params := dbgen.CreateDataSourceParams{
        UserID:       ds.UserID,
        Name:         ds.Name,
        Type:         ds.Type,
        Host:         ds.Host,
        Port:         int32(ds.Port),
        DatabaseName: ds.DatabaseName,
        Username:     ds.Username,
        PasswordEnc:  ds.PasswordEnc,
        Options:      toPgJSON(ds.Options),
    }
    out, err := r.queries.CreateDataSource(ctx, params)
    if err != nil {
        return nil, err
    }
    return r.toEntity(&out), nil
}

func (r *dataSourceRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.DataSource, error) {
    out, err := r.queries.GetDataSourceByID(ctx, dbgen.GetDataSourceByIDParams{ID: id, UserID: userID})
    if err != nil {
        return nil, err
    }
    return r.toEntity(&out), nil
}

func (r *dataSourceRepository) toEntity(v *dbgen.DataSource) *entities.DataSource {
    return &entities.DataSource{
        ID:           v.ID,
        UserID:       v.UserID,
        Name:         v.Name,
        Type:         v.Type,
        Host:         v.Host,
        Port:         int32(v.Port),
        DatabaseName: v.DatabaseName,
        Username:     v.Username,
        PasswordEnc:  v.PasswordEnc,
        Options:      fromPgJSON(v.Options),
        CreatedAt:    v.CreatedAt.Time,
        UpdatedAt:    v.UpdatedAt.Time,
    }
}
