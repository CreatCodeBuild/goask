package entity

// Tag represents a string that can be attached to a question to indicate the category.
type Tag string

type TagSet map[Tag]struct{}

func NewTagSet(elements ...Tag) TagSet {
	set := TagSet{}
	set.Add(elements...)
	return set
}

func (set TagSet) Contains(element Tag) bool {
	_, ok := set[element]
	return ok
}

func (set TagSet) Add(elements ...Tag) {
	for _, element := range elements {
		set[element] = struct{}{}
	}
}

func (set TagSet) Remove(elements ...Tag) {
	for _, element := range elements {
		delete(set, element)
	}
}

func (set TagSet) Union(set2 TagSet) TagSet {
	set3 := make(TagSet)
	set3.Add(set.Slice()...)
	set3.Add(set2.Slice()...)
	return set3
}

// Or returns a set set3 such that for each element in set3, this element is in either set or set2.
func (set TagSet) Or(set2 TagSet) TagSet {
	return set.Union(set2)
}

func (set TagSet) Intersection(set2 TagSet) TagSet {
	set3 := make(TagSet)
	for element := range set2 {
		if set.Contains(element) {
			set3.Add(element)
		}
	}
	return set3
}

// And returns a set set3 such that for each element in set3, this element is in both set and set2.
func (set TagSet) And(set2 TagSet) TagSet {
	return set.Intersection(set2)
}

func (set TagSet) Slice() []Tag {
	set2 := make([]Tag, len(set))
	i := 0
	for element := range set {
		set2[i] = element
		i++
	}
	return set2
}
