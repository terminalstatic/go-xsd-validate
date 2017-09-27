package xsdvalidate

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestXsdUrlHandlerPass(t *testing.T) {

	handler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}
func TestXsdUrlHandlerFail(t *testing.T) {
	handler, err := NewXsdHandlerUrl("examples/test1_fail.xsd", ParserDefault)
	t.Logf("Error: %s %s", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
	defer handler.Free()
}
func TestXmlMemHandlerPass(t *testing.T) {
	xmlFilePass, err := os.Open("examples/test1_pass.xml")
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		return
	}
	defer xmlFilePass.Close()

	inXml, _ := ioutil.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}

func TestXmlMemHandlerFail(t *testing.T) {
	xmlFilePass, err := os.Open("examples/test1_fail1.xml")
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		return
	}
	defer xmlFilePass.Close()

	inXml, _ := ioutil.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, ParserDefault)
	t.Logf("%s %s", t.Name(), err.Error())
	if err == nil {
		t.Log(err.Error())
	}
	defer handler.Free()
}

func TestValidateWithXsdHandlerPass(t *testing.T) {
	xsdhandler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("examples/test1_pass.xml")
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateWithXsdHandlerFail(t *testing.T) {
	xsdhandler, err := NewXsdHandlerUrl("examples/test1_split.xsd", ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}
	defer xsdhandler.Free()

	xmlFile, err := os.Open("examples/test1_fail2.xml")
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		return
	}
	defer xmlFile.Close()
	inXml, _ := ioutil.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, ParserDefault)
	if err != nil {
		t.Logf("%s %s", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ParserDefault)
	t.Logf("%s %s", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
}
