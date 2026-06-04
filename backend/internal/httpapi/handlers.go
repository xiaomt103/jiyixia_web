package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"jiyixia/member-system/internal/store"
)

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Username != a.cfg.AdminUsername || input.Password != a.cfg.AdminPassword {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := a.sessions.Create(r.Context(), input.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create session failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (a *API) listPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := a.db.ListPlans(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list plans failed")
		return
	}
	writeJSON(w, http.StatusOK, plans)
}

func (a *API) createPlan(w http.ResponseWriter, r *http.Request) {
	var input store.CreatePlanInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	plan, err := a.db.CreatePlan(r.Context(), input)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}

func (a *API) createMember(w http.ResponseWriter, r *http.Request) {
	var input store.CreateMemberInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	member, err := a.db.CreateMember(r.Context(), input)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, member)
}

func (a *API) getMember(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "/api/members/")
	if !ok {
		return
	}
	member, err := a.db.GetMember(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func (a *API) listMembers(w http.ResponseWriter, r *http.Request) {
	members, err := a.db.ListMembers(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, members)
}

func (a *API) updateMemberStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "/api/v1/admin/members/")
	if !ok {
		return
	}
	var input struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	member, err := a.db.UpdateMemberStatus(r.Context(), id, input.Status)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func pathID(w http.ResponseWriter, r *http.Request, prefix string) (int64, bool) {
	value := strings.TrimPrefix(r.URL.Path, prefix)
	value = strings.TrimSuffix(value, "/status")
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
