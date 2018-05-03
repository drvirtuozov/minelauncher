package events

var TaskProgress = make(chan ProgressBarFraction)

type ProgressBarFraction struct {
	Fraction float64
	Text     string
}
