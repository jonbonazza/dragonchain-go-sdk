package dragonchain

const (
	// ExecutionOrderSerial informs DragonChain to execute
	// multiple instances of a Contract serially.
	ExecutionOrderSerial = "serial"
	// ExecutionOrderParallel infoms DragonChain to execute
	// multiple instances of the Contract in paralell
	ExecutionOrderParallel = "parallel"
)

// ExecutionOrder informs DragonChain how to execute multiple
// instances of a smart contract.
type ExecutionOrder string

// ContractDefinition is used to create new smart contracts and
// defines how a contract will run on a DragonChain.
type ContractDefinition struct {
	Type           string            `json:"txn_type"`
	Order          ExecutionOrder    `json:"execution_order"`
	Image          string            `json:"image"`
	Command        string            `json:"cmd"`
	Args           []string          `json:"args"`
	Environment    map[string]string `json:"env"`
	Secrets        map[string]string `json:"secrets"`
	Seconds        string            `json:"seconds"`
	Cron           string            `json:"cron"`
	Authentication string            `json:"auth"`
}
