//go:build apitest
// +build apitest

package xsdvalidate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAddressUrlHandlerPass(t *testing.T) {
	Init()
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("./examples/test_address.xsd", ParsErrVerbose)
	if err != nil {
		fmt.Printf("%s %s\n", t.Name(), err.Error())
		t.Fail()
	}
	defer handler.Free()
}

func TestXsdUrlHandlerPass(t *testing.T) {
	Init()
	defer Cleanup()

	handler, err := NewXsdHandlerUrl("./examples/test1_split.xsd", ParsErrVerbose)
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
	fmt.Printf("Error OK:\n%s %s\n", t.Name(), err.Error())
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

	handler, err := NewXmlHandlerMem(inXml, ParsErrVerbose)
	if err == nil {
		t.Fail()
	} else {
		fmt.Printf("Error OK:\n%s %s\n", t.Name(), err.Error())
	}
	defer handler.Free()
}

func TestXmlMemHandlerLargeXml(t *testing.T) {
	Init()
	defer Cleanup()

	inXml := []byte(`<?xml version="1.0" encoding="UTF-8"?><root>` +
		strings.Repeat("<a>", 300) +
		`<data>` + strings.Repeat("A", 10*1024*1024+1) + `</data>` +
		strings.Repeat("</a>", 300) +
		`</root>`)

	handler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
	if err == nil {
		defer handler.Free()
		t.Fatal("expected parsing huge XML without ParsXmlHuge to fail")
	}

	handler, err = NewXmlHandlerMem(inXml, ParsXmlHuge)
	if err != nil {
		t.Fatalf("expected parsing huge XML with ParsXmlHuge to succeed: %v", err)
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
	fmt.Printf("Error OK:\n %s %s\n", t.Name(), err.Error())
	if err == nil {
		t.Fail()
	}
	if validationErr, ok := err.(ValidationError); ok {
		if len(validationErr.Errors) == 0 || validationErr.Errors[0].NodePath == "" {
			t.Fatalf("expected validation error node path, got %#v", validationErr.Errors)
		}
	} else {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestValidateWithXsdHandlerNodePath(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("testdata/nodepath.xsd", ParsErrVerbose)
	if err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	defer xsdhandler.Free()

	inXml := readTestFile(t, "testdata/nodepath_bad_nested_value.xml")
	xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
	if err != nil {
		t.Fatalf("failed to parse xml: %v", err)
	}
	defer xmlhandler.Free()

	err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
	validationErr := requireValidationError(t, err)
	requireStructError(t, validationErr, StructError{
		Line:     6,
		NodeName: "quantity",
		NodePath: "/order/items/item/quantity",
		Message:  "'many' is not a valid value",
	})
}

func TestValidateMemWithXsdHandlerNodePath(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("testdata/nodepath.xsd", ParsErrVerbose)
	if err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	defer xsdhandler.Free()

	tests := []struct {
		name string
		file string
		want StructError
	}{
		{
			name: "missing child",
			file: "testdata/nodepath_missing_child.xml",
			want: StructError{
				Line:     3,
				NodeName: "items",
				NodePath: "/order/items",
				Message:  "This element is not expected. Expected is ( customer )",
			},
		},
		{
			name: "nested value",
			file: "testdata/nodepath_bad_nested_value.xml",
			want: StructError{
				Line:     6,
				NodeName: "quantity",
				NodePath: "/order/items/item/quantity",
				Message:  "'many' is not a valid value",
			},
		},
		{
			name: "attribute",
			file: "testdata/nodepath_bad_attribute.xml",
			want: StructError{
				Line:     2,
				NodeName: "order",
				NodePath: "/order",
				Message:  "'bad-id' is not a valid value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := xsdhandler.ValidateMem(readTestFile(t, tt.file), ParsErrDefault)
			validationErr := requireValidationError(t, err)
			requireStructError(t, validationErr, tt.want)
		})
	}
}

func TestValidateMemWithXsdHandlerMultipleNodePaths(t *testing.T) {
	Init()
	defer Cleanup()

	xsdhandler, err := NewXsdHandlerUrl("testdata/nodepath.xsd", ParsErrVerbose)
	if err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}
	defer xsdhandler.Free()

	err = xsdhandler.ValidateMem(readTestFile(t, "testdata/nodepath_multiple_errors.xml"), ParsErrDefault)
	validationErr := requireValidationError(t, err)
	requireStructError(t, validationErr, StructError{
		Line:     2,
		NodeName: "order",
		NodePath: "/order",
		Message:  "'bad-id' is not a valid value",
	})
	requireStructError(t, validationErr, StructError{
		Line:     6,
		NodeName: "quantity",
		NodePath: "/order/items/item/quantity",
		Message:  "'many' is not a valid value",
	})
}

func readTestFile(t *testing.T, path string) []byte {
	t.Helper()

	in, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return in
}

func requireValidationError(t *testing.T, err error) ValidationError {
	t.Helper()

	if err == nil {
		t.Fatal("expected validation error")
	}
	validationErr, ok := err.(ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(validationErr.Errors) == 0 {
		t.Fatal("expected at least one structured validation error")
	}
	return validationErr
}

func requireStructError(t *testing.T, validationErr ValidationError, want StructError) {
	t.Helper()

	for _, got := range validationErr.Errors {
		if got.Line == want.Line && got.NodeName == want.NodeName && got.NodePath == want.NodePath && strings.Contains(got.Message, want.Message) {
			return
		}
	}
	t.Fatalf("expected validation error matching %#v, got %#v", want, validationErr.Errors)
}

func TestValidateMemWithXsdHandlerPass(t *testing.T) {
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

	err = xsdhandler.ValidateMem(inXml, ParsErrDefault)
	if err != nil {
		fmt.Printf("Error: %s %s\n", t.Name(), err.Error())
		t.Fail()
	}

}

func TestValidateMemWithXsdHandlerFail(t *testing.T) {
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
