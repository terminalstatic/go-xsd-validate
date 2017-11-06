package xsdvalidate_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/terminalstatic/go-xsd-validate"
)

// An example on how to use the package.
// In some situations, e.g. programatically looping over xml documents you might have to explicitly free the handler without defer. Calling xsdvalidate.Init() is only required once before you start parsing and validating, and xsdvalidate.Cleanup() respectively when finished.
func Example() {
	xsdvalidate.Init()
	defer xsdvalidate.Cleanup()
	xsdhandler, err := xsdvalidate.NewXsdHandlerUrl("examples/test1_split.xsd", xsdvalidate.ParsErrDefault)
	if err != nil {
		panic(err)
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("examples/test1_fail2.xml")
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()
	inXml, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}

	xmlhandler, err := xsdvalidate.NewXmlHandlerMem(inXml, xsdvalidate.ParsErrDefault)
	if err != nil {
		panic(err)
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, xsdvalidate.ValidErrDefault)
	if err != nil {
		fmt.Printf("Error in line: %d\n", err.(xsdvalidate.ValidationError).Line)
		fmt.Printf(err.(xsdvalidate.ValidationError).
			Format(func(ve xsdvalidate.ValidationError) string {
				return fmt.Sprintf("Error in Line: %d, Severity : %d, Message: %s\n", ve.Line, ve.Level, ve.Message)
			}))
		fmt.Println(err)
	}
	// Output:
	// Error in line: 3
	// Error in Line: 3, Severity : 2, Message: Element 'shipto': This element is not expected. Expected is ( orderperson ).
	// Element 'shipto': This element is not expected. Expected is ( orderperson ).
}
