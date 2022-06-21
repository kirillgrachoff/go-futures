package future

type Promise[T any] struct {
	value  chan<- T
	errors chan<- error
}

func (p Promise[T]) SetError(err error) {
	p.errors <- err
}

func (p Promise[T]) SetValue(value T) {
	p.value <- value
}

type Future[T any] struct {
	value  <-chan T
	errors <-chan error
}

func Async[T any](worker func() (T, error)) Future[T] {
    future, promise := MakeContract[T]()

    go func() {
        value, err := worker()
        if err != nil {
            promise.SetError(err)
        } else {
            promise.SetValue(value)
        }
    }()

    return future
}

func (f Future[T]) Get() (value T, err error) {
	select {
	case value = <-f.value:
	case err = <-f.errors:
	}
    return
}

func (f Future[T]) GetUnsafe() T {
	return <-f.value
}

type callback[T any] struct {
	success success[T]
	fail    fail
}

type success[T any] func(value T)
type fail func(err error)

func MakeContract[T any]() (future Future[T], promise Promise[T]) {
	value := make(chan T)
	errors := make(chan error)
	future.value = value
	future.errors = errors
	promise.value = value
	promise.errors = errors
	return
}

func (f Future[T]) SetCallback(cb callback[T]) {
	go func() {
		select {
		case value := <-f.value:
			cb.success(value)
		case err := <-f.errors:
			cb.fail(err)
		}
	}()
}

type mapFunction[T any] func(value T) (T, error)

func (f Future[T]) Map(mapFunction mapFunction[T]) Future[T] {
	future, promise := MakeContract[T]()

	f.SetCallback(callback[T]{
		success: func(value T) {
			val, err := mapFunction(value)
			if err != nil {
				promise.SetError(err)
			} else {
				promise.SetValue(val)
			}
		},
		fail: func(err error) {
			promise.SetError(err)
		},
	})

	return future
}

type recoverFunction[T any] func(err error) (T, error)

func (f Future[T]) Recover(recoverFunction recoverFunction[T]) Future[T] {
	future, promise := MakeContract[T]()

	f.SetCallback(callback[T]{
		success: func(value T) {
			promise.SetValue(value)
		},
		fail: func(err error) {
			value, nerr := recoverFunction(err)
			if nerr != nil {
				promise.SetError(nerr)
			} else {
				promise.SetValue(value)
			}
		},
	})

	return future
}
