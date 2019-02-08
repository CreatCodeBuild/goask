package fakeadapter

import (
	"goask/core/adapter"
	"goask/core/entity"

	"github.com/pkg/errors"
)

type QuestionDAO struct {
	data    *Data
	userDAO *UserDAO
}

func NewQuestionDAO(data *Data, userDAO *UserDAO) *QuestionDAO {
	return &QuestionDAO{data, userDAO}
}

func (d *QuestionDAO) QuestionByID(ID entity.ID) (entity.Question, error) {
	for _, q := range d.data.questions {
		if q.ID == ID {
			return q, nil
		}
	}
	return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: ID})
}

func (d *QuestionDAO) CreateQuestion(q entity.Question) (entity.Question, error) {
	_, err := d.userDAO.UserByID(q.AuthorID)
	if err != nil {
		return entity.Question{}, err
	}

	q.ID = entity.ID(len(d.data.questions) + 1)
	d.data.questions = append(d.data.questions, q)
	return d.data.questions[len(d.data.questions)-1], d.data.serialize()
}

func (d *QuestionDAO) UpdateQuestion(p entity.QuestionUpdate) (entity.Question, error) {
	if p.ID == 0 {
		return entity.Question{}, errors.New("ID can not be 0 nor absent")
	}
	for i, q := range d.data.questions {
		if q.ID == p.ID {
			if p.Content != nil {
				q.Content = *p.Content
			}
			if p.Title != nil {
				q.Title = *p.Title
			}
			d.data.questions[i] = q
			return q, d.data.serialize()
		}
	}
	return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: p.ID})
}

func (d *QuestionDAO) DeleteQuestion(userID entity.ID, questionID entity.ID) (entity.Question, error) {
	// todo: what is the semantics of deleting a question. Are the answers associated with it deleted as well?
	_, err := d.userDAO.UserByID(userID)
	if err != nil {
		return entity.Question{}, err
	}

	question, err := d.QuestionByID(questionID)
	if err != nil {
		return entity.Question{}, err
	}

	if question.AuthorID != userID {
		return entity.Question{}, errors.WithStack(&adapter.ErrQuestionMutationDenied{QuestionID: questionID, UserID: userID})
	}

	question, ok := d.data.questions.Pop(questionID)
	if !ok {
		return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	d.data.answers = d.data.answers.Filter(func(a entity.Answer) bool { return a.QuestionID != questionID })
	return question, nil
}

func (d *QuestionDAO) VoteCount(questionID entity.ID) (up, down int, err error) {
	_, ok := d.data.questions.Get(questionID)
	if !ok {
		return 0, 0, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	return d.data.questionVotes.Count(questionID)
}

func (d *QuestionDAO) VoteQuestion(userID, questionID entity.ID, voteType entity.VoteType) (entity.Vote, error) {

	_, ok := d.data.users.Get(userID)
	if !ok {
		return entity.Vote{}, errors.WithStack(&adapter.ErrUserNotFound{ID: userID})
	}

	_, ok = d.data.questions.Get(questionID)
	if !ok {
		return entity.Vote{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	return d.data.questionVotes.Update(userID, questionID, voteType)
}

func (d *QuestionDAO) Answers(questionID entity.ID) (ret []entity.Answer, err error) {
	for _, answer := range d.data.answers {
		if answer.QuestionID == questionID {
			ret = append(ret, answer)
		}
	}
	return
}

func (d *QuestionDAO) GetAuthor(questionID entity.ID) (entity.User, error) {
	question, ok := d.data.questions.Get(questionID)
	if !ok {
		return entity.User{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	user, ok := d.data.users.Get(question.AuthorID)
	if !ok {
		return user, errors.WithMessage(errors.WithStack(&adapter.ErrUserNotFound{ID: question.AuthorID}), "impossible")
	}
	return user, nil
}

func (d *QuestionDAO) Tags(questionID entity.ID) ([]entity.Tag, error) {
	// todo
	return nil, nil
}
