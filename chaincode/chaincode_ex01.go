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
	//"strconv"
	
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
	FilamentBatchId string `json:"filamentBatchId"`
	LedBatchId string `json:"ledBatchId"`
	CircuitBoardBatchId string `json:"circuitBoardBatchId"`
	WireBatchId string `json:"wireBatchId"`
	CasingBatchId string `json:"casingBatchId"`
	AdaptorBatchId string `json:"adaptorBatchId"`
	StickPodBatchId string `json:"stickPodBatchId"`
	ManufacturingPlant string `json:"manufacturingPlant"`
	AssemblyStatus string `json:"assemblyStatus"`
	AssemblyCreationDate string `json:"assemblyCreationDate"`
	AssemblyLastUpdatedOn string `json:"assemblyLastUpdateOn"`
	AssemblyCreatedBy string `json:"assemblyCreatedBy"`
	AssemblyLastUpdatedBy string `json:"assemblyLastUpdatedBy"`
	}

type AssemblyID_Holder struct {
	AssemblyIDs 	[]string `json:"assemblyIDs"`
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

/*	var _temp int;
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
*/	
	return nil, nil
}

//API to create an assembly
func (t *TnT) createAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 12 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 12. Got: %d.", len(args))
		}

		//var columns []shim.Column
		//_assemblyId:= rand.New(rand.NewSource(99)).Int31

		//Generate the AssemblyId
		//rand.Seed(time.Now().Unix())
		
		//_assemblyId := strconv.Itoa(rand.Int())
		_assemblyId := args[0]
		_deviceSerialNo:= args[1]
		_deviceType:=args[2]
		_filamentBatchId:=args[3]
		_ledBatchId:=args[4]
		_circuitBoardBatchId:=args[5]
		_wireBatchId:=args[6]
		_casingBatchId:=args[7]
		_adaptorBatchId:=args[8]
		_stickPodBatchId:=args[9]
		_manufacturingPlant:=args[10]
		_assemblyStatus:= args[11]

		_time:= time.Now().Local()

		_assemblyCreationDate := _time.Format("2006-01-02")
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		_assemblyCreatedBy := ""
		_assemblyLastUpdatedBy := ""

	//Checking if the Assembly already exists
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil { return nil, errors.New("Failed to get assembly Id") }
		if assemblyAsBytes != nil { return nil, errors.New("Assembly already exists") }

		//setting the AssemblyLine to create
		assem := AssemblyLine{}
		assem.AssemblyId = _assemblyId
		assem.DeviceSerialNo = _deviceSerialNo
		assem.DeviceType = _deviceType
		assem.FilamentBatchId = _filamentBatchId
		assem.LedBatchId = _ledBatchId
		assem.CircuitBoardBatchId = _circuitBoardBatchId
		assem.WireBatchId = _wireBatchId
		assem.CasingBatchId = _casingBatchId
		assem.AdaptorBatchId = _adaptorBatchId
		assem.StickPodBatchId = _stickPodBatchId
		assem.ManufacturingPlant = _manufacturingPlant
		assem.AssemblyStatus = _assemblyStatus
		assem.AssemblyCreationDate = _assemblyCreationDate
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assem.AssemblyCreatedBy = _assemblyCreatedBy
		assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

		bytes, err := json.Marshal(assem)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }

	/*
		// Holding the AssemblyIDs in State separately
		bytes, err = stub.GetState("Assemblies")
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

		var assemIDs AssemblyID_Holder

		err = json.Unmarshal(bytes, &assemIDs)
		if err != nil {	return nil, errors.New("Corrupt Assemblies record") }

		assemIDs.AssemblyIDs = append(assemIDs.AssemblyIDs, _assemblyId)
		
		bytes, err = json.Marshal(assemIDs)
		if err != nil { return nil, errors.New("Error creating Assembly_Holder record") }

		err = stub.PutState("Assemblies", bytes)
		if err != nil { return nil, errors.New("Unable to put the state") }
*/
		fmt.Println("Created Assembly successfully")
		
		return nil, nil

}

//Update Assembly based on Id (Now only status)
func (t *TnT) updateAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12.")
	} 
	
		_assemblyId := args[0]
		//_deviceSerialNo:= args[1] - No Change
		//_deviceType:=args[2] - No Change
		_filamentBatchId:=args[3]
		_ledBatchId:=args[4]
		_circuitBoardBatchId:=args[5]
		_wireBatchId:=args[6]
		_casingBatchId:=args[7]
		_adaptorBatchId:=args[8]
		_stickPodBatchId:=args[9]
		_manufacturingPlant:=args[10]
		_assemblyStatus:= args[11]
		
		_time:= time.Now().Local()
		//_assemblyCreationDate - No change
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		//_assemblyCreatedBy - No change
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assem := AssemblyLine{}
		json.Unmarshal(assemblyAsBytes, &assem)

		//update the AssemblyLine status
		assem.AssemblyId = _assemblyId
		//assem.DeviceSerialNo = _deviceSerialNo
		//assem.DeviceType = _deviceType
		assem.FilamentBatchId = _filamentBatchId
		assem.LedBatchId = _ledBatchId
		assem.CircuitBoardBatchId = _circuitBoardBatchId
		assem.WireBatchId = _wireBatchId
		assem.CasingBatchId = _casingBatchId
		assem.AdaptorBatchId = _adaptorBatchId
		assem.StickPodBatchId = _stickPodBatchId
		assem.ManufacturingPlant = _manufacturingPlant
		assem.AssemblyStatus = _assemblyStatus
		//assem.AssemblyCreationDate = _assemblyCreationDate
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		//assem.AssemblyCreatedBy = _assemblyCreatedBy
		assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

		
		bytes, err := json.Marshal(assem)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }

		return nil, nil
			
}


//get the Assembly against ID
func (t *TnT) getAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting AssemblyID to query")
	}

	_assemblyId := args[0]
	
	//get the var from chaincode state
	valAsbytes, err := stub.GetState(_assemblyId)									
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " +  _assemblyId  + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil	

}



func (t *TnT) getAssemblies(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	var assemblyIDs AssemblyID_Holder

	err = json.Unmarshal(bytes, &assemblyIDs)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }

	res2E:= []*AssemblyLine{}	

	for _, assemblyId := range assemblyIDs.AssemblyIDs {

		//Get the existing AssemblyLine
		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")	}

		res := new(AssemblyLine)
		json.Unmarshal(assemblyAsBytes, &res)

		// Append Assembly to Assembly Array
		res2E=append(res2E,res)
		}

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
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
