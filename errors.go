package xsdvalidate

// Common error
type CommonError struct {
	Message string
}

// Implementation of String Interface
func (e CommonError) String() string {
	return e.Message
}

// Implementation of Error Interface
func (e CommonError) Error() string {
	return e.String()
}

//Libxml2 Error
type Libxml2Error struct {
	CommonError
}

// XmlParser Error
type XmlParserError struct {
	CommonError
}

// XsdParser Error
type XsdParserError struct {
	CommonError
}

// The validation error, to access the fields use type assertion
type ValidationError struct {
	Code     int
	Message  string
	Level    int
	Line     int
	NodeName string
}

// Implementation of Stringer interface
func (ve ValidationError) String() string {
	return ve.Message
}

// Implementation of Error interface
func (ve ValidationError) Error() string {
	return ve.String()
}
