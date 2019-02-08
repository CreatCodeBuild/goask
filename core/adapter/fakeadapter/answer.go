package fakeadapter

import (
	"goask/core/adapter"
	"goask/core/entity"

	"github.com/pkg/errors"
)

// AnswerDAO implements adapter.AnswerDAO
type AnswerDAO struct {
	data *Data
}

func NewAnswerDAO(data *Data) *AnswerDAO {
	return &AnswerDAO{data}
}

// AnswersOfQuestion get answers belonging to the question.
func (d *AnswerDAO) AnswersOfQuestion(QuestionID entity.ID) (ret []entity.Answer) {
	for _, answer := range d.data.answers {
		if answer.QuestionID == QuestionID {
			ret = append(ret, answer)
		}
	}
	return
}

// CreateAnswer creates an answer.
func (d *AnswerDAO) CreateAnswer(QuestionID entity.ID, Content string, AuthorID entity.ID) (entity.Answer, error) {
	for _, q := range d.data.questions {
		if q.ID == QuestionID {
			answer := d.data.answers.Add(QuestionID, Content, AuthorID)
			return answer, d.data.serialize()
		}
	}
	return entity.Answer{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: QuestionID})
}

// AcceptAnswer accepts an answer.
func (d *AnswerDAO) AcceptAnswer(AnswerID entity.ID, UserID entity.ID) (entity.Answer, error) {

	// Find the question this answer belongs to
	answer, ok := d.data.answers.Get(AnswerID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrAnswerNotFound{ID: AnswerID})
	}

	q, ok := d.data.questions.Get(answer.QuestionID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrQuestionOfAnswerNotFound{QuestionID: answer.QuestionID, AnswerID: AnswerID})
	}

	// Find if this user is the author of the question this answer belongs to
	if q.AuthorID != UserID {
		return answer, errors.WithStack(&adapter.ErrUserIsNotAuthorOfQuestion{QuestionID: q.ID, UserID: UserID})
	}

	answer = d.data.answers.Accept(AnswerID)
	return answer, d.data.serialize()
}

// DeleteAnswer delete an answer
func (d *AnswerDAO) DeleteAnswer(AnswerID entity.ID, UserID entity.ID) (entity.Answer, error) {
	answer, ok := d.data.answers.Get(AnswerID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrAnswerNotFound{ID: AnswerID})
	}

	if answer.AuthorID != UserID {
		return answer, errors.WithStack(&adapter.ErrUserIsNotAuthorOfAnswer{AnswerID: answer.ID, UserID: UserID})
	}

	d.data.answers.Delete(answer.ID)
	return answer, d.data.serialize()
}
