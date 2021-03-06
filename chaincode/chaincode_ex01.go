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


//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert and Specific Assembly Statuses
//==============================================================================================================================
const   ASSEMBLYLINE_ROLE  		=	"assemblyline_role"
const   PACKAGELINE_ROLE   		=	"packageline_role"
const 	QA_VIEWER_ROLE 			= 	"qaviewer_role" // new role for viewing purpose
const   ASSEMBLYSTATUS_RFP   	=	"6" //Ready For Packaging"
const  	ASSEMBLYSTATUS_PKG 		=	"7" //Packaged" 
const  	ASSEMBLYSTATUS_CAN 		=	"8" //Cancelled"
const  	ASSEMBLYSTATUS_QAF 		=	"2" //QA Failed"
const   FIL_BATCH  				=	"FilamentBatchId"	
const   LED_BATCH  				=	"LedBatchId"
const   CIR_BATCH  				=	"CircuitBoardBatchId"
const   WRE_BATCH  				=	"WireBatchId"
const   CAS_BATCH  				=	"CasingBatchId"
const   ADP_BATCH  				=	"AdaptorBatchId"
const   STK_BATCH  				=	"StickPodBatchId"
const   HLD_ASSMB_TYP  			=	"HolderAssemblyId"
const 	CHG_ASSMB_TYP 			= 	"ChargerAssemblyId"


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
	AssemblyPackage string `json:"assemblyPackage"`
	AssemblyInfo1 string `json:"assemblyInfo1"`
	AssemblyInfo2 string `json:"assemblyInfo2"`
	//_assemblyPackage,_assemblyInfo1,_assemblyInfo2
	}


//AssemblyID Holder
type AssemblyID_Holder struct {
	AssemblyIDs 	[]string `json:"assemblyIDs"`
}

//AssemblyLine Holder
type AssemblyLine_Holder struct {
	AssemblyLines 	[]AssemblyLine `json:"assemblyLines"`
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
	PackageInfo1 string `json:"packageInfo1"`
	PackageInfo2 string `json:"packageInfo2"`
	}




/* Assembly Section */

type PackageCaseID_Holder struct {
	PackageCaseIDs 	[]string `json:"packageCaseIDs"`
}

//PackageLine Holder
type PackageLine_Holder struct {
	PackageLines 	[]PackageLine `json:"packageLines"`
}


//API to create an assembly
//"args": [ "ASM0101","DEV0101","HOLDER","FIL0002","LED0002","CIR0002","WIR0002","CAS0002","ADA0002","STK0002","MAN0002","1","20170608","aluser1"]
//_assemblyId,_deviceSerialNo,_deviceType,_filamentBatchId,_ledBatchId,_circuitBoardBatchId,_wireBatchId,_casingBatchId,_adaptorBatchId,_stickPodBatchId,_manufacturingPlant,_assemblyStatus _assemblyDate,_assemblyPackage,_assemblyInfo1,_assemblyInfo2 ,user_name
func (t *TnT) createAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 17. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

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
		_assemblyPackage:= args[13]
		_assemblyInfo1:= args[14]
		_assemblyInfo2:= args[15]
		_time:= time.Now().Local()

		_assemblyCreationDate := _time.Format("20060102150405")
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		_assemblyCreatedBy := user_name
		_assemblyLastUpdatedBy := user_name

	//Check Date
	if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
	//Checking if the Assembly already exists
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil { return nil, errors.New("Failed to get assembly Id") }
		if assemblyAsBytes != nil { return nil, errors.New("Assembly already exists") }



		/* AssemblyLine history -----------------Starts */
		var assemLine_HolderInit AssemblyLine_Holder

		assemLine_HolderKey := _assemblyId + "H" // Indicates history key
		bytesAssemblyLinesInit, err := json.Marshal(assemLine_HolderInit)
		if err != nil { return nil, errors.New("Error creating assemID_Holder record") }
		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLinesInit)
		/* AssemblyLine history -----------------Ends */


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
		assem.AssemblyPackage = _assemblyPackage
		assem.AssemblyInfo1 = _assemblyInfo1
		assem.AssemblyInfo2 = _assemblyInfo2
		
		

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

		/* AssemblyLine history ------------------------------------------Starts */
		bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

		var assemLine_Holder AssemblyLine_Holder

		err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

		assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the newly created AssemblyLine
		
		bytesAssemblyLines, err = json.Marshal(assemLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		/* AssemblyLine history ------------------------------------------Ends */
		
		fmt.Println("Created Assembly successfully")
		
		return nil, nil

}

//Update Assembly based on Id - All except AssemblyId, DeviceSerialNo,DeviceType and AssemblyCreationDate and AssemblyCreatedBy
//"args": [ "ASM0101","DEV0101","HOLDER","FIL0002","LED0002","CIR0002","WIR0002","CAS0002","ADA0002","STK0002","MAN0002","1","20170608","CASE0001","INFO1","INFO2"aluser1"]
//_assemblyId,_deviceSerialNo,_deviceType,_filamentBatchId,_ledBatchId,_circuitBoardBatchId,_wireBatchId,_casingBatchId,_adaptorBatchId,_stickPodBatchId,_manufacturingPlant,_assemblyStatus _assemblyDate,_assemblyPackage,_assemblyInfo1,_assemblyInfo2 ,user_name
func (t *TnT) updateAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {
		return nil, errors.New("Incorrect number of arguments. Expecting 17.")
	} 
	
	/* Access check -------------------------------------------- Starts*/
	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}	

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	
	/* Access check -------------------------------------------- Ends*/

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
		_assemblyPackage:= args[13]
		_assemblyInfo1:= args[14]
		_assemblyInfo2:= args[15]
		
		_time:= time.Now().Local()
		//_assemblyCreationDate - No change
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		//_assemblyCreatedBy - No change
		_assemblyLastUpdatedBy := user_name

		//Check Date
		if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
	
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
		assem.AssemblyPackage = _assemblyPackage
		assem.AssemblyInfo1 = _assemblyInfo1
		assem.AssemblyInfo2 = _assemblyInfo2


		
		bytes, err := json.Marshal(assem)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

		err = stub.PutState(_assemblyId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


		/* AssemblyLine history ------------------------------------------Starts */
		// assemLine_HolderKey := _assemblyId + "H" // Indicates history key
		assemLine_HolderKey := _assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

		var assemLine_Holder AssemblyLine_Holder

		err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

		assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the updated AssemblyLine
		
		bytesAssemblyLines, err = json.Marshal(assemLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		/* AssemblyLine history ------------------------------------------Ends */

		return nil, nil
			
}


//Update Assembly based on Id - AssemblyStatus
func (t *TnT) updateAssemblyStatusByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	} 
	
	/* Access check -------------------------------------------- Starts*/
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}	

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

		_assemblyId := args[0]
		_assemblyStatus:= args[1]
		
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		_assemblyLastUpdatedBy := user_name

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

		/* AssemblyLine history ------------------------------------------Starts */
		// assemLine_HolderKey := _assemblyId + "H" // Indicates history key
		assemLine_HolderKey := _assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

		var assemLine_Holder AssemblyLine_Holder

		err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

		assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the updated AssemblyLine
		
		bytesAssemblyLines, err = json.Marshal(assemLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

		err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		/* AssemblyLine history ------------------------------------------Ends */

		return nil, nil
			
}

