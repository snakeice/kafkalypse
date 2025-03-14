package app

type ResourceMetadata struct {
	Kind   string
	Render func()
}
