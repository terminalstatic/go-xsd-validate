// A simple package for xsd validation, uses
//
// The rationale behind this package is to preload xsd files and use their in-memory structure to validate incoming xml documents in a concurrent environment, eg. the post bodys of xml service endpoints, and return useful error messages when appropriate.
//
// This package is part of a rewrite for an online service I'm currently working on. As there are limited resources on how to handle xsd validation in go this might be useful for somebody.
//
// libxml2-dev is needed, below an example how to install the latest sources (Ubuntu, change prefix according to where libs and include files are located):
//  curl -sL ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz | tar -xzf -
//  cd ./libxml2-2.9.5/
//  ./configure --prefix=/usr  --enable-static --with-threads --with-history
//  make
//  sudo make install
package xsdvalidate

import "C"
import (
	"errors"

	"github.com/terminalstatic/go-xsd-validate/common"
	"github.com/terminalstatic/go-xsd-validate/libxml2"
)

// Handles schema parsing and validation, wraps CXsdHandler.
type XsdHandler struct {
	schemaPtr *libxml2.SchemaPtr
}

// Handles xml parsing, wraps CXmlHandler.
type XmlHandler struct {
	docPtr *libxml2.DocPtr
}

// Initialize the xml handler struct.
// Always use the Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXmlHandlerMem(inXml []byte, options common.Options) (*XmlHandler, error) {
	xPtr, err := libxml2.ParseXmlMem(inXml, options)
	docPtr := libxml2.DocPtr(xPtr)
	return &XmlHandler{&docPtr}, err
}

// Initialize the xml handler struct.
// Always use Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXsdHandlerUrl(url string, options common.Options) (*XsdHandler, error) {
	sPtr, err := libxml2.ParseUrlSchema(url, options)
	schemaPtr := libxml2.SchemaPtr(sPtr)
	return &XsdHandler{&schemaPtr}, err
}

// The validation method validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
// Both xmlHandler and xsdHandler have to be created first with the appropriate New... functions.
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options common.Options) error {
	if xsdHandler.schemaPtr == nil {
		return errors.New("Xsd handler not properly initialized, use 'New...'")

	}
	if xmlHandler.docPtr == nil {
		return errors.New("Xml handler not properly initialized, use 'New...'")
	}
	return libxml2.ValidateWithXsd(xmlHandler.docPtr, xsdHandler.schemaPtr)

}

// Frees the schemaPtr, call this when this handler is not needed anymore.
func (xsdHandler *XsdHandler) Free() {
	libxml2.FreeSchemaPtr(xsdHandler.schemaPtr)
}

// Frees the xml docPtr, call this when this handler is not needed anymore.
func (xmlHandler *XmlHandler) Free() {
	libxml2.FreeDocPtr(xmlHandler.docPtr)
}
