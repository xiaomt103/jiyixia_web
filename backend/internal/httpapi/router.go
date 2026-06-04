package httpapi

import (
	"net/http"
	"strings"

	"jiyixia/member-system/internal/config"
)

type API struct {
	db       MemberStore
	sessions SessionManager
	cfg      config.Config
}

func NewRouter(db MemberStore, sessions SessionManager, cfg config.Config) http.Handler {
	api := &API{db: db, sessions: sessions, cfg: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", api.health)
	mux.HandleFunc("GET /api/plans", api.listPlans)
	mux.HandleFunc("POST /api/members", api.createMember)
	mux.HandleFunc("GET /api/members/{id}", api.getMember)
	mux.HandleFunc("POST /api/admin/login", api.login)
	mux.Handle("/api/v1/admin/", api.requireAdmin(http.HandlerFunc(api.adminRouter)))
	return cors(mux)
}

func (a *API) adminRouter(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/admin/members":
		a.listMembers(w, r)
	case r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/status"):
		a.updateMemberStatus(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/admin/plans":
		a.listPlans(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/admin/plans":
		a.createPlan(w, r)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (a *API) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if _, ok := a.sessions.Validate(r.Context(), token); !ok {
			writeError(w, http.StatusUnauthorized, "admin token required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
