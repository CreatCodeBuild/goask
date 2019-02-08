package resolver

import (
	"goask/core/entity"
	"goask/log"
)

type Mutation struct {
	stdResolver
}

func NewMutation(stdResolver stdResolver) Mutation {
	return Mutation{stdResolver}
}

func (m *Mutation) QuestionMutation(args struct{ UserID int32 }) (QuestionMutation, error) {
	_, err := m.UserDAO.UserByID(entity.ID(args.UserID))
	if err != nil {
		return QuestionMutation{}, err
	}

	return QuestionMutation{
		stdResolver: m.stdResolver,
		userSession: UserSession{
			UserID: entity.ID(args.UserID),
		},
	}, nil
}

func (m *Mutation) Answer(args struct{ UserID int32 }) (AnswerMutation, error) {
	_, err := m.UserDAO.UserByID(entity.ID(args.UserID))
	if err != nil {
		return AnswerMutation{}, err
	}

	return AnswerMutation{
		stdResolver: stdResolver{
			QuestionDAO: m.QuestionDAO,
			AnswerDAO:   m.AnswerDAO,
			UserDAO:     m.UserDAO,
			log:         &log.Logger{},
		},
		userSession: UserSession{
			UserID: entity.ID(args.UserID),
		},
	}, nil
}

func (m *Mutation) User() (UserMutation, error) {
	return UserMutation{stdResolver: stdResolver{
		QuestionDAO: m.QuestionDAO,
		AnswerDAO:   m.AnswerDAO,
		UserDAO:     m.UserDAO,
		log:         &log.Logger{},
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

	q, err := m.QuestionDAO.CreateQuestion(
		entity.Question{
			Title:    args.Title,
			Content:  args.Content,
			AuthorID: m.userSession.UserID,
		},
	)

	return QuestionOne(q, m.stdResolver), err
}

// Update updates a question
func (m QuestionMutation) Update(input QuestionInput) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	input.QuestionUpdate.ID = entity.ID(input.ID)
	q, err := m.QuestionDAO.UpdateQuestion(input.QuestionUpdate)
	if err != nil {
		m.log.Error(err)
	}
	return QuestionOne(q, m.stdResolver), err
}

func (m QuestionMutation) Delete(args struct{ ID int32 }) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	question, err := m.QuestionDAO.DeleteQuestion(entity.ID(m.userSession.UserID), entity.ID(args.ID))
	return QuestionOne(question, m.stdResolver), err
}

func (m QuestionMutation) Vote(args struct {
	ID   int32
	Type string
}) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	_, err := m.QuestionDAO.VoteQuestion(m.userSession.UserID, entity.ID(args.ID), entity.VoteType(args.Type))
	if err != nil {
		return Question{}, err
	}

	question, err := m.QuestionDAO.QuestionByID(entity.ID(args.ID))
	return QuestionOne(question, m.stdResolver), err
}

type AnswerMutation struct {
	stdResolver
	userSession UserSession
}

func (m AnswerMutation) Create(args AnswerCreationInput) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	answer, err := m.AnswerDAO.CreateAnswer(entity.ID(args.QuestionID), args.Content, m.userSession.UserID)
	if err != nil {
		m.log.Error(err)
	}
	return AnswerOne(answer, m.stdResolver), err
}

func (m AnswerMutation) Accept(args struct{ AnswerID int32 }) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	an, err := m.AnswerDAO.AcceptAnswer(entity.ID(args.AnswerID), m.userSession.UserID)
	return AnswerOne(an, m.stdResolver), err
}

func (m AnswerMutation) Delete(args struct{ AnswerID int32 }) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	an, err := m.AnswerDAO.DeleteAnswer(entity.ID(args.AnswerID), m.userSession.UserID)
	return AnswerOne(an, m.stdResolver), err
}

type UserMutation struct {
	stdResolver
}

func (m UserMutation) Create(args struct{ Name string }) (User, error) {
	if err := m.check(); err != nil {
		return User{}, err
	}

	user, err := m.UserDAO.CreateUser(args.Name)
	return UserOne(user, m.stdResolver), err
}

type logger interface {
	Error(err error)
}
