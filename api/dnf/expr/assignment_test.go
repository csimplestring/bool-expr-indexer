package expr

import "testing"

func TestValidateAssignment(t *testing.T) {
	type args struct {
		a Assignment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{a: Assignment{
				{Name: "a", Value: "foo"},
			}},
			false,
		},
		{
			"",
			args{a: Assignment{
				{Name: "a", Value: "foo"},
				{Name: "b", Value: "bar"},
				{Name: "b", Value: "coo"},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateAssignment(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAssignment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
