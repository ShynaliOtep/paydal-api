package filesystem

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const filesystem_host = "localhost:8015"

func UploadFile(filename string, file io.ReadSeeker) (string, error) {
	fmt.Println(filename)
	if filename == "" {
		return filename, errors.New("Filename не передан")
	}

	// Создаем временный файл
	tempFile, err := os.CreateTemp("./uploads", "temp-*.tmp")
	if err != nil {
		return filename, err
	}
	defer tempFile.Close()
	// Копируем данные из тела запроса во временный файл
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return filename, err
	}

	// Получаем путь для сохранения файла
	savePath := fmt.Sprintf("./uploads/%s", filename)

	// Проверяем, существует ли папка для сохранения файла
	dir := filepath.Dir(savePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Если папка не существует, создаем ее
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return filename, err
		}
	}

	// Перемещаем временный файл в указанный путь
	err = os.Rename(tempFile.Name(), savePath)
	if err != nil {
		return filename, err
	}

	// Отправляем ответ об успешной загрузке файла
	return savePath, nil
}

//func UploadFile(w http.ResponseWriter, r *http.Request) {
//	file, handler, err := r.FormFile("file")
//	if err != nil {
//		fmt.Println("Error retrieving file:", err)
//		return
//	}
//	defer file.Close()
//
//	// Создаем временный файл для сохранения загруженного файла
//	tempFile, err := os.Create(filepath.Join("./uploads", handler.Filename))
//	if err != nil {
//		fmt.Println("Error creating file:", err)
//		return
//	}
//	defer tempFile.Close()
//
//	// Копируем содержимое файла во временный файл
//	_, err = io.Copy(tempFile, file)
//	if err != nil {
//		fmt.Println("Error copying file:", err)
//		return
//	}
//
//	fmt.Fprintf(w, "File uploaded successfully")
//}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	err := os.Remove(filepath.Join("./uploads", filename))
	if err != nil {
		fmt.Println("Error deleting file:", err)
		fmt.Fprintf(w, "Error deleting file: %v", err)
		return
	}
	fmt.Fprintf(w, "File deleted successfully")
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	file, err := os.Open(filepath.Join("./uploads", filename))
	if err != nil {
		fmt.Println("Error opening file:", err)
		fmt.Fprintf(w, "Error opening file: %v", err)
		return
	}
	defer file.Close()

	// Копируем содержимое файла в ответ
	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Println("Error copying file to response:", err)
		return
	}
}

func GetExtension(filename string) string {
	extension := filepath.Ext(filename)
	return extension
}

func GetFilenameFromPath(path string) string {
	extension := filepath.Base(path)
	return extension
}

func GenerateFilename() string {
	str := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(1000))
	b := []byte(str)
	return base64.URLEncoding.EncodeToString(b)
}

func GeneratePath(path string, ext string) string {
	imageFilename := fmt.Sprintf("%s/%s%s", path, GenerateFilename(), ext)
	return imageFilename
}

func GerFileUrl(path string) string {
	fmt.Println("PATH", path)
	return fmt.Sprintf("%s:%s%s", os.Getenv("HOST"), os.Getenv("PORT"), path[1:])
}
