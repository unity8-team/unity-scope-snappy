package interfaces

// PreviewGenerator is an interface to be implemented by any struct that wishes
// to provide previews for use in this scope.
type PreviewGenerator interface {
	Generate(receiver WidgetReceiver) error
}
