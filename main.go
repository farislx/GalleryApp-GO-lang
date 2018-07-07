package main

import (
    "encoding/json"
    "io/ioutil"
    "github.com/gorilla/mux"
    "log"
    "fmt"
	"io"
	"net/http"
	"os"
    "strings"
    "path/filepath"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)
type Shot struct {
    Id   int `json:"id,omitempty"`
    Images   *Image `json:"images,omitempty"`
    Title   string `json:"title,omitempty"`
    Description string   `json:"description,omitempty"`
}
type Image struct {
    Normal string `json:"normal,omitempty"`
    Teaser string `json:"teaser,omitempty"`
}

type Photo struct{
    Id int `json:"id,omitempty"`
    Title string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
    Filename string `json:"filename,omitempty"`
    
    
}
var shots []Shot
var photos []Photo

// Download all photos
func DownloadPhotos(w http.ResponseWriter, r *http.Request) {
    
	url := "https://api.dribbble.com/v2/user/shots?access_token=363804be630a4ae4fec2d66c6071c8ea92fd336ca671d8030781f897534d348a"
    
    res, err := http.Get(url)
    if err != nil {
        panic(err.Error())
    }
    
    body, err := ioutil.ReadAll(res.Body)

    if err != nil {
        panic(err.Error())
    }

    json.Unmarshal(body, &shots)
    fmt.Printf("Results: %v\n", shots)
    db, err := sql.Open("sqlite3", "./test.db")
    checkErr(err)

    for _, item := range shots {
            img := item.Images
            downloadFromUrl( img.Normal)
            
            stmt, err := db.Prepare("INSERT INTO imageinfo(title,description,filename) values(?,?,?)")
            checkErr(err)

            res, err := stmt.Exec(item.Title, item.Description, img.Normal)
            checkErr(err)

            affect, err := res.RowsAffected()
            checkErr(err)
            fmt.Println(affect)
            
     }
    db.Close()
}

func GetPhotoList(w http.ResponseWriter, r *http.Request) {

    params := mux.Vars(r)
    db, err := sql.Open("sqlite3", "./test.db")
    checkErr(err)
    fmt.Println("title", params["title"])
    fmt.Println("description", params["description"])

    
    rows, err := db.Query("SELECT * FROM imageinfo WHERE title LIKE '%"+params["title"]+"%' AND description LIKE '%"+params["description"]+"%' " )
    checkErr(err)
    var id int
    var title string
    var description string
    var filename string
    photos=photos[:0]
    for rows.Next() {
        err = rows.Scan(&id, &title, &description, &filename)
        checkErr(err)
        photos=append(photos, Photo{Id: id, Title: title, Description: description, Filename: filename})
    }

    rows.Close() //good habit to close
    db.Close()
    json.NewEncoder(w).Encode(photos)
}
//download a file
func downloadFromUrl(url string) {
	 tokens := strings.Split(url, "/")
     filename := tokens[len(tokens)-1]
     basepath :="photos"
	 fmt.Println("Downloading", url, "to", filename)
     if err := os.MkdirAll(basepath, 0777); err != nil {
         panic("Unable to create directory for file! - " + err.Error())
     }
     imagefile :=filepath.Join(basepath,filename)

     output, err := os.Create(imagefile)
	 if err != nil {
	 	fmt.Println("Error while creating", imagefile, "-", err)
	 	return
	 }
	 defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
    }
    fmt.Println(response, "response.")
    
	defer response.Body.Close()

	 n, err := io.Copy(output, response.Body)
	 if err != nil {
	 	fmt.Println("Error while downloading", url, "-", err)
	 	return
	 }

    fmt.Println(n, "bytes downloaded.")
    

}
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
// main function to boot up everything
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/download", DownloadPhotos).Methods("GET")
    router.HandleFunc("/photolist/{title}/{description}", GetPhotoList).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", router))
}