package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type account struct {
	Bank_ID string `json:"bank_ID"`
	Balance int    `json:"balance"`
	Name    string `json:"name"`
}

var indexes string = "index"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var empty []string
	indexAsbytes, _ := json.Marshal(empty)
	err := stub.PutState(indexes, indexAsbytes)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "make_account" {
		return t.make_account(stub, args)
	} else if function == "deposit" {
		return t.deposit(stub, args)
	} else if function == "withdrawal" {
		return t.withdrawal(stub, args)
	} else if function == "work" {
		return t.work(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) make_account(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the key and value to set")
	}

	str := `{"bank_ID": "` + args[0] + `", "balance": ` + args[1] + `, "name": "` + args[2] + `"}`
	var index []string

	key = args[0] //rename for funsies
	value = str
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	valAsbytes, err := stub.GetState(indexes)
	json.Unmarshal(valAsbytes, &index)
	index = append(index, args[0])
	indexAsbytes, _ := json.Marshal(index)
	err = stub.PutState(indexes, indexAsbytes)
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "seeAll" {
		return t.seeAll(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func thread(stub shim.ChaincodeStubInterface, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(10 * time.Second)
	_ = stub.PutState(key, []byte("love"))
}

func (t *SimpleChaincode) work(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var wg sync.WaitGroup
	_ = stub.PutState(args[0], []byte("I"))
	wg.Add(1)
	go t.thread(stub, args[1], &wg)
	wg.Wait()
	_ = stub.PutState(args[2], []byte("you"))
	return nil, nil
}

func (t *SimpleChaincode) deposit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	acc := account{}
	json.Unmarshal(valAsbytes, &acc)
	fmt.Println(acc)
	num, err := strconv.Atoi(args[1])

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	acc.Balance += num
	str := `{"bank_ID": "` + acc.Bank_ID + `", "balance": ` + strconv.Itoa(acc.Balance) + `, "name": "` + acc.Name + `"}`
	err = stub.PutState(key, []byte(str))

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return nil, nil

}

func (t *SimpleChaincode) withdrawal(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	acc := account{}
	json.Unmarshal(valAsbytes, &acc)
	fmt.Println(acc)
	num, err := strconv.Atoi(args[1])

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	acc.Balance -= num
	str := `{"bank_ID": "` + acc.Bank_ID + `", "balance": ` + strconv.Itoa(acc.Balance) + `, "name": "` + acc.Name + `"}`
	err = stub.PutState(key, []byte(str))

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return nil, nil

}

func (t *SimpleChaincode) seeAll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var index []string
	var allResults string
	if len(args) != 0 {
		return nil, errors.New("expecting 0 args")
	}
	valAsbytes, err := stub.GetState(indexes)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(valAsbytes, &index)
	for i := range index {
		oneResult, err := stub.GetState(index[i])
		if err != nil {
			return nil, errors.New("error!!")
		}
		allResults = allResults + string(oneResult[:])
	}
	return []byte(allResults), nil
}
