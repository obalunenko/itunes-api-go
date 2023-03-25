// Package itunes_search_go provides client for access to iTunes API.
// See https://affiliate.itunes.apple.com/resources/documentation/itunes-store-web-service-search-api
package itunes_search_go

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/obalunenko/itunes-api-go/internal"
	"github.com/obalunenko/itunes-api-go/option"
)

const (
	schema = "https"
	host   = "itunes.apple.com"

	lookup = "lookup"

	countryParam = "country"
	mediaParam   = "media"
	idParam      = "id"

	logPfx = "itunes-search-go: "
)

// Client is an contract for itunes api.
type Client interface {
	// Lookup returns LookupResponse for application by it's id.
	Lookup(ctx context.Context, id int, opts ...option.LookupOption) (LookupResponse, error)
}

type httpclient interface {
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	httpclient
}

// New returns itunes api client.
func New() Client {
	return client{
		httpclient: &http.Client{},
	}
}

func (c client) Lookup(ctx context.Context, id int, opts ...option.LookupOption) (LookupResponse, error) {
	opts = append(opts, option.WithID(id))

	params := newLookupParams(opts)

	uri := url.URL{
		Scheme:      schema,
		Opaque:      "",
		User:        nil,
		Host:        host,
		Path:        lookup,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}

	v := uri.Query()

	v.Add(idParam, strconv.Itoa(params.ID))
	v.Add(mediaParam, params.Media)
	v.Add(countryParam, params.Country)

	uri.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, uri.String(), http.NoBody)
	if err != nil {
		return LookupResponse{}, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return LookupResponse{}, fmt.Errorf("client do: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("%s: failed to close body: %v", logPfx, err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LookupResponse{}, fmt.Errorf("read body: %w", err)
	}

	res, err := UnmarshalLookupResponse(body)
	if err != nil {
		return LookupResponse{}, fmt.Errorf("response unmarshal: %w", err)
	}

	return res, nil
}

func newLookupParams(opts []option.LookupOption) internal.LookupParams {
	var p internal.LookupParams

	for _, opt := range opts {
		opt.Apply(&p)
	}

	return p
}
