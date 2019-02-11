package fakeadapter

import (
	"goask/core/entity"
)

type questionTag struct {
	tag        entity.Tag
	questionID entity.ID
}

type questionTagSet map[questionTag]struct{}

func (set questionTagSet) Contains(element questionTag) bool {
	_, ok := set[element]
	return ok
}

func (set questionTagSet) Add(elements ...questionTag) {
	for _, element := range elements {
		set[element] = struct{}{}
	}
}

func (set questionTagSet) Remove(elements ...questionTag) {
	for _, element := range elements {
		delete(set, element)
	}
}

func (set questionTagSet) Union(set2 questionTagSet) questionTagSet {
	set3 := make(questionTagSet)
	set3.Add(set.Slice()...)
	set3.Add(set2.Slice()...)
	return set3
}

// Or returns a set set3 such that for each element in set3, this element is in either set or set2.
func (set questionTagSet) Or(set2 questionTagSet) questionTagSet {
	return set.Union(set2)
}

func (set questionTagSet) Intersection(set2 questionTagSet) questionTagSet {
	set3 := make(questionTagSet)
	for element := range set2 {
		if set.Contains(element) {
			set3.Add(element)
		}
	}
	return set3
}

// And returns a set set3 such that for each element in set3, this element is in both set and set2.
func (set questionTagSet) And(set2 questionTagSet) questionTagSet {
	return set.Intersection(set2)
}

func (set questionTagSet) Slice() []questionTag {
	set2 := make([]questionTag, len(set))
	i := 0
	for element := range set {
		set2[i] = element
		i++
	}
	return set2
}

type Tags struct {
	tags questionTagSet
}

// GetQuestionIDs returns all questionIDs this tag is associated with.
// It turns nil if the tag is assocaited with 0 questionID or the tag doesn't exist.
// The implementation treat them the same.
func (t *Tags) GetQuestionIDs(tag entity.Tag) []entity.ID {
	var questionIDs []entity.ID
	for qt := range t.tags {
		if qt.tag == tag {
			questionIDs = append(questionIDs, qt.questionID)
		}
	}
	return questionIDs
}

func (t *Tags) GetTagsOfQuestion(questionID entity.ID) entity.TagSet {
	tags := entity.TagSet{}
	for qt := range t.tags {
		if qt.questionID == questionID {
			tags.Add(qt.tag)
		}
	}
	return tags
}

// UpdateQuestion replaces the tags of given question.
func (t *Tags) UpdateQuestion(questionID entity.ID, tags []entity.Tag) {
	// This implementation is not elegant
	tagSet := entity.TagSet{}
	tagSet.Add(tags...)

	toRemove := questionTagSet{}
	for qt := range t.tags {
		if qt.questionID == questionID {
			if !tagSet.Contains(qt.tag) {
				toRemove.Add(qt)
			}
		}
	}
	t.tags.Remove(toRemove.Slice()...)

	toAdd := questionTagSet{}
	for tag := range tagSet {
		toAdd.Add(questionTag{tag: tag, questionID: questionID})
	}
	t.tags.Add(toAdd.Slice()...)
}
