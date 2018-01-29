package sync

import (
	"net/http"
	"github.com/dataprism/dataprism-commons/utils"
	"github.com/dataprism/dataprism-sync-runtime/plugins"
)

type PluginRouter struct {
	registry *plugins.PluginRegistry
}

func NewPluginRouter(registry *plugins.PluginRegistry) (*PluginRouter) {
	return &PluginRouter{registry}
}

func (router *PluginRouter) List(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponse(w, router.registry.Plugins(), nil)
}

func (router *PluginRouter) ListInputs(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponse(w, router.registry.GetInputTypes(), nil)
}

func (router *PluginRouter) ListOutputs(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponse(w, router.registry.GetOutputTypes(), nil)
}