// +build apitest

package xsdv

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/terminalstatic/go-xsd-validate/libxml2"
)

func TestXsdUrlHandlerPass(t *testing.T) {
	l2 := libxml2.NewInit(10)
	defer l2.Shutdown()

	handler, err := NewXsdHandlerUrl("../examples/test1_split.xsd", libxml2.ParsErrDefault)
	if err != nil {
		fmt.Printf("%s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}
func TestXsdUrlHandlerFail(t *testing.T) {
	l2 := libxml2.NewInit(10)
	defer l2.Shutdown()

	handler, err := NewXsdHandlerUrl("../examples/test1_fail.xsd", libxml2.ParsErrVerbose)
	fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
	defer handler.Free()
}
func TestXmlMemHandlerPass(t *testing.T) {
	l2 := libxml2.NewInit(10)
	defer l2.Shutdown()

	xmlFilePass, err := os.Open("../examples/test1_pass.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFilePass.Close()

	inXml, _ := ioutil.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, libxml2.ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}

func TestXmlMemHandlerFail(t *testing.T) {
	l2 := libxml2.NewInit(10)
	defer l2.Shutdown()

	xmlFilePass, err := os.Open("../examples/test1_fail1.xml")
	if err != nil {
		panic(err)
	}
	defer xmlFilePass.Close()

	inXml, _ := ioutil.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, libxml2.ParsErrDefault)
	if err == nil {
		t.Fail()
	} else {
		fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	}
	defer handler.Free()
}

func TestValidateWithXsdHandlerPass(t *testing.T) {
	l2 := libxml2.NewInit(10)
	defer l2.Shutdown()

	xsdhandler, err := NewXsdHandlerUrl("../examples/test1_split.xsd", libxml2.ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("../examples/test1_pass.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, libxml2.ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, libxml2.ValidErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateWithXsdHandlerFail(t *testing.T) {
	l2 := libxml2.NewInit(10)
	defer l2.Shutdown()

	xsdhandler, err := NewXsdHandlerUrl("../examples/test1_split.xsd", libxml2.ParsErrVerbose)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("../examples/test1_fail2.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, libxml2.ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, libxml2.ValidErrDefault)
	fmt.Printf("Error OK: %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
}
