package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/dannypark95/ChicagoOnnuri/services"
)

// UploadPDF handles the uploading of a PDF file to AWS S3
func UploadPDF(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// validate token

	// Parse the multipart form data, setting a maximum file size limit
	err := r.ParseMultipartForm(32 << 20) // 32MB max file size
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	//Retrieve the uploaded file from the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the uploaded file has the correct content type (PDF)
	contentType := header.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		http.Error(w, "File must be a PDF", http.StatusBadRequest)
		return
	}

	// Construct the S3 object key (file name) using the original file name
	key := fmt.Sprintf("jubo/%s", header.Filename)

	// Upload the file to S3 and retrieve the URL of the uploaded file
	url, err := services.UploadToS3(file, key, contentType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the live jubo in S3 metadata
	err = services.WriteLiveJuboToS3(header.Filename)
	if err != nil {
		http.Error(w, "Error updating live jubo metadata", http.StatusInternalServerError)
		return
	}

	// Send the URL of the uploaded file as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"url":url})
}

// ListPDFs lists all PDF files in the "jubo" folder in the S3 bucket.
func ListPDFs(w http.ResponseWriter, r *http.Request) {
	// Read the live jubo from S3
	liveJubo, err := services.ReadLiveJuboFromS3()
	if err != nil {
		http.Error(w, "Error reading live jubo from S3", http.StatusInternalServerError)
		return
	}

	// Parse the live jubo from the metadata
	var liveJuboMap map[string]string
	err = json.Unmarshal([]byte(liveJubo), &liveJuboMap)
	if err != nil {
		http.Error(w, "Error parsing live jubo metadata", http.StatusInternalServerError)
		return
	}

	// Get all pdfs
	pdfURLs, err := services.ListPDFs()
	if err != nil {
		http.Error(w, "Error listing PDFs", http.StatusInternalServerError)
		return
	}

	// Prepare the response: a list of all PDFs with the 'live' PDF marked
	pdfList := make([]map[string]interface{}, 0, len(pdfURLs)) // Initialize with length 0
	for _, url := range pdfURLs {
		fileName := path.Base(url) // Extract the file name from the URL

		// Skip if it's a directory
		if fileName == "jubo%2F" {
			continue
		}

		// Remove "jubo%2F" from the filename
		fileName = strings.TrimPrefix(fileName, "jubo%2F")

		pdfList = append(pdfList, map[string]interface{}{
			"url": url,
			"name": fileName,
			"live": fileName == liveJuboMap["liveJubo"],
		})
	}

	// Send the list of PDFs as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pdfList)
}


// SetLiveJubo sets a specific PDF as the live Jubo
func SetLiveJubo(w http.ResponseWriter, r *http.Request) {
	// get filename from the request, you could use either query params or request body
	// for this example, let's assume you're using query params
	fmt.Printf("SetLiveJubo called")

	var request struct {
		Filename string `json:"filename"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
	}

	filename := request.Filename
	if filename == "" {
			http.Error(w, "Filename is missing", http.StatusBadRequest)
			return
	}

	// Call WriteLiveJuboToS3 with the key
	err = services.WriteLiveJuboToS3(filename)
	if err != nil {
		http.Error(w, "Error setting live jubo", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

// DeletePDF deletes a specific PDF file from the S3 bucket
func DeletePDF(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is DELETE
	if r.Method != http.MethodDelete {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
	}

	fmt.Println("Request is Delete")

	var request struct {
		Filename string `json:"filename"`
	}


	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
	}

	filename := request.Filename
	if filename == "" {
			http.Error(w, "Filename is missing", http.StatusBadRequest)
			return
	}

	// Construct the key for the S3 object
	key := fmt.Sprintf("jubo/%s", filename)

	// Delete the file from S3
	err = services.DeleteFromS3(key)
	if err != nil {
		http.Error(w, "Error deleting PDF file", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

// ShowJubo redirects to the live jubo URL on S3
func ShowJubo(w http.ResponseWriter, r *http.Request) {
	log.Printf("/jubo showing jubo")
	// Read the live jubo from S3
	liveJubo, err := services.ReadLiveJuboFromS3()
	if err != nil {
		http.Error(w, "Error reading live jubo from S3", http.StatusInternalServerError)
		return
	}

	// Parse the live jubo from the metadata
	var liveJuboMap map[string]string
	err = json.Unmarshal([]byte(liveJubo), &liveJuboMap)
	if err != nil {
		http.Error(w, "Error parsing live jubo metadata", http.StatusInternalServerError)
		return
	}

	// Construct the URL for the live jubo
	liveJuboURL := services.GetObjectURL("chicagoonnuri", "jubo/" + liveJuboMap["liveJubo"])

	// Issue a redirect to the S3 URL
	http.Redirect(w, r, liveJuboURL, http.StatusFound)
}
