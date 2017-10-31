// +build memtest

package xsdv_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/terminalstatic/go-xsd-validate/libxml2"
	"github.com/terminalstatic/go-xsd-validate/xsdv"
)

func TestMemParseXsd(t *testing.T) {

	l2 := libxml2.NewDefault()
	defer l2.Shutdown()

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < 1000000; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			//handler, err := xsdv.NewXsdHandlerUrl("../examples/test1_fail.xsd", libxml2.ParsErrVerbose)
			handler, err := xsdv.NewXsdHandlerUrl("../examples/test1_pass.xsd", libxml2.ParsErrVerbose)
			if err != nil {
				fmt.Println(err)
			}
			handler.Free()
			<-guard
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestMemParseXml(t *testing.T) {

	l2 := libxml2.NewDefault()
	defer l2.Shutdown()

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	xmlfile := "../examples/test1_fail1.xml"
	//xmlfile := "../examples/test1_pass.xml"

	fxml, err := os.Open(xmlfile)
	if err != nil {
		log.Printf("failed to open file: %s", err)
		return
	}
	defer fxml.Close()

	inXml, err := ioutil.ReadAll(fxml)
	if err != nil {
		log.Printf("failed to read file: %s", err)
		return
	}

	for i := 0; i < 10000000; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(inXml []byte) {
			xmlhandler, err := xsdv.NewXmlHandlerMem(inXml, libxml2.ParsErrDefault)
			//xmlhandler, err := xsdv.NewXmlHandlerMem(inXml, libxml2.ParsErrVerbose)
			if err != nil {
				fmt.Println(err)
			}
			xmlhandler.Free()
			<-guard
			wg.Done()
		}(inXml)
	}
	wg.Wait()
}
func TestMemValidate(t *testing.T) {

	l2 := libxml2.NewDefault()
	defer l2.Shutdown()

	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	xmlfile := "../examples/test1_fail2.xml"
	//xmlfile := "../examples/test1_pass.xml"

	fxml, err := os.Open(xmlfile)
	if err != nil {
		log.Printf("failed to open file: %s", err)
		return
	}
	defer fxml.Close()

	inXml, err := ioutil.ReadAll(fxml)
	if err != nil {
		log.Printf("failed to read file: %s", err)
		return
	}

	xsdhandler, err := xsdv.NewXsdHandlerUrl("../examples/test1_pass.xsd", libxml2.ParsErrDefault)
	if err != nil {
		panic(err)
	}

	defer xsdhandler.Free()

	wg.Add(1)
	for i := 0; i < 100000000; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(inXml []byte) {
			xmlhandler, err := xsdv.NewXmlHandlerMem(inXml, libxml2.ParsErrDefault)
			if err != nil {
				panic(err)
			}
			err = xsdhandler.Validate(xmlhandler, libxml2.ValidErrDefault)
			if err != nil {
				//log.Print(err)
			}
			xmlhandler.Free()
			<-guard
			wg.Done()
		}(inXml)
		/*if i > 0 && (i%1000000) == 0 {
			fmt.Println("Test paused ...\a\n")
			time.Sleep(30 * time.Second)
		}*/
	}
	wg.Wait()
}
