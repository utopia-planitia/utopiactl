package utopia

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVerifyTests(t *testing.T) {
	type args struct {
		directory       string
		services        []string
		testAllServices bool
	}
	tests := []struct {
		name    string
		args    args
		want    []error
		wantErr bool
	}{
		{
			name: "valid tests",
			args: args{
				directory:       "testdata/verify_tests",
				services:        []string{"make-tests-valid"},
				testAllServices: false,
			},
			want:    []error{},
			wantErr: false,
		},
		{
			name: "broken tests",
			args: args{
				directory:       "testdata/verify_tests",
				services:        []string{"make-tests-fail"},
				testAllServices: false,
			},
			want:    []error{},
			wantErr: true,
		},
		{
			name: "missing tests",
			args: args{
				directory:       "testdata/verify_tests",
				services:        []string{"make-tests-missing"},
				testAllServices: false,
			},
			want:    []error{},
			wantErr: false,
		},
		{
			name: "missing services",
			args: args{
				directory:       "testdata/verify_tests",
				services:        []string{},
				testAllServices: false,
			},
			want:    []error{},
			wantErr: true,
		},
		{
			name: "test all services",
			args: args{
				directory:       "testdata/verify_tests",
				services:        []string{"make-tests-fail", "make-tests-missing", "make-tests-fail2", "make-tests-valid"},
				testAllServices: true,
			},
			want: []error{
				fmt.Errorf("service make-tests-fail tests failed: exit status 2"),
				fmt.Errorf("service make-tests-fail2 tests failed: exit status 2"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyTests(tt.args.directory, tt.args.services, tt.args.testAllServices)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyTests() = %v, want %v", got, tt.want)
			}
		})
	}
}
