package expr

// Conjunction consists of a slice of Attributes, which are combined with 'AND' logic.
type Conjunction struct {
	ID         int
	Attributes []*Attribute
}

// GetKSize is the size of attributes in c, excluding any 'not-included' type attribute
func (c *Conjunction) GetKSize() int {
	ksize := 0
	for _, a := range c.Attributes {
		if a.Contains {
			ksize++
		}
	}
	return ksize
}

// NewConjunction creates a new Conjunction
func NewConjunction(ID int, attrs []*Attribute) *Conjunction {

	return &Conjunction{
		ID:         ID,
		Attributes: attrs,
	}
}
