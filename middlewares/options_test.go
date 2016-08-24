package middlewares

import (
	"testing"

	"golang.org/x/net/context"

	"tower.pro/renderer/options"
	"tower.pro/renderer/template"
)

func TestConstructOptions(t *testing.T) {

	m := &Middleware{
		Context: template.Context{"query": template.Context{"owner_id": "account.id"}},
	}

	opts, err := m.ConstructOptions()
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "account.id", "test")

	query, err := opts.Map(ctx, &options.Option{Name: "query"})
	if err != nil {
		t.Error(err)
	}

	if len(query) != 1 {
		t.Fatal("Expected query to be of length 1")
	}

	if acc, _ := query["owner_id"]; acc != "test" {
		t.Fatalf("Unexpected value: %q", acc)
	}

}
