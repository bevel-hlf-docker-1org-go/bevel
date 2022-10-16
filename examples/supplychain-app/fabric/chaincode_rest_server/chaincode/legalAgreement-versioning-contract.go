/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract contract for managing CRUD for Lglagrmt
type SmartContract struct {
	contractapi.Contract
}

// LglagrmtExists returns true when asset with given ID exists in world state
func (c *SmartContract) LegalAgreementVersioningExists(ctx contractapi.TransactionContextInterface, legalAgreementVersioningID string) (bool, error) {
	data, err := ctx.GetStub().GetState(legalAgreementVersioningID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateLglagrmt creates a new instance of Lglagrmt
func (c *SmartContract) CreateLegalAgreementVersioning(ctx contractapi.TransactionContextInterface, legalAgreementVersioningID string, legalAgreementVersioningContent string, legalAgreementVersioningVersion string) error {
	exists, err := c.LegalAgreementVersioningExists(ctx, legalAgreementVersioningID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The legal agreement versioning %s already exists", legalAgreementVersioningID)
	}

	contentHash := sha256.Sum256([]byte(legalAgreementVersioningContent))

	legalAgreementVersioning := new(LegalAgreementVersioning)
	legalAgreementVersioning.ID = legalAgreementVersioningID
	legalAgreementVersioning.Content = legalAgreementVersioningContent
	legalAgreementVersioning.ContentHash = fmt.Sprintf("%x", contentHash)
	legalAgreementVersioning.Timestamp = time.Now().UTC().Unix()
	legalAgreementVersioning.Version = legalAgreementVersioningVersion

	bytes, _ := json.Marshal(legalAgreementVersioning)

	return ctx.GetStub().PutState(legalAgreementVersioningID, bytes)
}

// ReadLglagrmt retrieves an instance of Lglagrmt from the world state
func (c *SmartContract) ReadLegalAgreementVersioning(ctx contractapi.TransactionContextInterface, legalAgreementVersioningID string) (*LegalAgreementVersioning, error) {
	exists, err := c.LegalAgreementVersioningExists(ctx, legalAgreementVersioningID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", legalAgreementVersioningID)
	}

	bytes, _ := ctx.GetStub().GetState(legalAgreementVersioningID)

	legalAgreementVersioning := new(LegalAgreementVersioning)

	err = json.Unmarshal(bytes, legalAgreementVersioning)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type Lglagrmt")
	}

	return legalAgreementVersioning, nil
}

// ReadLglagrmt retrieves an instance of Lglagrmt from the world state
func (c *SmartContract) ReadLatestLegalAgreementVersioning(ctx contractapi.TransactionContextInterface) (*LegalAgreementVersioning, error) {
	// Get iterator for all entries
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("Error getting state iterator: %s", err)
	}
	defer iterator.Close()

	latestLegalAgreementVersioning := new(LegalAgreementVersioning)
	for iterator.HasNext() {
		// Get the next item
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Error getting next item: %s", err)
		}

		// Unmarshal item
		legalAgreement := new(LegalAgreementVersioning)
		err = json.Unmarshal(item.Value, &legalAgreement)
		if err != nil && err.Error() != "Not a LegalAgreement" {
			return nil, fmt.Errorf("Error unmarshaling item: %s", err)
		}

		// Update latest version
		if latestLegalAgreementVersioning.Version < legalAgreement.Version {
			latestLegalAgreementVersioning = legalAgreement
		}
	}

	return latestLegalAgreementVersioning, nil
}

// ReadLglagrmt retrieves an instance of Lglagrmt from the world state
func (c *SmartContract) ReadAllLegalAgreementVersioning(ctx contractapi.TransactionContextInterface) ([]LegalAgreementVersioning, error) {
	// Get iterator for all entries
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("Error getting state iterator: %s", err)
	}
	defer iterator.Close()

	var allLegalAgreementVersioning []LegalAgreementVersioning
	for iterator.HasNext() {
		// Get the next item
		item, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Error getting next item: %s", err)
		}

		// Unmarshal item
		legalAgreement := new(LegalAgreementVersioning)
		err = json.Unmarshal(item.Value, &legalAgreement)
		if err != nil && err.Error() != "Not a LegalAgreement" {
			return nil, fmt.Errorf("Error unmarshaling item: %s", err)
		}

		allLegalAgreementVersioning = append(allLegalAgreementVersioning, *legalAgreement)
	}

	return allLegalAgreementVersioning, nil
}
