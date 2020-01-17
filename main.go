package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
)

func runPythonScript(scriptName string) {
	cmd := exec.Command("/usr/local/bin/python3", "../python/"+scriptName+".py")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "python ", "..\\python\\"+scriptName+".py")
	}

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

}

func runScriptHandler(w http.ResponseWriter, req *http.Request) {

	runPythonScript("transakcje_wzajemne")
	ff, _ := os.Open("todownload\\wynik.xlsx")
	FileHeader := make([]byte, 512)
	ff.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)
	//Get the file size
	FileStat, _ := ff.Stat()                           //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
	w.Header().Set("Content-Disposition", "attachment; filename=wynik.xlsx")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)
	ff.Seek(0, 0)
	io.Copy(w, ff)
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<h1>strona glowna</h1>`)
}

func getEtlLogFiles(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	itemID := req.FormValue("itemid")
	logID := req.FormValue("logid")
	// matches, err := filepath.Glob(`d:\tymczasowy\2020\01\2020-01-03\*` + itemID + `_lid_` + logID + `.txt`)
	matches, err := filepath.Glob(`/home/oracle/Octago_app/commands_log/*` + itemID + `_lid_` + logID + `.txt`)
	if err != nil {
		io.WriteString(w, `<h1>Nie znaleziono pliku</h1>`)
		return
	}
	// header := req.Header
	strHTML := ""
	/*
		for k, v := range header {
			strHTML += k + "<ul>"
			for _, e := range v {
				strHTML += "<li>" + e + "</li>"
			}
			strHTML += "</ul>"
		}
	*/

	if len(matches) != 0 {
		for _, f := range matches {
			strHTML += "<h3>Plik: " + filepath.Base(f) + "</h3>"
			strHTML += "<div>"
			file, err := os.Open(f)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			// b, _ := ioutil.ReadAll(file)
			// s := fmt.Sprintf("%s", b)

			sc := bufio.NewScanner(file)
			for sc.Scan() {
				strHTML += "<br/>" + sc.Text()
				// io.WriteString(w, "<br/>"+sc.Text())
			}
			strHTML += "<br/><br/><br/>"
			strHTML += "</div>"

			if err := sc.Err(); err != nil {
				log.Fatalf("scan file error: %v", err)
				return
			}
		}

		io.WriteString(w, strHTML)

		// io.WriteString(w, s)
	} else {
		io.WriteString(w, `<h1>Nie znaleziono pliku</h1>`)
	}

}

func getEtlAppLog(w http.ResponseWriter, req *http.Request) {
	ff, err := os.Open("/home/oracle/Octago_app/application_log/Octago.log")
	// ff, err := os.Open(`d:\tymczasowy\2020\01\2020-01-13\Octago.log`)
	if err != nil {
		io.WriteString(w, `<h1>Nie znaleziono pliku</h1>`)
		return
	}
	defer ff.Close()
	w.Header().Set("Content-Disposition", "attachment; filename=log_pijawki.txt")
	w.Header().Set("Content-Security-Policy", "frame-src: 'https://mis.centrala.bzwbk/'")
	w.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", req.Header.Get("Content-Length"))
	io.Copy(w, ff)

}

func startWebServer() {
	//  http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/tw", runScriptHandler)
	http.HandleFunc("/", index)
	http.HandleFunc("/etllogs", getEtlLogFiles)
	http.HandleFunc("/etlapplog", getEtlAppLog)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8000", nil)
}

func main() {
	startWebServer()
	// runPythonScript("transakcje_wzajemne")

}
