/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Transit struct {
	LocLatitude     string `json:"lat"`
	LocLongitude    string `json:"lon"`
	Time            string `json:"time"`
	HaulierReceptor string `json:"haulierreceptor"`
}

type Arrival struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

type Asset struct {
	Type     string    `json:"type"`
	Qty      string    `json:"qty"`
	Price    string    `json:"price"`
	DateL    string    `json:"datel"`
	Agent    string    `json:"agent"`
	Transits []Transit `json:"transits"`
	Arrivals []Arrival `json:"arrival"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryAllAssets" {
		return s.queryAllAssets(APIstub)
	} else if function == "queryAssets" {
		return s.queryAssets(APIstub)
	} else if function == "queryByAsset" {
		return s.queryByAsset(APIstub, args)
	} else if function == "buyAsset" {
		return s.buyAsset(APIstub, args)
	} else if function == "generateTransit" {
		return s.generateTransit(APIstub, args)
	} else if function == "arrival" {
		return s.arrival(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) buyAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	var transit = Transit{
		LocLatitude:     args[6],
		LocLongitude:    args[7],
		Time:            args[8],
		HaulierReceptor: args[5],
	}

	var asset = Asset{
		Type:     args[1],
		Qty:      args[2],
		Price:    args[3],
		DateL:    args[4],
		Agent:    args[5],
		Transits: []Transit{transit},
		Arrivals: nil,
	}
	assetAsBytes, _ := json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) generateTransit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var transit = Transit{
		LocLatitude:     args[1],
		LocLongitude:    args[2],
		Time:            args[3],
		HaulierReceptor: args[4],
	}

	assetAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get specified Asset")
	}

	asset := Asset{}
	json.Unmarshal(assetAsBytes, &asset)

	asset.Agent = args[4]

	asset.Transits = append(asset.Transits, transit)
	fmt.Println("!!! appended transit to Asset")

	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) arrival(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var arrival = Arrival{
		Date:   args[1],
		Status: args[2],
	}

	assetAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get specified Asset")
	}

	asset := Asset{}
	json.Unmarshal(assetAsBytes, &asset)

	asset.Arrivals = append(asset.Arrivals, arrival)
	fmt.Println("!!! appended arrival to Asset")

	assetAsBytes, _ = json.Marshal(asset)
	APIstub.PutState(args[0], assetAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllAssets(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "ASSET0"
	endKey := "ASSET999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllAssets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryAssets(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "ASSET0"
	endKey := "ASSET999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- queryAllAssets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryByAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Empty key. Expecting an Asset")
	}

	assetAsBytes, _ := APIstub.GetState(args[0])

	if len(assetAsBytes) == 0 {
		return shim.Error("Invalid key. Expecting an Asset")
	}

	return shim.Success(assetAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
