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

	handler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParsErrDefault)
	if err != nil {
		fmt.Printf("%s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}
func TestXsdUrlHandlerFail(t *testing.T) {
	Init()
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("examples/test1_fail.xsd", ParsErrVerbose)
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

	handler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
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

	handler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
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

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParsErrDefault)
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

	xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateWithXsdHandlerHttpPass(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("http://schemas.opengis.net/cat/csw/3.0/cswAll.xsd", ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("examples/test_csw.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrVerbose)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateWithXsdHandlerFail(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParsErrVerbose)
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

	xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
	fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
}

func TestIsInitialized(t *testing.T) {
	if err := Init(); err != nil {
		t.Fail()
	}
	if err := Init(); err == nil {
		t.Fail()
	}
	Cleanup()
	if err := Init(); err != nil {
		t.Fail()
	}
	Cleanup()
}
