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

func Test_addArm(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addARM
	checkInvoke(t, stub, [][]byte{[]byte("addARM"), []byte("ARM1"), []byte("My ARM")})

	checkState(t, stub, "ARM1", "ARM1")

}

func Test_addLaboratory(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addARM
	checkInvoke(t, stub, [][]byte{[]byte("addARM"), []byte("ARM1"), []byte("My ARM")})

	// addLaboratory
	checkInvoke(t, stub, [][]byte{[]byte("addLaboratory"), []byte("ARM1"), []byte("BAYER")})

	// addLaboratory
	checkInvoke(t, stub, [][]byte{[]byte("addLaboratory"), []byte("ARM1"), []byte("GLX")})

	checkState(t, stub, "ARM1", "ARM1")
	checkState(t, stub, "ARM1", "BAYER")
	checkState(t, stub, "ARM1", "GLX")

}

func Test_addLaboratoryError(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addARM
	checkInvoke(t, stub, [][]byte{[]byte("addARM"), []byte("ARM1"), []byte("My ARM")})

	// addLaboratory
	checkInvokeError(t, stub, [][]byte{[]byte("addLaboratory"), []byte("ARM1")})

}

func Test_addLaboratoryWithoutArmError(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addLaboratory
	checkInvokeError(t, stub, [][]byte{[]byte("addLaboratory"), []byte("ARM1")})

}

func Test_addPermission(t *testing.T) {
	scc := new(SmartContract)
	stub := shim.NewMockStub("ex01", scc)

	// addARM
	checkInvoke(t, stub, [][]byte{[]byte("addARM"), []byte("ARM1"), []byte("My ARM")})

	// addLaboratory
	checkInvoke(t, stub, [][]byte{[]byte("addLaboratory"), []byte("ARM1"), []byte("BAYER")})

	// addLaboratory
	checkInvoke(t, stub, [][]byte{[]byte("addLaboratory"), []byte("ARM1"), []byte("GLX")})

	// addMarketingAuthorization
	checkInvoke(t, stub, [][]byte{[]byte("addMarketingAuthorization"), []byte("ARM1"), []byte("BAYER"), []byte("Med1"), []byte("10/10/2018")})

	// addMarketingAuthorization
	checkInvoke(t, stub, [][]byte{[]byte("addMarketingAuthorization"), []byte("ARM1"), []byte("BAYER"), []byte("Med2"), []byte("10/12/2018")})

	checkState(t, stub, "ARM1", "ARM1")
	checkState(t, stub, "ARM1", "BAYER")
	checkState(t, stub, "ARM1", "GLX")
	checkState(t, stub, "ARM1", "Med1")
	checkState(t, stub, "ARM1", "Med2")

	checkQuery(t, stub, "queryByMarketingAuthorization", "ARM2")

// func checkQuery(t *testing.T, stub *shim.MockStub, tx string, name string, values ...string) {
// func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
// func checkState(t *testing.T, stub *shim.MockStub, name string, values ...string) {

}
