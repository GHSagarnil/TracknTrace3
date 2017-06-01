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
	AssemblyDate string `json:"assemblyDate"` // New
	AssemblyCreationDate string `json:"assemblyCreationDate"`
	AssemblyLastUpdatedOn string `json:"assemblyLastUpdateOn"`
	AssemblyCreatedBy string `json:"assemblyCreatedBy"`
	AssemblyLastUpdatedBy string `json:"assemblyLastUpdatedBy"`
	}


//AssemblyID Holder
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
	ShippingToAddress string `json:"shippingToAddress"`
	PackageCreationDate string `json:"packageCreationDate"`
	PackageLastUpdatedOn string `json:"packageLastUpdateOn"`
	PackageCreatedBy string `json:"packageCreatedBy"`
	PackageLastUpdatedBy string `json:"packageLastUpdatedBy"`
	}


//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const   ASSEMBLYLINE_ROLE  =  "assemblyline_role"
const   PACKAGELINE_ROLE   =  "packageline_role"

/* Assembly Section */

type PackageCaseID_Holder struct {
	PackageCaseIDs 	[]string `json:"packageCaseIDs"`
}

//API to create an assembly
func (t *TnT) createAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 14 {
			return nil, errors.New("Incorrect number of arguments. Expecting 14.")
		}
	/*
	user_name := args[13]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)

		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
		}
	}
	*/
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
		_assemblyDate:= args[12]

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
		assem.AssemblyDate = _assemblyDate
		assem.AssemblyCreationDate = _assemblyCreationDate
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assem.AssemblyCreatedBy = _assemblyCreatedBy
		assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

		bytes, err := json.Marshal(assem)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }

	/* GetAll changes-------------------------starts--------------------------*/
		// Holding the AssemblyIDs in State separately
		bytesAssemHolder, err := stub.GetState("Assemblies")
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

		var assemID_Holder AssemblyID_Holder

		err = json.Unmarshal(bytesAssemHolder, &assemID_Holder)
		if err != nil {	return nil, errors.New("Corrupt Assemblies record") }

		assemID_Holder.AssemblyIDs = append(assemID_Holder.AssemblyIDs, _assemblyId)
		
		bytesAssemHolder, err = json.Marshal(assemID_Holder)
		if err != nil { return nil, errors.New("Error creating Assembly_Holder record") }

		err = stub.PutState("Assemblies", bytesAssemHolder)
		if err != nil { return nil, errors.New("Unable to put the state") }
	/* GetAll changes---------------------------ends------------------------ */
		
		fmt.Println("Created Assembly successfully")
		
		return nil, nil

}

//Update Assembly based on Id - All except AssemblyId, DeviceSerialNo,DeviceType and AssemblyCreationDate and AssemblyCreatedBy
func (t *TnT) updateAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 13 {
		return nil, errors.New("Incorrect number of arguments. Expecting 13.")
	} 
	
		_assemblyId := args[0]
		_deviceSerialNo:= args[1]
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
		_assemblyDate:= args[12] 
		
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

		//update the AssemblyLine 
		//assem.AssemblyId = _assemblyId
		assem.DeviceSerialNo = _deviceSerialNo
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
		assem.AssemblyDate = _assemblyDate
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


//Update Assembly based on Id - AssemblyStatus
func (t *TnT) updateAssemblyStatusByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	} 
	
		_assemblyId := args[0]
		_assemblyStatus:= args[1]
		
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assem := AssemblyLine{}
		json.Unmarshal(assemblyAsBytes, &assem)

		//update the AssemblyLine status
		assem.AssemblyStatus = _assemblyStatus
		assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
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

//get all Assemblies
func (t *TnT) getAllAssemblies(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}

	user_name := args[0]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := stub.GetState(user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	var assemID_Holder AssemblyID_Holder

	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }

	res2E:= []*AssemblyLine{}	

	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the existing AssemblyLine
		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}

		if assemblyAsBytes != nil { 
		res := new(AssemblyLine)
		json.Unmarshal(assemblyAsBytes, &res)

		// Append Assembly to Assembly Array
		res2E=append(res2E,res)
		} // If ends
		} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}


