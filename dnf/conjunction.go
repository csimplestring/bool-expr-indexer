package dnf

// Attribute is the pair of key-value, e.g., age:10, representing 'belongs to'
// The value here is discrete, but for range values, like age < 40, we can convert it into multiple pairs:
// age < 40 ---> age = 10, age = 20, age = 30, the granularity is 10.
//
// Name and Values are internally represented as int value, to save memory space, there should be a mapping
// between the integer and human readable string, like:
// Name: browser -> 1, os -> 2
// Value: Chrome -> 1, Safari -> 2
//
// TODO: for longer range values such as dates, using a hierarchy structure to convert
type Attribute struct {
	Name     int
	Values   []int
	Contains bool
}

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
