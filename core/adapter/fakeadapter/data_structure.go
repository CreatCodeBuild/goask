package fakeadapter

import (
	"goask/core/adapter"
	"strings"

	"goask/core/entity"

	"github.com/pkg/errors"
)

func match(s1, s2 string) bool {
	return strings.Contains(s1, s2)
}

type Users []entity.User

func (u *Users) Get(userID entity.ID) (entity.User, bool) {
	for _, qu := range *u {
		if qu.ID == userID {
			return qu, true
		}
	}
	return entity.User{}, false
}

type Questions []entity.Question

func (q *Questions) Get(questionID entity.ID) (entity.Question, bool) {
	for _, qu := range *q {
		if qu.ID == questionID {
			return qu, true
		}
	}
	return entity.Question{}, false
}

func (q *Questions) Pop(questionID entity.ID) (entity.Question, bool) {
	for i, qu := range *q {
		if qu.ID == questionID {
			*q = q.Delete(i)
			return qu, true
		}
	}
	return entity.Question{}, false
}

func (q Questions) Delete(i int) Questions {
	return append(q[:i], q[i+1:]...)
}

type Answers []entity.Answer

func (a *Answers) Add(QuestionID entity.ID, Content string, AuthorID entity.ID) entity.Answer {
	// todo: serialize
	*a = append(*a, entity.Answer{
		ID:         entity.ID(len(*a) + 1),
		Content:    Content,
		QuestionID: QuestionID,
		AuthorID:   AuthorID,
	})
	return (*a)[len(*a)-1]
}

func (a *Answers) OfQuestion(questionID entity.ID) Answers {
	var ans Answers
	for _, answer := range *a {
		if answer.QuestionID == questionID {
			ans = append(ans, answer)
		}
	}
	return ans
}

func (a *Answers) Get(answerID entity.ID) (entity.Answer, bool) {
	for _, an := range *a {
		if an.ID == answerID {
			return an, true
		}
	}
	return entity.Answer{}, false
}

func (a *Answers) Accept(answerID entity.ID) entity.Answer {
	// todo: serialize
	for i := range *a {
		if (*a)[i].ID == answerID {
			(*a)[i].Accepted = true
			return (*a)[i]
		}
	}
	return entity.Answer{}
}

func (a *Answers) Filter(f func(entity.Answer) bool) Answers {
	var ret Answers
	for _, an := range *a {
		if f(an) {
			ret = append(ret, an)
		}
	}
	return ret
}

func (a *Answers) Delete(answerID entity.ID) {

	answers := a.Filter(func(a entity.Answer) bool {
		return a.ID != answerID
	})

	*a = answers
}

type voteTarget struct {
	userID, questionID entity.ID
}

type QuestionVotes map[voteTarget]entity.VoteType

func (a *QuestionVotes) Update(userID, questionID entity.ID, voteType entity.VoteType) (entity.Vote, error) {
	m := *a

	vt := voteTarget{
		userID:     userID,
		questionID: questionID,
	}

	vote := entity.Vote{
		UserID:     userID,
		QuestionID: questionID,
		Type:       voteType,
	}

	if vType, ok := m[vt]; ok {
		if vType == voteType {
			return entity.Vote{}, errors.WithStack(&adapter.ErrVoteDuplicated{Vote: vote})
		}
		m[vt] = voteType
	}

	m[vt] = voteType
	return vote, nil
}

func (a *QuestionVotes) Count(questionID entity.ID) (up int, down int, err error) {
	up = 0
	down = 0

	for _, vType := range *a {
		if vType == entity.UpVote() {
			up += 1
		} else if vType == entity.DownVote() {
			down += 1
		} else {
			err = errors.Errorf("impossible, vote type is %s", vType)
			return
		}
	}
	return
}
