package dragonchain

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Transaction represents a single transaction on a DragonChain blockchain.
type Transaction struct {
	DCRN    string            `json:"dcrn"`
	Version int               `json:"vesion"`
	Header  TransactionHeader `json:"header"`
	Payload string            `json:"payload"`
	Proof   TransactionProof  `json:"proof"`
}

// TransactionHeader contains metadata for a Transaction.
type TransactionHeader struct {
	// ID is the unique transaction ID.
	ID string `json:"txn_id"`
	// Type is the transaction type.
	//
	// For smart contract transactions, this will be the
	// name of the smart contract.
	//
	// For other transactions, this will be the type that
	// was provided when the transaction was created.
	Type string `json:"txn_type"`
	// DragonChainID is the unique ID of the DragonChain that
	// the transaction is associated with.
	DragonChainID string `json:"dc_id"`
	// Tag is an optional string that can be applied to a transaction.
	// Tags are useful when querying transactions to filter the results.
	Tag string `json:"tag"`
	// Timestamp is the epoch timestamp that the transaction was committed
	// to the DragonChain blockchain.
	Timestamp string `json:"timestamp"`
	// BlockID is the unique ID of the block on the DragonChain blockchain
	// that the transaction was committed to.
	BlockID string `json:"block_id"`

	Invoker string `json:"invoker"`
}

// TransactionProof contains the proof for the transaction.
type TransactionProof struct {
	Full     string `json:"full"`
	Stripped string `json:"stripped"`
}

// TransactionDefinition defines a transaction to be created on a DragonChain.
type TransactionDefinition struct {
	Version string `json:"version"`
	// Type is the transaction type of the transaction. For smart contract transactions,
	// this is the name of the smart contract. Otherwise, the type is arbitrary and is
	// only really useful for querying purposes.
	Type string `json:"txn_type"`
	// Payload is the payload that will be contained by the transaction.
	Payload interface{} `json:"payload"`
	// Tag is an optional string that will be associated with the transaction.
	// Tags are useful when querying transactions to filter results.
	Tag string `json:"tag,omitempty"`
}

// GetTransaction retrieves the transaction for the DragonChain with the provided id.
// An APIError error is returned if the request failed or a non-2xx status code was returned.
func (c *Client) GetTransaction(id string) (*Transaction, error) {
	var resp struct {
		HTTPResponse
		Response Transaction
	}
	if err := c.Get("/transaction/"+id, &resp); err != nil {
		return nil, &APIError{Err: err}
	}
	if !resp.OK {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}
	return &resp.Response, nil
}

type createTransactionResponse struct {
	ID string `json:"transaction_id"`
}

// CreateTransaction creates a new transaction on a DragonChain blockchain.
// The transaction is defined by def.
//
// An error is returned if the transaction could not be created. If the error is
// a result of the request failing or returned a non-2xx status code, the returned
// error will be of type APIError.
func (c *Client) CreateTransaction(def *TransactionDefinition) (string, error) {
	vers, err := strconv.Atoi(def.Version)
	if err != nil {
		return "", fmt.Errorf("invalid version %s: %s", def.Version, err)
	}
	if vers <= 0 {
		def.Version = "1"
	}
	b, err := json.Marshal(def)
	if err != nil {
		return "", fmt.Errorf("failed to JSON marshal transaction object: %s", err)
	}
	var resp struct {
		HTTPResponse
		Response createTransactionResponse
	}
	if err = c.Post("/transaction", b, &resp); err != nil {
		return "", &APIError{Err: err}
	}
	if !resp.OK {
		return "", &APIError{StatusCode: resp.StatusCode}
	}
	return resp.Response.ID, nil
}

// QueryTransactions queries all transactions on a DragonChain blockchain using the provided
// QueryOptions, returning a list of Transactions that match the query.
//
// An error is returned if the query could not be completed. If the error is a result of the
// HTTP request failing or returning a non-2xx status code, the returned error will be of
// type APIError.
func (c *Client) QueryTransactions(q *QueryOptions) ([]*Transaction, error) {
	var resp struct {
		HTTPResponse
		Response struct {
			Results []*Transaction
		}
	}
	url := "/transation"
	params := luceneQueryParams(q)
	if params != "" {
		url += "?" + params
	}
	if err := c.Get(url, &resp); err != nil {
		return nil, &APIError{Err: err}
	}
	if !resp.OK {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}
	results := resp.Response.Results
	if results == nil {
		results = make([]*Transaction, 0)
	}
	return results, nil
}

type bulkTransactionDef struct {
	Payload []*TransactionDefinition `json:"payload"`
}

// BulkCreateTransactions creates multiple transactions in a single API request. This is more
// efficient than creating a separate request for each transaction, and should be used whereever
// possible to optimize client CPU and network performance. It can also (though not always) result
// in lower costs when using cloud providers such as AWS or Google Cloud.
//
// It should be noted that the bulk creation of contracts is *not atomic* and as such, if one or more
// contracts cannot be created, those that were created will not be rolled back. Additionally, the
// overall request will be considered successful if at least one transaction was created successfully.
//
// An error is returned if the operation fails. If the error is a result of the HTTP request failing or
// a non-2xx status code being returend, the error will be of type APIError.
func (c *Client) BulkCreateTransactions(txs []*TransactionDefinition) ([]string, error) {
	payload := bulkTransactionDef{txs}
	b, err := json.Marshal(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal bulk transaction objects: %s", err)
	}
	var resp struct {
		HTTPResponse
		Response []struct {
			ID string `json:"transaction_id"`
		}
	}
	if err = c.Post("/transaction_bulk", b, &resp); err != nil {
		return nil, &APIError{Err: err}
	}
	if !resp.OK {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}
	ids := make([]string, len(resp.Response))
	for i, txResp := range resp.Response {
		ids[i] = txResp.ID
	}
	return ids, nil
}
