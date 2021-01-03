package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/signintech/gopdf"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"syscall/js"
)

//go:embed "font.ttf"
var font []byte

//go:embed "czech.png"
var flag []byte

var flags = map[string]string{
	"France": "https://upload.wikimedia.org/wikipedia/en/thumb/c/c3/Flag_of_France.svg/800px-Flag_of_France.svg.png",
	"Greece": "https://upload.wikimedia.org/wikipedia/commons/thumb/5/5c/Flag_of_Greece.svg/600px-Flag_of_Greece.svg.png",
	"Italy":  "https://upload.wikimedia.org/wikipedia/commons/thumb/0/03/Flag_of_Italy.svg/800px-Flag_of_Italy.svg.png",
	"Cyprus": "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d4/Flag_of_Cyprus.svg/800px-Flag_of_Cyprus.svg.png",
	"Spain":  "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9a/Flag_of_Spain.svg/750px-Flag_of_Spain.svg.png",
}

var flagData map[string][]byte

// downloadFlag
func downloadFlag(country, flag string) {
	fmt.Printf("downloading flag for %s\n", country)

	response, err := http.DefaultClient.Get(flag)
	if err != nil {
		fmt.Printf("ERROR(%s)", err.Error())
		return
	}

	defer func() {
		_ = response.Body.Close()
	}()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("ERROR(%s)", err.Error())
		return
	}

	flagData[country] = data
}

// generatePDF
func generatePDF() js.Func {
	var wg sync.WaitGroup

	for country, flag := range flags {
		wg.Add(1)
		go func(cc, ff string, lock *sync.WaitGroup) {
			defer lock.Done()
			downloadFlag(cc, ff)
		}(country, flag, &wg)
	}

	wg.Wait()

	fmt.Print("All flags downloaded\n")

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("generating pdf\n")

		var err error

		pdf := gopdf.GoPdf{}
		defer func() {
			_ = pdf.Close()
		}()

		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
		pdf.AddPage()

		fmt.Printf("adding ttf\n")

		fontReader := strings.NewReader(string(font))

		if err = pdf.AddTTFFontByReader("opensans", fontReader); err != nil {
			fmt.Printf("error adding ttf - %s", err.Error())
			return err
		}

		fmt.Printf("setting font\n")

		if err = pdf.SetFont("opensans", "", 16); err != nil {
			fmt.Printf("error setting font - %s", err.Error())
			return err
		}

		imageHolder, err := gopdf.ImageHolderByBytes(flag)

		if err != nil {
			fmt.Printf("error creating image holder by bytes - %s", err.Error())
			return err
		}

		err = pdf.ImageByHolder(imageHolder, 50, 100, nil)

		if err != nil {
			fmt.Printf("error loading image - %s", err.Error())
			return err
		}

		pdf.SetX(50)
		pdf.SetY(50)

		fmt.Printf("arguments length %d\n", len(args))
		fmt.Printf("arguments %+v\n", args)

		if len(args) == 0 {
			if err := pdf.Cell(nil, "Sample text"); err != nil {
				return err
			}
		} else {
			if err := pdf.Cell(nil, args[0].String()); err != nil {
				return err
			}
		}

		for country, data := range flagData {
			fmt.Printf("rendering flag for %s\n", country)

			pdf.AddPage()
			pdf.SetX(50)
			pdf.SetY(30)

			_ = pdf.Cell(nil, country)

			// draw flag
			imageHolder, err := gopdf.ImageHolderByBytes(data)
			if err != nil {
				fmt.Printf("ERROR(%s)\n", err.Error())
				break
			}

			err = pdf.ImageByHolder(imageHolder, 50, 70, nil)
			if err != nil {
				fmt.Printf("ERROR(%s)\n", err.Error())
				break
			}
		}

		return base64.StdEncoding.EncodeToString(pdf.GetBytesPdf())
	})
}

func main() {
	flagData = make(map[string][]byte)
	js.Global().Set("generatePDF", generatePDF())
	select {}
}
