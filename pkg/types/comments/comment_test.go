package comments

import (
	"reflect"
	"testing"

	markers2 "github.com/crossplane-contrib/terrajet/pkg/types/markers"

	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/crossplane-contrib/terrajet/pkg/config"
)

func TestComment_Build(t *testing.T) {
	tftag := "-"
	type args struct {
		text string
		opts []Option
	}
	type want struct {
		out   string
		mopts markers2.Options
		err   error
	}

	cases := map[string]struct {
		args
		want
	}{
		"OnlyTextNoMarker": {
			args: args{
				text: "hello world!",
			},
			want: want{
				out:   "// hello world!\n",
				mopts: markers2.Options{},
			},
		},
		"MultilineTextNoMarker": {
			args: args{
				text: `hello world!
this is a test
yes, this is a test`,
			},
			want: want{
				out: `// hello world!
// this is a test
// yes, this is a test
`,
				mopts: markers2.Options{},
			},
		},
		"TextWithTerrajetMarker": {
			args: args{
				text: `hello world!
+terrajet:crd:field:TFTag=-
`,
			},
			want: want{
				out: `// hello world!
// +terrajet:crd:field:TFTag=-
`,
				mopts: markers2.Options{
					TerrajetOptions: markers2.TerrajetOptions{
						FieldTFTag: &tftag,
					},
				},
			},
		},
		"TextWithOtherMarker": {
			args: args{
				text: `hello world!
+kubebuilder:validation:Required
`,
			},
			want: want{
				out: `// hello world!
// +kubebuilder:validation:Required
`,
				mopts: markers2.Options{},
			},
		},
		"CommentWithTerrajetOptions": {
			args: args{
				text: `hello world!`,
				opts: []Option{
					WithTFTag("-"),
				},
			},
			want: want{
				out: `// hello world!
// +terrajet:crd:field:TFTag=-
`,
				mopts: markers2.Options{
					TerrajetOptions: markers2.TerrajetOptions{
						FieldTFTag: &tftag,
					},
				},
			},
		},
		"CommentWithMixedOptions": {
			args: args{
				text: `hello world!`,
				opts: []Option{
					WithTFTag("-"),
					WithReferenceConfig(config.Reference{
						Type: reflect.TypeOf(Comment{}).String(),
					}),
				},
			},
			want: want{
				out: `// hello world!
// +terrajet:crd:field:TFTag=-
// +crossplane:generate:reference:type=comments.Comment
`,
				mopts: markers2.Options{
					TerrajetOptions: markers2.TerrajetOptions{
						FieldTFTag: &tftag,
					},
					CrossplaneOptions: markers2.CrossplaneOptions{
						Reference: config.Reference{
							Type: "comments.Comment",
						},
					},
				},
			},
		},
		"CommentWithUnsupportedTerrajetMarker": {
			args: args{
				text: `hello world!
+terrajet:crd:field:TFTag=-
+terrajet:unsupported:key=value
`,
			},
			want: want{
				err: errors.New("cannot parse as a terrajet prefix: +terrajet:unsupported:key=value"),
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c, gotErr := New(tc.text, tc.opts...)
			if diff := cmp.Diff(tc.want.err, gotErr, test.EquateErrors()); diff != "" {
				t.Fatalf("comment.New(...): -want error, +got error: %s", diff)
			}
			if gotErr != nil {
				return
			}
			if diff := cmp.Diff(tc.want.mopts, c.Options); diff != "" {
				t.Errorf("comment.New(...) opts = %v, want %v", c.Options, tc.want.mopts)
			}
			got := c.Build()
			if diff := cmp.Diff(tc.want.out, got); diff != "" {
				t.Errorf("Build() out = %v, want %v", got, tc.want.out)
			}
		})
	}
}