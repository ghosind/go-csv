package csv

type csvBuilder struct {
	comma   rune
	useCRLF bool
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