/* Package section*/

//API to create an Package
// Assemblies related to the package is updated with status = PACKAGED
func (t *TnT) createPackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 8 { return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8. Got: %d.", len(args)) }

		_caseId := args[0]
		_holderAssemblyId := args[1]
		_chargerAssemblyId := args[2]
		_packageStatus := args[3]
		_packagingDate := args[4]
		_shippingToAddress := args[5]
		// Status of associated Assemblies	
		_assemblyStatus:= args[6]

		_time:= time.Now().Local()

		_packageCreationDate := _time.Format("2006-01-02")
		_packageLastUpdatedOn := _time.Format("2006-01-02")
		_packageCreatedBy := ""
		_packageLastUpdatedBy := ""

	//Checking if the Package already exists
		packageAsBytes, err := stub.GetState(_caseId)
		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes != nil { return nil, errors.New("Package already exists") }

		//setting the Package to create
		pack := PackageLine{}
		pack.CaseId = _caseId
		pack.HolderAssemblyId = _holderAssemblyId
		pack.ChargerAssemblyId = _chargerAssemblyId
		pack.PackageStatus = _packageStatus
		pack.PackagingDate = _packagingDate
		pack.ShippingToAddress = _shippingToAddress
		pack.PackageCreationDate = _packageCreationDate
		pack.PackageLastUpdatedOn = _packageLastUpdatedOn
		pack.PackageCreatedBy = _packageCreatedBy
		pack.PackageLastUpdatedBy = _packageLastUpdatedBy

		bytes, err := json.Marshal(pack)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }

		err = stub.PutState(_caseId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }

		fmt.Println("Created Package successfully")

		//Update Holder Assemblies to Packaged status
		if 	len(_holderAssemblyId) > 0	{
		//_assemblyStatus:= "PACKAGED"
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyHolderAsBytes, err := stub.GetState(_holderAssemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyHolderAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assemHolder := AssemblyLine{}
		json.Unmarshal(assemblyHolderAsBytes, &assemHolder)

		//update the AssemblyLine status
		assemHolder.AssemblyStatus = _assemblyStatus
		assemHolder.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assemHolder.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

		
		bytesHolder, err := json.Marshal(assemHolder)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_holderAssemblyId, bytesHolder)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
		}

		//Update Charger Assemblies to Packaged status
		if 	len(_chargerAssemblyId) > 0		{
		//_assemblyStatus:= "PACKAGED"
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyChargerAsBytes, err := stub.GetState(_chargerAssemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyChargerAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assemCharger := AssemblyLine{}
		json.Unmarshal(assemblyChargerAsBytes, &assemCharger)

		//update the AssemblyLine status
		assemCharger.AssemblyStatus = _assemblyStatus
		assemCharger.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assemCharger.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

		
		bytesCharger, err := json.Marshal(assemCharger)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_chargerAssemblyId, bytesCharger)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
		}

	/* GetAll changes-------------------------starts--------------------------*/
		// Holding the PackageCaseIDs in State separately
		bytesPackageCaseHolder, err := stub.GetState("Packages")
		if err != nil { return nil, errors.New("Unable to get Packages") }

		var packageCaseID_Holder PackageCaseID_Holder

		err = json.Unmarshal(bytesPackageCaseHolder, &packageCaseID_Holder)
		if err != nil {	return nil, errors.New("Corrupt Packages record") }

		packageCaseID_Holder.PackageCaseIDs = append(packageCaseID_Holder.PackageCaseIDs, _caseId)
		
		bytesPackageCaseHolder, err = json.Marshal(packageCaseID_Holder)
		if err != nil { return nil, errors.New("Error creating PackageCaseID_Holder record") }

		err = stub.PutState("Packages", bytesPackageCaseHolder)
		if err != nil { return nil, errors.New("Unable to put the state") }
	/* GetAll changes---------------------------ends------------------------ */


		return nil, nil

}


