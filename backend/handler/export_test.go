// This file exposes internal logic from the handler package to external tests (ex. test files in the "handler_test"
// package). Functions should be added here instead of exporting them to just use in a test.

package handler

// Exposes generateHmacSha256Signature to external tests as GenerateHmacSha256Signature.
var GenerateHmacSha256Signature = generateHmacSha256Signature
