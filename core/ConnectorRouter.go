package core

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dataprism/dataprism-commons/utils"
	"io/ioutil"
	"encoding/json"
)

type ConnectorRouter struct {
	manager *ConnectorManager
}

func NewConnectorRouter(connectorManager *ConnectorManager) (*ConnectorRouter) {
	return &ConnectorRouter{manager:connectorManager}
}

func (router *ConnectorRouter) ListConnectors(w http.ResponseWriter, r *http.Request) {
	resp, err := router.manager.ListConnectors(r.Context())
	utils.HandleResponse(w, resp, err)
}

func (router *ConnectorRouter) GetConnector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := router.manager.GetConnector(r.Context(), id)
	utils.HandleResponse(w, resp, err)
}

func (router *ConnectorRouter) SetConnector(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var entity Connector
	err = json.Unmarshal(body, &entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := router.manager.SetConnector(r.Context(), &entity)
	utils.HandleResponse(w, response, err)
}
func (router *ConnectorRouter) RemoveConnector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := router.manager.RemoveConnector(r.Context(), id)
	utils.HandleStatus(w, 200, "Deleted", err)
}