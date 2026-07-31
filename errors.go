package xsdvalidate

import (
	"fmt"
	"strings"
)

// Common String and Error implementations.
type errorMessage struct {
	Message string
}

// Implementation of the Stringer Interface.
func (e errorMessage) String() string {
	return e.Message
}

// Implementation of the Error Interface.
func (e errorMessage) Error() string {
	return e.String()
}

// Libxml2Error is returned when a Libxm2 initialization error occured.
type Libxml2Error struct {
	errorMessage
}

// XmlParserError is returned when xml parsing caused error(s).
type XmlParserError struct {
	errorMessage
}

// XsdParserError is returned when xsd parsing caused a error(s).
type XsdParserError struct {
	errorMessage
}

// StructError is a subset of libxml2 xmlError struct.
type StructError struct {
	// Code is the libxml2 error code.
	Code int
	// Message is the validation error message returned by libxml2.
	Message string
	// Level is the libxml2 error level.
	Level int
	// Line is the line number associated with the validation error.
	Line int
	// NodeName is the XML node name associated with the validation error.
	NodeName string
	// NodePath is the XPath-style path to the XML node associated with the validation error.
	NodePath string
}

// ValidationError is returned when xsd validation caused an error, to access the fields of the Errors slice use type assertion (see example).
type ValidationError struct {
	Errors []StructError
}

// Implementation of the Stringer interface. Aggregates line numbers and messages of the Errors slice.
func (e ValidationError) String() string {
	var em string
	for _, eelem := range e.Errors {
		em = em + fmt.Sprintf("%d: %s\n", eelem.Line, eelem.Message)
	}
	return strings.TrimRight(em, "\n")
}

// Implementation of the Error interface.
func (e ValidationError) Error() string {
	return e.String()
}
