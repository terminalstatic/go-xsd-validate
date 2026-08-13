package xsdvalidate

/*
#cgo CFLAGS: -std=c99
#cgo pkg-config: libxml-2.0
#include "xsdvalidate_libxml2.h"
*/
import "C"
import (
	"runtime"
	"strings"
	"time"
	"unsafe"
)

// XsdHandler handles schema parsing and validation and wraps a pointer to libxml2's xmlSchemaPtr.
type XsdHandler struct {
	schemaPtr C.xmlSchemaPtr
}

// XmlHandler handles xml parsing and wraps a pointer to libxml2's xmlDocPtr.
type XmlHandler struct {
	docPtr C.xmlDocPtr
}

// Initializes the libxml2 parser, suggested for multithreading
func libXml2Init() {
	C.xsdValidateInit()
}

// Cleans up the libxml2 parser
func libXml2Cleanup() {
	C.xsdValidateCleanup()
}

func byteSlicePtr(b []byte) unsafe.Pointer {
	if len(b) == 0 {
		return nil
	}
	return unsafe.Pointer(&b[0])
}

// The helper function for parsing xml
func parseXmlMem(inXml []byte, options Options) (C.xmlDocPtr, error) {
	pRes, err := C.xsdValidateParseDoc(byteSlicePtr(inXml), C.int(len(inXml)), C.short(options))
	runtime.KeepAlive(inXml)

	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, XmlParserError{errorMessage{strings.Trim(rStr, "\n")}}
	}
	return pRes.docPtr, nil
}

// The helper function for parsing the schema
func parseUrlSchema(url string, options Options) (C.xmlSchemaPtr, error) {
	strUrl := C.CString(url)
	defer C.free(unsafe.Pointer(strUrl))

	pRes, err := C.xsdValidateParseUrlSchema(strUrl, C.short(options))
	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, XsdParserError{errorMessage{strings.Trim(rStr, "\n")}}
	}
	return pRes.schemaPtr, nil
}

// The helper function for parsing an in-memory schema
func parseMemSchema(xsd []byte, options Options) (C.xmlSchemaPtr, error) {
	pRes, err := C.xsdValidateParseMemSchema(byteSlicePtr(xsd), C.int(len(xsd)), C.short(options))
	runtime.KeepAlive(xsd)
	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, XsdParserError{errorMessage{strings.Trim(rStr, "\n")}}
	}
	return pRes.schemaPtr, nil
}

func handleErrArray(errSlice []C.struct_simpleXmlError) ValidationError {
	ve := ValidationError{make([]StructError, len(errSlice))}
	for i := 0; i < len(errSlice); i++ {
		ve.Errors[i] = StructError{
			Code:     int(errSlice[i].code),
			Message:  strings.Trim(C.GoString(errSlice[i].message), "\n"),
			Level:    int(errSlice[i].level),
			Line:     int(errSlice[i].line),
			NodeName: C.GoString(errSlice[i].node),
			NodePath: C.GoString(errSlice[i].nodePath)}
	}
	return ve

}

const maxSimpleXmlErrors = 1 << 20

func simpleXmlErrorSlice(errArr C.errArray) ([]C.struct_simpleXmlError, error) {
	if errArr.len > maxSimpleXmlErrors {
		return nil, Libxml2Error{errorMessage{"too many validation errors"}}
	}
	return unsafe.Slice(errArr.data, errArr.len), nil
}

// Helper function for validating given an xml document
func validateWithXsd(xmlHandler *XmlHandler, xsdHandler *XsdHandler) error {
	sErr, err := C.xsdValidateDoc(xmlHandler.docPtr, xsdHandler.schemaPtr)
	defer C.xsdValidateFreeErrArray(&sErr)
	if err != nil {
		errSlice, sliceErr := simpleXmlErrorSlice(sErr)
		if sliceErr != nil {
			return sliceErr
		}
		return handleErrArray(errSlice)
	}
	return nil
}

// Helper function for validating given an xml byte slice
func validateBufWithXsd(inXml []byte, options Options, xsdHandler *XsdHandler) error {
	sErr, err := C.xsdValidateBuf(byteSlicePtr(inXml), C.int(len(inXml)), C.short(options), xsdHandler.schemaPtr)
	runtime.KeepAlive(inXml)
	defer C.xsdValidateFreeErrArray(&sErr)
	if err != nil {
		errSlice, sliceErr := simpleXmlErrorSlice(sErr)
		if sliceErr != nil {
			return sliceErr
		}
		switch errSlice[0]._type {
		case C.VALIDATION_ERROR:
			return handleErrArray(errSlice)
		case C.XML_PARSER_ERROR:
			return XmlParserError{errorMessage{strings.Trim(C.GoString(errSlice[0].message), "\n")}}
		case C.LIBXML2_ERROR:
			return Libxml2Error{errorMessage{strings.Trim(C.GoString(errSlice[0].message), "\n")}}
		case C.XSD_PARSER_ERROR:
			return XsdParserError{errorMessage{strings.Trim(C.GoString(errSlice[0].message), "\n")}}
		default:
			return Libxml2Error{errorMessage{"Unknown error"}}
		}
		return ValidationError{}
	}
	return nil
}

// Wrapper for the xmlSchemaFree function
func freeSchemaPtr(xsdHandler *XsdHandler) {
	if xsdHandler != nil && xsdHandler.schemaPtr != nil {
		C.xmlSchemaFree(xsdHandler.schemaPtr)
		xsdHandler.schemaPtr = nil
	}
}

// Wrapper for the xmlFreeDoc function
func freeDocPtr(xmlHandler *XmlHandler) {
	if xmlHandler != nil && xmlHandler.docPtr != nil {
		C.xmlFreeDoc(xmlHandler.docPtr)
		xmlHandler.docPtr = nil
	}
}

// Ticker for gc
func gcTicker(d time.Duration, quit chan struct{}) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			runtime.GC()
			//C.malloc_trim(0)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
