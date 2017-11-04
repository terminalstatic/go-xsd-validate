// A go package for xsd validation, uses libxml2.
//
// The rationale behind this package is to preload xsd files and use their in-memory structure to validate incoming xml documents in a concurrent environment, eg. the post bodys of xml service endpoints, and return useful error messages when appropriate. Existing packages either didn't provide error details or got stuck under load.
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

// The type for parser/validation options.
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

// Initializes libxml2, suggested for multithreading, see http://xmlsoft.org/threads.html.
func Init() error {
	g.Lock()
	defer g.Unlock()
	if g.isInitialized() {
		return Libxml2Error{CommonError{"Libxml2 already initialized"}}
	}

	libXml2Init()
	g.setInitialized(true)
	return nil
}

// Initializes lbxml2 with a goroutine which trims memory and runs the go gc every d duration.
// Not required but can help to keep the memory footprint at bay when doing tons of validations.
func InitWithGc(d time.Duration) {
	Init()
	quit = make(chan struct{})
	go gcTicker(d, quit)
}

// Cleans up the libxml2 parser, use this when application ends or libxml2 is not needed anymore.
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

// Initialize the xml handler struct.
// Always use the Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error) {
	if !g.isInitialized() {
		return nil, Libxml2Error{CommonError{"Libxml2 not initialized"}}
	}

	xPtr, err := parseXmlMem(inXml, options)
	return &XmlHandler{xPtr}, err
}

// Initialize the xml handler struct.
// Always use Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error) {
	if !g.isInitialized() {
		return nil, Libxml2Error{CommonError{"Libxml2 not initialized"}}
	}
	sPtr, err := parseUrlSchema(url, options)
	return &XsdHandler{sPtr}, err
}

// The validation method validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
// Both xmlHandler and xsdHandler have to be created first.
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error {
	if !g.isInitialized() {
		return Libxml2Error{CommonError{"Libxml2 not initialized"}}
	}

	if xsdHandler == nil || xsdHandler.schemaPtr == nil {
		return XsdParserError{CommonError{"Xsd handler not properly initialized"}}

	}
	if xmlHandler == nil || xmlHandler.docPtr == nil {
		return XmlParserError{CommonError{"Xml handler not properly initialized"}}
	}
	return validateWithXsd(xmlHandler, xsdHandler)

}

// Frees the schemaPtr, call this when this handler is not needed anymore.
func (xsdHandler *XsdHandler) Free() {
	freeSchemaPtr(xsdHandler)
}

// Frees the xml docPtr, call this when this handler is not needed anymore.
func (xmlHandler *XmlHandler) Free() {
	freeDocPtr(xmlHandler)
}
