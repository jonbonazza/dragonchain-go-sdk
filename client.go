package dragonchain

import "fmt"

// HTTPResponse is a base response structure that is present
// in ever response from the DragonChain API, whether the
// response was successful or not.
type HTTPResponse struct {
	OK         bool `json:"ok"`
	StatusCode int  `json:"status"`
}

// Inquirer provides RESTful primitives for communicating with
// the DragonChain API server.
type Inquirer interface {
	Get(resource string, out interface{}) error
	Post(resource string, body []byte, out interface{}) error
	Put(resource string, body []byte, out interface{}) error
	Delete(resource string, out interface{}) error
}

// APIError is an error returned from the DragonChain API.
type APIError struct {
	// StatusCode is the status code of a response from
	// the DragonChain API server. StatusCode will only
	// be present when the request itself was successful,
	// but a non-2xx status code was returned.
	StatusCode int
	// Err is the actual error provided by the inquierer.
	// Err will only be available when the request itself
	// failed.
	Err error
}

// Error satisfies the error interface and returns a string
// containing the error message. This message will be the
// message from Err if it's available, otherwise, an error
// message contianing the status code will be returned.
func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	if e.StatusCode > 0 {
		return fmt.Sprintf("received status code %d from server", e.StatusCode)
	}
	return ""
}

// String statisfies the Stringer interface and is a proxy
// call to the Error method.
//
// This allows for more convenient usage with formatting functons,
// such as fmt.Sprintf and fmt.Errorf.
func (e *APIError) String() string {
	return e.Error()
}

// QueryOptions defines various parameters to be used when
// querying.
type QueryOptions struct {
	// QueryString is a Lucene query string that will be used to query.
	QueryString string
	// Sort is an optional string that defines how query results should be
	// sorted.
	Sort string
	// Offset declares the within the result array to return. All results before
	// this offset will be dropped before returning. This, along with Limit, allows
	// for server-side pagination.
	Offset int
	// Limit is the maximum number of results to return, starting with the provided Offset,
	// allowing for network-effecitent server-side pagination.
	Limit int
}

// Client is used to interact with a DragonChain using the DragonChain API.
type Client struct {
	Inquirer

	// Credentials are the credentials used to access a DragonChain via
	// the DragonChain API.
	Credentials *Credentials
}
