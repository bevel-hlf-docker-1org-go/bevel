/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testLegalAgreementVersioning := new(LegalAgreementVersioning)
	testLegalAgreementVersioning.ID = "001"
	testLegalAgreementVersioning.Content = "this is a content"
	testLegalAgreementVersioning.ContentHash = "this is a content"
	testLegalAgreementVersioning.Timestamp = time.Now().UTC().Unix()
	testLegalAgreementVersioning.Version = "1.0.0"
	lglagrmtBytes, _ := json.Marshal(testLegalAgreementVersioning)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "lglagrmtkey").Return(lglagrmtBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestLegalAgreementVersioningExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(SmartContract)

	exists, err = c.LegalAgreementVersioningExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.LegalAgreementVersioningExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.LegalAgreementVersioningExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateLegalAgreementVersioninging(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(SmartContract)

	err = c.CreateLegalAgreementVersioning(ctx, "statebad", "some value", "1.0.0")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateLegalAgreementVersioning(ctx, "existingkey", "some value", "1.0.0")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateLegalAgreementVersioning(ctx, "missingkey", "some value", "1.0.0")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadLegalAgreementVersioninging(t *testing.T) {
	var lglagrmt *LegalAgreementVersioning
	var err error

	ctx, _ := configureStub()
	c := new(SmartContract)

	lglagrmt, err = c.ReadLegalAgreementVersioning(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, lglagrmt, "should not return LegalAgreementVersioning when exists errors when reading")

	lglagrmt, err = c.ReadLegalAgreementVersioning(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, lglagrmt, "should not return LegalAgreementVersioning when key does not exist in world state when reading")

	lglagrmt, err = c.ReadLegalAgreementVersioning(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type LegalAgreementVersioning", "should error when data in key is not LegalAgreementVersioning")
	assert.Nil(t, lglagrmt, "should not return LegalAgreementVersioning when data in key is not of type LegalAgreementVersioning")

	lglagrmt, err = c.ReadLegalAgreementVersioning(ctx, "lglagrmtkey")
	expectedLegalAgreementVersioning := new(LegalAgreementVersioning)
	expectedLegalAgreementVersioning.ID = "001"
	expectedLegalAgreementVersioning.Content = "this is a content"
	expectedLegalAgreementVersioning.ContentHash = "this is a content"
	expectedLegalAgreementVersioning.Timestamp = time.Now().UTC().Unix()
	expectedLegalAgreementVersioning.Version = "1.0.0"
	assert.Nil(t, err, "should not return error when LegalAgreementVersioning exists in world state when reading")
	assert.Equal(t, expectedLegalAgreementVersioning, lglagrmt, "should return deserialized LegalAgreementVersioning from world state")
}
