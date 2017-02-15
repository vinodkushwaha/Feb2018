package main

import (

	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Customer Chaincode implementation
type CustomerChaincode struct {
}

var customerIndexTxStr = "_customerIndexTxStr"

type CustomerData struct{
	CUSTOMER_ID string `json:"CUSTOMER_ID"`
	CUSTOMER_NAME string `json:"CUSTOMER_NAME"`
	CUSTOMER_DOB string `json:"CUSTOMER_DOB"`
	CUSTOMER_KYC_FLAG string `json:"CUSTOMER_KYC_FLAG"`
	CUSTOMER_DOC string `json:"CUSTOMER_DOC"`
	}


func (t *CustomerChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	var err error
	// Initialize the chaincode
	// function, args := stub.GetFunctionAndParameters()
	fmt.Printf("Deployment of Customer ChainCode is completed\n")

	var emptyCustomerTxs []CustomerData
	jsonAsBytes, _ := json.Marshal(emptyCustomerTxs)
	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}

// Add customer data for the policy
func (t *CustomerChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.RegisterCustomer(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")

}

func (t *CustomerChaincode)  RegisterCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var CustomerDataObj CustomerData
	var CustomerDataList []CustomerData
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Need 5 arguments")
	}

	// Initialize the chaincode
	CustomerDataObj.CUSTOMER_ID = args[0]
	CustomerDataObj.CUSTOMER_NAME = args[1]
	CustomerDataObj.CUSTOMER_DOB = args[2]
	CustomerDataObj.CUSTOMER_KYC_FLAG = args[3]
	CustomerDataObj.CUSTOMER_DOC = args[4]

	fmt.Printf("Input from user:%s\n", CustomerDataObj)

	customerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return shim.Error("Failed to get customer transactions")
	}
	json.Unmarshal(customerTxsAsBytes, &CustomerDataList)

	CustomerDataList = append(CustomerDataList, CustomerDataObj)
	jsonAsBytes, _ := json.Marshal(CustomerDataList)

	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *CustomerChaincode) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var customer_name string // Entities
	var customer_id string
	var customer_dob string
	var err error
	var resAsBytes []byte

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 parameters to query")
	}

	customer_name = args[0]
	customer_id = args[1]
	customer_dob = args[2]

	return (t.GetCustomerDetails(stub, customer_name, customer_id, customer_dob))

	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + customer_name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(resAsBytes)
}

func (t *CustomerChaincode)  GetCustomerDetails(stub shim.ChaincodeStubInterface, customer_name string, customer_id string, customer_dob string) pb.Response {

	//var requiredObj CustomerData
	var objFound bool
	CustomerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return shim.Error("Failed to get Customer Transactions")
	}
	var CustomerTxObjects []CustomerData
	var CustomerTxObjects1 []CustomerData
	json.Unmarshal(CustomerTxsAsBytes, &CustomerTxObjects)
	length := len(CustomerTxObjects)
	fmt.Printf("Output from chaincode: %s\n", CustomerTxsAsBytes)

	if customer_id == "" {
		res, err := json.Marshal(CustomerTxObjects)
		if err != nil {
		return shim.Error("Failed to Marshal the required Obj")
		}
		return shim.Success(res)
	}

	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := CustomerTxObjects[i]
		if customer_id == obj.CUSTOMER_ID && customer_name == obj.CUSTOMER_NAME && customer_dob == obj.CUSTOMER_DOB {
			CustomerTxObjects1 = append(CustomerTxObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}

	if objFound {
		res, err := json.Marshal(CustomerTxObjects1)
		if err != nil {
		return shim.Error("Failed to Marshal the required Obj")
		}
		return shim.Success(res)
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return shim.Error("Failed to Marshal the required Obj")
		}
		return shim.Success(res)
	}
}

func main() {
	err := shim.Start(new(CustomerChaincode))
	if err != nil {
		fmt.Printf("Error starting Customer Simple chaincode: %s", err)
	}
}
