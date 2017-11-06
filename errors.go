package xsdvalidate

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

// Returned when initialization problems occured.
type Libxml2Error struct {
	errorMessage
}

// Returned when xml parsing caused a problem.
type XmlParserError struct {
	errorMessage
}

// Returned when xsd parsing caused a problem.
type XsdParserError struct {
	errorMessage
}

// Returned when validation caused a problem, to access the fields use type assertion.
type ValidationError struct {
	Code     int
	Message  string
	Level    int
	Line     int
	NodeName string
}

// Implementation of Stringer interface.
func (ve ValidationError) String() string {
	return ve.Message
}

// Implementation of Error interface.
func (ve ValidationError) Error() string {
	return ve.String()
}
