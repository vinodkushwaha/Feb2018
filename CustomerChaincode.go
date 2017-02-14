/*
Copyright Capgemini India. 2016 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
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


func (t *CustomerChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	// Initialize the chaincode

	fmt.Printf("Deployment of Customer ChainCode is completed\n")

	var emptyCustomerTxs []CustomerData
	jsonAsBytes, _ := json.Marshal(emptyCustomerTxs)
	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}


	return nil, nil
}

// Add customer data for the policy
func (t *CustomerChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == customerIndexTxStr {
		return t.RegisterCustomer(stub, args)
	}
	return nil, nil
}

func (t *CustomerChaincode)  RegisterCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var CustomerDataObj CustomerData
	var CustomerDataList []CustomerData
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Need 5 arguments")
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
		return nil, errors.New("Failed to get customer transactions")
	}
	json.Unmarshal(customerTxsAsBytes, &CustomerDataList)

	CustomerDataList = append(CustomerDataList, CustomerDataObj)
	jsonAsBytes, _ := json.Marshal(CustomerDataList)

	err = stub.PutState(customerIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *CustomerChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

	var customer_name string // Entities
	var customer_id string
	var customer_dob string
	var err error
	var resAsBytes []byte

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 parameters to query")
	}

	customer_name = args[0]
	customer_id = args[1]
	customer_dob = args[2]

	resAsBytes, err = t.GetCustomerDetails(stub, customer_name, customer_id, customer_dob)

	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *CustomerChaincode)  GetCustomerDetails(stub shim.ChaincodeStubInterface, customer_name string, customer_id string, customer_dob string) ([]byte, error) {

	//var requiredObj CustomerData
	var objFound bool
	CustomerTxsAsBytes, err := stub.GetState(customerIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	var CustomerTxObjects []CustomerData
	var CustomerTxObjects1 []CustomerData
	json.Unmarshal(CustomerTxsAsBytes, &CustomerTxObjects)
	length := len(CustomerTxObjects)
	fmt.Printf("Output from chaincode: %s\n", CustomerTxsAsBytes)

	if customer_id == "" {
		res, err := json.Marshal(CustomerTxObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
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
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}

func main() {
	err := shim.Start(new(CustomerChaincode))
	if err != nil {
		fmt.Printf("Error starting Customer Simple chaincode: %s", err)
	}
}
