package fakeadapter

import (
	"encoding/json"
	"goask/core/adapter"
	"io/ioutil"
	"os"
	"strings"

	"goask/core/entity"

	"github.com/pkg/errors"
)

func match(s1, s2 string) bool {
	return strings.Contains(s1, s2)
}

// Serializer: here in this file the implementer is defined in the same pkg. This is a bad practice. But it is ok for now.
type Serializer interface {
	Serialize([]byte) error
	Deserialize() ([]byte, error)
}

// FileSerializer implements Serializer with file io.
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
	if _, err := os.Stat(f.fileName); os.IsNotExist(err) {
		f, err := os.Create(f.fileName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		err = f.Close()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
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

// Data holds all data.
type Data struct {
	questions     Questions
	questionVotes QuestionVotes
	answers       Answers
	users         Users
	tags          Tags
	storage       Serializer
}

func NewData(storage Serializer) (*Data, error) {
	d := &Data{}
	if storage == nil {
		return nil, errors.New("storage == nil")
	}
	d.storage = storage
	d.questionVotes = make(QuestionVotes)
	d.tags = Tags{
		tags: make(questionTagSet),
	}

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

type Users []entity.User

func (u *Users) Get(userID entity.ID) (entity.User, bool) {
	for _, qu := range *u {
		if qu.ID == userID {
			return qu, true
		}
	}
	return entity.User{}, false
}

type Questions []entity.Question

func (a *Questions) Filter(f func(entity.Question) bool) Questions {
	var ret Questions
	for _, an := range *a {
		if f(an) {
			ret = append(ret, an)
		}
	}
	return ret
}

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
		ID:         entity.NewIDInt(len(*a) + 1),
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

func (a *Answers) Delete(answerID entity.ID) {

	answers := a.Filter(func(a entity.Answer) bool {
		return a.ID != answerID
	})

	*a = answers
}

type voteTarget struct {
	userID, questionID entity.ID
}

type QuestionVotes map[voteTarget]entity.VoteType

func (a *QuestionVotes) Update(userID, questionID entity.ID, voteType entity.VoteType) (entity.Vote, error) {
	m := *a

	vt := voteTarget{
		userID:     userID,
		questionID: questionID,
	}

	vote := entity.Vote{
		UserID:     userID,
		QuestionID: questionID,
		Type:       voteType,
	}

	if vType, ok := m[vt]; ok {
		if vType == voteType {
			return entity.Vote{}, errors.WithStack(&adapter.ErrVoteDuplicated{Vote: vote})
		}
		m[vt] = voteType
	}

	m[vt] = voteType
	return vote, nil
}

func (a *QuestionVotes) Count(questionID entity.ID) (up int, down int, err error) {
	up = 0
	down = 0

	for vTarget, vType := range *a {
		if vTarget.questionID == questionID {
			if vType == entity.UpVote() {
				up++
			} else if vType == entity.DownVote() {
				down++
			} else {
				err = errors.Errorf("impossible, vote type is %s", vType)
				return
			}
		}
	}
	return
}
