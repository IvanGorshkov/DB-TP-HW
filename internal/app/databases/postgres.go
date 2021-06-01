package databases

import (
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
)

type Postgres struct {
	postgresDatabase *pgx.ConnPool
}

func NewPostgres(dataSourceName string) (*Postgres, error) {
	pgxConnConfig, err := pgx.ParseConnectionString(dataSourceName)
	if err != nil {
		fmt.Println(err.Error())
	}
	pgxConnConfig.PreferSimpleProtocol = true

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     pgxConnConfig,
		MaxConnections: 200,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &Postgres{
		postgresDatabase: pool,
	}, nil
}

func (p *Postgres) GetDatabase() *pgx.ConnPool {
	return p.postgresDatabase
}

func (p *Postgres) Close() {
	p.postgresDatabase.Close()
}

func GetPostgresConfig() string {
//	return "host=localhost port=5432 user=postgres password=postgres dbname=formdb sslmode=disable"

	return "host=localhost port=5432 user=docker password=docker dbname=docker sslmode=disable"
}