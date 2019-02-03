package fakeadapter

import (
	"encoding/json"
	"goask/core/adapter"
	"goask/core/entity"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Data satisfied adapter.Data. It serializes to dist.
type Data struct {
	questions Questions
	answers   Answers
	users     []entity.User
}

func NewData() (Data, error) {
	d := Data{}
	err := d.deserialize()
	if err != nil {
		return d, err
	}
	return d, nil
}

type dataSerialization struct {
	Questions Questions
	Answers   Answers
	Users     []entity.User
}

var _ adapter.Data = &Data{}

func (d *Data) file() string {
	return "./data.json"
}

func (d *Data) serialize() error {
	data := dataSerialization{
		Questions: d.questions,
		Answers:   d.answers,
		Users:     d.users,
	}

	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return errors.WithStack(err)
	}

	err = ioutil.WriteFile(d.file(), b, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}

func (d *Data) deserialize() error {
	b, err := ioutil.ReadFile(d.file())
	if err != nil {
		return err
	}

	data := dataSerialization{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	d.questions = data.Questions
	d.answers = data.Answers
	d.users = data.Users
	return nil
}

func (d *Data) Questions(search *string) ([]entity.Question, error) {
	if search == nil {
		return d.questions, nil
	}
	ret := make([]entity.Question, 0)
	for _, q := range d.questions {
		if match(q.Content, *search) {
			ret = append(ret, q)
		}
	}
	return ret, nil
}

func (d *Data) QuestionByID(ID entity.ID) (entity.Question, error) {
	for _, q := range d.questions {
		if q.ID == ID {
			return q, nil
		}
	}
	return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: ID})
}

func (d *Data) QuestionsByUserID(ID entity.ID) ([]entity.Question, error) {
	var ret []entity.Question
	for _, q := range d.questions {
		if q.AuthorID == ID {
			ret = append(ret, q)
		}
	}
	return ret, nil
}

func (d *Data) CreateQuestion(q entity.Question) (entity.Question, error) {
	_, err := d.UserByID(q.AuthorID)
	if err != nil {
		return entity.Question{}, err
	}

	q.ID = entity.ID(len(d.questions) + 1)
	d.questions = append(d.questions, q)
	return d.questions[len(d.questions)-1], d.serialize()
}

func (d *Data) UpdateQuestion(p entity.QuestionUpdate) (entity.Question, error) {
	if p.ID == 0 {
		return entity.Question{}, errors.New("ID can not be 0 nor absent")
	}
	for i, q := range d.questions {
		if q.ID == p.ID {
			if p.Content != nil {
				q.Content = *p.Content
			}
			if p.Title != nil {
				q.Title = *p.Title
			}
			d.questions[i] = q
			return q, d.serialize()
		}
	}
	return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: p.ID})
}

func (d *Data) DeleteQuestion(userID entity.ID, questionID entity.ID) (entity.Question, error) {
	// todo: what is the semantics of deleting a question. Are the answers associated with it deleted as well?
	_, err := d.UserByID(userID)
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

	question, ok := d.questions.Pop(questionID)
	if !ok {
		return entity.Question{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	d.answers = d.answers.Filter(func(a entity.Answer) bool { return a.QuestionID != questionID })
	return question, nil
}

func (d *Data) AnswersOfQuestion(QuestionID entity.ID) (ret []entity.Answer) {
	for _, answer := range d.answers {
		if answer.QuestionID == QuestionID {
			ret = append(ret, answer)
		}
	}
	return
}

func (d *Data) CreateAnswer(QuestionID entity.ID, Content string, AuthorID entity.ID) (entity.Answer, error) {
	for _, q := range d.questions {
		if q.ID == QuestionID {
			answer := d.answers.Add(QuestionID, Content, AuthorID)
			return answer, d.serialize()
		}
	}
	return entity.Answer{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: QuestionID})
}

func (d *Data) AcceptAnswer(AnswerID entity.ID, UserID entity.ID) (entity.Answer, error) {

	// Find the question this answer belongs to
	answer, ok := d.answers.Get(AnswerID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrAnswerNotFound{ID: AnswerID})
	}

	q, ok := d.questions.Get(answer.QuestionID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrQuestionOfAnswerNotFound{QuestionID: answer.QuestionID, AnswerID: AnswerID})
	}

	// Find if this user is the author of the question this answer belongs to
	if q.AuthorID != UserID {
		return answer, errors.WithStack(&adapter.ErrUserIsNotAuthorOfQuestion{QuestionID: q.ID, UserID: UserID})
	}

	answer = d.answers.Accept(AnswerID)
	return answer, d.serialize()
}

func (d *Data) DeleteAnswer(ANswerID entity.ID, UserID entity.ID) (entity.Answer, error) {
	// todo
	return entity.Answer{}, nil
}

//

func (d *Data) UserByID(ID entity.ID) (entity.User, error) {
	for _, user := range d.users {
		if user.ID == ID {
			return user, nil
		}
	}
	return entity.User{}, errors.WithStack(&adapter.ErrUserNotFound{ID: ID})
}

func (d *Data) Users() ([]entity.User, error) {
	return d.users, nil
}

func (d *Data) CreateUser(name string) (entity.User, error) {
	user := entity.User{ID: entity.ID(len(d.users) + 1), Name: name}
	d.users = append(d.users, user)
	return user, d.serialize()
}

func match(s1, s2 string) bool {
	return strings.Contains(s1, s2)
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
