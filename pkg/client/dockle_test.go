package client_test

import (
	"fmt"
	"kube-dockle-exporter/pkg/client"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDockleResponseExtractImage(t *testing.T) {
	type want struct {
		first string
	}

	tests := []struct {
		name         string
		receiver     *client.DockleResponse
		want         want
		optsFunction func(interface{}) cmp.Option
	}{
		{
			func() string {
				_, _, line, _ := runtime.Caller(1)
				return fmt.Sprintf("L%d", line)
			}(),
			&client.DockleResponse{
				Target: "fake (fake)",
			},
			want{
				"fake",
			},
			func(got interface{}) cmp.Option {
				return nil
			},
		},
	}
	for _, tt := range tests {
		name := tt.name
		receiver := tt.receiver
		want := tt.want
		optsFunction := tt.optsFunction
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := receiver.ExtractImage()
			if diff := cmp.Diff(want.first, got, optsFunction(got)); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
