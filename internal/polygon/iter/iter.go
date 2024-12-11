package iter

import (
	"context"

	"github.com/jamesonhm/fingator/internal/uri"
)

type ListResponse interface {
	// NextPage returns a URL for the next page of list results.  The URL is returned by the API.
	NextPage() string
}

// Query defines a closure that iterators must implement. The implementation
// should include a call to the API and should return the API response
// with a different slice of results
type Query[T any] func(string) (ListResponse, []T, error)

type Iter[T any] struct {
	ctx   context.Context
	query Query[T]

	page    ListResponse
	item    T
	results []T

	err error
}

func NewIter[T any](ctx context.Context, path string, params any, query Query[T]) *Iter[T] {
	it := Iter[T]{
		ctx:   ctx,
		query: query,
	}

	uri := uri.New().EncodeParams(path, params)

	it.page, it.results, it.err = it.query(uri)
	return &it
}

func (it *Iter[T]) Next() bool {
	if it.err != nil {
		return false
	}

	if len(it.results) == 0 && it.page.NextPage() != "" {
		it.page, it.results, it.err = it.query(it.page.NextPage())
	}

	if it.err != nil || len(it.results) == 0 {
		return false
	}
	it.err = it.ctx.Err()
	if it.err != nil {
		return false
	}

	it.item = it.results[0]
	it.results = it.results[1:]
	return true
}

func (it *Iter[T]) Item() T {
	return it.item
}

func (it *Iter[T]) Err() error {
	return it.err
}
