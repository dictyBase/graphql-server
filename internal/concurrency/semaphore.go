package concurrency

// Semaphore provides a reusable mechanism for controlling concurrent access to a resource.
type Semaphore struct {
	sem chan struct{}
}

// NewSemaphore creates a new semaphore with the specified concurrency limit.
func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		sem: make(chan struct{}, limit),
	}
}

// Acquire acquires a semaphore permit. Blocks if no permits are available.
func (s *Semaphore) Acquire() {
	s.sem <- struct{}{}
}

// Release releases a semaphore permit.
func (s *Semaphore) Release() {
	<-s.sem
}

// Execute runs the provided function within the semaphore's permit.
// Automatically acquires a permit before executing the function and releases it afterward.
func (s *Semaphore) Execute(fn func()) {
	s.Acquire()
	defer s.Release()
	fn()
}