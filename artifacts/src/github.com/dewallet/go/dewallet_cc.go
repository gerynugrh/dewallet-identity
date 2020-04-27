package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("dewallet_chaincodes")

// DewalletChaincode is chaincode for dewallet operation
type DewalletChaincode struct {
}

// Identity saves the identity of user
// Data is an encrypted data of the user
// Data can only be decrypted by user private key
type Identity struct {
	Username  string `json:"username"`
	PublicKey string `json:"publicKey"`
	Data      string `json:"data"`
	Verified  string `json:"verified"`
}

// Init will initialize the chaincode
func (t *DewalletChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Initialize Dewallet Chaincode")
	return shim.Success(nil)
}

// Invoke will run the approriate function based on argument
func (t *DewalletChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("Invoking Dewallet Chaincode")

	function, args := stub.GetFunctionAndParameters()

	if function == "Register" {
		// Deletes an entity from its state
		return t.Register(stub, args)
	}

	if function == "UpdateUserData" {
		return t.UpdateUserData(stub, args)
	}

	if function == "GetPublicKey" {
		// queries an entity state
		return t.GetPublicKey(stub, args)
	}

	if function == "GetUserData" {
		return t.GetUserData(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'Register', 'GetPublicKey'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'Register', 'GetPublicKey'. But got: %v", args[0]))
}

// Register will add the user identity into blockchain
func (t *DewalletChaincode) Register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("Registering a member")

	var i Identity
	json.Unmarshal([]byte(args[0]), &i)

	iBytes, _ := json.Marshal(i)
	err := stub.PutState(i.Username, iBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(iBytes)
}

type updateUserDataRequest struct {
	Username string `json:"username"`
	Data string `json:"data"`
}

type updateUserDataResponse struct {
	Data string `json:"data"`
}

// UpdateUserData will query the blockchain
// and update the encrypted data
func (t *DewalletChaincode) UpdateUserData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("Updating data of user")

	var r updateUserDataRequest
	json.Unmarshal([]byte(args[0]), &r)

	iBytes, err := stub.GetState(r.Username)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if iBytes == nil {
		return shim.Error("Username not found")
	}

	var i Identity
	json.Unmarshal([]byte(iBytes), &i)
	i.Data = r.Data

	err = stub.PutState(i.Username, iBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(iBytes)
}

type getPublicKeyRequest struct {
	Username string `json:"username"`
}

type getPublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

// GetPublicKey will query the blockchain
// to get the public key of a username
func (t *DewalletChaincode) GetPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("Querying a member public key")

	var req getPublicKeyRequest
	json.Unmarshal([]byte(args[0]), &req)

	iBytes, err := stub.GetState(req.Username)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if iBytes == nil {
		return shim.Error("Username not found")
	}

	var i Identity
	json.Unmarshal([]byte(iBytes), &i)

	res := getPublicKeyResponse{
		PublicKey: i.PublicKey,
	}

	resBytes, _ := json.Marshal(res)

	return shim.Success(resBytes)
}

type getUserDataRequest struct {
	Username string `json:"username"`
}

type getUserDataResponse struct {
	Data string `json:"data"`
}

// GetUserData will query the blockchain
// and return encrypted data of a user
func (t *DewalletChaincode) GetUserData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("Querying a user data")

	var req getUserDataRequest
	json.Unmarshal([]byte(args[0]), &req)

	iBytes, err := stub.GetState(req.Username)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if iBytes == nil {
		return shim.Error("Username not found")
	}

	var i Identity
	json.Unmarshal([]byte(iBytes), &i)

	res := getUserDataResponse{
		Data: i.Data,
	}

	resBytes, _ := json.Marshal(res)

	return shim.Success(resBytes)
}

func main() {
	err := shim.Start(new(DewalletChaincode))
	if err != nil {
		logger.Errorf("Error starting Dewallet chaincode: %s", err)
	}
}
