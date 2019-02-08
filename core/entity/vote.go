package entity

// Vote is the vote of a question.
type Vote struct {
	UserID     ID
	QuestionID ID
	Type       VoteType // UP or DOWN
}

type VoteType string

func UpVote() VoteType {
	return "UP"
}

func DownVote() VoteType {
	return "DOWN"
}
