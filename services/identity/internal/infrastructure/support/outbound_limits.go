package support

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/rest/client"
)

const (
	maxOpenConns        = 20
	maxIdleConns        = 10
	connMaxIdleTime     = 5 * time.Minute
	outboundHTTPTimeout = 10 * time.Second
)

func ConfigureConnectionPool(dataSource *database.DataSource) error {
	sqlDB, err := dataSource.DB().DB()
	if err != nil {
		return fmt.Errorf("configure connection pool: %w", err)
	}
	applyConnectionPoolLimits(sqlDB)
	return nil
}

func applyConnectionPoolLimits(sqlDB *sql.DB) {
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)
}

func WithOutboundTimeout(restyClient *client.RestyClient) *client.RestyClient {
	restyClient.Client.SetTimeout(outboundHTTPTimeout)
	return restyClient
}
