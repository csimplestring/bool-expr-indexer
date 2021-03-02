package cnf

// Attribute is the pair of key-value, e.g., age:10, representing 'belongs to'
// The value here is discrete, but for range values, like age < 40, we can convert it into multiple pairs:
// age < 40 ---> age = 10, age = 20, age = 30, the granularity is 10.
// TODO: for longer range values such as dates, using a hierarchy structure to convert
type Attribute struct {
	Name    string
	Values  []string
	Include bool
}

// Disjunction consists of a slice of Attributes, which are combined with 'OR' logic.
type Disjunction struct {
	Attributes []*Attribute
}

func (d *Disjunction) hasExclude() bool {
	for _, a := range d.Attributes {
		if !a.Include {
			return true
		}
	}
	return false
}

// NewDisjunction creates a new Disjunction
func NewDisjunction(attrs []*Attribute) *Disjunction {
	return &Disjunction{
		Attributes: attrs,
	}
}

// CNF consists of a slice of Disjunction, combined with 'AND' logic
// CNF indexing algorithm does not need to split CNF into disjunctions, but be processed as a whole.
type CNF struct {
	ID           int64
	Disjunctions []*Disjunction
}

// NewCNF creates a new CNF
func NewCNF(ID int64, Disjunctions []*Disjunction) *CNF {
	return &CNF{
		ID:           ID,
		Disjunctions: Disjunctions,
	}
}

func (c *CNF) getKSize() int {
	kSize := 0
	for _, d := range c.Disjunctions {
		if !d.hasExclude() {
			kSize++
		}
	}
	return kSize
}
