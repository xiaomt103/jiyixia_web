package httpapi

import (
	"net/http"

	"jiyixia/member-system/internal/store"
)

func (a *API) listOperations(w http.ResponseWriter, r *http.Request) {
	records, err := a.db.ListOperations(r.Context(), r.URL.Query().Get("module"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (a *API) createOperation(w http.ResponseWriter, r *http.Request) {
	var input store.CreateOperationInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	record, err := a.db.CreateOperation(r.Context(), input)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (a *API) updateOperationStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r, "/api/v1/admin/operations/")
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
	record, err := a.db.UpdateOperationStatus(r.Context(), id, input.Status)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}
