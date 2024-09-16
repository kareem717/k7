package shared

// Injectable is the interface the represents a
// template that can be used to generate a file for an app
type Injectable struct {
	FilePath string
	Bytes    []byte
}

type Stringable interface {
	String() string
}

type AppTemplate struct {
	Templates  []Injectable
	// Extensions map[Stringable][]byte
}

// type Injectable interface {
// 	Template() AppTemplate
// 	GetGlobalTemplate(name Stringable) ([]byte, error)
// }

type GlobalFiles map[Stringable][]byte

