package griz

import (
	"context"
	"testing"

	"github.com/estenssoros/dasorm"
)

func TestReadSQL(t *testing.T) {
	db, err := dasorm.ConnectDB("dev-local")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.WithValue(context.Background(), "db", db.DB)
	query := `
	SELECT created_at
		, rv_month
		, rv_year
		, source
		, teu
		, spec
		, per_day
	FROM rv_data
	LIMIT 200
	`
	df, err := ReadSQL(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := 200, df.Len(); want != have {
		t.Errorf("have: %d, want: %d", have, want)
	}
	df.Head(10)
}
