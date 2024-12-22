package logger

type Option struct {
	IsEnable                     bool
	MaskingFields                []string
	MaskingType                  string
	AdditionalSkippedContentType []string
}
