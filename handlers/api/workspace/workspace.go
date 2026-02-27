package workspace

import (
	"encoding/json"
	"excalidraw-complete/core"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type saveRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func HandleList(store core.WorkspaceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := store.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, files)
	}
}

func HandleGet(store core.WorkspaceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		file, err := store.Find(r.Context(), id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		render.JSON(w, r, file)
	}
}

func HandleCreate(store core.WorkspaceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		var req saveRequest
		if err := json.Unmarshal(body, &req); err != nil || req.Name == "" {
			http.Error(w, "invalid body: name and data are required", http.StatusBadRequest)
			return
		}
		updatedBy := r.Header.Get("X-Forwarded-Email")
		id, err := store.Save(r.Context(), "", req.Name, []byte(req.Data), updatedBy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file, err := store.Find(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, file.WorkspaceFileMeta)
	}
}

func HandleUpdate(store core.WorkspaceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		var req saveRequest
		if err := json.Unmarshal(body, &req); err != nil || req.Name == "" {
			http.Error(w, "invalid body: name and data are required", http.StatusBadRequest)
			return
		}
		updatedBy := r.Header.Get("X-Forwarded-Email")
		if _, err := store.Save(r.Context(), id, req.Name, []byte(req.Data), updatedBy); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file, err := store.Find(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, file.WorkspaceFileMeta)
	}
}

func HandleDelete(store core.WorkspaceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := store.Delete(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
