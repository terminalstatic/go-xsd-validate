# xsdvalidate
[![GoDoc](https://godoc.org/github.com/terminalstatic/go-xsd-validate?status.svg)](https://godoc.org/github.com/terminalstatic/go-xsd-validate)

The goal of this package is to preload xsd files and use their in-memory representation to validate xml documents in a concurrent environment, eg. the post bodys of xml service endpoints and hand through libxml2 error messages. Similar packages I found on github either didn't provide error details or got stuck under load.

# Api Reference
[https://godoc.org/github.com/terminalstatic/go-xsd-validate](https://godoc.org/github.com/terminalstatic/go-xsd-validate)

# Install
libxml2-dev is needed, below an example how to install the latest sources as at the time of writing (Ubuntu, change prefix according to where libs and include files are located):                                                                                               

	curl -sL ftp://xmlsoft.org/libxml2/libxml2-2.9.7.tar.gz| tar -xzf -
	cd ./libxml2-2.9.7/
	./configure --prefix=/usr  --enable-static --with-threads --with-history
	make
	sudo make install
	
Go get the package:

	go get github.com/terminalstatic/go-xsd-validate
	
# Examples
Check [this](./examples/_server/simple/simple.go) for a simple example and [that](./examples/_server/simpler/simpler.go) for an even simpler one.

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

# Licence
[MIT](./LICENSE)
