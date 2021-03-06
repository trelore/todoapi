package internal

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

/* swagger:route DELETE /items/{itemID} item deleteItem
Deletes an item from the todo list given an ID.

	Consumes:
	- application/json

	Produces:
	- application/json

	Schemes: https

	Security:
		basicAuth:
			type: http
			scheme: basic

	Responses:
		default: genericError
		204:
*/
func (s server) DeleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = idParam{
		ID: ps.ByName("id"),
	}
	if _, err := uuid.Parse(params.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.db.Delete(params.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
