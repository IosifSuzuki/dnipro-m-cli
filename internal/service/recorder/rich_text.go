package recorder

type RichText struct {
	Value           string
	Link            string
	IsBold          bool
	BackgroundColor *Color
}

type Color struct {
	Red   float64
	Green float64
	Blue  float64
}
