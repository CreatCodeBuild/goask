package entity

type Question struct {
	ID       ID
	Content  string
	AuthorID ID // The ID of the user who created it.
	Title    string
}

type QuestionUpdate struct {
	ID      ID
	Title   *string
	Content *string
}

type Answer struct {
	ID         ID
	QuestionID ID
	Content    string
	AuthorID   ID // The ID of the user who created it.
	Accepted   bool
}

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

type User struct {
	ID   ID
	Name string
}

type ID int64

func (id ID) Equal(i int32) bool {
	return id == ID(i)
}
