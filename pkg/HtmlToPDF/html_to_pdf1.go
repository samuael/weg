package HtmlToPDF

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/Projects/ScientificNRS/pkg/Helper"

// 	"github.com/Projects/ScientificNRS/internal/pkg/entity"
// 	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
// )

// // GetThePdf function
// func GetThePdf(fileDirectory string) string {
// 	pdfg, erra := wkhtmltopdf.NewPDFGenerator()
// 	if erra != nil {
// 		fmt.Println("Error While Generating the Pdf ")
// 		return ""
// 	}
// 	pdfg.Dpi.Set(30)
// 	pdfg.ImageDpi.Set(30)
// 	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
// 	pdfg.Grayscale.Set(false)
// 	dir := strings.Trim(fileDirectory, "../") // the Trimmed Path of the  fiel
// 	pathDirectory, era := os.Getwd()          // this gives us the path to the main Directory
// 	if era != nil {
// 		fmt.Println("error while getting the path to the Main File ")
// 		return ""
// 	}
// 	// Trimming the cmd/DTIS/ suffix from the path
// 	pathDirectory = strings.TrimSuffix(pathDirectory, "cmd/DTIS")
// 	homepath := entity.FileSchema + pathDirectory + dir
// 	page := wkhtmltopdf.NewPage(homepath)
// 	// Set options for this page
// 	page.FooterRight.Set("[page]")
// 	page.FooterFontSize.Set(10)
// 	page.Zoom.Set(0.95)
// 	pdfg.AddPage(page)
// 	pdfCreationError := pdfg.Create()
// 	if pdfCreationError != nil {
// 		fmt.Println("Error while Creating the pdf file ")
// 		return ""
// 	}
// 	PathToPdfs := entity.PathToPdfs
// 	// Generating Random Name to Be Output Name
// 	OutPutname := PathToPdfs + Helper.GenerateRandomString(10, Helper.CHARACTERS) + ".pdf"
// 	writingError := pdfg.WriteFile(OutPutname)
// 	if writingError != nil {
// 		log.Println("Error While Writing the File to The Directory Name ", OutPutname)
// 		return ""
// 	}
// 	return OutPutname
// }

// // GetIDSPdf   function
// func GetIDSPdf(filePaths ...string) string {
// 	pdfg, erra := wkhtmltopdf.NewPDFGenerator()
// 	if erra != nil {
// 		fmt.Println("Error While Generating the Pdf ")
// 		return ""
// 	}
// 	pdfg.Dpi.Set(30)
// 	pdfg.ImageDpi.Set(30)
// 	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
// 	pdfg.Grayscale.Set(false)
// 	// dir := strings.Trim(fileDirectory, "../") // the Trimmed Path of the  fiel
// 	// pathDirectory, era := os.Getwd()          // this gives us the path to the main Directory
// 	// if era != nil {
// 	// fmt.Println("error while getting the path to the Main File ")
// 	// return ""
// 	// }
// 	// Trimming the cmd/DTIS/ suffix from the path
// 	// pathDirectory = strings.TrimSuffix(pathDirectory, "cmd/DTIS")
// 	// homepath := entity.FileSchema + pathDirectory + dir
// 	for _, path := range filePaths {
// 		page := wkhtmltopdf.NewPage(path)
// 		// Set options for this page
// 		page.FooterRight.Set("[page]")
// 		page.FooterFontSize.Set(10)
// 		page.Zoom.Set(1)
// 		pdfg.AddPage(page)
// 	}
// 	pdfCreationError := pdfg.Create()
// 	if pdfCreationError != nil {
// 		fmt.Println("Error while Creating the pdf file ")
// 		return ""
// 	}
// 	PathToPdfs := entity.PathToPdfs
// 	// Generating Random Name to Be Output Name
// 	OutPutname := PathToPdfs + Helper.GenerateRandomString(10, Helper.CHARACTERS) + ".pdf"
// 	writingError := pdfg.WriteFile(OutPutname)
// 	if writingError != nil {
// 		log.Println("Error While Writing the File to The Directory Name ", OutPutname)
// 		return ""
// 	}
// 	return OutPutname
// }

// // GetReciptThePdf function
// func GetReciptThePdf(fileDirectory string) string {
// 	pdfg, erra := wkhtmltopdf.NewPDFGenerator()
// 	if erra != nil {
// 		fmt.Println("Error While Generating the Pdf ")
// 		return ""
// 	}
// 	pdfg.Dpi.Set(30)
// 	pdfg.ImageDpi.Set(30)
// 	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
// 	pdfg.Grayscale.Set(false)
// 	dir := strings.Trim(fileDirectory, "../") // the Trimmed Path of the  fiel
// 	pathDirectory, era := os.Getwd()          // this gives us the path to the main Directory
// 	if era != nil {
// 		fmt.Println("error while getting the path to the Main File ")
// 		return ""
// 	}
// 	// Trimming the cmd/DTIS/ suffix from the path
// 	pathDirectory = strings.TrimSuffix(pathDirectory, "cmd/DTIS")
// 	homepath := entity.FileSchema + pathDirectory + dir
// 	page := wkhtmltopdf.NewPage(homepath)
// 	// Set options for this page
// 	page.FooterRight.Set("[page]")
// 	page.FooterFontSize.Set(10)
// 	page.Zoom.Set(0.25)
// 	pdfg.AddPage(page)
// 	pdfCreationError := pdfg.Create()
// 	if pdfCreationError != nil {
// 		fmt.Println("Error while Creating the pdf file ")
// 		return ""
// 	}
// 	PathToPdfs := entity.PathToPdfs
// 	// Generating Random Name to Be Output Name
// 	OutPutname := PathToPdfs + Helper.GenerateRandomString(10, Helper.CHARACTERS) + ".pdf"
// 	writingError := pdfg.WriteFile(OutPutname)
// 	if writingError != nil {
// 		log.Println("Error While Writing the File to The Directory Name ", OutPutname)
// 		return ""
// 	}
// 	return OutPutname
// }
