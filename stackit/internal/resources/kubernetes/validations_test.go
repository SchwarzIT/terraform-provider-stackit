package kubernetes

import (
	"testing"

	"github.com/Masterminds/semver"
)

func Test_maxVersionOption(t *testing.T) {
	type args struct {
		version        *semver.Version
		versionOptions []*semver.Version
	}
	tests := []struct {
		name string
		args args
		want *semver.Version
	}{
		{name: "patch version", args: args{version: semver.MustParse("1.18.0"), versionOptions: []*semver.Version{
			semver.MustParse("1.18.0"),
			semver.MustParse("1.18.1"),
		}}, want: semver.MustParse("1.18.0")},
		{name: "minor version", args: args{version: semver.MustParse("1.18"), versionOptions: []*semver.Version{
			semver.MustParse("1.18.1"),
			semver.MustParse("1.18.0"),
		}}, want: semver.MustParse("1.18.1")},
		{name: "minor version differs", args: args{version: semver.MustParse("1.18"), versionOptions: []*semver.Version{
			semver.MustParse("1.18.1"),
			semver.MustParse("1.18.0"),
			semver.MustParse("1.19.0"),
		}}, want: semver.MustParse("1.18.1")},
		{name: "regression", args: args{version: semver.MustParse("1.23"), versionOptions: []*semver.Version{
			semver.MustParse("1.23.3"),
			semver.MustParse("1.23.2"),
			semver.MustParse("1.23.1"),
		}}, want: semver.MustParse("1.23.3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint, err := toVersionConstraint(tt.args.version)
			if err != nil {
				t.Fatalf("toVersionConstraint() error = %v", err)
			}
			if got := maxVersionOption(constraint, tt.args.versionOptions); *got != *tt.want {
				t.Errorf("maxVersionOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
