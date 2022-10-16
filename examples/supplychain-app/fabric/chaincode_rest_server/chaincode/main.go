/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	smartContract := new(SmartContract)
	smartContract.Info.Version = "0.0.1"
	smartContract.Info.Description = "Smart Contract legal agreement versioning"
	smartContract.Info.License = new(metadata.LicenseMetadata)
	smartContract.Info.License.Name = "Apache-2.0"
	smartContract.Info.Contact = new(metadata.ContactMetadata)
	smartContract.Info.Contact.Name = "Oiga Technologies"

	chaincode, err := contractapi.NewChaincode(smartContract)
	chaincode.Info.Title = "chaincode chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from LglagrmtVersioningContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
