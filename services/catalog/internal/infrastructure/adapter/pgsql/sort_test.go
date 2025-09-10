package pgsql

import (
	"testing"

	"gorm.io/gen/field"
)

func TestParseSortTerm(t *testing.T) {
	cases := []struct {
		in       string
		wantName string
		wantDesc bool
	}{
		{"title", "title", false},
		{"title,desc", "title", true},
		{"title:desc", "title", true},
		{"title,asc", "title", false},
		{"  created_at ,  DESC ", "created_at", true},
		{"name:ASC", "name", false},
		{"name:whatever", "name", false},
	}
	for _, c := range cases {
		name, desc := parseSortTerm(c.in)
		if name != c.wantName || desc != c.wantDesc {
			t.Errorf("parseSortTerm(%q) = (%q, %v), want (%q, %v)", c.in, name, desc, c.wantName, c.wantDesc)
		}
	}
}

func TestBuildOrderExprs(t *testing.T) {
	title := field.NewString("work", "title")
	created := field.NewTime("work", "created_at")
	id := field.NewInt64("work", "id")
	whitelist := map[string]field.OrderExpr{
		"title":      title,
		"created_at": created,
	}

	t.Run("no sort falls back to id tiebreaker only", func(t *testing.T) {
		got := buildOrderExprs("", whitelist, id)
		if len(got) != 1 {
			t.Fatalf("len = %d, want 1 (id tiebreaker)", len(got))
		}
	})

	t.Run("unknown field is dropped, still id tiebreaker", func(t *testing.T) {
		got := buildOrderExprs("bogus,desc", whitelist, id)
		if len(got) != 1 {
			t.Fatalf("len = %d, want 1 (unknown dropped -> id only)", len(got))
		}
	})

	t.Run("whitelisted field plus id tiebreaker", func(t *testing.T) {
		got := buildOrderExprs("title,desc", whitelist, id)
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2 (title + id)", len(got))
		}
	})

	t.Run("multiple terms preserve order and append id last", func(t *testing.T) {
		got := buildOrderExprs("title;created_at,desc", whitelist, id)
		if len(got) != 3 {
			t.Fatalf("len = %d, want 3 (title, created_at, id)", len(got))
		}
	})
}
