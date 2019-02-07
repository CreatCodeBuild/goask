package adapter

import (
	"fmt"
	"goask/core/entity"
)

type ErrVoteDuplicated struct {
	Vote entity.Vote
}

func (e *ErrVoteDuplicated) Error() string {
	return fmt.Sprintf("user:%v has voted %v for question:%v", e.Vote.UserID, e.Vote.Type, e.Vote.QuestionID)
}
