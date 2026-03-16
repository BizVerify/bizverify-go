package bizverify

import "context"

const defaultPageSize = 50

// SearchIterator provides paginated access to search results using the Scanner pattern.
type SearchIterator struct {
	client    *httpClient
	params    SearchParams
	ctx       context.Context
	buffer    []SearchResult
	bufferIdx int
	offset    int
	total     int
	done      bool
	err       error
	fetched   bool
}

func newSearchIterator(ctx context.Context, c *httpClient, params SearchParams) *SearchIterator {
	return &SearchIterator{
		client: c,
		params: params,
		ctx:    ctx,
	}
}

// Next advances the iterator to the next result. Returns false when iteration is complete or on error.
func (it *SearchIterator) Next() bool {
	if it.done || it.err != nil {
		return false
	}

	if it.bufferIdx < len(it.buffer) {
		return true
	}

	if it.fetched && it.offset >= it.total {
		it.done = true
		return false
	}

	limit := defaultPageSize
	p := it.params
	p.Limit = &limit
	p.Offset = &it.offset

	var resp SearchResponse
	err := it.client.request(it.ctx, requestOptions{
		method: "POST",
		path:   "/v1/search",
		body:   p,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		it.err = err
		return false
	}

	it.fetched = true
	it.total = resp.Total
	it.buffer = resp.Results
	it.bufferIdx = 0
	it.offset += len(resp.Results)

	if len(resp.Results) == 0 {
		it.done = true
		return false
	}

	return true
}

// Value returns the current search result.
func (it *SearchIterator) Value() SearchResult {
	r := it.buffer[it.bufferIdx]
	it.bufferIdx++
	return r
}

// Err returns the error encountered during iteration, if any.
func (it *SearchIterator) Err() error {
	return it.err
}

// TotalResults returns the total number of results available.
func (it *SearchIterator) TotalResults() int {
	return it.total
}
