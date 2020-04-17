package main

import (
	"encoding/json"

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
