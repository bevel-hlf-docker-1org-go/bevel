/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// LglagrmtExists returns true when asset with given ID exists in world state
func (c *SmartContract) LegalAgreementSigninExists(ctx contractapi.TransactionContextInterface, legalAgreementSigninID string) (bool, error) {
	data, err := ctx.GetStub().GetState(legalAgreementSigninID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateLglagrmt creates a new instance of Lglagrmt
func (c *SmartContract) CreateSigninLegalAgreementVersion(
	ctx contractapi.TransactionContextInterface,
	legalAgreementSigninID string,
	legalAgreementSigninUserID string,
	legalAgreementVersioningID string,
	legalAgreementVersioningContentHash string,
	legalAgreementSigninAccepted bool) error {
	exists, err := c.LegalAgreementSigninExists(ctx, legalAgreementSigninID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The legal agreement signin %s already exists", legalAgreementSigninID)
	}

	legalAgreementVersioning, errLegalAgreementVersioning := c.ReadLegalAgreementVersioning(ctx, legalAgreementVersioningID)
	if errLegalAgreementVersioning != nil {
		return fmt.Errorf("%s", err)
	}

	if legalAgreementVersioning.ContentHash != legalAgreementVersioningContentHash {
		return fmt.Errorf("Content hash does not match with legal agreement")
	}

	legalAgreementSignin := new(LegalAgreementSignin)
	legalAgreementSignin.ID = legalAgreementSigninID
	legalAgreementSignin.UserID = legalAgreementSigninUserID
	legalAgreementSignin.LegalAgreementVersioningID = legalAgreementSigninUserID
	legalAgreementSignin.LegalAgreementVersioningContentHash = legalAgreementVersioningContentHash
	legalAgreementSignin.Accepted = legalAgreementSigninAccepted
	legalAgreementSignin.Timestamp = time.Now().UTC().Unix()

	bytes, _ := json.Marshal(legalAgreementSignin)

	return ctx.GetStub().PutState(legalAgreementSigninID, bytes)
}

// ReadLglagrmt retrieves an instance of Lglagrmt from the world state
func (c *SmartContract) ReadSigninLegalAgreementVersion(ctx contractapi.TransactionContextInterface, legalAgreementSigninID string) (*LegalAgreementSignin, error) {
	exists, err := c.LegalAgreementSigninExists(ctx, legalAgreementSigninID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The legal agreement signin %s does not exist", legalAgreementSigninID)
	}

	bytes, _ := ctx.GetStub().GetState(legalAgreementSigninID)

	legalAgreementSignin := new(LegalAgreementSignin)

	err = json.Unmarshal(bytes, legalAgreementSignin)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type Lglagrmt")
	}

	return legalAgreementSignin, nil
}
