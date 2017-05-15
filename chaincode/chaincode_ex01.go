/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"time"
	"strconv"
	
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

// TnT is a high level smart contract that collaborate together business artifact based smart contracts
type TnT struct {
}
var assemblyIndexStr = "_assemblyIndex" // Store Key value pair for Assembly
// Assembly Line Structure
type AssemblyLine struct{	
	AssemblyId string `json:"assemblyId"`
	DeviceSerialNo string `json:"deviceSerialNo"`
	DeviceType string `json:"deviceType"`
	//FilamentBatchId string `json:"filamentBatchId"`
	//LedBatchId string `json:"ledBatchId"`
	//CircuitBoardBatchId string `json:"circuitBoardBatchId"`
	//WireBatchId string `json:"wireBatchId"`
	//CasingBatchId string `json:"casingBatchId"`
	//AdaptorBatchId string `json:"adaptorBatchId"`
	//StickPodBatchId string `json:"stickPodBatchId"`
	//ManufacturingPlant string `json:"manufacturingPlant"`
	AssemblyStatus string `json:"assemblyStatus"`
	//AssemblyCreationDate string `json:"assemblyCreationDate"`
	AssemblyLastUpdatedOn string `json:"assemblyLastUpdateOn"`
	//AssemblyCreatedBy string `json:"assemblyCreatedBy"`
	//AssemblyLastUpdatedBy string `json:"assemblyLastUpdatedBy"`
	}

// Package Line Structure
type PackageLine struct{	
	CaseId string `json:"caseId"`
	HolderAssemblyId string `json:"holderAssemblyId"`
	ChargerAssemblyId string `json:"chargerAssemblyId"`
	PackageStatus string `json:"packageStatus"`
	PackagingDate string `json:"packagingDate"`
	PackageCreationDate string `json:"packagingCreationDate"`
	PackageLastUpdatedOn string `json:"packageLastUpdateOn"`
	ShippingToAddress string `json:"shippingToAddress"`
	PackageCreatedBy string `json:"packageCreatedBy"`
	PackageLastUpdatedBy string `json:"packageLastUpdatedBy"`
	}

// Init initializes the smart contracts
func (t *TnT) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var _temp int;
	var err error

	if len(args) != 1 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 1. Got: %d.", len(args))
		}

		// Initialize the chaincode
	_temp, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value ")
	}
	// Write the state to the ledger
	err = stub.PutState("12345678", []byte(strconv.Itoa(_temp)))				
	if err != nil {
		return nil, err
	}
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(assemblyIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}
//API to create an assembly
func (t *TnT) createAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
if len(args) != 4 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 11. Got: %d.", len(args))
		}

		//var columns []shim.Column
		//_assemblyId:= rand.New(rand.NewSource(99)).Int31

		//Generate the AssemblyId
		//rand.Seed(time.Now().Unix())
		
		//_assemblyId := strconv.Itoa(rand.Int())
		_assemblyId := args[0]
		_deviceSerialNo:= args[1]
		_deviceType:=args[2]
		//_FilamentBatchId:=args[2]
		//_LedBatchId:=args[3]
		//_CircuitBoardBatchId:=args[4]
		//_WireBatchId:=args[5]
		//_CasingBatchId:=args[6]
		//_AdaptorBatchId:=args[7]
		//_StickPodBatchId:=args[8]
		//_ManufacturingPlant:=args[9]
		_AssemblyStatus:= args[0]

		_time:= time.Now().Local()

		//_AssemblyCreationDate := _time.Format("2006-01-02")
		_AssemblyLastUpdateOn := _time.Format("2006-01-02")
		//_AssemblyCreatedBy := ""
		//_AssemblyLastUpdatedBy := ""

	//check if marble already exists
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil {
		return nil, errors.New("Failed to get assembly Id")
		}
		res := AssemblyLine{}
		json.Unmarshal(assemblyAsBytes, &res)
		if res.AssemblyId == _assemblyId{
		fmt.Println("This Assembly arleady exists: " + _assemblyId)
		fmt.Println(res);
		return nil, errors.New("This Assembly arleady exists")				//all stop an Assembly already exists
		}


		str := `{"assemblyId": "` + _assemblyId + `", "deviceSerialNo": "` + _deviceSerialNo + `", "deviceType": "` + _deviceType + `", "assemblyStatus": "`+ _AssemblyStatus +`", "assemblyLastUpdateOn": "` + _AssemblyLastUpdateOn + `"}`
		
		err = stub.PutState(_assemblyId, []byte(str))								//store assembly with id as key
		if err != nil {
		return nil, err
		}

		//get the assembly index
		assemblyAsBytes, err = stub.GetState(assemblyIndexStr)
		if err != nil {
			return nil, errors.New("Failed to get assembly index")
		}
		var assemblyIndex []string
		json.Unmarshal(assemblyAsBytes, &assemblyIndex)							//un stringify it aka JSON.parse()
		
		//append
		assemblyIndex = append(assemblyIndex, _assemblyId)								//add assembly id in Index list
		fmt.Println("! Assembly index: ", assemblyIndex)
		jsonAsBytes, _ := json.Marshal(assemblyIndex)
		err = stub.PutState(assemblyIndexStr, jsonAsBytes)						//store assembly

		fmt.Println("Create Assembly")
			
		return nil, nil

}

//Update Assembly based on Id (Now only status)
func (t *TnT) updateAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	} 
	
		_assemblyId := args[0]
		//_deviceSerialNo:= args[1]
		//_deviceType:=args[2]
		//_FilamentBatchId:=args[3]
		//_LedBatchId:=args[4]
		//_CircuitBoardBatchId:=args[5]
		//_WireBatchId:=args[6]
		//_CasingBatchId:=args[7]
		//_AdaptorBatchId:=args[8]
		//_StickPodBatchId:=args[9]
		//_ManufacturingPlant:=args[10]
		_AssemblyStatus:= args[1]
		//_AssemblyCreationDate := args[12]
		//_AssemblyCreatedBy :=  args[13]
		_time:= time.Now().Local()
		_AssemblyLastUpdateOn := _time.Format("2006-01-02")
		//_AssemblyLastUpdatedBy := ""
		str := `{ "assemblyStatus": "` + _AssemblyStatus + `", "assemblyLastUpdateOn": "` +  _AssemblyLastUpdateOn  + `"}`
		err := stub.PutState(_assemblyId, []byte(str))								//write the status into the chaincode state
		if err != nil {
		return nil, err
		}
	
		
	return nil, nil

}


//get the Assembly against ID
func (t *TnT) getAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting AssemblyID to query")
	}

	_assemblyId := args[0]
	
	valAsbytes, err := stub.GetState(_assemblyId)									//get the var from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " +  _assemblyId  + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil	

}



// Invoke callback representing the invocation of a chaincode
func (t *TnT) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")
	
	// Handle different functions
	if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "createAssembly" {
		fmt.Printf("Function is createAssembly")
		return t.createAssembly(stub, args)
	} else if function == "updateAssemblyByID" {
		fmt.Printf("Function is updateAssemblyByID")
		return t.updateAssemblyByID(stub, args)
	} 
	return nil, errors.New("Received unknown function invocation")
}


// query queries the chaincode
func (t *TnT) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function == "getAssemblyByID" { 
		t := TnT{}
		return t.getAssemblyByID(stub, args)
	}
	
	return nil, errors.New("Received unknown function query")
}

	func main() {
	err := shim.Start(new(TnT))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
