package stream

// The type of the elements in the stream
type e_string = string

// The type of the stream itself
type String func() (e_string, String)

func String_Empty() String {
	return nil
}
func String_Unit(a e_string) String {
	return func() (e_string, String) {
		return a, nil
	}
}
func String_FromSlice(slice []e_string) String {
	if len(slice) == 0 {
		return nil
	} else {
		return func() (e_string, String) {
			return slice[1], String_FromSlice(slice[1:])
		}
	}
}
func String_FromSet(m map[e_string]bool) String {
	slice := make([]e_string, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return String_FromSlice(slice)
}

////

func (es String) IsEmpty() bool {
	return es == nil
}

func (es String) Filtered(p func(e_string) bool) String {
	if es == nil {
		return nil
	} else {
		h, t := es()
		pass := p(h)
		for !pass && t != nil {
			h, t = t()
			pass = p(h)
		}
		if pass {
			return t.Filtered(p).PrecededBy(h)
		} else {
			return nil
		}
	}
}

func (es String) PrecededBy(a e_string) String {
	return func() (e_string, String) {
		return a, es
	}
}
func (es String) SuccedeedBy(a e_string) String {
	return es.FollowedBy(String_Unit(a))
}
func (as1 String) FollowedBy(as2 String) String {
	if as1 != nil {
		h, t := as1()
		return func() (e_string, String) { return h, t.FollowedBy(as2) }
	} else {
		return as2
	}
}
func (es1 String) IsEqualTo(es2 String) bool {
	if es1 == nil {
		return es2 == nil
	} else if es2 == nil {
		return false
	} else {
		h1, t1 := es1()
		h2, t2 := es2()
		return h1 == h2 && t1.IsEqualTo(t2)
	}
}

func (es String) AppendToSlice(s []e_string) []e_string {
	if es != nil {
		h, t := es()
		// All the following lines could be replaced by this << return t.AppendToSlice(append(s, h)) >> if the golang compiler supported tail recursion optimization.
		s = append(s, h)
		for t != nil {
			h, t = t()
			s = append(s, h)
		}
	}
	return s
}

func (es String) ToSlice(initialCapacity int) []e_string {
	slice := make([]e_string, 0, initialCapacity)
	return es.AppendToSlice(slice)
}
