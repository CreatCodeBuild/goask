package adaptertest

import (
	"fmt"
	"testing"

	"goask/core/adapter"
	"goask/core/entity"

	"github.com/stretchr/testify/require"
)

func All(t *testing.T, questionDAO adapter.QuestionDAO, answerDAO adapter.AnswerDAO, userDAO adapter.UserDAO, tagDAO adapter.TagDAO) {

	t.Run("create questions", func(t2 *testing.T) {
		_, err := questionDAO.CreateQuestion(entity.Question{AuthorID: "x"}, nil)
		require.EqualError(t, err, "user:x not found")

		user, err := userDAO.CreateUser("user 1")
		require.NoError(t, err)
		require.Equal(t, user.Name, "user 1")
		require.Equal(t, user.ID, entity.NewIDString("1"))

		question, err := questionDAO.CreateQuestion(entity.Question{AuthorID: "1"}, nil)
		require.NoError(t, err)
		require.Equal(t, entity.Question{AuthorID: "1", ID: "1"}, question)

		t.Run("create question with tag", func(t2 *testing.T) {
			q, err := questionDAO.CreateQuestion(entity.Question{AuthorID: "1"}, []entity.Tag{"Go1", "Go2"})
			require.NoError(t, err)

			tags, err := questionDAO.Tags(q.ID)
			require.NoError(t, err)
			require.Equal(t, entity.TagSet{"Go1": struct{}{}, "Go2": struct{}{}}, tags)
		})
	})

	t.Run("create answers", func(t2 *testing.T) {
		answer, err := answerDAO.CreateAnswer("1", "answer 1", "1")
		require.NoError(t, err)
		require.Equal(t, entity.Answer{ID: "1", QuestionID: "1", AuthorID: "1", Content: "answer 1"}, answer)

		t.Run("accept answers", func(t2 *testing.T) {
			acceptedAnswer, err := answerDAO.AcceptAnswer(answer.ID, answer.AuthorID)
			require.NoError(t, err)
			require.Equal(t,
				entity.Answer{ID: "1", QuestionID: "1", AuthorID: "1", Content: "answer 1", Accepted: true},
				acceptedAnswer)

			_, err = answerDAO.AcceptAnswer(answer.ID, "-1")
			require.EqualError(t, err, "user:-1 is no the author of question:1")

			t.Run("delete answers", func(t2 *testing.T) {
				deletedAnswer, err := answerDAO.DeleteAnswer(answer.ID, answer.AuthorID)
				require.NoError(t, err)
				require.Equal(t, deletedAnswer, acceptedAnswer)
			})
		})
	})

	t.Run("delete questions", func(t2 *testing.T) {
		_, err := questionDAO.DeleteQuestion("2", "1")
		require.EqualError(t, err, "user:2 not found")

		user, err := userDAO.CreateUser("user 2")
		require.NoError(t, err)

		_, err = questionDAO.DeleteQuestion(user.ID, "1")
		require.EqualError(t, err, "user:2 is not authorized to delete question:1")

		question, err := questionDAO.DeleteQuestion("1", "1")
		require.NoError(t, err)
		require.Equal(t, entity.Question{AuthorID: "1", ID: "1"}, question)

		_, err = questionDAO.QuestionByID("1")
		require.EqualError(t, err, "question:1 not found")

		answers, err := questionDAO.Answers("1") // The answers associated with this question should be deleted as well
		require.NoError(t, err)
		require.Empty(t, answers)
	})

	t.Run("vote a question", func(t2 *testing.T) {
		user, err := userDAO.CreateUser("user 2")
		require.NoError(t, err)

		question, err := questionDAO.CreateQuestion(entity.Question{AuthorID: user.ID}, nil)
		require.NoError(t, err)

		up, down, err := questionDAO.VoteCount(question.ID)
		require.NoError(t, err)
		require.Equal(t, 0, up)
		require.Equal(t, 0, down)

		vote, err := questionDAO.VoteQuestion(user.ID, question.ID, entity.UpVote())
		require.NoError(t, err)

		up, down, err = questionDAO.VoteCount(question.ID)
		require.NoError(t, err)
		require.Equal(t, 1, up)
		require.Equal(t, 0, down)

		require.Equal(t, user.ID, vote.UserID)
		require.Equal(t, question.ID, vote.QuestionID)
		require.Equal(t, entity.UpVote(), vote.Type)

		t.Run("The same use replace a vote with the opposite type for the same question", func(t2 *testing.T) {
			vote2, err := questionDAO.VoteQuestion(user.ID, question.ID, entity.DownVote())
			require.NoError(t, err)
			require.Equal(t, user.ID, vote2.UserID)
			require.Equal(t, question.ID, vote2.QuestionID)
			require.Equal(t, entity.DownVote(), vote2.Type)

			up, down, err = questionDAO.VoteCount(question.ID)
			require.NoError(t, err)
			require.Equal(t, 0, up)
			require.Equal(t, 1, down)
		})

		t.Run("The same user can't vote the same for the same question", func(t2 *testing.T) {
			_, err := questionDAO.VoteQuestion(user.ID, question.ID, entity.DownVote())
			require.Equal(t,
				fmt.Sprintf("user:%v has voted DOWN for question:%v", user.ID, question.ID),
				err.Error())
		})
	})

	t.Run("tag questions", func(t2 *testing.T) {
		user, err := userDAO.CreateUser("tagger")
		require.NoError(t, err)

		question1, err := questionDAO.CreateQuestion(entity.Question{AuthorID: user.ID}, nil)
		require.NoError(t, err)
		question2, err := questionDAO.CreateQuestion(entity.Question{AuthorID: user.ID}, nil)
		require.NoError(t, err)

		tags, err := questionDAO.Tags(question1.ID)
		require.NoError(t, err)
		require.Empty(t, tags)

		// Tag question1 with Python, Go
		_, err = questionDAO.UpdateQuestion(entity.QuestionUpdate{ID: question1.ID, Tags: []entity.Tag{"Python", "Go"}})
		require.NoError(t, err)

		tags, err = questionDAO.Tags(question1.ID)
		require.NoError(t, err)
		require.Equal(t, entity.NewTagSet("Python", "Go"), tags)

		// Tag question2 with Python
		_, err = questionDAO.UpdateQuestion(entity.QuestionUpdate{ID: question2.ID, Tags: []entity.Tag{"Python"}})
		require.NoError(t, err)

		tags, err = questionDAO.Tags(question2.ID)
		require.NoError(t, err)
		require.Equal(t, entity.NewTagSet("Python"), tags)

		// Assert tag Python has question1 and question2
		questions, err := tagDAO.Questions("Python")
		require.NoError(t, err)
		require.Equal(t, entity.NewQuestionSet(question1, question2), questions)

		// Assert tag Go has question1
		questions, err = tagDAO.Questions("Go")
		require.NoError(t, err)
		require.Equal(t, entity.NewQuestionSet(question1), questions)

		// Assert tag Java has no question
		questions, err = tagDAO.Questions("Java")
		require.NoError(t, err)
		require.Empty(t, questions)
	})

	t.Run("user question count", func(t2 *testing.T) {
		u := setupOneUserAndThreeQuestions(t, questionDAO, userDAO)
		count, err := userDAO.QuestionCount(u.ID)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})
}

func setupOneUserAndThreeQuestions(t *testing.T, qDAO adapter.QuestionDAO, uDAO adapter.UserDAO) entity.User {
	u, err := uDAO.CreateUser("Test User")
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		_, err = qDAO.CreateQuestion(entity.Question{
			AuthorID: u.ID,
		}, nil)
		require.NoError(t, err)
	}
	return u
}
