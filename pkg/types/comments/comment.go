package comments

import (
	"strings"

	markers2 "github.com/crossplane-contrib/terrajet/pkg/types/markers"

	"github.com/crossplane-contrib/terrajet/pkg/config"
)

// Option is a comment option
type Option func(*Comment)

// WithReferenceConfig returns a comment options with the given reference config
func WithReferenceConfig(cfg config.Reference) Option {
	return func(c *Comment) {
		c.Reference = cfg
	}
}

// WithTFTag returns a comment options with input tf tag
func WithTFTag(s string) Option {
	return func(c *Comment) {
		c.FieldTFTag = &s
	}
}

// New returns a Comment by parsing Terrajet markers as Options
func New(text string, opts ...Option) (*Comment, error) {
	to := markers2.TerrajetOptions{}
	co := markers2.CrossplaneOptions{}

	rawLines := strings.Split(strings.TrimSpace(text), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, rl := range rawLines {
		rl = strings.TrimSpace(rl)
		if rl == "" {
			lines = append(lines, rl)
			continue
		}
		// Only add raw marker line if not processed as an option (e.g. if it is
		// not a known marker.) Known markers will still be printed as
		// comments while building from options.
		parsed, err := markers2.ParseAsTerrajetOption(&to, rl)
		if err != nil {
			return nil, err
		}
		if parsed {
			continue
		}

		lines = append(lines, rl)
	}

	c := &Comment{
		Text: strings.Join(lines, "\n"),
		Options: markers2.Options{
			TerrajetOptions:   to,
			CrossplaneOptions: co,
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

// Comment represents a comment with text and supported marker options.
type Comment struct {
	Text string
	markers2.Options
}

// String returns a string representation of this Comment (no "// " prefix)
func (c *Comment) String() string {
	if c.Text == "" {
		return c.Options.String()
	}
	return c.Text + "\n" + c.Options.String()
}

// Build builds comments by adding comment prefix ("// ") to each line of
// the string representation of this Comment.
func (c *Comment) Build() string {
	all := strings.ReplaceAll("// "+c.String(), "\n", "\n// ")
	return strings.TrimSuffix(all, "// ")
}