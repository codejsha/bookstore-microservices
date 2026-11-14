package keycloak

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/identity/internal/config"
)

const (
	stubRealm        = "test"
	stubClientID     = "identity-client"
	stubClientSecret = "s3cr3t-client"
	stubAdminUser    = "realm-admin"
	stubAdminPass    = "s3cr3t-password"
	stubAccessToken  = "stub-access-token"
)

type stubRequest struct {
	Method      string
	Path        string
	Query       url.Values
	Auth        string
	ContentType string
	Body        string
}

func (r stubRequest) form(t *testing.T) url.Values {
	t.Helper()
	values, err := url.ParseQuery(r.Body)
	if err != nil {
		t.Fatalf("parse form body %q: %v", r.Body, err)
	}
	return values
}

type stubReply struct {
	Status int
	Body   string
}

type keycloakStub struct {
	URL string

	mu       sync.Mutex
	replies  map[string]stubReply
	respond  func(stubRequest) *stubReply
	requests []stubRequest
}

func newKeycloakStub(t *testing.T) *keycloakStub {
	t.Helper()
	stub := &keycloakStub{replies: map[string]stubReply{}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		stub.mu.Lock()
		stub.requests = append(stub.requests, stubRequest{
			Method:      r.Method,
			Path:        r.URL.Path,
			Query:       r.URL.Query(),
			Auth:        r.Header.Get("Authorization"),
			ContentType: r.Header.Get("Content-Type"),
			Body:        string(body),
		})
		reply, ok := stub.replies[r.Method+" "+r.URL.Path]
		if stub.respond != nil {
			if scripted := stub.respond(stub.requests[len(stub.requests)-1]); scripted != nil {
				reply, ok = *scripted, true
			}
		}
		stub.mu.Unlock()
		if !ok {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(reply.Status)
		_, _ = w.Write([]byte(reply.Body))
	}))
	t.Cleanup(srv.Close)
	stub.URL = srv.URL
	return stub
}

func (s *keycloakStub) on(method, path string, status int, body string) *keycloakStub {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replies[method+" "+path] = stubReply{Status: status, Body: body}
	return s
}

func (s *keycloakStub) recorded() []stubRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]stubRequest(nil), s.requests...)
}

func (s *keycloakStub) recordedFor(method, path string) []stubRequest {
	var out []stubRequest
	for _, r := range s.recorded() {
		if r.Method == method && r.Path == path {
			out = append(out, r)
		}
	}
	return out
}

func (s *keycloakStub) only(t *testing.T) stubRequest {
	t.Helper()
	requests := s.recorded()
	if len(requests) != 1 {
		t.Fatalf("requests = %+v, want exactly 1", requests)
	}
	return requests[0]
}

func stubConfig(baseURL string) *config.Config {
	return &config.Config{
		App: &sharedconfig.AppConfig{},
		Keycloak: &sharedconfig.KeycloakConfig{
			Url:                baseURL,
			Realm:              stubRealm,
			ClientId:           stubClientID,
			ClientSecret:       stubClientSecret,
			RealmAdminUsername: stubAdminUser,
			RealmAdminPassword: stubAdminPass,
		},
	}
}

func unreachableKeycloakURL(t *testing.T) string {
	t.Helper()
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()
	return srv.URL
}

func newStubRestyClient() *resty.Client {
	return resty.New().SetTimeout(2 * time.Second)
}

func freshTokenHelper() *AdminTokenHelper {
	return &AdminTokenHelper{accessToken: stubAccessToken, expiresAt: time.Now().Add(time.Hour)}
}
