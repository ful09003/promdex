package internal

//PromdexFlavor holds data pertaining to Promdex 'flavor' text
type PromdexFlavor struct {
	CtxString string `json:"context"` //String representing the context (flavor) for a particular Prom metric
}

//NewFlavor takes a string containing metric context and returns a PromdexFlavor
func NewFlavor(c string) PromdexFlavor {
	var f PromdexFlavor
	f.CtxString = c

	return f
}
