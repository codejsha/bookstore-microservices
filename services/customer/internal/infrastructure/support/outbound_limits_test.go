package support

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/codejsha/shared-library-go/pkg/rest/client"
)

type unreachableConnector struct{}

func (unreachableConnector) Connect(context.Context) (driver.Conn, error) {
	return nil, errors.New("no database in unit tests")
}

func (unreachableConnector) Driver() driver.Driver { return nil }

func TestApplyConnectionPoolLimits_freshPool_capsOpenConnections(t *testing.T) {
	sqlDB := sql.OpenDB(unreachableConnector{})
	t.Cleanup(func() { _ = sqlDB.Close() })

	applyConnectionPoolLimits(sqlDB)

	if got := sqlDB.Stats().MaxOpenConnections; got != maxOpenConns {
		t.Fatalf("MaxOpenConnections = %d, want %d", got, maxOpenConns)
	}
}

func TestWithOutboundTimeout_sharedRestyClient_setsHTTPClientTimeout(t *testing.T) {
	restyClient := WithOutboundTimeout(client.NewRestyClient())

	if got := restyClient.Client.GetClient().Timeout; got != outboundHTTPTimeout {
		t.Fatalf("resty timeout = %v, want %v", got, outboundHTTPTimeout)
	}
}