//API to update an Package
// Assemblies related to the package is updated with status sent as parameter
func (t *TnT) updatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 7 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 7. Got: %d.", len(args))
		}

		
		_caseId := args[0]
		//_holderAssemblyId := args[1]
		//_chargerAssemblyId := args[2]
		_packageStatus := args[3]
		//_packagingDate := args[4]
		_shippingToAddress := args[5]
		// Status of associated Assemblies	
		_assemblyStatus := args[6]

		_time:= time.Now().Local()

		//_packageCreationDate := _time.Format("2006-01-02")
		_packageLastUpdatedOn := _time.Format("2006-01-02")
		//_packageCreatedBy := ""
		_packageLastUpdatedBy := ""


	//Checking if the Package already exists
		packageAsBytes, err := stub.GetState(_caseId)
		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes == nil { return nil, errors.New("Package doesn't exists") }

		//setting the Package to update
		pack := PackageLine{}
		json.Unmarshal(packageAsBytes, &pack)

		//pack.CaseId = _caseId
		//pack.HolderAssemblyId = _holderAssemblyId
		//pack.ChargerAssemblyId = _chargerAssemblyId
		pack.PackageStatus = _packageStatus
		//pack.PackagingDate = _packagingDate
		pack.ShippingToAddress = _shippingToAddress
		//pack.PackageCreationDate = _packageCreationDate
		pack.PackageLastUpdatedOn = _packageLastUpdatedOn
		//pack.PackageCreatedBy = _packageCreatedBy
		pack.PackageLastUpdatedBy = _packageLastUpdatedBy

		// Getting associate Assembly IDs
		_holderAssemblyId := pack.HolderAssemblyId
		_chargerAssemblyId := pack.ChargerAssemblyId


		bytes, err := json.Marshal(pack)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }

		err = stub.PutState(_caseId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }

		fmt.Println("Created Package successfully")

		//Update Holder Assemblies status
		if 	len(_holderAssemblyId) > 0	{
		//_assemblyStatus:= "PACKAGED"
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyHolderAsBytes, err := stub.GetState(_holderAssemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyHolderAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assemHolder := AssemblyLine{}
		json.Unmarshal(assemblyHolderAsBytes, &assemHolder)

		//update the AssemblyLine status
		assemHolder.AssemblyStatus = _assemblyStatus
		assemHolder.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assemHolder.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

		
		bytesHolder, err := json.Marshal(assemHolder)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_holderAssemblyId, bytesHolder)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
		}

		//Update Charger Assemblies status
		if 	len(_chargerAssemblyId) > 0		{
		//_assemblyStatus:= "PACKAGED"
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("2006-01-02")
		_assemblyLastUpdatedBy := ""

		//get the Assembly
		assemblyChargerAsBytes, err := stub.GetState(_chargerAssemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyChargerAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assemCharger := AssemblyLine{}
		json.Unmarshal(assemblyChargerAsBytes, &assemCharger)

		//update the AssemblyLine status
		assemCharger.AssemblyStatus = _assemblyStatus
		assemCharger.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
		assemCharger.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
		
		bytesCharger, err := json.Marshal(assemCharger)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_chargerAssemblyId, bytesCharger)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }
		}

		return nil, nil

}


//get the Package against ID
func (t *TnT) getPackageByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting CaseId to query")
	}

	_caseId := args[0]
	
	//get the var from chaincode state
	valAsbytes, err := stub.GetState(_caseId)									
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " +  _caseId  + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil	

}

//get all Packages
func (t *TnT) getAllPackages(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	
	if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1.")
		}
	user_name := args[0]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
		}
	}
	
	
	bytesPackageCaseHolder, err := stub.GetState("Packages")
	if err != nil { return nil, errors.New("Unable to get Packages") }

	var packageCaseID_Holder PackageCaseID_Holder

	err = json.Unmarshal(bytesPackageCaseHolder, &packageCaseID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }

	res2E:= []*PackageLine{}	

	for _, caseId := range packageCaseID_Holder.PackageCaseIDs {

		//Get the existing AssemblyLine
		packageAsBytes, err := stub.GetState(caseId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}

		if packageAsBytes != nil { 
		res := new(PackageLine)
		json.Unmarshal(packageAsBytes, &res)

		// Append Assembly to Assembly Array
		res2E=append(res2E,res)
		} // If ends
		} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}


