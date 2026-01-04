package csv

type csvBuilder struct {
	comma    rune
	useCRLF  bool
	noHeader bool
}

func newCSVBuilder(opts ...CSVOption) *csvBuilder {
	builder := &csvBuilder{
		comma:    ',',
		useCRLF:  false,
		noHeader: false,
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

type CSVOption func(*csvBuilder)

// WithComma sets the field delimiter rune for the CSV encoder/decoder.
func WithComma(r rune) CSVOption {
	return func(cb *csvBuilder) {
		cb.comma = r
	}
}

// WithCRLF sets whether to use \r\n as the line terminator for the CSV encoder/decoder.
func WithCRLF(useCRLF bool) CSVOption {
	return func(cb *csvBuilder) {
		cb.useCRLF = useCRLF
	}
}

// WithNoHeader sets whether to omit the header row in the CSV encoding.
func WithNoHeader(noHeader bool) CSVOption {
	return func(cb *csvBuilder) {
		cb.noHeader = noHeader
	}
}
