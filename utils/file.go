package utils

import (
	"fmt"
	"image"
	"math"
	"mime/multipart"
	"os"

	_ "image/png"

	"github.com/disintegration/imaging"
)

var DefaultExtension = []string{
	".jpg",
	".jpeg",
	".png",
	".svg",
}

func ValidateExtension(extension string, extLists []string) (string, bool) {
	var accetExtension string
	for _, ext := range extLists {
		if ext == extension {
			return "", true
		}
		accetExtension += " " + ext + ","
	}
	return "Image must be accept only" + accetExtension, false
}

func ValidateFileSize(size int64, maxSize int64) (string, bool) {
	calculateMaxSize := float64(size / 1000 / 1000)
	if math.Ceil(calculateMaxSize) > float64(maxSize) {
		return fmt.Sprintf("Image must be less than or equal %d MB only", maxSize), false
	}
	return "", true
}

func ImageSaver(imageFile *multipart.FileHeader, path, imageName string, size map[string]interface{}) error {
	file, err := imageFile.Open()
	if err != nil {
		return err
	}
	image, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	baseDir, err := generateBaseDirSubImage(path)
	if err != nil {
		return err
	}

	err = createDirectoryIfNotExists(baseDir + "/original/")
	if err != nil {
		return err
	}

	err = imaging.Save(image, baseDir+"/original/"+imageName)
	if err != nil {
		return err
	}

	err = resizeImageAndSave(image, size["l"].(int), baseDir, imageName, "large")
	if err != nil {
		return err
	}

	err = resizeImageAndSave(image, size["m"].(int), baseDir, imageName, "medium")
	if err != nil {
		return err
	}

	err = resizeImageAndSave(image, size["s"].(int), baseDir, imageName, "small")
	if err != nil {
		return err
	}

	return nil
}

func RemoveImageAllResolution(path, imageName string) error {
	dir := []string{
		"original",
		"large",
		"medium",
		"small",
	}
	for _, subDir := range dir {
		baseDir, err := generateBaseDirSubImage(path)
		if err != nil {
			return err
		}
		baseDir = baseDir + "/" + subDir + "/" + imageName
		if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
			os.Remove(baseDir)
		}
	}
	return nil
}

func resizeImageAndSave(image image.Image, width int, dir, imageName, resolution string) error {
	baseImageDir := dir + "/" + resolution + "/"
	err := createDirectoryIfNotExists(baseImageDir)
	if err != nil {
		return err
	}

	fileName := baseImageDir + imageName
	dstImage := imaging.Resize(image, width, 0, imaging.Lanczos)
	err = imaging.Save(dstImage, fileName)
	if err != nil {
		removePreviousImageBeforeError(dir, resolution, imageName)
		return err
	}
	return nil
}

func removePreviousImageBeforeError(dir, resolution, imageName string) {
	path := map[string][]string{
		"large":  {"original"},
		"medium": {"original", "large"},
		"small":  {"original", "large", "medium"},
	}
	for _, result := range path[resolution] {
		problemDir := dir + "/" + result + "/" + imageName
		if _, err := os.Stat(problemDir); !os.IsNotExist(err) {
			os.Remove(problemDir)
		}
	}
}

func createDirectoryIfNotExists(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		errMkDir := os.Mkdir(dir, 0755)
		if errMkDir != nil {
			return errMkDir
		}
	}

	return nil
}

func generateBaseDirSubImage(path string) (string, error) {
	baseDir, err := GetBaseDirectory()
	return baseDir + "/public/images/" + path, err
}
