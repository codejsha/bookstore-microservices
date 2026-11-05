package support

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type utcClockRow struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func TestUseUtcClock_NonUtcLocalZone_AutoTimestampsInUtc(t *testing.T) {
	local := time.Local
	time.Local = time.FixedZone("KST", 9*60*60)
	t.Cleanup(func() { time.Local = local })
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	useUtcClock(db)

	for name, vars := range map[string][]any{
		"create": db.Create(&utcClockRow{Name: "a"}).Statement.Vars,
		"save":   db.Save(&utcClockRow{Id: 1, Name: "a"}).Statement.Vars,
		"update": db.Model(&utcClockRow{Id: 1}).Updates(map[string]any{"name": "b"}).Statement.Vars,
	} {
		stamped := 0
		for _, v := range vars {
			if ts, ok := v.(time.Time); ok && !ts.IsZero() {
				stamped++
				if ts.Location() != time.UTC {
					t.Errorf("%s: auto timestamp must be UTC, got %v", name, ts)
				}
			}
		}
		if stamped == 0 {
			t.Errorf("%s: no auto timestamp bound", name)
		}
	}
}
