package xsdvalidate

// Common error for default String and Error implementations.
type CommonError struct {
	Message string
}

// Implementation of Stringer Interface
func (e CommonError) String() string {
	return e.Message
}

// Implementation of Error Interface
func (e CommonError) Error() string {
	return e.String()
}

// Returned when initialization problems occured.
type Libxml2Error struct {
	CommonError
}

// Returned when xml parsing caused a problem.
type XmlParserError struct {
	CommonError
}

// Returned when xsd parsing caused a problem.
type XsdParserError struct {
	CommonError
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
