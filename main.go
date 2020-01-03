package main

 

import (

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

               param := req.FormValue("q")

               matches, err := filepath.Glob(`d:\tymczasowy\2020\01\2020-01-03\fastExportLogFile_1578044073686_pid_142_iid_772_lid_51238.txt`)

 

               if err != nil {

                              io.WriteString(w, `<h1>Nie znaleziono pliku</h1>`)

               }

 

               header := req.Header

 

               strHTML := ""

               for k, v := range header {

                              strHTML += k + "<ul>"

                              for _, e := range v {

                                            strHTML += "<li>" + e + "</li>"

                              }

                              strHTML += "</ul>"

               }

 

               io.WriteString(w, `<h1>strona glowna</h1>`)

               io.WriteString(w, param)

               io.WriteString(w, strHTML)

 

               if len(matches) != 0 {

                              io.WriteString(w, "<h1>"+matches[0]+"</h1>")

               }

}

 

func startWebServer() {

               //  http.Handle("/", http.FileServer(http.Dir(".")))

               http.HandleFunc("/tw", runScriptHandler)

               http.HandleFunc("/", index)

               http.HandleFunc("/etllogs", getEtlLogFiles)

               http.Handle("/favicon.ico", http.NotFoundHandler())

               http.ListenAndServe(":8000", nil)

}

 

func main() {

               startWebServer()

               // runPythonScript("transakcje_wzajemne")

}
