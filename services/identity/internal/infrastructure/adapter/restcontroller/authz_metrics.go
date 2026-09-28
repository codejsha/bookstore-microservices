package restcontroller

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/codejsha/bookstore-microservices/identity/internal/domain/constant"
)

const (
	authzOutcomeAllow = "allow"
	authzOutcomeDeny  = "deny"

	authzReasonMissingBearer = "missing_bearer"
)

type authzMetrics struct {
	decisions metric.Int64Counter
}

func newAuthzMetrics(provider metric.MeterProvider) *authzMetrics {
	counter, err := provider.Meter(string(constant.MeterNameAuthz)).Int64Counter(
		constant.MetricAuthzDecisions,
		metric.WithDescription("Authorization decisions returned by /internal/authz, by outcome and reason"),
		metric.WithUnit("{decision}"),
	)
	if err != nil {
		logrus.WithError(err).Warn("authz decisions counter unavailable, recording disabled")
		counter = noop.Int64Counter{}
	}
	return &authzMetrics{decisions: counter}
}

func (m *authzMetrics) allow(ctx context.Context, reason string) {
	m.record(ctx, authzOutcomeAllow, reason)
}

func (m *authzMetrics) deny(ctx context.Context, reason string) {
	m.record(ctx, authzOutcomeDeny, reason)
}

func (m *authzMetrics) record(ctx context.Context, outcome, reason string) {
	m.decisions.Add(ctx, 1, metric.WithAttributes(
		attribute.String("outcome", outcome),
		attribute.String("reason", reason),
	))
}
