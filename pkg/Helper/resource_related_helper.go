package Helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/samuael/Project/Weg/internal/pkg/entity"
	etc "github.com/samuael/Project/Weg/pkg/ethiopianCalendar"
)

// EncodeToHLS function
func EncodeToHLS(fileDirectory string, destination string, randomDirectory string) (string, string, error) {
	var cmd *exec.Cmd
	dir, _ := os.Getwd()
	thePath := strings.TrimSuffix(dir, "/cmd/DTIS")                // path to Project
	theFileDirectory := strings.TrimPrefix(fileDirectory, "../..") // path to the File to Be COnverted from the Project
	os.Mkdir(thePath+strings.TrimPrefix(destination, "../..")+randomDirectory, 0700|os.ModeSticky)
	destinationFileDirectory := thePath + strings.TrimPrefix(destination, "../..") + randomDirectory
	destinationFileDirectoryLong := thePath + strings.TrimPrefix(destination, "../..") + "/" + randomDirectory + "/" + "index.m3u8"
	// fmt.Println(destinationFileDirectory)
	if runtime.GOOS == "windows" {
		cmd = exec.Command(
			thePath+"/internal/app/ffmpeg_windows.exe", "-i", thePath+theFileDirectory, "-profile:v", "baseline", "-level", "3.0",
			"-s", "640x360", "-start_number", "0", "-hls_time", "15",
			"-hls_list_size", "0", "-f", "hls", destinationFileDirectoryLong,
		)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command(
			thePath+"/internal/app/ffmpeg", "-i", thePath+theFileDirectory, "-profile:v", "baseline", "-level", "3.0",
			"-s", "640x360", "-start_number", "0", "-hls_time", "15",
			"-hls_list_size", "0", "-f", "hls", destinationFileDirectoryLong,
		)
	}
	cmd.Dir = dir
	erra := cmd.Start()
	if erra != nil {
		_ = os.RemoveAll(destination + randomDirectory)
		return "", erra.Error(), erra
	}
	cutted := strings.TrimSuffix(dir, "cmd/DTIS")
	meinRoot := strings.TrimSuffix(strings.TrimPrefix(destinationFileDirectory, cutted), ".m3u8") + "/stream/"
	// fmt.Println("Cutted ", cutted, "  Main : ", meinRoot)
	returnRoot := strings.TrimPrefix(meinRoot, "web/templates/Source/Resources")
	return strings.TrimPrefix(returnRoot, "/"), "Succesfully Created ", erra
}

// GetHLSFolderNumberName  function
// Example The Path is Like this  -->  "/media/1593087506/stream/"
func GetHLSFolderNumberName(value string) string {
	strips := strings.Split(value, "/")
	return strips[2]
}

// GetFirstFrameOfVideo  funection
func GetFirstFrameOfVideo(sourceDirectory string) (newImageDirectory string) {
	width := 640
	height := 360
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("../../internal/app/ffmpeg_windows.exe", "-i", sourceDirectory, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("../../internal/app/ffmpeg", "-i", sourceDirectory, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	}

	// cmd := exec.Command("ffmpeg", "-i", sourceDirectory, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	if cmd.Run() != nil {
		fmt.Println("could not generate frame")
		return ""
	}
	// Do something with buffer, which contains a JPEG image
	OutPutDirectory := "Source/Resources/FirstFrames/" + strconv.Itoa(int(etc.NewDate(0).Unix)) + ".jpg"
	era := ioutil.WriteFile(entity.PathToTemplates+OutPutDirectory, buffer.Bytes(), 0700)
	if era != nil {
		fmt.Println(era.Error())
		return ""
	}
	// fmt.Println("The Finale Image Directory Is : ", OutPutDirectory)
	return OutPutDirectory
}

// IsImage boolean returning function
// // taking the the File name as an Input
// func IsImage(filename string) bool {
// 	extension := GetExtension(filename)
// 	for _, value := range entity.PICTURES {
// 		if extension == value {
// 			return true
// 		}
// 	}
// 	return false
// }
