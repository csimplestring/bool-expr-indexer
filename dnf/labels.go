package dnf

// Label is a simple k/v pair: like <age:30>
type Label struct {
	Name  int
	Value int
}

// Assignment is a slice of Label, equals to 'assignment S' in the paper
type Assignment []Label
