package testData

import (
	"testing"

	"github.com/pedramktb/schwarzit-probearbeit/pkg/postgres"
)

func Local_Test_Setup_DB(t *testing.T) {
	db := postgres.Test_Create_DB("localhost", "5432", "core")
	MigrateTestData(db)
}
