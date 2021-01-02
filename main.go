package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/signintech/gopdf"
	"strings"
	"syscall/js"
)

//go:embed "font.ttf"
var font []byte

//go:embed "czech.png"
var flag []byte

// generatePDF
func generatePDF() js.Func {
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

		pdf.AddPage()

		return base64.StdEncoding.EncodeToString(pdf.GetBytesPdf())
	})
}

func main() {
	js.Global().Set("generatePDF", generatePDF())
	select {}
}
