package dragonchain

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	// ExecutionOrderSerial informs DragonChain to execute
	// multiple instances of a Contract serially.
	ExecutionOrderSerial = "serial"
	// ExecutionOrderParallel infoms DragonChain to execute
	// multiple instances of the Contract in paralell
	ExecutionOrderParallel = "parallel"

	// ContractTypeTransaction informs DragonChain that a contract
	// should be executed as a single transaction.
	ContractTypeTransaction = "transaction"
	// ContractTypeCron informs DragonChain that a contract should
	// be executed on a cron.
	ContractTypeCron = "cron"
)

// ExecutionOrder informs DragonChain how to execute multiple
// instances of a smart contract.
type ExecutionOrder string

// ContractType informs DragonChain of the type of contract.
type ContractType string

// ContractStatus represents the status of a contract.
type ContractStatus string

// ContractOrigin represents the contract's origin.
type ContractOrigin string

// ContractRuntime represents the runtime that will be used
// when executing the contract.
type ContractRuntime string

// ContractDesiredState represents the desired execution state
// of the contract.
type ContractDesiredState string

// Contract is a smart contract running on a DragonChain.
type Contract struct {
	DCRN              string            `json:"dcrn"`
	Version           string            `json:"version"`
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Status            ContractStatus    `json:"status"`
	CustomEnvironment map[string]string `json:"custom_environment_variables"`
	Origin            ContractOrigin    `json:"origin"`
	Runtime           ContractRuntime   `json:"runtime"`
	Type              string            `json:"sc_type"`
	Code              string            `json:"code"`
	S3Bucket          string            `json:"s3_bucket"`
	S3Path            string            `json:"s3_path"`
	Serial            bool              `json:"is_serial"`
}

// ContractDefinition is used to create new smart contracts and
// defines how a contract will run on a DragonChain.
type ContractDefinition struct {
	Version        string               `json:"version,omitempty"`
	Type           string               `json:"txn_type,omitempty"`
	Order          ExecutionOrder       `json:"execution_order,omitempty"`
	Image          string               `json:"image,omitempty"`
	Command        string               `json:"cmd,omitempty"`
	DesiredState   ContractDesiredState `json:"desired_state,omitempty"`
	Args           []string             `json:"args,omitempty"`
	Environment    map[string]string    `json:"env,omitempty"`
	Secrets        map[string]string    `json:"secrets,omitempty"`
	Seconds        int                  `json:"seconds,omitempty"`
	Cron           string               `json:"cron,omitempty"`
	Authentication string               `json:"auth,omitempty"`
}

// Contract retrieves the contract with the given id from a DragonChain.
//
// An error is returned if the contract could not be retrieved. The error
// will be an APIError.
func (c *Client) Contract(id string) (*Contract, error) {
	var resp struct {
		HTTPResponse
		Response Contract
	}
	if err := c.Get("/contract/"+id, &resp); err != nil {
		return nil, &APIError{Err: err}
	}
	if !resp.OK {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}
	return &resp.Response, nil
}

// QueryContracts queries all contracts on a DragonChain using the provided
// QueryOptions, returning a list of Contracts that match the query.
//
// An error is returned if the query could not be completed. If the error is a
// result of the HTTP request failing or returning a non-2xx status code, the
// returned error will be an APIError.
func (c *Client) QueryContracts(q *QueryOptions) ([]*Contract, error) {
	var resp struct {
		HTTPResponse
		Response struct {
			Results []*Contract
		}
	}
	url := "/contract"
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
		results = make([]*Contract, 0)
	}
	return results, nil
}

// UpdateContract updates the smart contract running on a DragonChain with the provided id.
//
// The provided ContractDefinition is used to define which fields should be updated. Only
// fields that are set in the ContractDefinition will be updated and all other fields will
// be left untouched.
//
// An error is returned if the contract could not be updated. If the error is a result of the
// HTTP request failing or returning a non-2xx status code, the error will be an APIError.
// If the error is a result of the contract failing to update on the server, the error will be
// ErrContractUpdateFailed.
func (c *Client) UpdateContract(id string, update *ContractDefinition) error {
	b, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to JSON marshal contract definition : %s", err)
	}
	var resp HTTPResponse
	if err = c.Put("/contract/"+id, b, &resp); err != nil {
		return &APIError{Err: err}
	}
	if !resp.OK {
		return &APIError{StatusCode: resp.StatusCode}
	}
	return nil
}

// DeleteContract deletes a smart contract from the DragonChain server with the provided id.
//
// An error is returned if the contract could not be deleted. If the error is a result of the
// HTTP request failing or returning a non-2xx status code, the error will be an APIError.
// Otherwise, if the error is a result of the server failing to delete the contract, the error
// will be ErrContractDeleteFailed.
func (c *Client) DeleteContract(id string) error {
	var resp HTTPResponse
	if err := c.Delete("/contract/"+id, &resp); err != nil {
		return &APIError{Err: err}
	}
	if !resp.OK {
		return &APIError{StatusCode: resp.StatusCode}
	}
	return nil
}

// CreateContract creates a new smart contract for a DragonChain defined by def.
//
// An error is returned if the contract could not be created. If the error is a
// result of the HTTP request failing or returning a non-2xx status code, the
// error will be an APIError.
func (c *Client) CreateContract(def *ContractDefinition) error {
	vers, err := strconv.Atoi(def.Version)
	if err != nil {
		return fmt.Errorf("invalid version %s: %s", def.Version, err)
	}
	if vers <= 0 {
		def.Version = "3"
	}
	b, err := json.Marshal(def)
	if err != nil {
		return fmt.Errorf("failed to JSON marshal contract definition: %s", err)
	}
	var resp HTTPResponse
	if err = c.Post("/contract", b, &resp); err != nil {
		return &APIError{Err: err}
	}
	if !resp.OK {
		return &APIError{StatusCode: resp.StatusCode}
	}
	return nil
}
