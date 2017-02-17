package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mwheatley3/ww/server/httproutes"
	"github.com/mwheatley3/ww/server/personal/api/api"
)

func (s *Server) getUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var (
		meUser = p.ByName("userID") == "me"
		resp   = httproutes.JSON(s.logger, w)
	)

	// handle other users later!
	if !meUser {
		resp.Error(api.ErrInvalidUserID, http.StatusNotFound)
		return
	}

	u, err := s.service.GetUser(getSession(r).UserID)

	switch err {
	case nil:
		resp.Success(u)
	case api.ErrInvalidUserID:
		resp.Error(err, http.StatusNotFound)
	default:
		resp.Error(err, http.StatusInternalServerError)
	}
}
