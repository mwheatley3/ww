package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mwheatley3/ww/server/httproutes"
	"github.com/mwheatley3/ww/server/httpserver"
	"github.com/mwheatley3/ww/server/personal/api/api"
	"github.com/satori/go.uuid"
)

type session struct {
	UserID uuid.UUID
}

var sesContextKey = "github.com/conversable/woodhouse/server/portal"

func withSession(r *http.Request, v *session) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), &sesContextKey, v))
}

func getSession(r *http.Request) *session {
	return r.Context().Value(&sesContextKey).(*session)
}

func (s *Server) authMiddleware() httpserver.MiddlewareFunc {
	return httpserver.MiddlewareFunc(
		func(method, path string, h httpserver.Handler) httpserver.Handler {
			if !strings.HasPrefix(path, "/api/") || (method == "POST" && path == "/api/auth") {
				return h
			}

			return httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				var ses session
				if err := s.cookie.GetSession(r, &ses); err != nil {
					s.logger.Errorf("Cookie get error: %s", err)
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}

				if ses.UserID == uuid.Nil {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}

				r = withSession(r, &ses)
				h.Handle(w, r, p)
			})
		})

}

func (s *Server) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		resp   = httproutes.JSON(s.logger, w)
		params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
	)

	if err := httproutes.DecodeJSON(r, &params); err != nil {
		s.logger.Errorf("Decode json error: %s", err)
		resp.Error(nil, http.StatusBadRequest)
		return
	}

	// email check takes no time, but pw check actually does work
	// since we don't expose which check failed to the UI, don't give
	// a time based hint about which check failed.
	st := time.Now()
	user, err := s.service.AuthenticateUser(params.Email, params.Password)
	dur := time.Since(st)

	if diff := 2*time.Second - dur; diff > 0 {
		time.Sleep(diff)
	}

	switch err {
	case nil:
		err = s.cookie.SetSession(w, &session{
			UserID: user.ID,
		})

		if err != nil {
			s.logger.Errorf("Set session err: %s", err)
			resp.Error(api.ErrUnknownError, http.StatusInternalServerError)
		} else {
			resp.Success(user)
		}
	case api.ErrInvalidPassword, api.ErrInvalidEmail:
		resp.Error(fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
	default:
		resp.Error(err, http.StatusInternalServerError)
	}
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := httproutes.JSON(s.logger, w)
	s.cookie.SetSession(w, nil)

	resp.Success(true)
}
