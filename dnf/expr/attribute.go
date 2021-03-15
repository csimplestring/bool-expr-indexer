package expr

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
	Name     string
	Values   []string
	Contains bool
}
