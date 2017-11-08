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

// ValidationError is returned when xsd validation caused an error, to access the fields use type assertion (see example).
type ValidationError struct {
	Code     int
	Message  string
	Level    int
	Line     int
	NodeName string
}

// Implementation of the Stringer interface.
func (e ValidationError) String() string {
	return e.Message
}

// Implementation of the Error interface.
func (e ValidationError) Error() string {
	return e.String()
}
