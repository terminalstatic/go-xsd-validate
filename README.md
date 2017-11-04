

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
* [type CommonError](#CommonError)
  * [func (e CommonError) Error() string](#CommonError.Error)
  * [func (e CommonError) String() string](#CommonError.String)
* [type Libxml2Error](#Libxml2Error)
* [type Options](#Options)
* [type ValidationError](#ValidationError)
  * [func (ve ValidationError) Error() string](#ValidationError.Error)
  * [func (ve ValidationError) String() string](#ValidationError.String)
* [type XmlHandler](#XmlHandler)
  * [func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error)](#NewXmlHandlerMem)
  * [func (xmlHandler *XmlHandler) Free()](#XmlHandler.Free)
* [type XmlParserError](#XmlParserError)
* [type XsdHandler](#XsdHandler)
  * [func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error)](#NewXsdHandlerUrl)
  * [func (xsdHandler *XsdHandler) Free()](#XsdHandler.Free)
  * [func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error](#XsdHandler.Validate)
* [type XsdParserError](#XsdParserError)

#### <a name="pkg-examples">Examples</a>
* [Package](#example_)

#### <a name="pkg-files">Package files</a>
[errors.go](/src/github.com/terminalstatic/go-xsd-validate/errors.go) [libxml2.go](/src/github.com/terminalstatic/go-xsd-validate/libxml2.go) [validate_xsd.go](/src/github.com/terminalstatic/go-xsd-validate/validate_xsd.go) 





## <a name="Cleanup">func</a> [Cleanup](/src/target/validate_xsd.go?s=2337:2351#L81)
``` go
func Cleanup()
```
Cleans up the libxml2 parser, use this when application ends or libxml2 is not needed anymore.



## <a name="Init">func</a> [Init](/src/target/validate_xsd.go?s=1755:1772#L60)
``` go
func Init() error
```
Initializes libxml2, suggested for multithreading, see <a href="http://xmlsoft.org/threads.html">http://xmlsoft.org/threads.html</a>.



## <a name="InitWithGc">func</a> [InitWithGc](/src/target/validate_xsd.go?s=2143:2175#L74)
``` go
func InitWithGc(d time.Duration)
```
Initializes lbxml2 with a goroutine which trims memory and runs the go gc every d duration.
Not required but can help to keep the memory footprint at bay when doing tons of validations.




## <a name="CommonError">type</a> [CommonError](/src/target/errors.go?s=83:126#L4)
``` go
type CommonError struct {
    Message string
}
```
Common error for default String and Error implementations.










### <a name="CommonError.Error">func</a> (CommonError) [Error](/src/target/errors.go?s=265:300#L14)
``` go
func (e CommonError) Error() string
```
Implementation of Error Interface




### <a name="CommonError.String">func</a> (CommonError) [String](/src/target/errors.go?s=168:204#L9)
``` go
func (e CommonError) String() string
```
Implementation of Stringer Interface




## <a name="Libxml2Error">type</a> [Libxml2Error](/src/target/errors.go?s=375:416#L19)
``` go
type Libxml2Error struct {
    CommonError
}
```
Returned when initialization problems occured.










## <a name="Options">type</a> [Options](/src/target/validate_xsd.go?s=1239:1257#L44)
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










## <a name="ValidationError">type</a> [ValidationError](/src/target/errors.go?s=689:796#L34)
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










### <a name="ValidationError.Error">func</a> (ValidationError) [Error](/src/target/errors.go?s=943:983#L48)
``` go
func (ve ValidationError) Error() string
```
Implementation of Error interface.




### <a name="ValidationError.String">func</a> (ValidationError) [String](/src/target/errors.go?s=839:880#L43)
``` go
func (ve ValidationError) String() string
```
Implementation of Stringer interface.




## <a name="XmlHandler">type</a> [XmlHandler](/src/target/libxml2.go?s=7218:7264#L309)
``` go
type XmlHandler struct {
    // contains filtered or unexported fields
}
```
Handles xml parsing, wraps a pointer to libxml2's xmlDocPtr.







### <a name="NewXmlHandlerMem">func</a> [NewXmlHandlerMem](/src/target/validate_xsd.go?s=2679:2752#L95)
``` go
func NewXmlHandlerMem(inXml []byte, options Options) (*XmlHandler, error)
```
Initialize the xml handler struct.
Always use the Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XmlHandler.Free">func</a> (\*XmlHandler) [Free](/src/target/validate_xsd.go?s=4262:4298#L139)
``` go
func (xmlHandler *XmlHandler) Free()
```
Frees the xml docPtr, call this when this handler is not needed anymore.




## <a name="XmlParserError">type</a> [XmlParserError](/src/target/errors.go?s=465:508#L24)
``` go
type XmlParserError struct {
    CommonError
}
```
Returned when xml parsing caused a problem.










## <a name="XsdHandler">type</a> [XsdHandler](/src/target/libxml2.go?s=7100:7152#L304)
``` go
type XsdHandler struct {
    // contains filtered or unexported fields
}
```
Handles schema parsing and validation, wraps a pointer to libxml2's xmlSchemaPtr.







### <a name="NewXsdHandlerUrl">func</a> [NewXsdHandlerUrl](/src/target/validate_xsd.go?s=3119:3190#L107)
``` go
func NewXsdHandlerUrl(url string, options Options) (*XsdHandler, error)
```
Initialize the xml handler struct.
Always use Free() method when done using this handler or memory will be leaking.
The go garbage collector will not collect the allocated resources.





### <a name="XsdHandler.Free">func</a> (\*XsdHandler) [Free](/src/target/validate_xsd.go?s=4117:4153#L134)
``` go
func (xsdHandler *XsdHandler) Free()
```
Frees the schemaPtr, call this when this handler is not needed anymore.




### <a name="XsdHandler.Validate">func</a> (\*XsdHandler) [Validate](/src/target/validate_xsd.go?s=3543:3628#L117)
``` go
func (xsdHandler *XsdHandler) Validate(xmlHandler *XmlHandler, options Options) error
```
The validation method validates an xmlHandler against an xsdHandler and returns the libxml2 validation error text.
Both xmlHandler and xsdHandler have to be created first.




## <a name="XsdParserError">type</a> [XsdParserError](/src/target/errors.go?s=557:600#L29)
``` go
type XsdParserError struct {
    CommonError
}
```
Returned when xsd parsing caused a problem.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
