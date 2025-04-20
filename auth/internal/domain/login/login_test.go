package login

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

var (
	userId   = valueobjects.NewUserID(uuid.New())
	password = "testpassword"
)

func TestLogin(t *testing.T) {
	t.Run("Creating a login - success", func(t *testing.T) {
		login, err := CreateLogin(userId, password)
		assert.Nil(t, err)
		assert.NotNil(t, login)
		assert.Equal(t, userId, login.UserId())
		assert.Equal(t, password, login.Password())
	})

	t.Run("Creating a login - invalid user id", func(t *testing.T) {
		_, err := CreateLogin(valueobjects.NewUserID(uuid.Nil), password)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrCreateLogin)
	})

	t.Run("Creating a login - invalid password", func(t *testing.T) {
		_, err := CreateLogin(userId, "")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrCreateLogin)
	})

}
