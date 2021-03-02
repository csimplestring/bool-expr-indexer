package dnf

// Label is a simple k/v pair: like <age:30>
type Label struct {
	Name  string
	Value string
}

// Assignment is a slice of Label, equals to 'assignment S' in the paper
type Assignment []Label
