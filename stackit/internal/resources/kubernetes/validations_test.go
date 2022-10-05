package kubernetes

import (
	"github.com/SchwarzIT/community-stackit-go-client/pkg/api/v1/kubernetes/options"
	"testing"
)

func Test_maxVersionOption(t *testing.T) {
	type args struct {
		version        string
		versionOptions []options.KubernetesVersion
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "patch version", args: args{version: "1.18.0", versionOptions: []options.KubernetesVersion{
			{Version: "1.18.0"},
			{Version: "1.18.1"},
		}}, want: "1.18.0"},
		{name: "minor version", args: args{version: "1.18", versionOptions: []options.KubernetesVersion{
			{Version: "1.18.1"},
			{Version: "1.18.0"},
		}}, want: "1.18.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxVersionOption(tt.args.version, tt.args.versionOptions); got != tt.want {
				t.Errorf("maxVersionOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
