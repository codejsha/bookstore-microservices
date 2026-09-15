package di

import (
	"testing"

	pkgconfig "github.com/codejsha/shared-library-go/pkg/config"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure"
)

func TestModule_fullGraph_validates(t *testing.T) {
	err := fx.ValidateApp(
		fx.Supply(&pkgconfig.PreConfig{}, &pkgconfig.Metadata{}),
		Module,
		fx.Invoke(func(*infrastructure.Infra) {}),
	)
	if err != nil {
		t.Fatalf("fx.ValidateApp: %v", err)
	}
}
