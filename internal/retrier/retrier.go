package retrier

type Retrier interface {
	Retry(func() error) (int, error)
}
