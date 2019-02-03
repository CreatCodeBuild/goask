package adaptertest

import (
	"goask/core/adapter"
	"goask/core/entity"
	"testing"

	"github.com/stretchr/testify/require"
)

func Data(t *testing.T, data adapter.Data) {

	t.Run("create questions", func(t *testing.T) {
		_, err := data.CreateQuestion(entity.Question{})
		require.EqualError(t, err, "user:0 not found")

		user, err := data.CreateUser("user 1")
		require.NoError(t, err)
		require.Equal(t, user.Name, "user 1")
		require.Equal(t, user.ID, entity.ID(1))

		question, err := data.CreateQuestion(entity.Question{AuthorID: 1})
		require.NoError(t, err)
		require.Equal(t, entity.Question{AuthorID: 1, ID: 1}, question)
	})

	t.Run("create answers", func(t *testing.T) {
		answer, err := data.CreateAnswer(1, "answer 1", 1)
		require.NoError(t, err)
		require.Equal(t, entity.Answer{ID: 1, QuestionID: 1, AuthorID: 1, Content: "answer 1"}, answer)

		t.Run("accept answers", func(t *testing.T) {
			acceptedAnswer, err := data.AcceptAnswer(answer.ID, answer.AuthorID)
			require.NoError(t, err)
			require.Equal(t,
				entity.Answer{ID: 1, QuestionID: 1, AuthorID: 1, Content: "answer 1", Accepted: true},
				acceptedAnswer)

			_, err = data.AcceptAnswer(answer.ID, -1)
			require.EqualError(t, err, "user:-1 is no the author of question:1")

			t.Run("delete answers", func(t *testing.T) {
				deletedAnswer, err := data.DeleteAnswer(answer.ID, answer.AuthorID)
				require.NoError(t, err)
				require.Equal(t, deletedAnswer, acceptedAnswer)
			})
		})
	})

	t.Run("delete questions", func(t *testing.T) {
		_, err := data.DeleteQuestion(2, 1)
		require.EqualError(t, err, "user:2 not found")

		user, err := data.CreateUser("user 2")
		require.NoError(t, err)

		_, err = data.DeleteQuestion(user.ID, 1)
		require.EqualError(t, err, "user:2 is not authorized to delete question:1")

		question, err := data.DeleteQuestion(1, 1)
		require.NoError(t, err)
		require.Equal(t, entity.Question{AuthorID: 1, ID: 1}, question)

		_, err = data.QuestionByID(1)
		require.EqualError(t, err, "question:1 not found")

		answers := data.AnswersOfQuestion(1) // The answers associated with this question should be deleted as well
		require.Empty(t, answers)
	})
}
