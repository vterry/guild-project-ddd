package session

import (
	"time"

	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

type SessionBuilder struct {
	Session
}

func NewSessionBuilder() *SessionBuilder {
	return &SessionBuilder{}
}

func (b *SessionBuilder) WithSessionId(sessionId SessionID) *SessionBuilder {
	b.SessionID = sessionId
	return b
}

func (b *SessionBuilder) WithUserId(userId valueobjects.UserID) *SessionBuilder {
	b.userId = userId
	return b
}

func (b *SessionBuilder) WithJWTToken(jwtToken JWTToken) *SessionBuilder {
	b.jwtToken = jwtToken
	return b
}

func (b *SessionBuilder) WithRenewToken(renewToken JWTToken) *SessionBuilder {
	b.renewToken = renewToken
	return b
}

func (b *SessionBuilder) WithCSRFToken(csrfToken CSRFToken) *SessionBuilder {
	b.csrfToken = csrfToken
	return b
}

func (b *SessionBuilder) WithExpirationAt(expiresAt time.Time) *SessionBuilder {
	b.expiresAt = expiresAt
	return b
}

func (b *SessionBuilder) WithRevokeStatus(revokedStatus bool) *SessionBuilder {
	b.revoked = revokedStatus
	return b
}

func (b *SessionBuilder) Build() Session {
	return Session{
		SessionID:  b.SessionID,
		userId:     b.userId,
		jwtToken:   b.jwtToken,
		renewToken: b.renewToken,
		csrfToken:  b.csrfToken,
		expiresAt:  b.expiresAt,
		revoked:    b.revoked,
	}
}
