package dnf

// Attribute is the pair of key-value, e.g., age:10, representing 'belongs to'
// The value here is discrete, but for range values, like age < 40, we can convert it into multiple pairs:
// age < 40 ---> age = 10, age = 20, age = 30, the granularity is 10.
// TODO: for longer range values such as dates, using a hierarchy structure to convert
type Attribute struct {
	Name     string
	Values   []string
	Contains bool
}

// Conjunction consists of a slice of Attributes, which are combined with 'AND' logic.
type Conjunction struct {
	ID         int
	Attributes []*Attribute
	kSize      int
}

// GetKSize is the size of attributes in c, excluding any 'not-included' type attribute
func (c *Conjunction) GetKSize() int {
	return c.kSize
}

// NewConjunction creates a new Conjunction
func NewConjunction(ID int, attrs []*Attribute) *Conjunction {
	ksize := 0
	for _, a := range attrs {
		if a.Contains {
			ksize++
		}
	}
	return &Conjunction{
		ID:         ID,
		Attributes: attrs,
		kSize:      ksize,
	}
}
