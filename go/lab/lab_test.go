package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

////////////////// Util Methods //////////////////

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, values ...string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	for _, v := range values {
		if !strings.Contains(string(bytes), v) {
			fmt.Println("State value", name, "was not", v, "as expected")
			t.FailNow()
		}
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, tx string, name string, values ...string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(tx), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", tx, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", tx, "failed to get value")
		t.FailNow()
	}
	for _, v := range values {
		if !strings.Contains(string(res.Payload), v) {
			fmt.Println("State value", name, "was not", v, "as expected")
			t.FailNow()
		}
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvokeError(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.ERROR {
		fmt.Println("Invoke", args, "success", string(res.Message))
		t.FailNow()
	}
}

////////////////// Tests //////////////////

func Test_givenANewLaboratoryWhenAddLaboratoryThenLaboratoryIsPersisted(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addLaboratory LabXXX, 15/03/2018, 1st Street, ARM
	checkInvoke(t, stub, [][]byte{[]byte("addLaboratory"), []byte("LabXXX"), []byte("15/03/2018"), []byte("1st Street"), []byte("ARM")})

	checkState(t, stub, "LabXXX", "address", "1st Street", "armOwner", "ARM")
}

func Test_givenANewLaboratoryWhenAddLaboratoryWithOneParamThenError(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addLaboratory LabXXX
	checkInvokeError(t, stub, [][]byte{[]byte("addLaboratory"), []byte("LabXXX")})

}
