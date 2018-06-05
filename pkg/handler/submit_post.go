package handler

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/techmexdev/socialnet"
)

func (h *handler) SubmitPost(w http.ResponseWriter, r *http.Request) {
	un := r.Context().Value(ctxUsnKey).(string)
	var post socialnet.Post

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, `Request body must be in the form: 
		{"title": "I am the Walrus!",
		"body": "I am he as you are he as you are me, and we are all together ♫"}`, http.StatusBadRequest)
		return
	}

	post.CreatedAt = time.Now()
	post.Author = un

	createdPost, err := h.postSvc.Store.Create(post)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, createdPost, http.StatusCreated)
}
