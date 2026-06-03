package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"jiyixia/member-system/internal/config"
	"jiyixia/member-system/internal/store"
)

func TestHealthResponse(t *testing.T) {
	handler := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if res.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected cors header")
	}
}

func TestAdminMembersRequiresToken(t *testing.T) {
	handler := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/members", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", res.Code)
	}
}

func TestAdminLoginAndListMembers(t *testing.T) {
	handler := newTestRouter()
	login := `{"username":"admin","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", strings.NewReader(login))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "test-token") {
		t.Fatalf("expected login token, status=%d body=%s", res.Code, res.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/admin/members", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), "Alice") {
		t.Fatalf("expected member list, status=%d body=%s", res.Code, res.Body.String())
	}
}

func newTestRouter() http.Handler {
	cfg := config.Config{AdminUsername: "admin", AdminPassword: "secret"}
	return NewRouter(fakeStore{}, fakeSessions{}, cfg)
}

type fakeSessions struct{}

func (fakeSessions) Create(context.Context, string) (string, error) { return "test-token", nil }
func (fakeSessions) Validate(_ context.Context, token string) (string, bool) {
	return "admin", token == "test-token"
}

type fakeStore struct{}

func (fakeStore) ListPlans(context.Context) ([]store.Plan, error) {
	return []store.Plan{{ID: 1, Name: "月度会员", DurationDay: 30}}, nil
}

func (fakeStore) CreatePlan(_ context.Context, input store.CreatePlanInput) (store.Plan, error) {
	return store.Plan{ID: 2, Name: input.Name, PriceCents: input.PriceCents, DurationDay: input.DurationDay}, nil
}

func (fakeStore) CreateMember(_ context.Context, input store.CreateMemberInput) (store.Member, error) {
	return store.Member{ID: 9, Name: input.Name, Phone: input.Phone, Email: input.Email, PlanID: input.PlanID, Status: "pending"}, nil
}

func (fakeStore) GetMember(context.Context, int64) (store.Member, error) {
	return testMember(), nil
}

func (fakeStore) ListMembers(context.Context, string) ([]store.Member, error) {
	return []store.Member{testMember()}, nil
}

func (fakeStore) UpdateMemberStatus(_ context.Context, id int64, status string) (store.Member, error) {
	member := testMember()
	member.ID = id
	member.Status = status
	return member, nil
}

func testMember() store.Member {
	now := time.Date(2026, 6, 3, 8, 0, 0, 0, time.UTC)
	return store.Member{ID: 1, Name: "Alice", Phone: "13800000000", Email: "alice@example.com", PlanID: 1, PlanName: "月度会员", Status: "active", CreatedAt: now, UpdatedAt: now}
}