//Update Assembly Info2 - HashCode based on Id 
// Parameters = ASM0001, HASCODE, USERNAME
func (t *TnT) updateAssemblyInfo2ByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	} 
	
	/* Access check -------------------------------------------- Starts*/
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}	

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

		_assemblyId := args[0]
		_assemblyInfo2:= args[1]
		
		_time:= time.Now().Local()
		_assemblyLastUpdatedOn := _time.Format("20060102150405")
		_assemblyLastUpdatedBy := user_name

		//get the Assembly
		assemblyAsBytes, err := stub.GetState(_assemblyId)
		if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
		if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

		assem := AssemblyLine{}
		json.Unmarshal(assemblyAsBytes, &assem)

		// Update Assembly only if the hashcode is not already set
		if len(assem.AssemblyInfo2) == 0 {
			//update the AssemblyLine Info2
			assem.AssemblyInfo2 = _assemblyInfo2
			assem.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
			assem.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy

			
			bytes, err := json.Marshal(assem)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

			err = stub.PutState(_assemblyId, bytes)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }

			/* AssemblyLine history ------------------------------------------Starts */
			// For HashCode update don't store an Assembly History but update the last History with Info2
			// assemLine_HolderKey := _assemblyId + "H" // Indicates history key
			assemLine_HolderKey := _assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
			bytesAssemblyLines, err := stub.GetState(assemLine_HolderKey)
			if err != nil { return nil, errors.New("Unable to get Assemblies") }

			var assemLine_Holder AssemblyLine_Holder

			err = json.Unmarshal(bytesAssemblyLines, &assemLine_Holder)
			if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

			//assemLine_Holder.AssemblyLines = append(assemLine_Holder.AssemblyLines, assem) //appending the updated AssemblyLine

			//Overwrite exisitng last element with the updated element - Don't apend but overwrite
			latestIndex := len(assemLine_Holder.AssemblyLines)
			assemLine_Holder.AssemblyLines[latestIndex-1] = assem
			
			bytesAssemblyLines, err = json.Marshal(assemLine_Holder)
			if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

			err = stub.PutState(assemLine_HolderKey, bytesAssemblyLines)
			if err != nil { return nil, errors.New("Unable to put the state") }
			/* AssemblyLine history ------------------------------------------Ends */
		} // AssemblyInfo2 lenght check ends	

		return nil, nil
			
}

//get the Assembly against ID
func (t *TnT) getAssemblyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments to query")
	}

	_assemblyId := args[0]
	user_name:= args[1]	
	/* Access check -------------------------------------------- Starts*/
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not an AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

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
	
	/* Access check -------------------------------------------- Starts*/
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
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

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

//get all Assemblies based on Type & BatchNo
func (t *TnT) getAssembliesByBatchNumber(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	_batchType:= args[0]
	_batchNumber:= args[1]
	_assemblyFlag:= 0

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

			//Check the filter condition
			if 		   _batchType == FIL_BATCH					&&
						res.FilamentBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			} else if  _batchType == LED_BATCH					&&
						res.LedBatchId == _batchNumber			{ 
						_assemblyFlag = 1
			} else if  _batchType == CIR_BATCH					&&
						res.CircuitBoardBatchId == _batchNumber	{ 
						_assemblyFlag = 1
			} else if  _batchType == WRE_BATCH					&&
						res.WireBatchId == _batchNumber			{ 
						_assemblyFlag = 1
			} else if  _batchType == CAS_BATCH					&&
						res.CasingBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			} else if  _batchType == ADP_BATCH					&&
						res.AdaptorBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			} else if  _batchType == STK_BATCH					&&
						res.StickPodBatchId == _batchNumber		{ 
						_assemblyFlag = 1
			}
			

			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)
			if _assemblyFlag == 1 {
				res2E=append(res2E,res)
			}
		} // If ends
		//re-setting the flag to 0
		_assemblyFlag = 0
	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}

//get all Assemblies based on FromDate & ToDate
func (t *TnT) getAssembliesByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64
	
	_fromDate, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	_assemblyFlag:= 0

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64
	
	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }

	res2E:= []*AssemblyLine{}	

	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the existing AssemblyLine History
		assemblyAsBytes, err := stub.GetState(assemblyId)
		if err != nil { return nil, errors.New("Failed to get Assembly")}
		if assemblyAsBytes == nil { return nil, errors.New("Failed to get Assembly")}

		res := new(AssemblyLine)
		json.Unmarshal(assemblyAsBytes, &res)

		//fmt.Printf("%T, %v\n", _fromDate, _fromDate)
		//fmt.Printf("%T, %v\n", _toDate, _toDate)
		//if _fromDate == _toDate { return nil, errors.New("Failed to get Assembly")}
		
		//Check the filter condition YYYYMMDDHHMMSS
		if len(res.AssemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}
		if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err != nil { errors.New ("Error in converting AssemblyDate to int64")}
		if	_assemblyDateInt64 >= _fromDate		&&
			_assemblyDateInt64 <= _toDate		{ 
			_assemblyFlag = 1
		} 
					
		// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)
		if _assemblyFlag == 1 {
			res2E=append(res2E,res)
		}
	//re-setting the flag and AssemblyDate
		_assemblyFlag = 0
		_assemblyDateInt64 = 0
	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}

//get all Assemblies based on Type & BatchNo & From & To Date
func (t *TnT) getAssembliesByBatchNumberAndByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 5 {
			return nil, errors.New("Incorrect number of arguments. Expecting 5.")
		}
	user_name := args[4]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	_batchType:= args[0]
	_batchNumber:= args[1]
	_assemblyFlag:= 0

	_fromDate, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}



	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64

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

			//Check the filter condition
			if len(res.AssemblyDate) == 14 {
				if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err == nil { 
					if	_assemblyDateInt64 >= _fromDate		&&
						_assemblyDateInt64 <= _toDate		{
							if 		   _batchType == FIL_BATCH					&&
										res.FilamentBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == LED_BATCH					&&
										res.LedBatchId == _batchNumber			{ 
										_assemblyFlag = 1
							} else if  _batchType == CIR_BATCH					&&
										res.CircuitBoardBatchId == _batchNumber	{ 
										_assemblyFlag = 1
							} else if  _batchType == WRE_BATCH					&&
										res.WireBatchId == _batchNumber			{ 
										_assemblyFlag = 1
							} else if  _batchType == CAS_BATCH					&&
										res.CasingBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == ADP_BATCH					&&
										res.AdaptorBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == STK_BATCH					&&
										res.StickPodBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							}
						}// from date and to date check
				}// if date parse
			}// if date lenght
			

			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)
			if _assemblyFlag == 1 {
				res2E=append(res2E,res)
			}
		} // If ends
		//re-setting the flag to 0
		_assemblyFlag = 0
		_assemblyDateInt64 = 0
	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}


