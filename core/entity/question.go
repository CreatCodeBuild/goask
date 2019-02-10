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
	Tags    []Tag
}

type QuestionSet map[Question]struct{}

func (set QuestionSet) Contains(element Question) bool {
	_, ok := set[element]
	return ok
}

func (set QuestionSet) Add(elements ...Question) {
	for _, element := range elements {
		set[element] = struct{}{}
	}
}

func (set QuestionSet) Remove(elements ...Question) {
	for _, element := range elements {
		delete(set, element)
	}
}

func (set QuestionSet) Union(set2 QuestionSet) QuestionSet {
	set3 := make(QuestionSet)
	set3.Add(set.Slice()...)
	set3.Add(set2.Slice()...)
	return set3
}

// Or returns a set set3 such that for each element in set3, this element is in either set or set2.
func (set QuestionSet) Or(set2 QuestionSet) QuestionSet {
	return set.Union(set2)
}

func (set QuestionSet) Intersection(set2 QuestionSet) QuestionSet {
	set3 := make(QuestionSet)
	for element := range set2 {
		if set.Contains(element) {
			set3.Add(element)
		}
	}
	return set3
}

// And returns a set set3 such that for each element in set3, this element is in both set and set2.
func (set QuestionSet) And(set2 QuestionSet) QuestionSet {
	return set.Intersection(set2)
}

func (set QuestionSet) Slice() []Question {
	set2 := make([]Question, len(set))
	i := 0
	for element := range set {
		set2[i] = element
		i++
	}
	return set2
}
