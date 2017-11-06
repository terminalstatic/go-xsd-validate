// +build memtest

package xsdvalidate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const iterations = 5000000
const maxGoroutines = 100

func TestMemParseXsd(t *testing.T) {
	fmt.Println("Now Running TestMemParseXsd")
	InitWithGc(time.Duration(30) * time.Second)

	defer Cleanup()

	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < iterations; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			//handler, err := NewXsdHandlerUrl("examples/test1_fail.xsd", ParsErrVerbose)
			handler, err := NewXsdHandlerUrl("examples/test1_pass.xsd", ParsErrDefault)
			if err != nil {
				//fmt.Println(err)
			}
			handler.Free()
			<-guard
			wg.Done()
		}()
	}
	wg.Wait()
}
func TestMemParseXml(t *testing.T) {
	fmt.Println("Now Running TestMemParseXml")
	InitWithGc(time.Duration(30) * time.Second)

	defer Cleanup()

	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	xmlfile := "examples/test1_fail1.xml"
	//xmlfile := "examples/test1_pass.xml"

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

	for i := 0; i < iterations; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(inXml []byte) {
			xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
			//xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrVerbose)
			if err != nil {
				//fmt.Println(err)
			}
			xmlhandler.Free()
			<-guard
			wg.Done()
		}(inXml)
	}
	wg.Wait()
}
func TestMemValidate(t *testing.T) {
	fmt.Println("Now Running TestMemValidate")
	InitWithGc(time.Duration(30) * time.Second)

	defer Cleanup()

	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	xmlfile := "examples/test1_fail2.xml"
	//xmlfile := "examples/test1_pass.xml"

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

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_pass.xsd", ParsErrDefault)
	if err != nil {
		panic(err)
	}

	defer xsdhandler.Free()

	for i := 0; i < iterations; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func(inXml []byte) {
			xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
			if err != nil {
				//panic(err)
			}
			err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
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
func TestMemAltValidate(t *testing.T) {
	fmt.Println("Now Running TestMemAltValidate")
	InitWithGc(time.Duration(30) * time.Second)

	//fmt.Println(syscall.Gettid())
	defer Cleanup()

	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	xmlfile1 := "examples/test1_fail2.xml"

	fxml1, err := os.Open(xmlfile1)
	if err != nil {
		log.Printf("failed to open file: %s", err)
		return
	}
	defer fxml1.Close()

	inXml1, err := ioutil.ReadAll(fxml1)
	if err != nil {
		log.Printf("failed to read file: %s", err)
		return
	}

	xmlfile2 := "examples/test1_fail3.xml"

	fxml2, err := os.Open(xmlfile2)
	if err != nil {
		log.Printf("failed to open file: %s", err)
		return
	}
	defer fxml2.Close()

	inXml2, err := ioutil.ReadAll(fxml2)
	if err != nil {
		log.Printf("failed to read file: %s", err)
		return
	}

	xsdhandler, err := NewXsdHandlerUrl("examples/test1_pass.xsd", ParsErrDefault)
	if err != nil {
		panic(err)
	}

	defer xsdhandler.Free()

	for i := 0; i < iterations; i++ {
		guard <- struct{}{}
		wg.Add(1)

		var inXml []byte
		if i%2 == 0 {
			inXml = inXml1
		} else {
			inXml = inXml2
		}
		go func(inXml []byte, i int) {
			xmlhandler, err := NewXmlHandlerMem(inXml, ParsErrDefault)
			if err != nil {
				//panic(err)
			}
			//start := time.Now()
			err = xsdhandler.Validate(xmlhandler, ValidErrDefault)
			if err != nil {
				if i%2 == 1 {
					if !strings.Contains(err.Error(), "Element 'name1'") {
						panic(err)
					}
				} else {
					if !strings.Contains(err.Error(), "Element 'shipto'") {
						panic(err)
					}
				}
				//log.Print(err)
			}
			//elapsed := time.Since(start)
			//log.Printf("Validation took %s", elapsed)
			xmlhandler.Free()
			<-guard
			wg.Done()
		}(inXml, i)
	}
	wg.Wait()
}
