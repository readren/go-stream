package templates

// DO NOT COPY contained lines - BEGIN. They exist to make the compiler happy.

type someTypeE struct{}

// DO NOT COPY contained lines - END

// The type of the elements in the stream
type e_type = someTypeE

func type_equality(e1, e2 e_type) bool {
	return e1 == e2
}

// The type of the stream whose elements are of type `e_type`
type EStream func() (e_type, EStream)

func EStream_Empty() EStream {
	return nil
}

func EStream_Single(e e_type) EStream {
	return func() (e_type, EStream) {
		return e, nil
	}
}

func EStream_Forever(e e_type) EStream {
	return func() (e_type, EStream) {
		return e, EStream_Forever(e)
	}
}

func EStream_FromSlice(slice []e_type) EStream {
	if len(slice) == 0 {
		return nil
	} else {
		return func() (e_type, EStream) {
			return slice[0], EStream_FromSlice(slice[1:])
		}
	}
}

func EStream_FromSet(m map[e_type]bool) EStream {
	slice := make([]e_type, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return EStream_FromSlice(slice)
}

////

func (es EStream) IsEmpty() bool {
	return es == nil
}

func (es EStream) TakeWhile(indexBase int, p func(elem e_type, index int) bool) EStream {
	if es == nil {
		return nil
	} else {
		h, t := es()
		if p(h, indexBase) {
			return func() (e_type, EStream) {
				return h, t.TakeWhile(indexBase+1, p)
			}
		} else {
			return nil
		}
	}
}

func (es EStream) DropWhile(indexBase int, p func(elem e_type, index int) bool) EStream {
	for es != nil {
		h, t := es()
		if !p(h, indexBase) {
			return es
		}
		indexBase += 1
		es = t
	}
	return nil
}

func (es EStream) Filtered(p func(elem e_type) bool) EStream {
	var h e_type
	for es != nil {
		h, es = es()
		if p(h) {
			return func() (e_type, EStream) {
				return h, es.Filtered(p)
			}
		}
	}
	return nil
}

func (es EStream) PrecededBy(a e_type) EStream {
	return func() (e_type, EStream) {
		return a, es
	}
}

func (es EStream) SuccedeedBy(a e_type) EStream {
	return es.FollowedBy(EStream_Single(a))
}

func (es1 EStream) FollowedBy(es2 EStream) EStream {
	if es1 != nil {
		return func() (e_type, EStream) {
			h, t := es1()
			return h, t.FollowedBy(es2)
		}
	} else {
		return es2
	}
}

func (es EStream) ForAll(p func(e_type) bool) bool {
	z := true
	var h e_type
	for es != nil && z {
		h, es = es()
		z = p(h)
	}
	return z
}

func (es EStream) ForAny(p func(e_type) bool) bool {
	return es.ForAll(func(e e_type) bool {
		return !p(e)
	})
}

func (es1 EStream) IsEqualTo(es2 EStream) bool {
	var h1, h2 e_type
	for es1 != nil && es2 != nil {
		h1, es1 = es1()
		h2, es2 = es2()
		if !type_equality(h1, h2) {
			return false
		}
	}
	return es1 == nil && es2 == nil
}

func (es EStream) AppendToSlice(s []e_type) []e_type {
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

func (es EStream) ToSlice(initialCapacity int) []e_type {
	slice := make([]e_type, 0, initialCapacity)
	return es.AppendToSlice(slice)
}