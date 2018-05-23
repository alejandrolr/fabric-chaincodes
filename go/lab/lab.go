package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract defines laboratory transactions
type SmartContract struct {
}

type Order struct {
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Quantity    int64  `json:"quantity"`
	DateCreated string `json:"datecreated"`
	DateSent    string `json:"datesent"`
	DateArrival string `json:"datearrival"`
	DateCancelled string `json:"datecancelled"`
	SentFlag    string `json:"sentflag"`
}

type Pharmacy struct {
	Pharmacy string  `json:"pharmacy"`
	Order    []Order `json:"order"`
}

// MarketingAuthorization defines a marketing authorization in order to produce a medicine
type MarketingAuthorization struct {
	Medicine    string `json:"medicine"`
	CreatedDate string `json:"createdDate"`
}

// Laboratory defines a company wich produces medicines
type Laboratory struct {
	LaboratoryName         string                   `json:"laboratoryName"`
	CreatedDate            string                   `json:"createdDate"`
	Address                string                   `json:"address"`
	ARMOwner               string                   `json:"armOwner"`
	MarketingAuthorization []MarketingAuthorization `json:"authorizations"`
	Pharmacy               []Pharmacy               `json:"pharmacy"`
}

var logger = *shim.NewLogger("PHALogger")

// Init is called during Instantiate transaction
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Printf("SmartContract has been instantiated \n")
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "addLaboratory" {
		return s.addLaboratory(APIstub, args)
	} else if function == "queryLabByARM" {
		return s.queryLabByARM(APIstub, args)
	} else if function == "createMarketingAuthorization" {
		return s.createMarketingAuthorization(APIstub, args)
	} else if function == "addMedicineOrder" {
		return s.addMedicineOrder(APIstub, args)
	} else if function == "SendOrder" {
		return s.SendOrder(APIstub, args)
	} else if function == "queryByLab" {
		return s.queryByLab(APIstub, args)
	} else if function == "queryLabsJSON" {
		return s.queryLabsJSON(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// ./executeTransaction.sh '{"Args":["addMedicineOrder", "BAYER", "FarmaciaAluche", "IBUPROFENO", "IBUPROFENODESC", "7"]}' labcc
// ./executeTransaction.sh '{"Args":["createMedicineOrder", "FarmaciaAluche", "IBUPROFENO", "IBUPROFENODESC", "7", "BAYER"]}' phacc
func (s *SmartContract) addMedicineOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	if len(args[0]) == 0 {
		return shim.Error("Empty key. Expecting a LAB")
	}

	current_time := time.Now().Local()
	str := current_time.Format("02/01/2006")

	quantity, _ := strconv.ParseInt(args[4], 10, 64)
	var order = Order{
		Name:        args[2],
		Desc:        args[3],
		Quantity:    quantity,
		DateCreated: str,
		DateSent:    "",
		DateArrival: "",
		DateCancelled: "",
		SentFlag:    "",
	}

	labAsBytes, _ := APIstub.GetState(args[0])
	laboratory := Laboratory{}
	json.Unmarshal(labAsBytes, &laboratory)

	existe := 0
	for _, pha := range laboratory.Pharmacy {
		if pha.Pharmacy == args[1] {
			existe = 1
		}
	}

	if existe == 0 { // crea una nueva estructura farmacia e incorpora el pedido
		var pharmacy = Pharmacy{
			Pharmacy: args[1],
			Order:    []Order{order},
		}

		laboratory.Pharmacy = append(laboratory.Pharmacy, pharmacy)

		labAsBytes, _ := json.Marshal(laboratory)
		APIstub.PutState(args[0], labAsBytes)
		fmt.Println("!!! appended PHA")

	} else { // recorrer farmacias y append en la que aplique
		l := 0
		for i, pha := range laboratory.Pharmacy {
			if pha.Pharmacy == args[1] {
				l = i
				break
			}
		}

		laboratory.Pharmacy[l].Order = append(laboratory.Pharmacy[l].Order, order)

		labAsBytes, _ = json.Marshal(laboratory)
		APIstub.PutState(args[0], labAsBytes)
		fmt.Println("!!! appended order to PHA")
	}

	return shim.Success(nil)
}

// ./executeTransaction.sh '{"Args":["SendOrder", "FarmaciaAluche", "01/07/2018"]}' labcc
// ./executeTransaction.sh '{"Args":["create", "FarmaciaAluche", "IBUPROFENO", "IBUPROFENODESC", "7", "FECHA"]}' phacc
// ./executeTransaction.sh '{"Args":["create", "BAYERN", FarmaciaAluche", "IBUPROFENO", "IBUPROFENODESC", "7", "FECHA"]}' phacc
func (s *SmartContract) SendOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	labAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get specified Lab")
	}

	labStruct := Laboratory{}
	json.Unmarshal(labAsBytes, &labStruct)

	foundPharma := false
	foundOrder := false
	pahrmaIndex := 0
	orderIndex := 0
	quantity, _ := strconv.ParseInt(args[4], 10, 64)
	for i, pharma := range labStruct.Pharmacy {
		if pharma.Pharmacy == args[1] {
			foundPharma = true
			pahrmaIndex = i
			for j, order := range pharma.Order {
				if (order.Name == args[2] && order.Desc == args[3] && order.Quantity == quantity && order.DateCancelled == "") {
					foundOrder = true
					orderIndex = j
					break
				}
			}
			break
		}
	}
	current_time := time.Now().Local()
	str := current_time.Format("02/01/2006")
	labStruct.Pharmacy[pahrmaIndex].Order[orderIndex].SentFlag = "true"
	labStruct.Pharmacy[pahrmaIndex].Order[orderIndex].DateSent = str

	if (foundPharma) {
		if (foundOrder) {
			labAsBytes, _ := json.Marshal(labStruct)
			APIstub.PutState(args[0], labAsBytes)
			return shim.Success(nil)
		} else {
			return shim.Error("Failed to get specified Order")
		}
	} else {
		return shim.Error("Failed to get specified Pharma")
	}
}

