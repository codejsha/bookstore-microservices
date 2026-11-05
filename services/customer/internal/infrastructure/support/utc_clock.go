package support

import (
	"time"

	"github.com/codejsha/shared-library-go/pkg/database"
	"gorm.io/gorm"
)

func ConfigureUtcClock(dataSource *database.DataSource) {
	useUtcClock(dataSource.DB())
}

func useUtcClock(db *gorm.DB) {
	db.Config.NowFunc = func() time.Time { return time.Now().UTC() }
}
