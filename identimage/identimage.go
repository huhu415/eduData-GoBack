package identimage

type IdentImage interface {
	Identify(base64Image *string) (string, error)
}
