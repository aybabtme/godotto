package do

import (
	"net/http"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/julienschmidt/httprouter"
)

type handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params)

type dropletsAPI struct {
	db *database
}

func (api *dropletsAPI) create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := new(struct {
		Name   string `json:"name"`
		Region string `json:"region"`
		Size   string `json:"size"`
		Image  string `json:"image"`
	})
	decodeReq(r, &req)
	switch {
	case req.Name == "":
		throwBadRequest("name is required")
	case req.Region == "":
		throwBadRequest("region is required")
	case req.Size == "":
		throwBadRequest("size is required")
	case req.Image == "":
		throwBadRequest("image slug is required")
	}

	d, err := api.db.Droplets().Create(
		req.Name,
		req.Region,
		req.Size,
		req.Image,
	)
	if err != nil {
		throwInternalError("querying database: %v", err)
	}

	respond(w, http.StatusCreated, struct {
		Droplet *godo.Droplet `json:"droplet"`
		Links   *godo.Links   `json:"links"`
	}{
		Droplet: d,
	})
}

func (api *dropletsAPI) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		throwBadRequest("not an integer: %v", err)
	}
	if err := api.db.Droplets().Delete(id); err != nil {
		throwInternalError("querying database: %v", err)
	}
}
