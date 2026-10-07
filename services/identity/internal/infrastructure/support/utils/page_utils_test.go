package utils

import (
	"math"
	"testing"

	"github.com/codejsha/shared-library-go/pkg/pagination"
)

func i32p(v int32) *int32   { return &v }
func strp(v string) *string { return &v }

func TestRowOffsetPage(t *testing.T) {
	cases := []struct {
		name       string
		in         pagination.PageOption
		wantSize   int32
		wantOffset int32
	}{
		{"firstPage_offsetZero", pagination.NewPageOption(i32p(20), i32p(0), nil), 20, 0},
		{"pageUnset_offsetZero", pagination.NewPageOption(i32p(20), nil, nil), 20, 0},
		{"secondPage_offsetIsOneSize", pagination.NewPageOption(i32p(20), i32p(1), nil), 20, 20},
		{"thirdPage_offsetIsTwoSizes", pagination.NewPageOption(i32p(5), i32p(3), nil), 5, 15},
		{"negativePage_offsetZero", pagination.NewPageOption(i32p(20), i32p(-3), nil), 20, 0},
		{"sizeUnset_defaultSize", pagination.NewPageOption(nil, i32p(4), nil), 20, 80},
		{"sizeZero_defaultSize", pagination.NewPageOption(i32p(0), i32p(1), nil), 20, 20},
		{"sizeAboveMax_cappedSize", pagination.NewPageOption(i32p(500), i32p(2), nil), 100, 200},
		{"offsetOverflow_cappedAtMaxInt32", pagination.NewPageOption(i32p(100), i32p(math.MaxInt32), nil), 100, math.MaxInt32},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RowOffsetPage(c.in)
			if got.GetSize() != c.wantSize {
				t.Errorf("size = %d, want %d", got.GetSize(), c.wantSize)
			}
			if got.GetPage() != c.wantOffset {
				t.Errorf("offset = %d, want %d", got.GetPage(), c.wantOffset)
			}
		})
	}
}

func TestRowOffsetPage_sortGiven_keepsSort(t *testing.T) {
	in := pagination.NewPageOption(i32p(10), i32p(2), strp("email,asc;createdAt,desc"))
	got := RowOffsetPage(in)
	if got.GetSort() != "email,asc;createdAt,desc" {
		t.Errorf("sort = %q, want %q", got.GetSort(), "email,asc;createdAt,desc")
	}
	if in.GetPage() != 2 {
		t.Errorf("input page = %d, want 2 left untouched", in.GetPage())
	}
}

func TestWithDefaultSort(t *testing.T) {
	cases := []struct {
		name string
		in   pagination.PageOption
		want string
	}{
		{"sortUnset_defaultSort", pagination.NewPageOption(i32p(20), i32p(0), nil), "createdat,desc;uid,asc"},
		{"sortEmpty_defaultSort", pagination.NewPageOption(i32p(20), i32p(0), strp("")), "createdat,desc;uid,asc"},
		{"sortGiven_keepsSort", pagination.NewPageOption(i32p(20), i32p(0), strp("email,asc")), "email,asc"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := WithDefaultSort(c.in).GetSort(); got != c.want {
				t.Errorf("sort = %q, want %q", got, c.want)
			}
		})
	}
}