// ./executeTransaction.sh '{"Args":["createMarketingAuthorization", "OWNER1", "BAYER", "IBUPROFENO", "01/07/2018"]}' labcc
func (s *SmartContract) createMarketingAuthorization(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Expecting 4 {OWNER, LAB, MEDICINE, DATE}")
	}

	chainCodeToCall := "arm"
	channelID := "mychannel"
	f := "addMarketingAuthorization"

	invokeArgs := toChaincodeArgs(f, args[0], args[1], args[2], args[3])
	response := APIstub.InvokeChaincode(chainCodeToCall, invokeArgs, channelID)
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to invoke armcc. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}

	fmt.Printf("Invoke armcc successful. Got response %s", string(response.Payload))

	return shim.Success(response.Payload)
}

// ./executeTransaction.sh '{"Args":["addLaboratory", "BAYER", "01/03/2018", "calle de BAYER", "OWNER01"]}' labcc
func (s *SmartContract) addLaboratory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var lab = Laboratory{
		LaboratoryName:         args[0],
		CreatedDate:            args[1],
		Address:                args[2],
		ARMOwner:               args[3],
		MarketingAuthorization: nil,
		Pharmacy:               nil,
	}

	// TODO check lab already exists

	labAsBytes, _ := json.Marshal(lab)
	APIstub.PutState(args[0], labAsBytes)

	return shim.Success(nil)
}

// ./executeTransaction.sh '{"Args":["queryByLab", "BAYER"]}' labcc
func (s *SmartContract) queryByLab(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Empty key. Expecting a LAB")
	}

	labAsBytes, _ := APIstub.GetState(args[0])
	if len(labAsBytes) == 0 {
		return shim.Error("Invalid key. Expecting a LAB")
	}

	labStruct := Laboratory{}
	json.Unmarshal(labAsBytes, &labStruct)
	labAsBytes, _ = json.Marshal(labStruct)
	return shim.Success(labAsBytes)
}

func (s *SmartContract) queryLabsJSON(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Empty key. Expecting a LAB")
	}

	labAsBytes, _ := APIstub.GetState(args[0])
	if len(labAsBytes) == 0 {
		return shim.Error("Invalid key. Expecting a LAB")
	}

	labStruct := Laboratory{}
	json.Unmarshal(labAsBytes, &labStruct)

	var buffer bytes.Buffer
	buffer.WriteString("{")

	buffer.WriteString("\"LaboratoryName\":")
	buffer.WriteString("\"")
	buffer.WriteString(labStruct.LaboratoryName)
	buffer.WriteString("\"")

	buffer.WriteString("\"CreatedDate\":")
	buffer.WriteString("\"")
	buffer.WriteString(labStruct.CreatedDate)
	buffer.WriteString("\"")

	buffer.WriteString("\"Address\":")
	buffer.WriteString("\"")
	buffer.WriteString(labStruct.Address)
	buffer.WriteString("\"")

	buffer.WriteString("\"ARMOwner\":")
	buffer.WriteString("\"")
	buffer.WriteString(labStruct.ARMOwner)
	buffer.WriteString("\"")

	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())

}

// ./executeTransaction.sh '{"Args":["queryLabByARM", "BAYER"]}' labcc /// CouchDB !!!!!!!
// _______________________________________________________________________________________

func (s *SmartContract) queryLabByARM(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Empty key. Expecting an Asset")
	}

	lab := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"lab\",\"ARMowner\":\"%s\"}}", lab)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

}
