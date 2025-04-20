package session

type Repository interface {
	Save(session Session) error
	Update(session Session) (*Session, error)
	FindSessionByID(sessionID SessionID) (*Session, error)
}
