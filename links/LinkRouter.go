package links

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dataprism/dataprism-commons/utils"
	"io/ioutil"
	"encoding/json"
)

type LinkRouter struct {
	manager *LinkManager
}

func NewRouter(linkManager *LinkManager) (*LinkRouter) {
	return &LinkRouter{manager:linkManager}
}

func (router *LinkRouter) ListLinks(w http.ResponseWriter, r *http.Request) {
	resp, err := router.manager.ListLinks(r.Context())
	utils.HandleResponse(w, resp, err)
}

func (router *LinkRouter) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := router.manager.GetLink(r.Context(), id)
	utils.HandleResponse(w, resp, err)
}

func (router *LinkRouter) SetLink(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var entity Link
	err = json.Unmarshal(body, &entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := router.manager.SetLink(r.Context(), &entity)
	utils.HandleResponse(w, response, err)
}
func (router *LinkRouter) RemoveLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := router.manager.RemoveLink(r.Context(), id)
	utils.HandleStatus(w, 200, "Deleted", err)
}