package main

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"io"
	"log"
	"net/http"
	"os"
)

/*func Disembed(src string) ([]byte, error) {
	if src == "astilectron.zip" {
		return nil, nil
	} else if src == "electron.zip" {
		return nil, nil
	} else {
		return nil, errors.New("Wrong argument!")
	}
}*/


type W struct {
	io.Writer
}

func (w *W) Write(p []byte) (n int, err error) {
	for _, b := range p {
		_, err = fmt.Fprintf(w.Writer, "\\x%02x", b)
	}
	n = len(p)
	return
}

func main() {
	for _, o := range astilectron.ValidOSes() {
		for _, a := range []string{"amd64", "386"} {
			if o == "darwin" && a == "386" {
				continue
			}
			file, err := os.Create(fmt.Sprintf("astilectron_vendor_%s_%s.go", o, a))
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString("// Generated by go-astilectron-bindata\n\n")
			file.WriteString("// +build " + o + "\n")
			file.WriteString("// +build " + a + "\n")
			file.WriteString(`
package main
import "errors"
`)
			file.WriteString(`func Disembed(src string) ([]byte, error) {
	if src == "astilectron.zip" {
		return []byte("`)
			resp, err := http.Get(astilectron.AstilectronDownloadSrc())
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(&W{Writer: file}, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString(`"), nil
	} else if src == "electron.zip" {
		return []byte("`)
			resp, err = http.Get(astilectron.ElectronDownloadSrc(o, a))
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(&W{Writer: file}, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString(`"), nil
	} else {
		return nil, errors.New("Wrong argument!")
	}
}`)
		}
	}
}
