package sqlx_store

import (
	"github.com/gin-contrib/sessions"
	ginSessions "github.com/gin-contrib/sessions"
	"net/http"
)

type Store struct {
}

func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (st *Store) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Options(options ginSessions.Options) {
	//TODO implement me
	panic("implement me")
}
