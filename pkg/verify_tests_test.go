package utopia

import "testing"

func TestVerifyTests(t *testing.T) {
	type args struct {
		directory string
		services  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid tests",
			args: args{
				directory: "testdata/verify_tests",
				services:  []string{"make-tests-valid"},
			},
			wantErr: false,
		},
		{
			name: "broken tests",
			args: args{
				directory: "testdata/verify_tests",
				services:  []string{"make-tests-fail"},
			},
			wantErr: true,
		},
		{
			name: "missing tests",
			args: args{
				directory: "testdata/verify_tests",
				services:  []string{"make-tests-missing"},
			},
			wantErr: false,
		},
		{
			name: "missing services",
			args: args{
				directory: "testdata/verify_tests",
				services:  []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyTests(tt.args.directory, tt.args.services); (err != nil) != tt.wantErr {
				t.Errorf("VerifyTests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
