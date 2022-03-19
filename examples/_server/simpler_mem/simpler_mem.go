// A simpler standalone example for xsd validation and http
package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/terminalstatic/go-xsd-validate"
)

//go:embed address.xsd
var xsd []byte

var xsdHandler *xsdvalidate.XsdHandler

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("content-type", "application/xml; charset=utf-8")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%s<error><![CDATA[%s]]></error>", xml.Header, err))
		return
	}

	err = xsdHandler.ValidateMem(body, xsdvalidate.ParsErrVerbose)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%s<error><![CDATA[%s]]></error>", xml.Header, err))
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("%s<no-error>No errors</no-error>", xml.Header))
}

func main() {
	addr := ":9999"
	xsdvalidate.Init()
	defer xsdvalidate.Cleanup()
	var err error
	xsdHandler, err = xsdvalidate.NewXsdHandlerMem(xsd, xsdvalidate.ParsErrDefault)
	defer xsdHandler.Free()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/xsd", handler)
	fmt.Printf("Starting http server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
