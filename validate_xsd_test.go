package xsdvalidate

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestAddressUrlHandlerPass(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("./examples/test_address.xsd", ParsErrVerbose)
	if err != nil {
		fmt.Printf("%s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}

func TestXsdUrlHandlerPass(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("./examples/test1_split.xsd", ParsErrVerbose)
	if err != nil {
		fmt.Printf("%s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}
func TestXsdUrlHandlerFail(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("examples/test1_fail.xsd", ParsErrVerbose)
	fmt.Printf("Error OK:\n%s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
	defer handler.Free()
}
func TestXmlMemHandlerPass(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
	defer Cleanup()

	xmlFilePass, err := os.Open("examples/test1_pass.xml")
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		return
	}
	defer xmlFilePass.Close()

	inXml, _ := io.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}

func TestXmlMemHandlerFail(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
	defer Cleanup()

	xmlFilePass, err := os.Open("examples/test1_fail1.xml")
	if err != nil {
		panic(err)
	}
	defer xmlFilePass.Close()

	inXml, _ := io.ReadAll(xmlFilePass)

	handler, err := NewXmlHandlerMem(inXml, ParsErrVerbose)
	if err == nil {
		t.Fail()
	} else {
		fmt.Printf("Error OK:\n%s %s\n", t.Name(), err.Error())
	}
	defer handler.Free()
}

func TestValidateWithXsdHandlerPass(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
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
	inXml, _ := io.ReadAll(xmlFile)

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
	err := Init()
	if err != nil {
		panic(err)
	}
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
	inXml, _ := io.ReadAll(xmlFile)

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
	err := Init()
	if err != nil {
		panic(err)
	}
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
	inXml, _ := io.ReadAll(xmlFile)

	xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
	fmt.Printf("Error OK:\n %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
}

func TestValidateMemWithXsdHandlerPass(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
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
	inXml, _ := io.ReadAll(xmlFile)

	err = xsdhandler.ValidateMem(inXml, ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateMemWithXsdHandlerFail(t *testing.T) {
	err := Init()
	if err != nil {
		panic(err)
	}
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
	inXml, _ := io.ReadAll(xmlFile)

	err = xsdhandler.ValidateMem(inXml, ParsErrDefault)
	fmt.Printf("Error OK:\n %s %s\n", t.Name(), err.Error())
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
