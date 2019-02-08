package adaptertest

import (
	"fmt"
	"testing"

	"goask/core/adapter"
	"goask/core/entity"

	"github.com/stretchr/testify/require"
)

func All(t *testing.T, questionDAO adapter.QuestionDAO, answerDAO adapter.AnswerDAO, userDAO adapter.UserDAO) {

	t.Run("create questions", func(t *testing.T) {
		_, err := questionDAO.CreateQuestion(entity.Question{})
		require.EqualError(t, err, "user:0 not found")

		user, err := userDAO.CreateUser("user 1")
		require.NoError(t, err)
		require.Equal(t, user.Name, "user 1")
		require.Equal(t, user.ID, entity.ID(1))

		question, err := questionDAO.CreateQuestion(entity.Question{AuthorID: 1})
		require.NoError(t, err)
		require.Equal(t, entity.Question{AuthorID: 1, ID: 1}, question)
	})

	t.Run("create answers", func(t *testing.T) {
		answer, err := answerDAO.CreateAnswer(1, "answer 1", 1)
		require.NoError(t, err)
		require.Equal(t, entity.Answer{ID: 1, QuestionID: 1, AuthorID: 1, Content: "answer 1"}, answer)

		t.Run("accept answers", func(t *testing.T) {
			acceptedAnswer, err := answerDAO.AcceptAnswer(answer.ID, answer.AuthorID)
			require.NoError(t, err)
			require.Equal(t,
				entity.Answer{ID: 1, QuestionID: 1, AuthorID: 1, Content: "answer 1", Accepted: true},
				acceptedAnswer)

			_, err = answerDAO.AcceptAnswer(answer.ID, -1)
			require.EqualError(t, err, "user:-1 is no the author of question:1")

			t.Run("delete answers", func(t *testing.T) {
				deletedAnswer, err := answerDAO.DeleteAnswer(answer.ID, answer.AuthorID)
				require.NoError(t, err)
				require.Equal(t, deletedAnswer, acceptedAnswer)
			})
		})
	})

	t.Run("delete questions", func(t *testing.T) {
		_, err := questionDAO.DeleteQuestion(2, 1)
		require.EqualError(t, err, "user:2 not found")

		user, err := userDAO.CreateUser("user 2")
		require.NoError(t, err)

		_, err = questionDAO.DeleteQuestion(user.ID, 1)
		require.EqualError(t, err, "user:2 is not authorized to delete question:1")

		question, err := questionDAO.DeleteQuestion(1, 1)
		require.NoError(t, err)
		require.Equal(t, entity.Question{AuthorID: 1, ID: 1}, question)

		_, err = questionDAO.QuestionByID(1)
		require.EqualError(t, err, "question:1 not found")

		answers := answerDAO.AnswersOfQuestion(1) // The answers associated with this question should be deleted as well
		require.Empty(t, answers)
	})

	t.Run("vote a question", func(t *testing.T) {
		user, err := userDAO.CreateUser("user 2")
		require.NoError(t, err)

		question, err := questionDAO.CreateQuestion(entity.Question{AuthorID: user.ID})
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

		t.Run("The same use replace a vote with the opposite type for the same question", func(t *testing.T) {
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

		t.Run("The same user can't vote the same for the same question", func(t *testing.T) {
			_, err := questionDAO.VoteQuestion(user.ID, question.ID, entity.DownVote())
			require.Equal(t,
				fmt.Sprintf("user:%d has voted DOWN for question:%d", user.ID, question.ID),
				err.Error())
		})
	})
}
