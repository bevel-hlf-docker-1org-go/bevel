/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// Lglagrmt stores a value
type LegalAgreementSignin struct {
	ID                                  string `json:"ID"`
	UserID                              string `json:"userID"`
	LegalAgreementVersioningID          string `json:"legalAgreementVersioningID"`
	LegalAgreementVersioningContentHash string `json:"legalAgreementVersioningContentHash"`
	Accepted                            bool   `json:"accepted"`
	Timestamp                           int64  `json:"timestamp"`
}
