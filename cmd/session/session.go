package session

import "restforavito/cmd/generate"

type userSession struct {
	Username string
}

type Session struct {
	data map[string]*userSession
}

func NewSession() *Session {
	s := new(Session)
	s.data = make(map[string]*userSession)
	return s
}

func (s *Session) Init(username string) string {
	sessionId := generate.GenerateCookieToken()
	data := &userSession{Username: username}
	s.data[sessionId] = data
	return sessionId
}
