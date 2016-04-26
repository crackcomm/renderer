package middlewares

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/crackcomm/renderer/options"
)

func TestConstructOptions(t *testing.T) {

	m := &Middleware{
		Context: options.Options{"query": map[string]interface{}{"account_id": "account.id"}},
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

	if acc, _ := query["account_id"]; acc != "test" {
		t.Fatalf("Unexpected value: %q", acc)
	}

}
