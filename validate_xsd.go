// A simple package for xsd validation, wraps libxml2.
//
// The rationale behind this package is to provide a way to preload xsd files and use their memory representations to validate
// incoming xml documents in a concurrent environment, eg. the post bodys of xml service endpoints.
//
// Currently this package aims to be a proof of concept and is not used in production.
// Nevertheless there are limited resources on how to handle xsd validation in go, so this might be helpful for someone.
//
// As libxml2-dev is needed, as a reference how I installed the latest sources on my system (Ubuntu):
//  curl -sL ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz | tar -xzf -
//  cd ./libxml2-2.9.5/
//  ./configure --prefix=/usr  --enable-static --with-threads --with-history
//  make
//  sudo make install
package xsdvalidate

import "C"
import (
	"errors"
	"log"
)

type Param uint8

//Parameter constant(s) for possible future use, e.g. parser error detail level.
const (
	ParserDefault Param = 1 << iota
)

// Initializes the libxml2 parser, suggested for multithreading.
func Init() {
	log.Printf("Initializing the libxml2 parser")
	libXml2Init()
}

// Cleans up the libxml2 parser, use this when application ends or parser is not needed anymore.
func Cleanup() {
	log.Printf("Cleaning up the libxml2 parser")
	libXml2Cleanup()
}

// Initialize the xml handler struct.
// Always use the Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXmlHandlerMem(inXml []byte, param Param) (*XmlHandler, error) {
	xPtr, err := parseXmlMem(inXml)
	return &XmlHandler{xPtr}, err
}

// Initialize the xml handler struct.
// Always use Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXsdHandlerUrl(url string, param Param) (*XsdHandler, error) {
	sPtr, err := parseUrlSchema(url)
	return &XsdHandler{sPtr}, err
}

// The validation method validates an xmlHandler against an xsdHandler.
// Both xmlHandler and xsdHandler have to be created first with the appropriate New... functions.
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, param Param) error {
	if xsdHandler == nil || xsdHandler.schemaPtr == nil {
		return errors.New("Xsd handler not properly initialized, use 'New...'")

	}
	if xsdHandler == nil || xmlHandler.docPtr == nil {
		return errors.New("Xml handler not properly initialized, use 'New...'")
	}
	return validateWithXsd(xmlHandler, xsdHandler)

}

// Frees the schemaPtr, call this when handler is not needed anymore.
func (xsdHandler *XsdHandler) Free() {
	freeSchemaPtr(xsdHandler)
}

// Frees the xml docPtr, call this when handler is not needed anymore.
func (xmlHandler *XmlHandler) Free() {
	freeDocPtr(xmlHandler)
}
