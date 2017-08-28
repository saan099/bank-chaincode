package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
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
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}

	var empty []string
	indexAsbytes, _ := json.Marshal(empty)
	err := stub.PutState(indexes, indexAsbytes)

	if err != nil {
		return shim.Error(fmt.Sprintf("error yo"))
	}

	return shim.Success(nil)
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub)
	} else if function == "make_account" {
		return t.make_account(stub, args)
	} else if function == "deposit" {
		return t.deposit(stub, args)
	} else if function == "withdrawal" {
		return t.withdrawal(stub, args)
	} else if function == "work" {
		return t.work(stub, args)
	} else if function == "check" {
		return t.check(stub, args)
	} else if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return shim.Error(fmt.Sprintf("No function called"))
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) make_account(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}

	str := `{"bank_ID": "` + args[0] + `", "balance": ` + args[1] + `, "name": "` + args[2] + `"}`
	var index []string

	key = args[0] //rename for funsies
	value = str
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	valAsbytes, err := stub.GetState(indexes)
	json.Unmarshal(valAsbytes, &index)
	index = append(index, args[0])
	indexAsbytes, _ := json.Marshal(index)
	err = stub.PutState(indexes, indexAsbytes)
	return shim.Success(nil)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "seeAll" {
		return t.seeAll(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return shim.Error(fmt.Sprintf("Wrong number of arguments"))
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}

	return shim.Success(valAsbytes)
}

func thread(stub shim.ChaincodeStubInterface, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Second)
	}

	_ = stub.PutState(key, []byte("love"))
}

func (t *SimpleChaincode) work(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var wg sync.WaitGroup
	_ = stub.PutState(args[0], []byte("I"))
	wg.Add(1)
	go thread(stub, args[1], &wg)
	wg.Wait()
	_ = stub.PutState(args[2], []byte("you"))
	return shim.Success(nil)
}

func (t *SimpleChaincode) check(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	time.Sleep(20 * time.Second)
	_ = stub.PutState("some", []byte("yo yo"))

	return shim.Success(nil)

}

func (t *SimpleChaincode) deposit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	acc := account{}
	json.Unmarshal(valAsbytes, &acc)
	fmt.Println(acc)
	num, err := strconv.Atoi(args[1])

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	acc.Balance += num
	str := `{"bank_ID": "` + acc.Bank_ID + `", "balance": ` + strconv.Itoa(acc.Balance) + `, "name": "` + acc.Name + `"}`
	err = stub.PutState(key, []byte(str))

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	return shim.Success(nil)

}

func (t *SimpleChaincode) withdrawal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	acc := account{}
	json.Unmarshal(valAsbytes, &acc)
	fmt.Println(acc)
	num, err := strconv.Atoi(args[1])

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	acc.Balance -= num
	str := `{"bank_ID": "` + acc.Bank_ID + `", "balance": ` + strconv.Itoa(acc.Balance) + `, "name": "` + acc.Name + `"}`
	err = stub.PutState(key, []byte(str))

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	return shim.Success(nil)

}

func (t *SimpleChaincode) seeAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var index []string
	var allResults string
	if len(args) != 0 {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	valAsbytes, err := stub.GetState(indexes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Wrong number of arguments"))
	}
	json.Unmarshal(valAsbytes, &index)
	for i := range index {
		oneResult, err := stub.GetState(index[i])
		if err != nil {
			return shim.Error(fmt.Sprintf("Wrong number of arguments"))
		}
		allResults = allResults + string(oneResult[:])
	}
	return shim.Success([]byte(allResults))
}
