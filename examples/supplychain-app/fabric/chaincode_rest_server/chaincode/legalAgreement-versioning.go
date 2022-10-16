/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// LegalAgreementVersioning stores a value
type LegalAgreementVersioning struct {
	ID          string `json:"ID"`
	Content     string `json:"content"`
	ContentHash string `json:"hash"`
	Timestamp   int64  `json:"timestamp"`
	Version     string `json:"version"`
}
