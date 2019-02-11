package resolver

import (
	"goask/core/entity"
)

// UserSession is used for authentication and authorization
type UserSession struct {
	UserID entity.ID
}
