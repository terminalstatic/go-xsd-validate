package libxml2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/terminalstatic/go-xsd-validate/libxml2"
)

func TestInit(t *testing.T) {
	libxml2 := libxml2.New()
	libxml2.Init(2)
	time.Sleep(10 * time.Second)
	libxml2.Reset(5)
	time.Sleep(15 * time.Second)
	libxml2.Reset(1)
	time.Sleep(3 * time.Second)
	libxml2.Shutdown()
	time.Sleep(3 * time.Second)
	fmt.Println("Done")
}
