// +build apitest

package xsdvalidate

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestXsdUrlHandlerPass(t *testing.T) {
	Init()
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParserDefault)
	if err != nil {
		fmt.Printf("%s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}
func TestXsdUrlHandlerFail(t *testing.T) {
	Init()
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("examples/test1_fail.xsd", ParserDefault)
	fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
	defer handler.Free()
}
func TestXmlMemHandlerPass(t *testing.T) {
	Init()
	defer Cleanup()

	xmlFilePass, err := os.Open("examples/test1_pass.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFilePass.Close()

	inXml, _ := ioutil.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}

func TestXmlMemHandlerFail(t *testing.T) {
	Init()
	defer Cleanup()

	xmlFilePass, err := os.Open("examples/test1_fail1.xml")
	if err != nil {
		panic(err)
	}
	defer xmlFilePass.Close()

	inXml, _ := ioutil.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err == nil {
		t.Fail()
	} else {
		fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	}
	defer handler.Free()
}

func TestValidateWithXsdHandlerPass(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParserDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("examples/test1_pass.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ParserDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateWithXsdHandlerFail(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParserDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("examples/test1_fail2.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ParserDefault)
	fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
}
