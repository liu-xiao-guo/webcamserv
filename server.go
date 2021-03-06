package main

import (
	"net/http"
	"fmt"
	"os/exec"
	"log"
	"io"
	"os"
	"io/ioutil"	
	)
	
var quality = "40"
var resolution = "320x240"

func main() {	
	mux := http.NewServeMux()
	
	path := os.Getenv("SNAP_DATA");
	fmt.Println("path: ", path);
	
	mux.Handle("/", http.FileServer(http.Dir(path)))
	mux.HandleFunc("/takepic", takePicture)
	mux.HandleFunc("/getpic", getPicture)
	mux.HandleFunc("/hello", handleHello)
	mux.HandleFunc("/setquality", handleSetQuality)	
	mux.HandleFunc("/setresolution", handleResolution)	

	log.Println("Starting webserver on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal("http.ListendAndServer() failed with %s\n", err)
	}
}

func getPicture(w http.ResponseWriter, r *http.Request) {
	log.Println("entering getPicture")
	// cmd := exec.Command("fswebcam", "-")
	cmd := exec.Command("fswebcam", "--jpeg", quality, "-p", "YUYV", "-r", resolution, "-")
	// cmd := exec.Command("fswebcam", "--jpeg", "40", "-")
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
			
	bytes, err := ioutil.ReadAll(stdout);
		
	if  err != nil {
		log.Fatal(err)
	}

			
	if err := cmd.Wait(); err != nil {
		fmt.Println("it comes here", err.Error())
		log.Println("it is error!")
		log.Fatal(err)
	}	
		
	// After all of the data has been read, now, let's print out
	// fmt.Println("Number of bytes: ", len(bytes))
		
	// fmt.Fprintf(w, string(bytes))		
	w.Header().Set("Content-Length", fmt.Sprint(len(bytes)))
    w.Header().Set("Content-Type", "image/jpeg")
	
  	if _, err := w.Write(bytes); err != nil {
       fmt.Println("unable to write image.")
    }
    
}


func handleHello(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/" {
//		http.NotFound(w, r)
//		return
//	}

	fmt.Fprintf(w, "Hello World, xiaoguo\n")
}

func takePicture(res http.ResponseWriter, req *http.Request) {
	fmt.Println("GET params:", req.URL.Query());

	res.Header().Set("Content-Type", "text/plain")
	
	name := req.URL.Query().Get("name")
	fmt.Println("name: ", name);
		
	if ( name == "") {
		name = "shot.jpeg"
	}
		
	pwd := os.Getenv("pwd");
	fmt.Println("pwd: ", pwd);
	
	fmt.Println("Going to launch program")
	
	app_path := os.Getenv("SNAP")
	app_path += "/bin/capture"
	
	fmt.Println("app_path: ", app_path)
			
	cmd := exec.Command(app_path, name)
	err1 := cmd.Run()

	if err1 != nil {
		fmt.Println("error:", err1.Error())
		log.Fatal(err1)
	}
	
	io.WriteString(res, name );
				
//	res.Write([]byte("Picture is taken.\n"))
}	

func handleSetQuality(res http.ResponseWriter, req *http.Request) {
	quality = req.URL.Query().Get("quality")
	fmt.Println("Going to set the Quality: ", quality)
	fmt.Fprintf(res, "The quality is set to: "+quality)
}

func handleResolution(res http.ResponseWriter, req *http.Request) {
	resolution = req.URL.Query().Get("resolution")
	fmt.Println("Going to set the resolution: ", resolution)
	fmt.Fprintf(res, "The resolution is set to: " + resolution)
}
