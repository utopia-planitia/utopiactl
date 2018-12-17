package utopia

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMakefile(t *testing.T) {

	err := generateMakefile("testdata/makefile")
	if err != nil {
		t.Errorf("failed to generate makefile: %v", err)
		return
	}

	golden, err := ioutil.ReadFile("testdata/makefile/Makefile.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
		return
	}

	result, err := ioutil.ReadFile("testdata/makefile/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
		return
	}
	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
		return
	}

	defer os.Remove("testdata/makefile/Makefile")
}

func TestMakefileEmpty(t *testing.T) {

	err := generateMakefile("testdata/makefile-empty")
	if err != nil {
		t.Errorf("failed to generate makefile: %v", err)
		return
	}

	golden, err := ioutil.ReadFile("testdata/makefile-empty/Makefile.golden")
	if err != nil {
		t.Errorf("failed to read golden state: %v", err)
		return
	}

	result, err := ioutil.ReadFile("testdata/makefile-empty/Makefile")
	if err != nil {
		t.Errorf("failed to read result: %v", err)
		return
	}
	if bytes.Compare(result, golden) != 0 {
		t.Errorf("Makefile was incorrect, got: %+s, want: %+s.", result, golden)
		return
	}

	defer os.Remove("testdata/makefile-empty/Makefile")
}

func Test_moveStorageToFirst(t *testing.T) {
	type args struct {
		services []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "storage is missing",
			args: args{
				services: []string{"a", "bb", "cee"},
			},
			want: []string{"a", "bb", "cee"},
		},
		{
			name: "storage is in the beginning",
			args: args{
				services: []string{"storage", "a", "bb", "cee"},
			},
			want: []string{"storage", "a", "bb", "cee"},
		},
		{
			name: "storage is at the end",
			args: args{
				services: []string{"a", "bb", "cee", "storage"},
			},
			want: []string{"storage", "a", "bb", "cee"},
		},
		{
			name: "storage is in between",
			args: args{
				services: []string{"a", "bb", "storage", "cee"},
			},
			want: []string{"storage", "a", "bb", "cee"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := moveStorageToFirst(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("moveStorageToFirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
