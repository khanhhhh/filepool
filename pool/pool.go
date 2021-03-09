package pool

// Pool :
type Pool interface {
	Upload()
	Download()
	Clean()
}
