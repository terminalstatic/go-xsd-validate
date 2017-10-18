

# xsdvalidate
`import "github.com/terminalstatic/go-xsd-validate"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
A simple package for xsd validation, uses libxml2.

The rationale behind this package is to preload xsd files and use their in-memory structure to validate incoming xml documents in a concurrent environment, eg. the post bodys of xml service endpoints, and return useful error messages when appropriate. Existing packages either didn't provide error details or got stuck under load.

libxml2-dev is needed, below an example how to install the latest sources (Ubuntu, change prefix according to where libs and include files are located):


	curl -sL <a href="ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz">ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz</a> | tar -xzf -
	cd ./libxml2-2.9.5/
	./configure --prefix=/usr  --enable-static --with-threads --with-history
	make
	sudo make install




## <a name="pkg-index">Index</a>
* [func Cleanup()](#Cleanup)
* [func Init()](#Init)
* [type Options](#Options)
* [type XmlHandler](#XmlHandler)
  * [func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error)](#NewXmlHandlerMem)
  * [func (xmlHandler *XmlHandler) Free()](#XmlHandler.Free)
* [type XsdHandler](#XsdHandler)
  * [func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error)](#NewXsdHandlerUrl)
  * [func (xsdHandler *XsdHandler) Free()](#XsdHandler.Free)
  * [func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error](#XsdHandler.Validate)

#### <a name="pkg-examples">Examples</a>
An example on how to use the package.
In some situations, e.g. programatically looping over xml documents you might have to explicitly free the handler without defer.
You prabably want to call xsdvalidate.Init() and xsdvalidate.Cleanup() only once after app start and before app end.

	import (
		"fmt"
		"io/ioutil"
		"os"
		"github.com/terminalstatic/go-xsd-validate"
	)

	func Example() {
		xsdvalidate.Init()
		defer xsdvalidate.Cleanup()
		xsdhandler, err := xsdvalidate.NewXsdHandlerUrl("examples/test1_split.xsd", xsdvalidate.ParsErrDefault)
		if err != nil {
			panic(err)
		}
		defer xsdhandler.Free()

		xmlFile, err := os.Open("examples/test1_fail2.xml")
		if err != nil {
			panic(err)
		}
		defer xmlFile.Close()
		inXml, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			panic(err)
		}

		xmlhandler, err := xsdvalidate.NewXmlHandlerMem(inXml, xsdvalidate.ParsErrDefault)
		if err != nil {
			panic(err)
		}
		defer xmlhandler.Free()

		err = xsdhandler.Validate(xmlhandler, xsdvalidate.ValidErrDefault)
		if err != nil {
			fmt.Println(err)
		}
		// Output:
		// Validation error:
		// Element 'shipto': This element is not expected. Expected is ( orderperson ).
	}

#### <a name="pkg-files">Package files</a>
[libxml2.go](libxml2.go) [validate_xsd.go](validate_xsd.go) 





## <a name="Cleanup">func</a> [Cleanup](validate_xsd.go?s=1652:1666#L42)
``` go
func Cleanup()
```
Cleans up the libxml2 parser, use this when application ends or libxml2 is not needed anymore.



## <a name="Init">func</a> [Init](validate_xsd.go?s=1475:1486#L36)
``` go
func Init()
```
Initializes the libxml2 parser, suggested for multithreading, see (http://xmlsoft.org/threads.html).




## <a name="Options">type</a> [Options](validate_xsd.go?s=1042:1060#L22)
``` go
type Options int16
```
The type for parser/validation options.

``` go
const (
    ParsErrDefault Options = 1 << iota // Default parser error output
    ParsErrVerbose                     // Verbose parser error output, considerably slower!
)
```
The parser options, ParsErrVerbose will slow down parsing considerably!


``` go
const (
    ValidErrDefault Options = 1 << iota // Default validation error output
)
```
Validation options for possible future enhancements.
``` go
type XmlHandler struct {
    // contains filtered or unexported fields
}
```
Handles xml parsing, wraps a pointer to libxml2's xmlDocPtr.







### <a name="NewXmlHandlerMem">func</a> [NewXmlHandlerMem](validate_xsd.go?s=1932:2005#L50)
``` go
func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error)
```
Initialize the xml handler struct.
Always use the Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XmlHandler.Free">func</a> (\*XmlHandler) [Free](validate_xsd.go?s=3265:3301#L83)
``` go
func (xmlHandler *XmlHandler) Free()
```
Frees the xml docPtr, call this when this handler is not needed anymore.




## <a name="XsdHandler">type</a> [XsdHandler](libxml2.go?s=5651:5703#L242)
``` go
type XsdHandler struct {
    // contains filtered or unexported fields
}
```
Handles schema parsing and validation, wraps a pointer to libxml2's xmlSchemaPtr.







### <a name="NewXsdHandlerUrl">func</a> [NewXsdHandlerUrl](validate_xsd.go?s=2276:2347#L58)
``` go
func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error)
```
Initialize the xml handler struct.
Always use Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XsdHandler.Free">func</a> (\*XsdHandler) [Free](validate_xsd.go?s=3120:3156#L78)
``` go
func (xsdHandler *XsdHandler) Free()
```
Frees the schemaPtr, call this when this handler is not needed anymore.




### <a name="XsdHandler.Validate">func</a> (\*XsdHandler) [Validate](validate_xsd.go?s=2643:2728#L65)
``` go
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error
```
The validation method validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
Both xmlHandler and xsdHandler have to be created first with the appropriate New... functions.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
