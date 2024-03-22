package domain

type Labels []string

func (l Labels) Contains(label string) bool {
	for _, v := range l {
		if v == label {
			return true
		}
	}
	return false
}

type ParsedEvent struct {
	Branch string
	Labels Labels
}
