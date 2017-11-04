package xsdvalidate_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/terminalstatic/go-xsd-validate"
)

// An example on how to use the package.
// In some situations, e.g. programatically looping over xml documents you might have to explicitly free the handler without defer. You prabably want to call xsdvalidate.Init() and xsdvalidate.Cleanup() only once after app start and before app end.
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
		fmt.Println(err)
	}
	// Output:
	// Element 'shipto': This element is not expected. Expected is ( orderperson ).
}