//get all Assemblies History based on FromDate & ToDate
func (t *TnT) getAssembliesHistoryByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64
	
	_fromDate, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	_assemblyFlag:= 0

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64
	
	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }

	// Array of filtered Assemblies
	res2E:= []AssemblyLine{}	
	// Filtered Assembly
	//res := new(AssemblyLine)
	
	//Looping through the array of assemblyids
	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the AssemblyLine History for each AssemblyID
		assemLine_HolderKey := assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLinesHistoryByID, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesAssemblyLinesHistoryByID") }

		var assemLineHistory_Holder AssemblyLine_Holder

		err = json.Unmarshal(bytesAssemblyLinesHistoryByID, &assemLineHistory_Holder)
		if err != nil {	return nil, errors.New("Corrupt assemLineHistory_Holder record") }

		//Looping through the array of assemblies
		for _, res := range assemLineHistory_Holder.AssemblyLines {
		
			//Check the filter condition YYYYMMDDHHMMSS
			/*
			if len(res.AssemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}
			if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err != nil { errors.New ("Error in converting AssemblyDate to int64")}
			*/
			//Skip if not a valid date YYYYMMDDHHMMSS
			if len(res.AssemblyDate) == 14 {
				if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err == nil { 
					if	_assemblyDateInt64 >= _fromDate		&&
						_assemblyDateInt64 <= _toDate		{ 
						_assemblyFlag = 1
					} 
				}
			}
						
			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)
			if _assemblyFlag == 1 {
				res2E=append(res2E,res)
			}
			
			//re-setting the flag and AssemblyDate
				_assemblyFlag = 0
				_assemblyDateInt64 = 0
		} // For assemLineHistory_Holder.AssemblyLines ends
	} // For assemID_Holder.AssemblyIDs ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}


//get all Assemblies History based on Type & BatchNo & From & To Date
func (t *TnT) getAssembliesHistoryByBatchNumberAndByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 5 {
			return nil, errors.New("Incorrect number of arguments. Expecting 4.")
		}
	user_name := args[4]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not an AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64
	
	_batchType:= args[0]
	_batchNumber:= args[1]
	_assemblyFlag:= 0

	_fromDate, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}

	bytes, err := stub.GetState("Assemblies")
	if err != nil { return nil, errors.New("Unable to get Assemblies") }

	var assemID_Holder AssemblyID_Holder
	var _assemblyDateInt64 int64
	
	err = json.Unmarshal(bytes, &assemID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Assemblies") }

	// Array of filtered Assemblies
	res2E:= []AssemblyLine{}	
	// Filtered Assembly
	//res := new(AssemblyLine)
	
	//Looping through the array of assemblyids
	for _, assemblyId := range assemID_Holder.AssemblyIDs {

		//Get the AssemblyLine History for each AssemblyID
		assemLine_HolderKey := assemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
		bytesAssemblyLinesHistoryByID, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesAssemblyLinesHistoryByID") }

		var assemLineHistory_Holder AssemblyLine_Holder

		err = json.Unmarshal(bytesAssemblyLinesHistoryByID, &assemLineHistory_Holder)
		if err != nil {	return nil, errors.New("Corrupt assemLineHistory_Holder record") }

		//re-setting the flag and AssemblyDate
		_assemblyFlag = 0
		_assemblyDateInt64 = 0

		//Looping through the array of assemblies to check if the filter condition matches - then consider the Assembly for response (latest status only)
		for _, res := range assemLineHistory_Holder.AssemblyLines {
		
			//Check the filter condition YYYYMMDDHHMMSS
			/*
			if len(res.AssemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}
			if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err != nil { errors.New ("Error in converting AssemblyDate to int64")}
			*/
			//Skip if not a valid date YYYYMMDDHHMMSS
			if len(res.AssemblyDate) == 14 {
				if _assemblyDateInt64, err = strconv.ParseInt(res.AssemblyDate, 10, 64); err == nil { 
					if	_assemblyDateInt64 >= _fromDate		&&
						_assemblyDateInt64 <= _toDate		{ 
							if 		   _batchType == FIL_BATCH					&&
										res.FilamentBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == LED_BATCH					&&
										res.LedBatchId == _batchNumber			{ 
										_assemblyFlag = 1
							} else if  _batchType == CIR_BATCH					&&
										res.CircuitBoardBatchId == _batchNumber	{ 
										_assemblyFlag = 1
							} else if  _batchType == WRE_BATCH					&&
										res.WireBatchId == _batchNumber			{ 
										_assemblyFlag = 1
							} else if  _batchType == CAS_BATCH					&&
										res.CasingBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == ADP_BATCH					&&
										res.AdaptorBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							} else if  _batchType == STK_BATCH					&&
										res.StickPodBatchId == _batchNumber		{ 
										_assemblyFlag = 1
							}
					} // Date check
				}// Date parse
			}// Date lenght check
						
			// Append Assembly current status to Assembly Selection Array if the flag is 1 (indicates valid for filter criteria)
			if _assemblyFlag == 1 {
				latestIndex := len(assemLineHistory_Holder.AssemblyLines)
				latestRes := assemLineHistory_Holder.AssemblyLines[latestIndex-1]
				res2E=append(res2E,latestRes)
				break // break the for loop as selected Assembly has been added to the list
			}
			
			//re-setting the flag and AssemblyDate
				_assemblyFlag = 0
				_assemblyDateInt64 = 0
		} // For assemLineHistory_Holder.AssemblyLines ends
	} // For assemID_Holder.AssemblyIDs ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}

// All AssemblyLine history
func (t *TnT) getAssemblyLineHistoryByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments to query")
	}

	_assemblyId := args[0]
	user_name:= args[1]	
	/* Access check -------------------------------------------- Starts*/
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not an AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	assemLine_HolderKey := _assemblyId + "H" // Indicates history key

	bytesAssemLineHolder, err := stub.GetState(assemLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get Assemblies") }

	return bytesAssemLineHolder, nil	

}

/* Package section*/

