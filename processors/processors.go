package processors

type Tracker interface {
	TrackChanges() ([]byte, error)
}
