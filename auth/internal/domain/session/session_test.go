package session

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

var (
	userId     = valueobjects.NewUserID(uuid.New())
	jwtToken   = JWTToken("jwtToken")
	renewToken = JWTToken("renewJwtToken")
	csrfToken  = CSRFToken("csrfToken")
	expiresAt  = time.Now().Add(24 * time.Hour)
)

func TestSession(t *testing.T) {

	t.Run("valid creation input", func(t *testing.T) {
		session, err := Create(userId, jwtToken, renewToken, csrfToken, expiresAt)
		assert.Nil(t, err)
		assert.Equal(t, session.userId, userId)
		assert.Equal(t, session.jwtToken, jwtToken)
		assert.Equal(t, session.renewToken, renewToken)
		assert.Equal(t, session.csrfToken, csrfToken)
		assert.Equal(t, session.expiresAt, expiresAt)
		assert.False(t, session.revoked)
	})

	t.Run("invalid userId input", func(t *testing.T) {
		session, err := Create(valueobjects.UserID{}, jwtToken, renewToken, csrfToken, expiresAt)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrCreateSession)
	})

	t.Run("invalid jwt input", func(t *testing.T) {
		session, err := Create(userId, "", renewToken, csrfToken, expiresAt)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrCreateSession)
	})

	t.Run("invalid renew input", func(t *testing.T) {
		session, err := Create(userId, jwtToken, "", csrfToken, expiresAt)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrCreateSession)
	})

	t.Run("invalid csrf input", func(t *testing.T) {
		session, err := Create(userId, jwtToken, renewToken, "", expiresAt)
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrCreateSession)
	})

	t.Run("invalid expiresAt input", func(t *testing.T) {
		session, err := Create(userId, jwtToken, renewToken, csrfToken, time.Time{})
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrCreateSession)
	})

	t.Run("test revoked - correct status", func(t *testing.T) {
		session, err := Create(userId, jwtToken, renewToken, csrfToken, expiresAt)
		assert.Nil(t, err)
		session.Revoke()
		assert.True(t, session.IsRevoked())
	})

	t.Run("test revoked - invalid status", func(t *testing.T) {
		session, _ := Create(userId, jwtToken, renewToken, csrfToken, expiresAt)
		session.revoked = true
		assert.ErrorIs(t, session.Revoke(), ErrSessionAlreadyRevoke)
	})

	t.Run("test revoked - expired session", func(t *testing.T) {

		session, err := Create(userId, jwtToken, renewToken, csrfToken, time.Now().Add(time.Duration(-24)*time.Hour))
		assert.Nil(t, session)
		assert.ErrorIs(t, err, ErrCreateSession, ErrExpiresAtInThePast)
	})

	t.Run("test not expired", func(t *testing.T) {
		session, err := Create(userId, jwtToken, renewToken, csrfToken, expiresAt)
		assert.Nil(t, err)
		assert.False(t, session.IsExpired())
	})

	t.Run("test expired", func(t *testing.T) {
		session, err := Create(userId, jwtToken, renewToken, csrfToken, expiresAt)
		session.expiresAt = session.expiresAt.Add(time.Duration(-24) * time.Hour)
		assert.Nil(t, err)
		assert.True(t, session.IsExpired())
	})
}