//API to create an Package
// Assemblies related to the package is updated with status = PACKAGED
func (t *TnT) createPackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 10 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 10. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[9]
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
	/* Access check -------------------------------------------- Ends*/
	
		_caseId := args[0]
		_holderAssemblyId := args[1]
		_chargerAssemblyId := args[2]
		_packageStatus := args[3]
		_packagingDate := args[4]
		_shippingToAddress := args[5]
		// Status of associated Assemblies	
		_assemblyStatus:= args[6]
		_packageInfo1:= args[7]
		_packageInfo2:= args[8]

		_time:= time.Now().Local()

		_packageCreationDate := _time.Format("20060102150405")
		_packageLastUpdatedOn := _time.Format("20060102150405")
		_packageCreatedBy := user_name
		_packageLastUpdatedBy := user_name

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
		pack.PackageInfo1 = _packageInfo1
		pack.PackageInfo2 = _packageInfo2

		bytes, err := json.Marshal(pack)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }

		err = stub.PutState(_caseId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }


		/* PackageLine history -----------------Starts */
		// Initialises the PackageLine_Holder
		var packLine_HolderInit PackageLine_Holder

		packLine_HolderKey := _caseId + "H" // Indicates history key
		bytesPackLinesInit, err := json.Marshal(packLine_HolderInit)
		if err != nil { return nil, errors.New("Error creating packLine_HolderInit record") }
		err = stub.PutState(packLine_HolderKey, bytesPackLinesInit)
		/* PackageLine history -----------------Ends */

		/* PackageLine history ------------------------------------------Starts */
		//packLine_HolderKey := _caseId + "H" // Indicates history key

		bytesPackageLines, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesPackageLines") }

		var packLine_Holder PackageLine_Holder

		err = json.Unmarshal(bytesPackageLines, &packLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt bytesPackageLines record") }

		packLine_Holder.PackageLines = append(packLine_Holder.PackageLines, pack) //appending the newly created pack
		
		bytesPackageLines, err = json.Marshal(packLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

		err = stub.PutState(packLine_HolderKey, bytesPackageLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		/* PackageLine history ------------------------------------------Ends */


		fmt.Println("Created Package successfully")

		//Update Holder Assemblies to Packaged status
		if 	len(_holderAssemblyId) > 0	{
			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageCreatedBy
			_assemblyPackage:= _caseId // Keeping reference
			
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
			assemHolder.AssemblyPackage = _assemblyPackage
			assemHolder.AssemblyInfo2 = "" // to reset the hascode to be updated later as part of package hash code update
			//assemHolder.AssemblyInfo2 = _packageInfo2// specia case to store the transaction hash - This will never be the case on creation (only true for update) hence commented
			
			bytesHolder, err := json.Marshal(assemHolder)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

			err = stub.PutState(_holderAssemblyId, bytesHolder)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


			/* AssemblyLine history ------------------------------------------Starts */
			holderAssemLine_HolderKey := _holderAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
			bytesHolderAssemblyLines, err := stub.GetState(holderAssemLine_HolderKey)
			if err != nil { return nil, errors.New("Unable to get Assemblies") }

			var holderAssemLine_Holder AssemblyLine_Holder

			err = json.Unmarshal(bytesHolderAssemblyLines, &holderAssemLine_Holder)
			if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

			holderAssemLine_Holder.AssemblyLines = append(holderAssemLine_Holder.AssemblyLines, assemHolder) //appending the updated AssemblyLine
			
			bytesHolderAssemblyLines, err = json.Marshal(holderAssemLine_Holder)
			if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

			err = stub.PutState(holderAssemLine_HolderKey, bytesHolderAssemblyLines)
			if err != nil { return nil, errors.New("Unable to put the state") }
			/* AssemblyLine history ------------------------------------------Ends */


			}

		//Update Charger Assemblies to Packaged status
		if 	len(_chargerAssemblyId) > 0		{
			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageCreatedBy
			_assemblyPackage:= _caseId // Keeping reference

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
			assemCharger.AssemblyPackage = _assemblyPackage
			assemCharger.AssemblyInfo2 = "" // to reset the hascode to be updated later as part of package hash code update

			
			bytesCharger, err := json.Marshal(assemCharger)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

			err = stub.PutState(_chargerAssemblyId, bytesCharger)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


			/* AssemblyLine history ------------------------------------------Starts */
			chargerAssemLine_HolderKey := _chargerAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
			bytesChargerAssemblyLines, err := stub.GetState(chargerAssemLine_HolderKey)
			if err != nil { return nil, errors.New("Unable to get Assemblies") }

			var chargerAssemLine_Holder AssemblyLine_Holder

			err = json.Unmarshal(bytesChargerAssemblyLines, &chargerAssemLine_Holder)
			if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

			chargerAssemLine_Holder.AssemblyLines = append(chargerAssemLine_Holder.AssemblyLines, assemCharger) //appending the updated AssemblyLine
			
			bytesChargerAssemblyLines, err = json.Marshal(chargerAssemLine_Holder)
			if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

			err = stub.PutState(chargerAssemLine_HolderKey, bytesChargerAssemblyLines)
			if err != nil { return nil, errors.New("Unable to put the state") }
			/* AssemblyLine history ------------------------------------------Ends */

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

	if len(args) != 10 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 10. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[9]
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
	/* Access check -------------------------------------------- Ends*/
		
		_caseId := args[0]
		//_holderAssemblyId := args[1]
		//_chargerAssemblyId := args[2]
		_packageStatus := args[3]
		_packagingDate := args[4]
		_shippingToAddress := args[5]
		// Status of associated Assemblies	
		_assemblyStatus := args[6]
		_packageInfo1:= args[7]
		_packageInfo2:= args[8]

		_time:= time.Now().Local()

		//_packageCreationDate := _time.Format("2006-01-02")
		_packageLastUpdatedOn := _time.Format("20060102150405")
		//_packageCreatedBy := ""
		_packageLastUpdatedBy := user_name


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
		pack.PackagingDate = _packagingDate
		pack.ShippingToAddress = _shippingToAddress
		//pack.PackageCreationDate = _packageCreationDate
		pack.PackageLastUpdatedOn = _packageLastUpdatedOn
		//pack.PackageCreatedBy = _packageCreatedBy
		pack.PackageLastUpdatedBy = _packageLastUpdatedBy
		pack.PackageInfo1 = _packageInfo1
		pack.PackageInfo2 = _packageInfo2

		// Getting associate Assembly IDs
		_holderAssemblyId := pack.HolderAssemblyId
		_chargerAssemblyId := pack.ChargerAssemblyId


		bytes, err := json.Marshal(pack)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }

		err = stub.PutState(_caseId, bytes)
		if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }


		/* PackageLine history ------------------------------------------Starts */
		packLine_HolderKey := _caseId + "H" // Indicates history key

		bytesPackageLines, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesPackageLines") }

		var packLine_Holder PackageLine_Holder

		err = json.Unmarshal(bytesPackageLines, &packLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt bytesPackageLines record") }

		packLine_Holder.PackageLines = append(packLine_Holder.PackageLines, pack) //appending the newly created pack
		
		bytesPackageLines, err = json.Marshal(packLine_Holder)
		if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

		err = stub.PutState(packLine_HolderKey, bytesPackageLines)
		if err != nil { return nil, errors.New("Unable to put the state") }
		/* PackageLine history ------------------------------------------Ends */

		fmt.Println("Created Package successfully")

		//Update Holder Assemblies status
		if 	len(_holderAssemblyId) > 0	{
			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageLastUpdatedBy
			_assemblyPackage:= _caseId // Keeping reference

			//get the Assembly
			assemblyHolderAsBytes, err := stub.GetState(_holderAssemblyId)
			if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
			if assemblyHolderAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

			assemHolder := AssemblyLine{}
			json.Unmarshal(assemblyHolderAsBytes, &assemHolder)

			// Don't update assembly if there is no chnage in status
			// Update only when status moves say from Packaged -> Cancelled	
			if assemHolder.AssemblyStatus != _assemblyStatus {

				//update the AssemblyLine status
				assemHolder.AssemblyStatus = _assemblyStatus
				assemHolder.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
				assemHolder.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
				assemHolder.AssemblyPackage = _assemblyPackage
				assemHolder.AssemblyInfo2 = "" // to reset the hascode to be updated later as part of package hash code update
				
				bytesHolder, err := json.Marshal(assemHolder)
				if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

				err = stub.PutState(_holderAssemblyId, bytesHolder)
				if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


				/* AssemblyLine history ------------------------------------------Starts */
				holderAssemLine_HolderKey := _holderAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
				bytesHolderAssemblyLines, err := stub.GetState(holderAssemLine_HolderKey)
				if err != nil { return nil, errors.New("Unable to get Assemblies") }

				var holderAssemLine_Holder AssemblyLine_Holder

				err = json.Unmarshal(bytesHolderAssemblyLines, &holderAssemLine_Holder)
				if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

				holderAssemLine_Holder.AssemblyLines = append(holderAssemLine_Holder.AssemblyLines, assemHolder) //appending the updated AssemblyLine
				
				bytesHolderAssemblyLines, err = json.Marshal(holderAssemLine_Holder)
				if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

				err = stub.PutState(holderAssemLine_HolderKey, bytesHolderAssemblyLines)
				if err != nil { return nil, errors.New("Unable to put the state") }
				/* AssemblyLine history ------------------------------------------Ends */
			}// Change of Status ends	

		}

		//Update Charger Assemblies status
		if 	len(_chargerAssemblyId) > 0		{
			//_assemblyStatus:= "PACKAGED"
			_time:= time.Now().Local()
			_assemblyLastUpdatedOn := _time.Format("20060102150405")
			_assemblyLastUpdatedBy := _packageLastUpdatedBy
			_assemblyPackage:= _caseId // Keeping reference

			//get the Assembly
			assemblyChargerAsBytes, err := stub.GetState(_chargerAssemblyId)
			if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
			if assemblyChargerAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

			assemCharger := AssemblyLine{}
			json.Unmarshal(assemblyChargerAsBytes, &assemCharger)

			// Don't update assembly if there is no chnage in status
			// Update only when status moves say from Packaged -> Cancelled	
			if assemCharger.AssemblyStatus != _assemblyStatus {
				//update the AssemblyLine status
				assemCharger.AssemblyStatus = _assemblyStatus
				assemCharger.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
				assemCharger.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
				assemCharger.AssemblyPackage = _assemblyPackage
				assemCharger.AssemblyInfo2 = "" // to reset the hascode to be updated later as part of package hash code update

				bytesCharger, err := json.Marshal(assemCharger)
				if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

				err = stub.PutState(_chargerAssemblyId, bytesCharger)
				if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


				/* AssemblyLine history ------------------------------------------Starts */
				chargerAssemLine_HolderKey := _chargerAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
				bytesChargerAssemblyLines, err := stub.GetState(chargerAssemLine_HolderKey)
				if err != nil { return nil, errors.New("Unable to get Assemblies") }

				var chargerAssemLine_Holder AssemblyLine_Holder

				err = json.Unmarshal(bytesChargerAssemblyLines, &chargerAssemLine_Holder)
				if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

				chargerAssemLine_Holder.AssemblyLines = append(chargerAssemLine_Holder.AssemblyLines, assemCharger) //appending the updated AssemblyLine
				
				bytesChargerAssemblyLines, err = json.Marshal(chargerAssemLine_Holder)
				if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

				err = stub.PutState(chargerAssemLine_HolderKey, bytesChargerAssemblyLines)
				if err != nil { return nil, errors.New("Unable to put the state") }
				/* AssemblyLine history ------------------------------------------Ends */
			}// Check if status changes

		}

		return nil, nil

}


//API to update an Package HashCode
// Assemblies related to the package is updated with hashcode 
// Parameters: CAS0001, HASHCODE, USERNAME
func (t *TnT) updatePackageInfo2ById(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[2]
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
	/* Access check -------------------------------------------- Ends*/
		
		_caseId := args[0]
		_packageInfo2:= args[1]

		_time:= time.Now().Local()
		//_packageCreationDate := _time.Format("2006-01-02")
		_packageLastUpdatedOn := _time.Format("20060102150405")
		//_packageCreatedBy := ""
		_packageLastUpdatedBy := user_name


	//Checking if the Package exists
		packageAsBytes, err := stub.GetState(_caseId)
		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes == nil { return nil, errors.New("Package doesn't exists") }

		//setting the Package to update
		pack := PackageLine{}
		json.Unmarshal(packageAsBytes, &pack)

		// Update only when PackageInfo2 is not set to avoid uncenessary duplicate updates
		if len(pack.PackageInfo2) == 0 {
			pack.PackageLastUpdatedOn = _packageLastUpdatedOn
			pack.PackageLastUpdatedBy = _packageLastUpdatedBy
			pack.PackageInfo2 = _packageInfo2
			
			// Getting associate Assembly IDs
			_holderAssemblyId := pack.HolderAssemblyId
			_chargerAssemblyId := pack.ChargerAssemblyId


			bytes, err := json.Marshal(pack)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Package record: %s", err); return nil, errors.New("Error converting Package record") }

			err = stub.PutState(_caseId, bytes)
			if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Package record: %s", err); return nil, errors.New("Error storing Package record") }


			/* PackageLine history ------------------------------------------Starts */
			packLine_HolderKey := _caseId + "H" // Indicates history key

			bytesPackageLines, err := stub.GetState(packLine_HolderKey)
			if err != nil { return nil, errors.New("Unable to get bytesPackageLines") }

			var packLine_Holder PackageLine_Holder

			err = json.Unmarshal(bytesPackageLines, &packLine_Holder)
			if err != nil {	return nil, errors.New("Corrupt bytesPackageLines record") }

			//packLine_Holder.PackageLines = append(packLine_Holder.PackageLines, pack) //appending the newly created pack
			//Overwrite exisitng last element with the updated element - Don't apend but overwrite
			latestIndex := len(packLine_Holder.PackageLines)
			packLine_Holder.PackageLines[latestIndex-1] = pack


			bytesPackageLines, err = json.Marshal(packLine_Holder)
			if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

			err = stub.PutState(packLine_HolderKey, bytesPackageLines)
			if err != nil { return nil, errors.New("Unable to put the state") }
			/* PackageLine history ------------------------------------------Ends */

			fmt.Println("Updated Package successfully")

			//Update Holder Assemblies status
			if 	len(_holderAssemblyId) > 0	{
				//_assemblyStatus:= "PACKAGED"
				_time:= time.Now().Local()
				_assemblyLastUpdatedOn := _time.Format("20060102150405")
				_assemblyLastUpdatedBy := _packageLastUpdatedBy
				//_assemblyPackage:= _caseId // Keeping reference
				_assemblyInfo2:= _packageInfo2 // same hashcode as used for package update

				//get the Assembly
				assemblyHolderAsBytes, err := stub.GetState(_holderAssemblyId)
				if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
				if assemblyHolderAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

				assemHolder := AssemblyLine{}
				json.Unmarshal(assemblyHolderAsBytes, &assemHolder)

				// Don't update assembly if there is no chnage in status
				// Update only when AssemblyInfo2 is not set to avoid uncenessary duplicate updates
				if len(assemHolder.AssemblyInfo2) == 0 {

					//update the AssemblyLine status
					//assemHolder.AssemblyStatus = _assemblyStatus
					assemHolder.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
					assemHolder.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
					//assemHolder.AssemblyPackage = _assemblyPackage
					assemHolder.AssemblyInfo2 = _assemblyInfo2

					
					bytesHolder, err := json.Marshal(assemHolder)
					if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

					err = stub.PutState(_holderAssemblyId, bytesHolder)
					if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


					/* AssemblyLine history ------------------------------------------Starts */
					holderAssemLine_HolderKey := _holderAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
					bytesHolderAssemblyLines, err := stub.GetState(holderAssemLine_HolderKey)
					if err != nil { return nil, errors.New("Unable to get Assemblies") }

					var holderAssemLine_Holder AssemblyLine_Holder

					err = json.Unmarshal(bytesHolderAssemblyLines, &holderAssemLine_Holder)
					if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

					//holderAssemLine_Holder.AssemblyLines = append(holderAssemLine_Holder.AssemblyLines, assemHolder) //appending the updated AssemblyLine

					//Overwrite exisitng last element with the updated element - Don't apend but overwrite
					latestIndex := len(holderAssemLine_Holder.AssemblyLines)
					holderAssemLine_Holder.AssemblyLines[latestIndex-1] = assemHolder
					
					bytesHolderAssemblyLines, err = json.Marshal(holderAssemLine_Holder)
					if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

					err = stub.PutState(holderAssemLine_HolderKey, bytesHolderAssemblyLines)
					if err != nil { return nil, errors.New("Unable to put the state") }
					/* AssemblyLine history ------------------------------------------Ends */
				}// len(assemHolder.AssemblyInfo2) > 0 	

			}

			//Update Charger Assemblies status
			if 	len(_chargerAssemblyId) > 0		{
				//_assemblyStatus:= "PACKAGED"
				_time:= time.Now().Local()
				_assemblyLastUpdatedOn := _time.Format("20060102150405")
				_assemblyLastUpdatedBy := _packageLastUpdatedBy
				//_assemblyPackage:= _caseId // Keeping reference
				_assemblyInfo2:= _packageInfo2 // same hashcode as used for package update

				//get the Assembly
				assemblyChargerAsBytes, err := stub.GetState(_chargerAssemblyId)
				if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
				if assemblyChargerAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

				assemCharger := AssemblyLine{}
				json.Unmarshal(assemblyChargerAsBytes, &assemCharger)

				// Don't update assembly if there is no chnage in status
				// Update only when AssemblyInfo2 is not set to avoid uncenessary duplicate updates
				if len(assemCharger.AssemblyInfo2) == 0 {
					//update the AssemblyLine status
					//assemCharger.AssemblyStatus = _assemblyStatus
					assemCharger.AssemblyLastUpdatedOn = _assemblyLastUpdatedOn
					assemCharger.AssemblyLastUpdatedBy = _assemblyLastUpdatedBy
					//assemCharger.AssemblyPackage = _assemblyPackage
					assemCharger.AssemblyInfo2 = _assemblyInfo2

					bytesCharger, err := json.Marshal(assemCharger)
					if err != nil { fmt.Printf("SAVE_CHANGES: Error converting Assembly record: %s", err); return nil, errors.New("Error converting Assembly record") }

					err = stub.PutState(_chargerAssemblyId, bytesCharger)
					if err != nil { fmt.Printf("SAVE_CHANGES: Error storing Assembly record: %s", err); return nil, errors.New("Error storing Assembly record") }


					/* AssemblyLine history ------------------------------------------Starts */
					chargerAssemLine_HolderKey := _chargerAssemblyId + "H" // Indicates History Key for Assembly with ID = _assemblyId
					bytesChargerAssemblyLines, err := stub.GetState(chargerAssemLine_HolderKey)
					if err != nil { return nil, errors.New("Unable to get Assemblies") }

					var chargerAssemLine_Holder AssemblyLine_Holder

					err = json.Unmarshal(bytesChargerAssemblyLines, &chargerAssemLine_Holder)
					if err != nil {	return nil, errors.New("Corrupt AssemblyLines record") }

					//chargerAssemLine_Holder.AssemblyLines = append(chargerAssemLine_Holder.AssemblyLines, assemCharger) //appending the updated AssemblyLine
					
					//Overwrite exisitng last element with the updated element - Don't apend but overwrite
					latestIndex := len(chargerAssemLine_Holder.AssemblyLines)
					chargerAssemLine_Holder.AssemblyLines[latestIndex-1] = assemCharger
					
					bytesChargerAssemblyLines, err = json.Marshal(chargerAssemLine_Holder)
					if err != nil { return nil, errors.New("Error creating AssemblyLine_Holder record") }

					err = stub.PutState(chargerAssemLine_HolderKey, bytesChargerAssemblyLines)
					if err != nil { return nil, errors.New("Unable to put the state") }
					/* AssemblyLine history ------------------------------------------Ends */
				}// len(assemCharger.AssemblyInfo2) > 0
			} //len(_chargerAssemblyId) > 0	
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


	/* Access check -------------------------------------------- Starts*/
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
		if user_role != PACKAGELINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

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


//All Validators to be called before Invoke

// Validator before createAssembly invoke call
func (t *TnT) validateCreateAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 17. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied, not an AssemblyLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	
	//Checking if the Assembly already exists
	_assemblyId := args[0]
	assemblyAsBytes, err := stub.GetState(_assemblyId)
	if err != nil { return nil, errors.New("Failed to get assembly Id") }
	if assemblyAsBytes != nil { return nil, errors.New("Assembly already exists") }
	
	//Check Date
	_assemblyDate:= args[12]
	if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
		
	
	//No validation error proceed to call Invoke command
	return nil, nil
}

// Validator before createAssembly invoke call
func (t *TnT) validateUpdateAssembly(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 17 {
			return nil, errors.New("Incorrect number of arguments. Expecting 17.")
		} 
	

	_assemblyId := args[0]
	_assemblyStatus:= args[11]

	//get the Assembly
	assemblyAsBytes, err := stub.GetState(_assemblyId)
	if err != nil {	return nil, errors.New("Failed to get assembly Id")	}
	if assemblyAsBytes == nil { return nil, errors.New("Assembly doesn't exists") }

	//Check Date
	_assemblyDate:= args[12]
	if len(_assemblyDate) != 14 {return nil, errors.New("AssemblyDate must be 14 digit datetime field.")}	
	

	assem := AssemblyLine{}
	json.Unmarshal(assemblyAsBytes, &assem)


	/* Access check -------------------------------------------- Starts*/
	user_name := args[16]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)

		if user_role != ASSEMBLYLINE_ROLE {
			return nil, errors.New("Permission denied, not an AssemblyLine Role")
		}
		// AssemblyLine can't edit an Assembly in certain statuses
		if (user_role 			== ASSEMBLYLINE_ROLE 		&&
		(assem.AssemblyStatus 	== ASSEMBLYSTATUS_RFP || assem.AssemblyStatus == ASSEMBLYSTATUS_CAN)) 		{
			return nil, errors.New("Permission denied for AssemblyLine Role to update Assembly if status = 'Ready For Packaging'")
		}
		
		// AssemblyLine user can't move an AssemblyLine from QA Failed to Ready For packaging status
		if (user_role 			== ASSEMBLYLINE_ROLE 		&&
		assem.AssemblyStatus 	== ASSEMBLYSTATUS_QAF 		&&
		_assemblyStatus		 	== ASSEMBLYSTATUS_RFP) 		{
			return nil, errors.New("Permission denied for updating AssemblyLine with status = 'QA Failed' to 'Ready For Packaging' status")
		}

		// AssemblyLine user can't move an AssemblyLine to "Packaged" status directly; It is internally done in packaging line
		if (user_role 			== ASSEMBLYLINE_ROLE 		&&
		_assemblyStatus		 	== ASSEMBLYSTATUS_PKG) 		{
			return nil, errors.New("Permission denied for updating AssemblyLine to status 'Packaged'")
	}		
		
	}
	/* Access check -------------------------------------------- Ends*/	
	//No validation error proceed to call Invoke command
	return nil, nil
}

// Validator before createPackage invoke call
func (t *TnT) validateCreatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

		if len(args) != 10 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 10. Got: %d.", len(args))
		}

	/* Access check -------------------------------------------- Starts*/
	user_name := args[9]
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
	/* Access check -------------------------------------------- Ends*/
	
	//Checking if the Package already exists
		_caseId := args[0]
		packageAsBytes, err := stub.GetState(_caseId)
		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes != nil { return nil, errors.New("Package already exists") }
	
	//No validation error proceed to call Invoke command
	return nil, nil
}

