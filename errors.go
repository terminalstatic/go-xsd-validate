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

// Returned when initialization an initialization error occured.
type Libxml2Error struct {
	errorMessage
}

// Returned when xml parsing caused error(s).
type XmlParserError struct {
	errorMessage
}

// Returned when xsd parsing caused a error(s).
type XsdParserError struct {
	errorMessage
}

// Returned when validation caused an error, to access the fields use type assertion (see example).
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
