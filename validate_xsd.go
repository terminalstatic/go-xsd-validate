// Package xsdvalidate is a go package for xsd validation that utilizes libxml2.
//
// The goal of this package is to preload xsd files and use their in-memory representation to validate xml documents in a concurrent environment, eg. the post bodys of xml service endpoints and hand through libxml2 error messages. Similar packages on github either didn't provide error details or got stuck under load.
//
// libxml2-dev is needed, below an example how to install the latest sources as at the time of writing (Ubuntu, change prefix according to where libs and include files are located):
//  curl -sL ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz | tar -xzf -
//  cd ./libxml2-2.9.5/
//  ./configure --prefix=/usr  --enable-static --with-threads --with-history
//  make
//  sudo make install
package xsdvalidate

import "C"
import (
	"sync"
	"sync/atomic"
	"time"
)

type guard struct {
	sync.Mutex
	initialized uint32
}

func (guard *guard) isInitialized() bool {
	if atomic.LoadUint32(&guard.initialized) == 0 {
		return false
	}
	return true
}

func (guard *guard) setInitialized(b bool) {
	switch b {
	case true:
		atomic.StoreUint32(&guard.initialized, 1)
	case false:
		atomic.StoreUint32(&guard.initialized, 0)
	}
}

var g guard

// Options type for parser/validation options.
type Options int16

// The parser options, ParsErrVerbose will slow down parsing considerably!
const (
	ParsErrDefault Options = 1 << iota // Default parser error output
	ParsErrVerbose                     // Verbose parser error output, considerably slower!
)

// Validation options for possible future enhancements.
const (
	ValidErrDefault Options = 1 << iota // Default validation error output
)

var quit chan struct{}

// Init initializes libxml2, see http://xmlsoft.org/threads.html.
func Init() error {
	g.Lock()
	defer g.Unlock()
	if g.isInitialized() {
		return Libxml2Error{errorMessage{"Libxml2 already initialized"}}
	}

	libXml2Init()
	g.setInitialized(true)
	return nil
}

// InitWithGc initializes lbxml2 with a goroutine that trims memory and runs the go gc every d duration.
// Not required but can help to keep the memory footprint at bay when doing tons of validations.
func InitWithGc(d time.Duration) {
	Init()
	quit = make(chan struct{})
	go gcTicker(d, quit)
}

// Cleanup cleans up libxml2 memory and finishes gc goroutine when running.
func Cleanup() {
	g.Lock()
	defer g.Unlock()
	libXml2Cleanup()
	g.setInitialized(false)
	if quit != nil {
		quit <- struct{}{}
		quit = nil
	}
}

// NewXmlHandlerMem creates a xml handler struct.
// Always use the Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error) {
	if !g.isInitialized() {
		return nil, Libxml2Error{errorMessage{"Libxml2 not initialized"}}
	}

	xPtr, err := parseXmlMem(inXml, options)
	return &XmlHandler{xPtr}, err
}

// NewXsdHandlerUrl creates a xsd handler struct.
// Always use Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error) {
	g.Lock()
	defer g.Unlock()
	if !g.isInitialized() {
		return nil, Libxml2Error{errorMessage{"Libxml2 not initialized"}}
	}
	sPtr, err := parseUrlSchema(url, options)
	return &XsdHandler{sPtr}, err
}

// Validate validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
// Both xmlHandler and xsdHandler have to be created first.
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error {
	if !g.isInitialized() {
		return Libxml2Error{errorMessage{"Libxml2 not initialized"}}
	}

	if xsdHandler == nil || xsdHandler.schemaPtr == nil {
		return XsdParserError{errorMessage{"Xsd handler not properly initialized"}}

	}
	if xmlHandler == nil || xmlHandler.docPtr == nil {
		return XmlParserError{errorMessage{"Xml handler not properly initialized"}}
	}
	return validateWithXsd(xmlHandler, xsdHandler)

}

// Free frees the wrapped schemaPtr, call this when this handler is not needed anymore.
func (xsdHandler *XsdHandler) Free() {
	freeSchemaPtr(xsdHandler)
}

// Free frees the wrapped xml docPtr, call this when this handler is not needed anymore.
func (xmlHandler *XmlHandler) Free() {
	freeDocPtr(xmlHandler)
}
