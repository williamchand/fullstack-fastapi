package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/crypto"
)

type DataSourceService struct {
	cfg       *config.Config
	dsRepo    repositories.DataSourceRepository
	aiRepo    repositories.AICredentialRepository
	txManager repositories.TransactionManager
}

func NewDataSourceService(cfg *config.Config, dsRepo repositories.DataSourceRepository, aiRepo repositories.AICredentialRepository, tx repositories.TransactionManager) *DataSourceService {
	return &DataSourceService{cfg: cfg, dsRepo: dsRepo, aiRepo: aiRepo, txManager: tx}
}

func (s *DataSourceService) GetByID(ctx context.Context, userID uuid.UUID, dataSourceID uuid.UUID) (*entities.DataSource, error) {
	return s.dsRepo.GetByID(ctx, dataSourceID, userID)
}

func (s *DataSourceService) CreateDataSource(ctx context.Context, userID uuid.UUID, name, typ, host string, port int32, dbname, username, password string, options map[string]any) (*entities.DataSource, error) {
	key := []byte(strings.TrimSpace(s.cfg.Security.CredentialEncryptionKey))
	enc, err := crypto.Encrypt(key, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	ds := &entities.DataSource{
		UserID:       userID,
		Name:         name,
		Type:         typ,
		Host:         host,
		Port:         port,
		DatabaseName: dbname,
		Username:     username,
		PasswordEnc:  enc,
		Options:      options,
	}
	return s.dsRepo.Create(ctx, ds)
}

func (s *DataSourceService) TestConnection(ctx context.Context, ds *entities.DataSource) (bool, string) {
	dsn := s.buildDSN(ds)
	driver := s.sqlDriver(ds.Type)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return false, err.Error()
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return false, err.Error()
	}
	return true, "ok"
}

func (s *DataSourceService) IntrospectSchema(ctx context.Context, ds *entities.DataSource) (string, string, error) {
	dsn := s.buildDSN(ds)
	driver := s.sqlDriver(ds.Type)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	switch strings.ToLower(ds.Type) {
	case "postgres":
		rows, err := db.QueryContext(ctx, `SELECT table_schema, table_name, column_name, data_type FROM information_schema.columns WHERE table_schema NOT IN ('pg_catalog','information_schema') ORDER BY table_schema, table_name, ordinal_position`)
		if err != nil {
			return "", "", err
		}
		defer rows.Close()
		var b strings.Builder
		for rows.Next() {
			var sch, t, col, dt string
			if err := rows.Scan(&sch, &t, &col, &dt); err != nil {
				return "", "", err
			}
			b.WriteString(fmt.Sprintf("%s.%s.%s:%s\n", sch, t, col, dt))
		}
		return b.String(), "schema introspection", nil
	case "mysql":
		rows, err := db.QueryContext(ctx, `SELECT table_schema, table_name, column_name, data_type FROM information_schema.columns WHERE table_schema NOT IN ('mysql','information_schema','performance_schema','sys') ORDER BY table_schema, table_name, ordinal_position`)
		if err != nil {
			return "", "", err
		}
		defer rows.Close()
		var b strings.Builder
		for rows.Next() {
			var sch, t, col, dt string
			if err := rows.Scan(&sch, &t, &col, &dt); err != nil {
				return "", "", err
			}
			b.WriteString(fmt.Sprintf("%s.%s.%s:%s\n", sch, t, col, dt))
		}
		return b.String(), "schema introspection", nil
	}
	return "", "", fmt.Errorf("unsupported type")
}

func (s *DataSourceService) SetAICredential(ctx context.Context, userID uuid.UUID, provider, apiKey string) error {
	key := []byte(strings.TrimSpace(s.cfg.Security.CredentialEncryptionKey))
	enc, err := crypto.Encrypt(key, []byte(apiKey))
	if err != nil {
		return err
	}
	_, err = s.aiRepo.Upsert(ctx, &entities.AICredential{UserID: userID, Provider: provider, APIKeyEnc: enc})
	return err
}

func (s *DataSourceService) AskQuestion(ctx context.Context, ds *entities.DataSource, provider, question string) (string, string, string, error) {
	_ = os.Getenv("AI_MODEL")
	return "not implemented", "", "{}", nil
}

func (s *DataSourceService) sqlDriver(typ string) string {
	switch strings.ToLower(typ) {
	case "postgres":
		return "pgx"
	case "mysql":
		return "mysql"
	}
	return ""
}

func (s *DataSourceService) buildDSN(ds *entities.DataSource) string {
	key := []byte(strings.TrimSpace(s.cfg.Security.CredentialEncryptionKey))
	pwdBytes, _ := crypto.Decrypt(key, ds.PasswordEnc)
	pwd := string(pwdBytes)
	switch strings.ToLower(ds.Type) {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", ds.Username, pwd, ds.Host, ds.Port, ds.DatabaseName)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", ds.Username, pwd, ds.Host, ds.Port, ds.DatabaseName)
	}
	return ""
}
