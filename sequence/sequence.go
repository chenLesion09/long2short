package sequence

type Sequence interface {
	Next() (uint64, error)
}
