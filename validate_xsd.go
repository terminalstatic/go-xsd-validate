// A simple package for xsd validation, uses libxml2.
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
	"errors"
	"sync/atomic"
	"time"
)

var isInitialized uint32

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
	//runtime.GOMAXPROCS(2)
	if atomic.LoadUint32(&isInitialized) == 0 {
		libXml2Init()
		atomic.StoreUint32(&isInitialized, 1)
	} else {
		return errors.New("Libxml2 already initialized")
	}
	return nil
}

// Initializes lbxml2 with a goroutine which trims memory and runs the go gc every d duration.
// Helps to keep the memory footprint reasonable when doing millions of operations.
func InitWithGc(d time.Duration) {
	quit = make(chan struct{})
	go gcTicker(d, quit)
}

// Cleans up the libxml2 parser, use this when application ends or libxml2 is not needed anymore.
func Cleanup() {
	libXml2Cleanup()
	atomic.StoreUint32(&isInitialized, 0)
	if quit != nil {
		quit <- struct{}{}
		quit = nil
	}
}

// Initialize the xml handler struct.
// Always use the Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error) {
	if atomic.LoadUint32(&isInitialized) == 0 {
		return nil, errors.New("Libxml2 not initialized")
	}

	xPtr, err := parseXmlMem(inXml, options)
	return &XmlHandler{xPtr}, err
}

// Initialize the xml handler struct.
// Always use Free() method when done using this handler or memory will be leaking.
// The go garbage collector will not collect the allocated resources.
func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error) {
	if atomic.LoadUint32(&isInitialized) == 0 {
		return nil, errors.New("Libxml2 not initialized")
	}
	sPtr, err := parseUrlSchema(url, options)
	return &XsdHandler{sPtr}, err
}

// The validation method validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
// Both xmlHandler and xsdHandler have to be created first with the appropriate New... functions.
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error {
	if atomic.LoadUint32(&isInitialized) == 0 {
		return errors.New("Libxml2 not initialized")
	}
	if xsdHandler == nil || xsdHandler.schemaPtr == nil {
		return errors.New("Xsd handler not properly initialized")

	}
	if xmlHandler == nil || xmlHandler.docPtr == nil {
		return errors.New("Xml handler not properly initialized")
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
