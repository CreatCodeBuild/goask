package fakeadapter

import (
	"encoding/json"
	"goask/core/adapter"
	"goask/core/entity"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type Serializer interface {
	Serialize([]byte) error
	Deserialize() ([]byte, error)
}

type FileSerializer struct {
	fileName string
}

func NewFileSerializer(file string) FileSerializer {
	return FileSerializer{fileName: file}
}

func (f FileSerializer) Serialize(b []byte) error {
	err := ioutil.WriteFile(f.fileName, b, os.ModePerm)
	return errors.WithStack(err)
}

func (f FileSerializer) Deserialize() ([]byte, error) {
	b, err := ioutil.ReadFile(f.fileName)
	return b, errors.WithStack(err)
}

type BufferSerializer struct {
	data []byte
}

func (s BufferSerializer) Serialize(b []byte) error {
	s.data = b
	return nil
}

func (s BufferSerializer) Deserialize() ([]byte, error) {
	return s.data, nil
}

// Data satisfied adapter.Data. It serializes to dist.
type Data struct {
	questions     Questions
	questionVotes QuestionVotes
	answers       Answers
	users         Users
	storage       Serializer
}

func NewData(storage Serializer) (*Data, error) {
	d := &Data{}
	if storage == nil {
		return nil, errors.New("storage == nil")
	}
	d.storage = storage
	d.questionVotes = make(QuestionVotes)

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

	return d.storage.Serialize(b)
}

func (d *Data) deserialize() error {
	b, err := d.storage.Deserialize()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}

	data := dataSerialization{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return errors.WithStack(err)
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

func (d *Data) VoteCount(questionID entity.ID) (up, down int, err error) {
	_, ok := d.questions.Get(questionID)
	if !ok {
		return 0, 0, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	return d.questionVotes.Count(questionID)
}

func (d *Data) VoteQuestion(userID, questionID entity.ID, voteType entity.VoteType) (entity.Vote, error) {

	_, ok := d.users.Get(userID)
	if !ok {
		return entity.Vote{}, errors.WithStack(&adapter.ErrUserNotFound{ID: userID})
	}

	_, ok = d.questions.Get(questionID)
	if !ok {
		return entity.Vote{}, errors.WithStack(&adapter.ErrQuestionNotFound{ID: questionID})
	}

	return d.questionVotes.Update(userID, questionID, voteType)
}

// Answers

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

func (d *Data) DeleteAnswer(AnswerID entity.ID, UserID entity.ID) (entity.Answer, error) {
	answer, ok := d.answers.Get(AnswerID)
	if !ok {
		return answer, errors.WithStack(&adapter.ErrAnswerNotFound{ID: AnswerID})
	}

	if answer.AuthorID != UserID {
		return answer, errors.WithStack(&adapter.ErrUserIsNotAuthorOfAnswer{AnswerID: answer.ID, UserID: UserID})
	}

	d.answers.Delete(answer.ID)
	return answer, d.serialize()
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
