

# xsdvalidate
`import "/home/ubuntu/source/go/src/github.com/terminalstatic/go-xsd-validate"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
A simple package for xsd validation, wraps libxml2.

The rationale behind this package is to provide a way to preload xsd files and use their memory representations to validate
incoming xml documents in a concurrent environment, eg. the post bodys of xml service endpoints.

Currently this package aims to be a proof of concept and is not used in production.
Nevertheless there are limited resources on how to handle xsd validation in go, so this might be helpful for someone.

As libxml2-dev is needed, as a reference how I installed the latest sources on my system (Ubuntu):


	curl -sL <a href="ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz">ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz</a> | tar -xzf -
	cd ./libxml2-2.9.5/
	./configure --prefix=/usr  --enable-static --with-threads --with-history
	make
	sudo make install




## <a name="pkg-index">Index</a>
* [func Cleanup()](#Cleanup)
* [func Init()](#Init)
* [type Param](#Param)
* [type XmlHandler](#XmlHandler)
  * [func NewXmlHandlerMem(inXml []byte, param Param) (*XmlHandler, error)](#NewXmlHandlerMem)
  * [func (xmlHandler *XmlHandler) Free()](#XmlHandler.Free)
* [type XsdHandler](#XsdHandler)
  * [func NewXsdHandlerUrl(url string, param Param) (*XsdHandler, error)](#NewXsdHandlerUrl)
  * [func (xsdHandler *XsdHandler) Free()](#XsdHandler.Free)
  * [func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, param Param) error](#XsdHandler.Validate)

#### <a name="pkg-examples">Examples</a>
An example on how to use the api.
Always bear in mind to free the handlers, the go gc will not collect those.
On the other hand you prabably want to call xsdvalidate.Init() and xsdvalidate.Cleanup() only once after app start and before app end.


	import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/terminalstatic/go-xsd-validate"
	)

	func Example() {
		xsdvalidate.Init()
		defer xsdvalidate.Cleanup()
		xsdhandler, err := xsdvalidate.NewXsdHandlerUrl("examples/test1_split.xsd", xsdvalidate.ParserDefault)
		if err != nil {
			panic(err)
		}
		defer xsdhandler.Free()
		xmlFile, err := os.Open("examples/test1_pass.xml")
		if err != nil {
			panic(err)
		}
		defer xmlFile.Close()
		inXml, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			panic(err)
		}
		xmlhandler, err := xsdvalidate.NewXmlHandlerMem(inXml, xsdvalidate.ParserDefault)
		if err != nil {
			panic(err)
		}
		defer xmlhandler.Free()
		err = xsdhandler.Validate(xmlhandler, xsdvalidate.ParserDefault)
		if err != nil {
			panic(err)
		}
		fmt.Println("Validation OK")
		// Output: Validation OK
	}

#### <a name="pkg-files">Package files</a>
[libxml2.go](/src/target/libxml2.go) [validate_xsd.go](/src/target/validate_xsd.go) 





## <a name="Cleanup">func</a> [Cleanup](/src/target/validate_xsd.go?s=1251:1265#L37)
``` go
func Cleanup()
```
Cleans up the libxml2 parser, use this when application ends or parser is not needed anymore.



## <a name="Init">func</a> [Init](/src/target/validate_xsd.go?s=1075:1086#L31)
``` go
func Init()
```
Initializes the libxml2 parser, suggested for multithreading.




## <a name="Param">type</a> [Param](/src/target/validate_xsd.go?s=867:883#L23)
``` go
type Param uint8
```

``` go
const (
    ParserDefault Param = 1 << iota
)
```
Parameter constant(s) for possible future use, e.g. parser error detail level.










## <a name="XmlHandler">type</a> [XmlHandler](/src/target/libxml2.go?s=3907:3953#L174)
``` go
type XmlHandler struct {
    // contains filtered or unexported fields
}
```
Handles xml parsing, wraps a pointer to libxml2's xmlDocPtr.







### <a name="NewXmlHandlerMem">func</a> [NewXmlHandlerMem](/src/target/validate_xsd.go?s=1531:1600#L45)
``` go
func NewXmlHandlerMem(inXml []byte, param Param) (*XmlHandler, error)
```
Initialize the xml handler struct.
Always use the Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XmlHandler.Free">func</a> (\*XmlHandler) [Free](/src/target/validate_xsd.go?s=2778:2814#L78)
``` go
func (xmlHandler *XmlHandler) Free()
```
Frees the xml docPtr, call this when handler is not needed anymore.




## <a name="XsdHandler">type</a> [XsdHandler](/src/target/libxml2.go?s=3789:3841#L169)
``` go
type XsdHandler struct {
    // contains filtered or unexported fields
}
```
Handles schema parsing and validation, wraps a pointer to libxml2's xmlSchemaPtr.







### <a name="NewXsdHandlerUrl">func</a> [NewXsdHandlerUrl](/src/target/validate_xsd.go?s=1862:1929#L53)
``` go
func NewXsdHandlerUrl(url string, param Param) (*XsdHandler, error)
```
Initialize the xml handler struct.
Always use Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XsdHandler.Free">func</a> (\*XsdHandler) [Free](/src/target/validate_xsd.go?s=2638:2674#L73)
``` go
func (xsdHandler *XsdHandler) Free()
```
Frees the schemaPtr, call this when handler is not needed anymore.




### <a name="XsdHandler.Validate">func</a> (\*XsdHandler) [Validate](/src/target/validate_xsd.go?s=2170:2251#L60)
``` go
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, param Param) error
```
The validation method validates an xmlHandler against an xsdHandler.
Both xmlHandler and xsdHandler have to be created first with the appropriate New... functions.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