// Validator before updateAssembly invoke call
func (t *TnT) validateUpdatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 10 {
				return nil, fmt.Errorf("Incorrect number of arguments. Expecting 10. Got: %d.", len(args))
			}

		/* Access check -------------------------------------------- Starts*/
		user_name := args[9]
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
		/* Access check -------------------------------------------- Ends*/
			
		//Checking if the Package already exists
		_caseId := args[0]
		packageAsBytes, err := stub.GetState(_caseId)
		if err != nil { return nil, errors.New("Failed to get Package") }
		if packageAsBytes == nil { return nil, errors.New("Package doesn't exists") }
		
	//No validation error proceed to call Invoke command
	return nil, nil
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


// All PackageLine history
func (t *TnT) getPackageLineHistoryByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments to query")
	}

	_caseId := args[0]
	user_name:= args[1]	
	/* Access check -------------------------------------------- Starts*/
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not PakagingLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/


	packLine_HolderKey := _caseId + "H" // Indicates history key
	bytesPackLineHolder, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get PackageLine history") }

	return bytesPackLineHolder, nil	

}
// Search Package
//get all Packages based on Assembly Id
func (t *TnT) getPackagesByAssemblyId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not PakagingLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	_assemblyType:= args[0]
	_assemblyId := args[1]
	_packageFlag:= 0

	bytes, err := stub.GetState("Packages")
	if err != nil { return nil, errors.New("Unable to get Packages") }

	var packageCaseID_Holder PackageCaseID_Holder

	err = json.Unmarshal(bytes, &packageCaseID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Packages") }

	res2E:= []*PackageLine{}	

	for _, caseId := range packageCaseID_Holder.PackageCaseIDs {

		//Get the existing Packages
		packageAsBytes, err := stub.GetState(caseId)
		if err != nil { return nil, errors.New("Failed to get Package")}

		if packageAsBytes != nil { 
			res := new(PackageLine)
			json.Unmarshal(packageAsBytes, &res)

			//Check the filter condition
			if 		   _assemblyType == HLD_ASSMB_TYP	&&
						res.HolderAssemblyId == _assemblyId		{ 
						_packageFlag = 1
			} else if  _assemblyType == CHG_ASSMB_TYP	&&
						res.ChargerAssemblyId == _assemblyId	{ 
						_packageFlag = 1
			}
			

			// Append Assembly to Assembly Array if the flag is 1 (indicates valid for filter criteria)
			if _packageFlag == 1 {
				res2E=append(res2E,res)
			}
		} // If ends
		//re-setting the flag to 0
		_packageFlag = 0
	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}
//get all Packages based on FromDate & ToDate and AssemblyId
func (t *TnT) getPackagesByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64
	
	_fromDate, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	_packageFlag:= 0

	bytes, err := stub.GetState("Packages")
	if err != nil { return nil, errors.New("Unable to get Packages") }

	var packageCaseID_Holder PackageCaseID_Holder
	var _packageDateInt64 int64
	
	err = json.Unmarshal(bytes, &packageCaseID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Packages") }

	res2E:= []*PackageLine{}	

	for _, caseId := range packageCaseID_Holder.PackageCaseIDs {

		//Get the existing Package Line History
		packageAsBytes, err := stub.GetState(caseId)
		if err != nil { return nil, errors.New("Failed to get Case")}
		if packageAsBytes == nil { return nil, errors.New("Failed to get AsseCasembly")}

		res := new(PackageLine)
		json.Unmarshal(packageAsBytes, &res)

		//fmt.Printf("%T, %v\n", _fromDate, _fromDate)
		//fmt.Printf("%T, %v\n", _toDate, _toDate)
		//if _fromDate == _toDate { return nil, errors.New("Failed to get Assembly")}
		
		//Check the filter condition YYYYMMDDHHMMSS
		if len(res.PackagingDate) != 14 {return nil, errors.New("PackagingDate must be 14 digit datetime field.")}
		if _packageDateInt64, err = strconv.ParseInt(res.PackagingDate, 10, 64); err != nil { errors.New ("Error in converting PackagingDate to int64")}
		if	_packageDateInt64 >= _fromDate		&&
			_packageDateInt64 <= _toDate		{ 
			_packageFlag = 1
		} 
					
		// Append Package Case to Package Array if the flag is 1 (indicates valid for filter criteria)
		if _packageFlag == 1 {
			res2E=append(res2E,res)
		}
	//re-setting the flag and PackageCreationDate
		_packageFlag = 0
		_packageDateInt64 = 0
	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
}

//get all Package based on AssemblyID & From & To Date
func (t *TnT) getPackageByAssemblyIdAndByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	/* Access check -------------------------------------------- Starts*/
	if len(args) != 5 {
			return nil, errors.New("Incorrect number of arguments. Expecting 5.")
		}
	user_name := args[4]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not PackageLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64
	_assemblyType:= args[0]
	_assemblyId := args[1]
	_fromDate, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	_packageFlag:= 0

	bytes, err := stub.GetState("Packages")
	if err != nil { return nil, errors.New("Unable to get Packages") }

	var packageCaseID_Holder PackageCaseID_Holder
	var _packageDateInt64 int64
	
	err = json.Unmarshal(bytes, &packageCaseID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Packages") }

	res2E:= []*PackageLine{}	

	for _, caseId := range packageCaseID_Holder.PackageCaseIDs {

		//Get the existing Package Line History
		packageAsBytes, err := stub.GetState(caseId)
		if err != nil { return nil, errors.New("Failed to get Case")}
		if packageAsBytes == nil { return nil, errors.New("Failed to get AsseCasembly")}

		res := new(PackageLine)
		json.Unmarshal(packageAsBytes, &res)

		//fmt.Printf("%T, %v\n", _fromDate, _fromDate)
		//fmt.Printf("%T, %v\n", _toDate, _toDate)
		//if _fromDate == _toDate { return nil, errors.New("Failed to get Assembly")}
		
		//Check the filter condition YYYYMMDDHHMMSS
		if len(res.PackagingDate) != 14 {return nil, errors.New("PackagingDate must be 14 digit datetime field.")}
		if _packageDateInt64, err = strconv.ParseInt(res.PackagingDate, 10, 64); err != nil { errors.New ("Error in converting PackagingDate to int64")}
		if	_packageDateInt64 >= _fromDate		&&
			_packageDateInt64 <= _toDate{ 
				if  _assemblyType == HLD_ASSMB_TYP	&&
				res.HolderAssemblyId == _assemblyId	{ 
				_packageFlag = 1
			} else if  _assemblyType == CHG_ASSMB_TYP &&
						res.ChargerAssemblyId == _assemblyId{ 
						_packageFlag = 1
			}
		} 

		// Append Package Case to Package Array if the flag is 1 (indicates valid for filter criteria)
		if _packageFlag == 1 {
			res2E=append(res2E,res)
		}
	//re-setting the flag and PackageCreationDate
		_packageFlag = 0
		_packageDateInt64 = 0
	} // For ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil

}
//get all Packages History based on FromDate & ToDate
func (t *TnT) getPackagesHistoryByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	/* Access check -------------------------------------------- Starts*/
	if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3.")
		}
	user_name := args[2]
	if len(user_name) == 0 { return nil, errors.New("User name supplied as empty") }

	if len(user_name) > 0 {
		ecert_role, err := t.get_ecert(stub, user_name)
		if err != nil {return nil, errors.New("userrole couldn't be retrieved")}
		if ecert_role == nil {return nil, errors.New("username not defined")}

		user_role := string(ecert_role)
		if user_role != PACKAGELINE_ROLE &&
			user_role != QA_VIEWER_ROLE {
			return nil, errors.New("Permission denied not PackagingLine Role")
		}
	}
	/* Access check -------------------------------------------- Ends*/

	// YYYYMMDDHHMMSS (e.g. 20170612235959) handled as Int64
	//var _fromDate int64
	//var _toDate int64
	
	_fromDate, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting FromDate to int64")}
	
	_toDate, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil { return nil, errors.New ("Error in converting ToDate to int64")}
	
	_packageFlag:= 0

	bytes, err := stub.GetState("Packages")
	if err != nil { return nil, errors.New("Unable to get Packages") }

	var packageCaseID_Holder PackageCaseID_Holder
	var _packageDateInt64 int64
	
	err = json.Unmarshal(bytes, &packageCaseID_Holder)
	if err != nil {	return nil, errors.New("Corrupt Packages") }

	// Array of filtered Package Line
	res2E:= []PackageLine{}	

	
	//Looping through the array of packageCaseId
	for _, caseId := range packageCaseID_Holder.PackageCaseIDs {

		//Get the AssemblyLine History for each AssemblyID
		packLine_HolderKey := caseId + "H" // Indicates history key
		bytesPackageHistoryLines, err := stub.GetState(packLine_HolderKey)
		if err != nil { return nil, errors.New("Unable to get bytesPackageHistoryLines") }

		var packLine_Holder PackageLine_Holder

		err = json.Unmarshal(bytesPackageHistoryLines, &packLine_Holder)
		if err != nil {	return nil, errors.New("Corrupt History Holder record") }

		//Looping through the array of assemblies
		for _, res := range packLine_Holder.PackageLines {
		
			
			//Skip if not a valid date YYYYMMDDHHMMSS
			if len(res.PackagingDate) == 14 {
				if _packageDateInt64, err = strconv.ParseInt(res.PackagingDate, 10, 64); err == nil { 
					if	_packageDateInt64 >= _fromDate		&&
						_packageDateInt64 <= _toDate		{ 
						_packageFlag = 1
					} 
				}
			}
						
			// Append AssembPackagely to Package Array if the flag is 1 (indicates valid for filter criteria)
			if _packageFlag == 1 {
				res2E=append(res2E,res)
			}
			
			//re-setting the flag and PackagingDate
				_packageFlag = 0
				_packageDateInt64 = 0
		} // For packLine_Holder.PackageLines ends
	} // For packageCaseID_Holder.PackageCaseIDs ends

    mapB, _ := json.Marshal(res2E)
    //fmt.Println(string(mapB))
	return mapB, nil
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

	var assemID_Holder AssemblyID_Holder
	bytesAssembly, err := json.Marshal(assemID_Holder)
    if err != nil { return nil, errors.New("Error creating assemID_Holder record") }
	err = stub.PutState("Assemblies", bytesAssembly)

	var packageCaseID_Holder PackageCaseID_Holder
	bytesPackage, err := json.Marshal(packageCaseID_Holder)
    if err != nil { return nil, errors.New("Error creating packageCaseID_Holder record") }
	err = stub.PutState("Packages", bytesPackage)
	
	/* GetAll changes---------------------------ends------------------------ */

	// creating minimum default user and roles
	//"AssemblyLine_User1","assemblyline_role","PackageLine_User1", "packageline_role"
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
	} else if function == "updateAssemblyInfo2ByID" {
		fmt.Printf("Function is updateAssemblyInfo2ByID")
		return t.updateAssemblyInfo2ByID(stub, args)
	} else if function == "updatePackageInfo2ById" {
		fmt.Printf("Function is updatePackageInfo2ById")
		return t.updatePackageInfo2ById(stub, args)
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
		t := TnT{}
		return t.get_ecert(stub, args[0])
	} else if function == "validateCreateAssembly" {
		t := TnT{}
		return t.validateCreateAssembly(stub, args)
	} else if function == "validateUpdateAssembly" {
		t := TnT{}
		return t.validateUpdateAssembly(stub, args)
	} else if function == "validateCreatePackage" {
		t := TnT{}
		return t.validateCreatePackage(stub, args)
	} else if function == "validateUpdatePackage" {
		t := TnT{}
		return t.validateUpdatePackage(stub, args)
	} else if function == "getAssemblyLineHistoryByID" {
		t := TnT{}
		return t.getAssemblyLineHistoryByID(stub, args)
	} else if function == "getPackageLineHistoryByID" {
		t := TnT{}
		return t.getPackageLineHistoryByID(stub, args)
	} else if function == "getAssembliesByBatchNumber" {
		t := TnT{}
		return t.getAssembliesByBatchNumber(stub, args)
	} else if function == "getAssembliesByDate" {
		t := TnT{}
		return t.getAssembliesByDate(stub, args)
	} else if function == "getAssembliesHistoryByDate" {
		t := TnT{}
		return t.getAssembliesHistoryByDate(stub, args)
	} else if function == "getAssembliesByBatchNumberAndByDate" {
		t := TnT{}
		return t.getAssembliesByBatchNumberAndByDate(stub, args)
	} else if function == "getAssembliesHistoryByBatchNumberAndByDate" {
		t := TnT{}
		return t.getAssembliesHistoryByBatchNumberAndByDate(stub, args)
	} else if function == "getPackagesByAssemblyId" {
		t := TnT{}
		return t.getPackagesByAssemblyId(stub, args)
	} else if function == "getPackagesByDate" {
		t := TnT{}
		return t.getPackagesByDate(stub, args)
	} else if function == "getPackageByAssemblyIdAndByDate" {
		t := TnT{}
		return t.getPackageByAssemblyIdAndByDate(stub, args)
	} else if function == "getPackagesHistoryByDate" {
		t := TnT{}
		return t.getPackagesHistoryByDate(stub, args)
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