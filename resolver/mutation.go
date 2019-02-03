package resolver

import (
	"goask/core/adapter"
	"goask/core/entity"
	"goask/log"
)

type Mutation struct {
	Data adapter.Data
}

func (m *Mutation) QuestionMutation(args struct{ UserID int32 }) (QuestionMutation, error) {
	_, err := m.Data.UserByID(entity.ID(args.UserID))
	if err != nil {
		return QuestionMutation{}, err
	}

	return QuestionMutation{
		stdResolver: stdResolver{
			data: m.Data,
			log:  &log.Logger{},
		},
		userSession: UserSession{
			UserID: entity.ID(args.UserID),
		},
	}, nil
}

func (m *Mutation) Answer(args struct{ UserID int32 }) (AnswerMutation, error) {
	_, err := m.Data.UserByID(entity.ID(args.UserID))
	if err != nil {
		return AnswerMutation{}, err
	}

	return AnswerMutation{
		stdResolver: stdResolver{
			data: m.Data,
			log:  &log.Logger{},
		},
		userSession: UserSession{
			UserID: entity.ID(args.UserID),
		},
	}, nil
}

func (m *Mutation) User() (UserMutation, error) {
	return UserMutation{stdResolver: stdResolver{
		data: m.Data,
		log:  &log.Logger{},
	}}, nil
}

// QuestionMutation resolves all mutations of questions.
type QuestionMutation struct {
	stdResolver
	userSession UserSession
}

// Create creates a question.
func (m QuestionMutation) Create(args struct{ Title, Content string }) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	q, err := m.data.CreateQuestion(
		entity.Question{
			Title:    args.Title,
			Content:  args.Content,
			AuthorID: m.userSession.UserID,
		},
	)

	return QuestionOne(q, m.data), err
}

// Update updates a question
func (m QuestionMutation) Update(input QuestionInput) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	input.QuestionUpdate.ID = entity.ID(input.ID)
	q, err := m.data.UpdateQuestion(input.QuestionUpdate)
	if err != nil {
		m.log.Error(err)
	}
	return QuestionOne(q, m.data), err
}

func (m QuestionMutation) Delete(args struct{ ID int32 }) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	question, err := m.data.DeleteQuestion(entity.ID(m.userSession.UserID), entity.ID(args.ID))
	return QuestionOne(question, m.data), err
}

type AnswerMutation struct {
	stdResolver
	userSession UserSession
}

func (m AnswerMutation) Create(args AnswerCreationInput) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	answer, err := m.data.CreateAnswer(entity.ID(args.QuestionID), args.Content, m.userSession.UserID)
	if err != nil {
		m.log.Error(err)
	}
	return Answer{entity: answer, data: m.data}, err
}

func (m AnswerMutation) Accept(args struct{ AnswerID int32 }) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	an, err := m.data.AcceptAnswer(entity.ID(args.AnswerID), m.userSession.UserID)
	return AnswerOne(an, m.data), err
}

func (m AnswerMutation) Delete(args struct{ AnswerID int32 }) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	an, err := m.data.DeleteAnswer(entity.ID(args.AnswerID), m.userSession.UserID)
	return AnswerOne(an, m.data), err
}

type UserMutation struct {
	stdResolver
}

func (m UserMutation) Create(args struct{ Name string }) (User, error) {
	if err := m.check(); err != nil {
		return User{}, err
	}

	user, err := m.data.CreateUser(args.Name)
	return UserOne(user, m.data), err
}

type logger interface {
	Error(err error)
}
