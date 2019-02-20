package resolver

import (
	"goask/core/entity"

	"github.com/graph-gophers/graphql-go"
)

func (m *Mutation) QuestionMutation(args struct{ UserID graphql.ID }) (QuestionMutation, error) {
	_, err := m.UserDAO.UserByID(ToEntityID(args.UserID))
	if err != nil {
		return QuestionMutation{}, err
	}

	return QuestionMutation{
		stdResolver: m.stdResolver,
		userSession: UserSession{
			UserID: ToEntityID(args.UserID),
		},
	}, nil
}

func (m *Mutation) Answer(args struct{ UserID graphql.ID }) (AnswerMutation, error) {
	_, err := m.UserDAO.UserByID(ToEntityID(args.UserID))
	if err != nil {
		return AnswerMutation{}, err
	}

	return AnswerMutation{
		stdResolver: m.stdResolver,
		userSession: UserSession{
			UserID: ToEntityID(args.UserID),
		},
	}, nil
}

func (m *Mutation) User() (UserMutation, error) {
	return UserMutation{stdResolver: m.stdResolver}, m.check()
}

// QuestionMutation resolves all mutations of questions.
type QuestionMutation struct {
	stdResolver
	userSession UserSession
}

// Create creates a question.
func (m QuestionMutation) Create(args struct {
	Title, Content string
	Tags           *[]entity.Tag
}) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	var tags []entity.Tag
	if args.Tags != nil {
		tags = *args.Tags
	}

	q, err := m.QuestionDAO.CreateQuestion(
		entity.Question{
			Title:    args.Title,
			Content:  args.Content,
			AuthorID: m.userSession.UserID,
		},
		tags,
	)

	return QuestionOne(q, m.stdResolver), err
}

// Update updates a question
func (m QuestionMutation) Update(input QuestionInput) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	input.QuestionUpdate.ID = ToEntityID(input.ID)
	q, err := m.QuestionDAO.UpdateQuestion(input.QuestionUpdate)
	if err != nil {
		m.log.Error(err)
	}
	return QuestionOne(q, m.stdResolver), err
}

func (m QuestionMutation) Delete(args struct{ ID graphql.ID }) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	question, err := m.QuestionDAO.DeleteQuestion(m.userSession.UserID, ToEntityID(args.ID))
	return QuestionOne(question, m.stdResolver), err
}

func (m QuestionMutation) Vote(args struct {
	ID   graphql.ID
	Type string
}) (Question, error) {
	if err := m.check(); err != nil {
		return Question{}, err
	}

	_, err := m.QuestionDAO.VoteQuestion(m.userSession.UserID, ToEntityID(args.ID), entity.VoteType(args.Type))
	if err != nil {
		return Question{}, err
	}

	question, err := m.QuestionDAO.QuestionByID(ToEntityID(args.ID))
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

	answer, err := m.AnswerDAO.CreateAnswer(ToEntityID(args.QuestionID), args.Content, m.userSession.UserID)
	if err != nil {
		m.log.Error(err)
	}
	return AnswerOne(answer, m.stdResolver), err
}

func (m AnswerMutation) Accept(args struct{ AnswerID graphql.ID }) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	an, err := m.AnswerDAO.AcceptAnswer(ToEntityID(args.AnswerID), m.userSession.UserID)
	return AnswerOne(an, m.stdResolver), err
}

func (m AnswerMutation) Delete(args struct{ AnswerID graphql.ID }) (Answer, error) {
	if err := m.check(); err != nil {
		return Answer{}, err
	}

	an, err := m.AnswerDAO.DeleteAnswer(ToEntityID(args.AnswerID), m.userSession.UserID)
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
