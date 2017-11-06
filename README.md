

# xsdvalidate
`import "github.com/terminalstatic/go-xsd-validate"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
A go package for xsd validation, uses libxml2.

The rationale behind this package is to preload xsd files and use their in-memory structure to validate incoming xml documents in a concurrent environment, eg. the post bodys of xml service endpoints, and return useful error messages when appropriate. Existing packages either didn't provide error details or got stuck under load.

libxml2-dev is needed, below an example how to install the latest sources (Ubuntu, change prefix according to where libs and include files are located):


	curl -sL <a href="ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz">ftp://xmlsoft.org/libxml2/libxml2-2.9.5.tar.gz</a> | tar -xzf -
	cd ./libxml2-2.9.5/
	./configure --prefix=/usr  --enable-static --with-threads --with-history
	make
	sudo make install




## <a name="pkg-index">Index</a>
* [func Cleanup()](#Cleanup)
* [func Init() error](#Init)
* [func InitWithGc(d time.Duration)](#InitWithGc)
* [type Libxml2Error](#Libxml2Error)
  * [func (e Libxml2Error) Error() string](#Libxml2Error.Error)
  * [func (e Libxml2Error) String() string](#Libxml2Error.String)
* [type Options](#Options)
* [type ValidationError](#ValidationError)
  * [func (ve ValidationError) Error() string](#ValidationError.Error)
  * [func (ve ValidationError) String() string](#ValidationError.String)
* [type XmlHandler](#XmlHandler)
  * [func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error)](#NewXmlHandlerMem)
  * [func (xmlHandler *XmlHandler) Free()](#XmlHandler.Free)
* [type XmlParserError](#XmlParserError)
  * [func (e XmlParserError) Error() string](#XmlParserError.Error)
  * [func (e XmlParserError) String() string](#XmlParserError.String)
* [type XsdHandler](#XsdHandler)
  * [func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error)](#NewXsdHandlerUrl)
  * [func (xsdHandler *XsdHandler) Free()](#XsdHandler.Free)
  * [func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error](#XsdHandler.Validate)
* [type XsdParserError](#XsdParserError)
  * [func (e XsdParserError) Error() string](#XsdParserError.Error)
  * [func (e XsdParserError) String() string](#XsdParserError.String)

#### <a name="pkg-examples">Examples</a>
An example on how to use the package. Init() is only required once before parsing and validating, and Cleanup() respectively when finished.

Code:

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
	    fmt.Printf("Error in line: %d\n", err.(xsdvalidate.ValidationError).Line)
	    fmt.Println(err)
	}

#### <a name="pkg-files">Package files</a>
[errors.go](./errors.go) [libxml2.go](./libxml2.go) [validate_xsd.go](./validate_xsd.go) 





## <a name="Cleanup">func</a> [Cleanup](./validate_xsd.go?s=2278:2292#L81)
``` go
func Cleanup()
```
Cleans up libxml2 memory and finishes gc goroutine when running.



## <a name="Init">func</a> [Init](./validate_xsd.go?s=1725:1742#L60)
``` go
func Init() error
```
Initializes libxml2, see <a href="http://xmlsoft.org/threads.html">http://xmlsoft.org/threads.html</a>.



## <a name="InitWithGc">func</a> [InitWithGc](./validate_xsd.go?s=2114:2146#L74)
``` go
func InitWithGc(d time.Duration)
```
Initializes lbxml2 with a goroutine which trims memory and runs the go gc every d duration.
Not required but can help to keep the memory footprint at bay when doing tons of validations.




## <a name="Libxml2Error">type</a> [Libxml2Error](./errors.go?s=370:412#L19)
``` go
type Libxml2Error struct {
    // contains filtered or unexported fields
}
```
Returned when initialization problems occured.










### <a name="Libxml2Error.Error">func</a> (Libxml2Error) [Error](./errors.go?s=259:295#L14)
``` go
func (e Libxml2Error) Error() string
```
Implementation of the Error Interface.




### <a name="Libxml2Error.String">func</a> (Libxml2Error) [String](./errors.go?s=156:193#L9)
``` go
func (e Libxml2Error) String() string
```
Implementation of the Stringer Interface.




## <a name="Options">type</a> [Options](./validate_xsd.go?s=1239:1257#L44)
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










## <a name="ValidationError">type</a> [ValidationError](./errors.go?s=687:794#L34)
``` go
type ValidationError struct {
    Code     int
    Message  string
    Level    int
    Line     int
    NodeName string
}
```
Returned when validation caused a problem, to access the fields use type assertion.










### <a name="ValidationError.Error">func</a> (ValidationError) [Error](./errors.go?s=941:981#L48)
``` go
func (ve ValidationError) Error() string
```
Implementation of Error interface.




### <a name="ValidationError.String">func</a> (ValidationError) [String](./errors.go?s=837:878#L43)
``` go
func (ve ValidationError) String() string
```
Implementation of Stringer interface.




## <a name="XmlHandler">type</a> [XmlHandler](./libxml2.go?s=7211:7257#L308)
``` go
type XmlHandler struct {
    // contains filtered or unexported fields
}
```
Handles xml parsing, wraps a pointer to libxml2's xmlDocPtr.







### <a name="NewXmlHandlerMem">func</a> [NewXmlHandlerMem](./validate_xsd.go?s=2620:2693#L95)
``` go
func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error)
```
Initialize the xml handler struct.
Always use the Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XmlHandler.Free">func</a> (\*XmlHandler) [Free](./validate_xsd.go?s=4191:4227#L139)
``` go
func (xmlHandler *XmlHandler) Free()
```
Frees the xml docPtr, call this when this handler is not needed anymore.




## <a name="XmlParserError">type</a> [XmlParserError](./errors.go?s=461:505#L24)
``` go
type XmlParserError struct {
    // contains filtered or unexported fields
}
```
Returned when xml parsing caused a problem.










### <a name="XmlParserError.Error">func</a> (XmlParserError) [Error](./errors.go?s=259:295#L14)
``` go
func (e XmlParserError) Error() string
```
Implementation of the Error Interface.




### <a name="XmlParserError.String">func</a> (XmlParserError) [String](./errors.go?s=156:193#L9)
``` go
func (e XmlParserError) String() string
```
Implementation of the Stringer Interface.




## <a name="XsdHandler">type</a> [XsdHandler](./libxml2.go?s=7093:7145#L303)
``` go
type XsdHandler struct {
    // contains filtered or unexported fields
}
```
Handles schema parsing and validation, wraps a pointer to libxml2's xmlSchemaPtr.







### <a name="NewXsdHandlerUrl">func</a> [NewXsdHandlerUrl](./validate_xsd.go?s=3061:3132#L107)
``` go
func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error)
```
Initialize the xml handler struct.
Always use Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XsdHandler.Free">func</a> (\*XsdHandler) [Free](./validate_xsd.go?s=4046:4082#L134)
``` go
func (xsdHandler *XsdHandler) Free()
```
Frees the schemaPtr, call this when this handler is not needed anymore.




### <a name="XsdHandler.Validate">func</a> (\*XsdHandler) [Validate](./validate_xsd.go?s=3469:3554#L117)
``` go
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error
```
This validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
Both xmlHandler and xsdHandler have to be created first.




## <a name="XsdParserError">type</a> [XsdParserError](./errors.go?s=554:598#L29)
``` go
type XsdParserError struct {
    // contains filtered or unexported fields
}
```
Returned when xsd parsing caused a problem.










### <a name="XsdParserError.Error">func</a> (XsdParserError) [Error](./errors.go?s=259:295#L14)
``` go
func (e XsdParserError) Error() string
```
Implementation of the Error Interface.




### <a name="XsdParserError.String">func</a> (XsdParserError) [String](./errors.go?s=156:193#L9)
``` go
func (e XsdParserError) String() string
```
Implementation of the Stringer Interface.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
