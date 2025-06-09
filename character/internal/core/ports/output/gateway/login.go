package gateway

import (
	"context"

	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
)

type Login interface {
	IsLoginValid(ctx context.Context, loginId login.LoginID) (bool, error)
}
