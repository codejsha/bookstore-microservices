package di

import (
	"testing"

	pkgconfig "github.com/codejsha/shared-library-go/pkg/config"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure"
)

func TestModule_suppliedConfig_validatesGraph(t *testing.T) {
	err := fx.ValidateApp(
		fx.Supply(&pkgconfig.PreConfig{}, &pkgconfig.Metadata{}),
		Module,
		fx.Invoke(func(*infrastructure.Infra) {}),
	)
	if err != nil {
		t.Fatalf("fx.ValidateApp = %v, want nil", err)
	}
}
