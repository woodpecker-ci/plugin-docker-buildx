package plugin

import (
	"reflect"
	"testing"
)

func Test_stripTagPrefix(t *testing.T) {
	tests := []struct {
		Before string
		After  string
	}{
		{"refs/tags/1.0.0", "1.0.0"},
		{"refs/tags/v1.0.0", "1.0.0"},
		{"v1.0.0", "1.0.0"},
	}

	for _, test := range tests {
		got, want := stripTagPrefix(test.Before), test.After
		if got != want {
			t.Errorf("Got tag %s, want %s", got, want)
		}
	}
}

func TestDefaultTags(t *testing.T) {
	tests := []struct {
		DefaultTag string
		Before     string
		After      []string
	}{
		{"latest", "", []string{"latest"}},
		{"latest", "refs/heads/master", []string{"latest"}},
		{"latest", "refs/tags/0.9.0", []string{"0.9", "0.9.0"}},
		{"latest", "refs/tags/1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"latest", "refs/tags/v1.0.0", []string{"1", "1.0", "1.0.0"}},
		{"latest", "refs/tags/v1.0.0-alpha.1", []string{"1.0.0-alpha.1"}},
	}

	for _, test := range tests {
		tags, err := DefaultTags(test.Before, test.DefaultTag)
		if err != nil {
			t.Error(err)
			continue
		}
		got, want := tags, test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got tag %v, want %v", got, want)
		}
	}
}

func TestDefaultTagsError(t *testing.T) {
	tests := []struct {
		DefaultTag string
		Before     string
	}{
		{
			DefaultTag: "latest",
			Before:     "refs/tags/x1.0.0",
		},
		{
			DefaultTag: "latest",
			Before:     "refs/tags/20190203",
		},
	}

	for _, test := range tests {
		_, err := DefaultTags(test.Before, test.DefaultTag)
		if err == nil {
			t.Errorf("Expect tag error for %s", test)
		}
	}
}

func TestDefaultTagSuffix(t *testing.T) {
	tests := []struct {
		Name       string
		Before     string
		Suffix     string
		After      []string
		DefaultTag string
	}{
		{
			Name:       "Default tag without suffix",
			DefaultTag: "latest",
			After:      []string{"latest"},
		},
		{
			Name:       "Overridden default tag without suffix",
			DefaultTag: "next",
			After:      []string{"next"},
		},
		{
			Name:       "Generate version",
			DefaultTag: "latest",
			Before:     "refs/tags/v1.0.0",
			After: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		{
			Name:       "Generate version with overridden default tag",
			DefaultTag: "next",
			Before:     "refs/tags/v1.0.0",
			After: []string{
				"1",
				"1.0",
				"1.0.0",
			},
		},
		{
			Name:       "Default tag with suffix (linux-amd64)",
			DefaultTag: "latest",
			Suffix:     "linux-amd64",
			After:      []string{"linux-amd64"},
		},
		{
			Name:       "Overridden default tag with suffix (linux-amd64)",
			DefaultTag: "next",
			Suffix:     "linux-amd64",
			After:      []string{"linux-amd64"},
		},
		{
			Name:       "Generate version with suffix (linux-amd64)",
			DefaultTag: "latest",
			Before:     "refs/tags/v1.0.0",
			Suffix:     "linux-amd64",
			After: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
		{
			Name:       "Generate version with suffix (linux-amd64) and overridden default tag (next)",
			DefaultTag: "next",
			Before:     "refs/tags/v1.0.0",
			Suffix:     "linux-amd64",
			After: []string{
				"1-linux-amd64",
				"1.0-linux-amd64",
				"1.0.0-linux-amd64",
			},
		},
		{
			Name:       "Default tag with suffix (nanoserver)",
			DefaultTag: "latest",
			Suffix:     "nanoserver",
			After:      []string{"nanoserver"},
		},
		{
			Name:       "Overridden default tag with suffix (nanoserver)",
			DefaultTag: "next",
			Suffix:     "nanoserver",
			After:      []string{"nanoserver"},
		},
		{
			Name:       "Generate version with suffix (nanoserver)",
			DefaultTag: "latest",
			Before:     "refs/tags/v1.9.2",
			Suffix:     "nanoserver",
			After: []string{
				"1-nanoserver",
				"1.9-nanoserver",
				"1.9.2-nanoserver",
			},
		},
		{
			Name:       "Generate version with suffix (nanoserver) and overridden default tag (next)",
			DefaultTag: "latest",
			Before:     "refs/tags/v1.9.2",
			Suffix:     "nanoserver",
			After: []string{
				"1-nanoserver",
				"1.9-nanoserver",
				"1.9.2-nanoserver",
			},
		},
		{
			Name:       "Generate version with suffix (zero-padded version)",
			DefaultTag: "latest",
			Before:     "refs/tags/v18.06.0",
			Suffix:     "nanoserver",
			After: []string{
				"18-nanoserver",
				"18.06-nanoserver",
				"18.06.0-nanoserver",
			},
		},
		{
			Name:       "Generate version with suffix (zero-padded version) with overridden default tag (next)",
			DefaultTag: "next",
			Before:     "refs/tags/v18.06.0",
			Suffix:     "nanoserver",
			After: []string{
				"18-nanoserver",
				"18.06-nanoserver",
				"18.06.0-nanoserver",
			},
		},
	}

	for _, test := range tests {
		tag, err := DefaultTagSuffix(test.Before, test.DefaultTag, test.Suffix)
		if err != nil {
			t.Error(err)
			continue
		}
		got, want := tag, test.After
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%q. Got tag %v, want %v", test.Name, got, want)
		}
	}
}

func Test_stripHeadPrefix(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				ref: "refs/heads/master",
			},
			want: "master",
		},
	}
	for _, tt := range tests {
		if got := stripHeadPrefix(tt.args.ref); got != tt.want {
			t.Errorf("stripHeadPrefix() = %v, want %v", got, tt.want)
		}
	}
}

func TestUseDefaultTag(t *testing.T) {
	type args struct {
		ref           string
		defaultBranch string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "latest tag for default branch",
			args: args{
				ref:           "refs/heads/master",
				defaultBranch: "master",
			},
			want: true,
		},
		{
			name: "build from tags",
			args: args{
				ref:           "refs/tags/v1.0.0",
				defaultBranch: "master",
			},
			want: true,
		},
		{
			name: "skip build for not default branch",
			args: args{
				ref:           "refs/heads/develop",
				defaultBranch: "master",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		if got := UseDefaultTag(tt.args.ref, tt.args.defaultBranch); got != tt.want {
			t.Errorf("%q. UseDefaultTag() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_isSingleTag(t *testing.T) {
	tests := []struct {
		Tag     string
		IsValid bool
	}{
		{"latest", true},
		{" latest", false},
		{"LaTest__Hi", true},
		{"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ__.-0123456789", true},
		{"_wierd.but-ok1", true},
		{"latest ", false},
		{"latest,next", false},
		// more tests to be added, once the validation is more powerful
	}

	for _, test := range tests {
		valid := isSingleTag(test.Tag)
		if valid != test.IsValid {
			t.Errorf("Tag verification '%s' tag %v, want %v", test.Tag, valid, test.IsValid)
		}
	}
}