//AllAssemblyIDS
//get the all Assembly IDs from AssemblyID_Holder - To Test only
func (t *TnT) getAllAssemblyIDs(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting zero argument to query")
	}

	bytesAssemHolder, err := stub.GetState("Assemblies")
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

	return bytesAssemHolder, nil	

}

//AllPackageCaseIDs
//get the all Package CaseIDs from PackageCaseID_Holder - To Test only
func (t *TnT) getAllPackageCaseIDs(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting zero argument to query")
	}

	bytesPackageCaseHolder, err := stub.GetState("Packages")
		if err != nil { return nil, errors.New("Unable to get Packages") }

	return bytesPackageCaseHolder, nil	

}

//Security & Access

//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *TnT) get_ecert(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {

	ecert, err := stub.GetState(name)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	return ecert, nil
}

//==============================================================================================================================
//	 add_ecert - Adds a new ecert and user pair to the table of ecerts
//==============================================================================================================================

func (t *TnT) add_ecert(stub shim.ChaincodeStubInterface, name string, ecert string) ([]byte, error) {

	err := stub.PutState(name, []byte(ecert))

	if err == nil {
		return nil, errors.New("Error storing eCert for user " + name + " identity: " + ecert)
	}

	return nil, nil

}


/*Standard Calls*/

// Init initializes the smart contracts
func (t *TnT) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	/* GetAll changes-------------------------starts--------------------------*/
	/*
		var assemID_Holder AssemblyID_Holder
		// Adding a dummy assembly to test
		if len(args) != 0 {
				_assemblyId := args[0]
				assemID_Holder.AssemblyIDs = append(assemID_Holder.AssemblyIDs, _assemblyId)
			}

		bytesAssembly, err := json.Marshal(assemID_Holder)
		if err != nil { return nil, errors.New("Error creating assemID_Holder record") }
		err = stub.PutState("Assemblies", bytesAssembly)

		var packageCaseID_Holder PackageCaseID_Holder
		// Adding a dummy package to test
		if len(args) != 0 {
				_packageId := args[0]
				packageCaseID_Holder.PackageCaseIDs = append(packageCaseID_Holder.PackageCaseIDs, _packageId)
			}
		bytesPackage, err := json.Marshal(packageCaseID_Holder)
		if err != nil { return nil, errors.New("Error creating packageCaseID_Holder record") }
		err = stub.PutState("Packages", bytesPackage)
	*/	
	/* GetAll changes---------------------------ends------------------------ */

	// creating minimum default user and roles
	//
	for i:=0; i < len(args); i=i+2 {
		t.add_ecert(stub, args[i], args[i+1])
	}


	return nil, nil

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
	}  else if function == "createPackage" {
		fmt.Printf("Function is createPackage")
		return t.createPackage(stub, args)
	} else if function == "updatePackage" {
		fmt.Printf("Function is updatePackage")
		return t.updatePackage(stub, args)
	} else if function == "add_ecert" {
		fmt.Printf("Function is add_ecert")
		return t.add_ecert(stub, args[0], args[1])
	} 
	return nil, errors.New("Received unknown function invocation")
}

// query queries the chaincode
func (t *TnT) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")

	if function == "getAssemblyByID" { 
		t := TnT{}
		return t.getAssemblyByID(stub, args)
	} else if function == "getPackageByID" { 
		t := TnT{}
		return t.getPackageByID(stub, args)
	} else if function == "getAllAssemblies" { 
		t := TnT{}
		return t.getAllAssemblies(stub, args)
	} else if function == "getAllPackages" { 
		t := TnT{}
		return t.getAllPackages(stub, args)
	} else if function == "getAllAssemblyIDs" { 
		t := TnT{}
		return t.getAllAssemblyIDs(stub, args)
	} else if function == "getAllPackageCaseIDs" { 
		t := TnT{}
		return t.getAllPackageCaseIDs(stub, args)
	} else if function == "get_ecert" {
		return t.get_ecert(stub, args[0])
	} 


	return nil, errors.New("Received unknown function query")
}

//main function
func main() {
	err := shim.Start(new(TnT))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
