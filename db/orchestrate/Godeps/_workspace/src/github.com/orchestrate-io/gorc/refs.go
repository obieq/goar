// Copyright 2014, Orchestrate.IO, Inc.

package gorc

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Holds results returned from a ref list.
type RefResults struct {
	Count   uint64      `json:"count"`
	Results []RefResult `json:"results"`
	Next    string      `json:"next,omitempty"`
}

// An individual ref result.
type RefResult struct {
	Path     Path            `json:"path"`
	RawValue json.RawMessage `json:"value,omitempty"`
	RefTime  uint64          `json:"reftime"`
}

// Get a collection-key pair's value at a specific ref.
func (c *Client) GetRef(collection, key, ref string) (*KVResult, error) {
	return c.GetPath(&Path{Collection: collection, Key: key, Ref: ref})
}

// List the refs of a value in time order with the specified page size
// optionally retrieving values.
func (c *Client) ListRefs(collection, key string, limit int, values bool) (*RefResults, error) {
	queryVariables := url.Values{
		"limit":  []string{strconv.Itoa(limit)},
		"values": []string{strconv.FormatBool(values)},
	}

	trailingUri := collection + "/" + key + "/refs/?" + queryVariables.Encode()

	return c.doListRefs(trailingUri)
}

// List the refs of a value in time order with the specified page size
// optionally retrieving values starting at the specified offset.
func (c *Client) ListRefsFromOffset(collection, key string, limit int, values bool, offset int) (*RefResults, error) {
	queryVariables := url.Values{
		"limit":  []string{strconv.Itoa(limit)},
		"values": []string{strconv.FormatBool(values)},
		"offset": []string{strconv.Itoa(offset)},
	}

	trailingUri := collection + "/" + key + "/refs/?" + queryVariables.Encode()

	return c.doListRefs(trailingUri)
}

// Get the page of ref list results that follow the provided set.
func (c *Client) ListRefsGetNext(results *RefResults) (*RefResults, error) {
	return c.doListRefs(results.Next[4:])
}

// Execute a ref list operation.
func (c *Client) doListRefs(trailingUri string) (*RefResults, error) {
	resp, err := c.doRequest("GET", trailingUri, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	decoder := json.NewDecoder(resp.Body)
	result := new(RefResults)
	if err := decoder.Decode(result); err != nil {
		return result, err
	}

	return result, nil
}

// Check if there is a subsequent page of ref list results.
func (r *RefResults) HasNext() bool {
	return r.Next != ""
}

// Marshall the value of a RefResult into the provided object.
func (r *RefResult) Value(value interface{}) error {
	return json.Unmarshal(r.RawValue, value)
}

// Determines if the given ref represents a deletion.
func (r *RefResult) IsDeleted() bool {
	return r.Path.Tombstone
}
