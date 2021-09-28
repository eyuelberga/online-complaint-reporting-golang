package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		log.Println("DONE!")
		return
	}

	log.Printf("File upload in progress: %d\n", pr.BytesRead)
}

func upload(w http.ResponseWriter, r *http.Request) (string, bool) {
	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("[Handler:Upload]  %s", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return "", false
	}

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		if fileHeader.Size > 1024*1024 {
			http.Error(w, fmt.Sprintf("The uploaded file is too big: %s. Please upload a PDF less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return "", true
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", true
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", true
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "application/pdf" {
			http.Error(w, "The provided file format is not allowed. Please upload a PDF file", http.StatusBadRequest)
			return "", true
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", true
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", true
		}
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		f, err := os.Create(fmt.Sprintf("./uploads/%s", filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return "", true
		}

		defer f.Close()

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(f, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return "", true
		}
		log.Println("Upload successful")
		return filename, false
	}
	return "", false

}

// func upload(w http.ResponseWriter, r *http.Request) (string, bool) {
// 	if err := r.ParseMultipartForm(32 << 20); err != nil {
// 		log.Printf("[Handler:Upload]  %s", err)
// 		return "", false
// 	}
// 	file, _, err := r.FormFile("file")
// 	if err != nil {

// 		log.Printf("[Handler:Upload]  %s", err)
// 		return "", false
// 	}
// 	defer file.Close()
// 	//docs tell that it take only first 512 bytes into consideration
// 	// buff := make([]byte, 512)
// 	// if _, err = file.Read(buff); err != nil {
// 	// 	log.Printf("[Handler:Upload]  %s", err)
// 	// 	return "", false
// 	// }
// 	// if http.DetectContentType(buff) != "application/pdf" {
// 	// 	http.Error(w, "File not PDF", http.StatusBadRequest)
// 	// 	return "", true
// 	// }
// 	// upload file to directory
// 	filename := uuid.NewV4().String() + ".pdf"
// 	f, err := os.OpenFile("uploads/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
// 	if err != nil {
// 		log.Printf("[Handler:Upload]  %s", err)
// 		http.Error(w, "Could not complete Upload", http.StatusInternalServerError)
// 		return "", true
// 	}
// 	defer f.Close()
// 	io.Copy(f, file)

// 	return filename, false
// }

func downloadHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	fileBytes, err := ioutil.ReadFile("uploads/" + vars["file"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(fileBytes)

}
