package postgres

import (
	"context"
	"fmt"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pedramktb/schwarzit-probearbeit/migration"
	v1Migration "github.com/pedramktb/schwarzit-probearbeit/migration/v1"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/testcontainers/testcontainers-go"
	postgresC "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_Create_Container() (container testcontainers.Container, ip, port string) {
	ctx := context.Background()

	postgresContainer, err := postgresC.Run(ctx, "postgres:latest",
		postgresC.WithUsername("postgres"),
		postgresC.WithPassword("example"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Minute)),
	)

	if err != nil {
		panic(err)
	}

	ip, err = postgresContainer.Host(ctx)
	if err != nil {
		panic(err)
	}
	natPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	return postgresContainer, ip, natPort.Port()
}

func Test_Create_DB(ip, port, dbName string) *gorm.DB {
	db := create(fmt.Sprintf("postgres://postgres:example@%s:%s", ip, port))

	err := db.Exec(fmt.Sprintf("CREATE DATABASE %q", dbName)).Error
	if err != nil {
		panic(err)
	}

	Close(db)

	db = create(fmt.Sprintf("postgres://postgres:example@%s:%s/%s", ip, port, dbName))

	fx.New(
		fx.Provide(func() *gorm.DB { return db }),
		v1Migration.FXV1TestMigrationProvide,
		fx.Invoke(func(m migration.Migrator) {
			m.Migrate(context.Background())
		}),
	)

	return db
}

func Test_Drop_DB(db *gorm.DB, ip, port, dbName string) {
	Close(db)

	db = create(fmt.Sprintf("postgres://postgres:example@%s:%s", ip, port))
	db.Exec(fmt.Sprintf("DROP DATABASE %q WITH (force)", dbName))
	Close(db)
}
