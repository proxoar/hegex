package hegex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHegex(t *testing.T) {
	type args struct {
		expression string
		s          string
		template   string
	}
	type want struct {
		matched     bool
		substituted string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"1", args{"{site}.example.com", "book.example.com", "/my-{site}"}, want{true, "/my-book"}},
		{"1.1", args{"{site}.example.com", ".example.com", "/my-{site}"}, want{false, ""}},
		{"1.2", args{"{site}.example.com", "example.com", "/my-{site}"}, want{false, ""}},
		{"1.3", args{"{site}.example.com", "", "/my-{site}"}, want{false, ""}},
		{"2", args{"{site}.example.com", "example.com", "/my-{site}"}, want{false, ""}},
		{"2", args{"{site}.example.com", "a.b.example.com", "/my-{site}"}, want{false, ""}},
		{"3", args{"{site[x|y]}.example.com", "x.example.com", "/my-{site}"}, want{true, "/my-x"}},
		{"3.1", args{"{site[xxa|yya|abc]}.example.com", "abc.example.com", "/my-{site}"}, want{true, "/my-abc"}},
		{"4", args{"{site[x|y]}.example.com", "b.example.com", "/my-{site}"}, want{false, ""}},
		{"4", args{"{}.example.com", "xy.example.com", "/{}z"}, want{true, "/xyz"}},
		{"5", args{"*.example.com", ".example.com", "/my-*"}, want{true, "/my-"}},
		{"6", args{"*.example.com", ".example.com", "/my-*-*"}, want{true, "/my--"}},
		{"7", args{"*.example.com", "abc.example.com", "/my-*"}, want{true, "/my-abc"}},
		{"8", args{"*.example.com", "abc.example.com", "/my-*-*"}, want{true, "/my-abc-abc"}},
		{"9", args{"*.example.com", "a.b.c.example.com", "/my/*"}, want{true, "/my/a.b.c"}},
		{"10", args{"*.example.com", "example.com", "/my/*"}, want{false, ""}},
		{"11", args{"*.example.com", "example.x.com", "/my/*"}, want{false, ""}},
		{"12", args{"*.example.com", "", "/my/*"}, want{false, ""}},

		{"a1", args{"/{media}/size", "/text/size", "/my-{media}"}, want{true, "/my-text"}},
		{"a1", args{"/{media}/size", "/small-video/size", "/my-{media}"}, want{true, "/my-small-video"}},
		{"a2", args{"/{media}/size", "//size", "/my-{media}"}, want{false, ""}},
		{"a2", args{"/{media}/size", "/size", "/my-{media}"}, want{false, ""}},
		{"a2", args{"/{media}/size", "/a/bsize", "/my-{media}"}, want{false, ""}},
		{"a3", args{"/{media[video|image]}/size", "/video/size", "/my-{media}"}, want{true, "/my-video"}},
		{"a3", args{"/{media[video|image]}/size", "/jpeg/size", "/my-{media}"}, want{false, ""}},
		{"a3", args{"/{media[video|image]}/size", "//size", "/my-{media}"}, want{false, ""}},
		{"a3", args{"/{media[video|image]}/size", "/size", "/my-{media}"}, want{false, ""}},
		{"a5", args{"/*/size", "//size", "/my-*"}, want{true, "/my-"}},
		{"a6", args{"/*/size", "//size", "/my-*-*"}, want{true, "/my--"}},
		{"a7", args{"/*/size", "/abc/size", "/my-*"}, want{true, "/my-abc"}},
		{"a8", args{"/*/size", "/abc/size", "/my-*-*"}, want{true, "/my-abc-abc"}},
		{"a9", args{"/*/size", "/a/b/c/size", "/my/*"}, want{true, "/my/a/b/c"}},
		{"a10", args{"/*/size", "/size", "/my/*"}, want{false, ""}},
		{"a11", args{"/*/size", "/anysize", "/my/*"}, want{false, ""}},
		{"a12", args{"/*/size", "", "/my/*"}, want{false, ""}},
		{"a13", args{"/*size", "/size", "/my/*"}, want{true, "/my/"}},
		{"a14", args{"/*size", "/0size", "/my/*"}, want{true, "/my/0"}},
		{"a15", args{"/*", "/yes", "/my/*"}, want{true, "/my/yes"}},
		{"a16", args{"/*", "/yes/no", "/my/*"}, want{true, "/my/yes/no"}},
		{"a17", args{"/size*", "/size", "/my/*"}, want{true, "/my/"}},
		{"a18", args{"/size*", "/size1900", "/my/*"}, want{true, "/my/1900"}},
		{"a19", args{"/size*", "/size1900/2000", "/my/*"}, want{true, "/my/1900/2000"}},

		{"b1", args{"/{name[ab|cd]}/*", "/ab/health", "/*"}, want{true, "/health"}},
		{"b2", args{"/home/assets//.yml", "/home/assets//.yml", ""}, want{true, ""}},
		{"b2", args{"/home/assets/*/.yml", "/home/assets//.yml", "/*"}, want{true, "/"}},
		{"b2", args{"/home/assets/*/.yml", "/home/assets//.yml", "/home/assets/config.yml"}, want{true, "/home/assets/config.yml"}},

		// example in README
		{"c1", args{"/*/**", "/path/data", "/*/to/**"}, want{true, "/path/to/data"}},
		{"c1", args{"/*/**", "/path/to/data", "/*/to/**"}, want{true, "/path/to/to/data"}},
		{"c1", args{"/***/abc/*/**def", "/s03/abc/s01/s02def", "/*/to/**/***"}, want{true, "/s01/to/s02/s03"}},
		{"c1", args{"/***/a", "/a", ""}, want{false, ""}},
		{"c1", args{"/***/a", "//a", ""}, want{true, ""}},
		{"c2", args{"/path/*/api.{postfix[json|yml]}", "/path/doc/api.json", "{postfix}-*"}, want{true, "json-doc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compile := MustCompile(tt.args.expression)
			assert.Equal(t, tt.want.matched, compile.MatchString(tt.args.s))
			substituted, ok := compile.MatchAndSubstitute(tt.args.s, tt.args.template)
			assert.Equal(t, tt.want.matched, ok)
			if ok {
				assert.Equal(t, tt.want.substituted, substituted)
			} else {
				assert.Empty(t, substituted)
			}
		})
	}
}
