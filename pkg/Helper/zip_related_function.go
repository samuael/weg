// Package Helper for Handling Zip Related function
package Helper

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

// AddFileToZip  function
func AddFileToZip(zipWriter *zip.Writer, InputFileName, OutputFileName string) error {
	// InputFileName = InputFileName + ".jpg
	if !strings.HasSuffix(InputFileName, ".jpg") {
		InputFileName = InputFileName + ".jpg"
	}
	fileToZip, err := os.Open(InputFileName)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = OutputFileName
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// GetExtension function to return the extension of the File Input FileName
func GetExtension(Filename string) string {
	fileSlice := strings.Split(Filename, ".")
	if len(fileSlice) >= 1 {
		return fileSlice[len(fileSlice)-1]
	}
	return ""
}

// JPEGFileName function
func JPEGFileName(filename string) string {
	filenameSlice := strings.Split(filename, ".")
	if len(filenameSlice) > 1 {
		filenames := strings.Join(filenameSlice[:len(filenameSlice)-1], "")
		filenames += ".jpg"
		return filenames
	}
	return filename + ".jpg"
}
