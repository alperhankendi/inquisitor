package server

import (
	"fmt"
	"net/http"
	"io"
	"log"
	"os"
	"encoding/json"
	"time"

	//"github.com/nytimes/gziphandler"

	"github.com/serainville/inquisitor/variables"
	"github.com/serainville/inquisitor/models"
)

type Config struct {
	ip string
	port string
	TLSCertFile string
	TLSKeyFile string
	UseTLS bool
	Database
}
type Database struct {
	host string
	port string
	user string
	password string
}

func Init(
	bindIP string,
	bindPort string,
	TLSCertFile string,
	TLSKeyFile string,
	UseTLS bool) *Config {

	return &Config{ip: bindIP, port: bindPort, TLSCertFile: TLSCertFile, TLSKeyFile: TLSKeyFile, UseTLS: UseTLS}

}

func StartServer(c *Config) bool {

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/api/v1/metrics", receiveMetrics)
	http.HandleFunc("/api/v1/apm", receiveAPM)

	if !c.UseTLS { 
		fmt.Println("Starting " + variables.AppName + " server...")
		fmt.Println("Version: " + variables.Version)
		fmt.Println("Commit: " + variables.CommitID)
		fmt.Println("")
		fmt.Println("WARNING: Server is not secure!")
		fmt.Println("         TLS should be enabled to encrypt communications.")
		fmt.Println("         Do not run in production environment!")
		fmt.Println("")
		fmt.Printf("Listening on http://%s:%s\n\n", c.ip, c.port)
		http.ListenAndServe(c.ip + ":" + c.port, nil)
	} else {
		fmt.Println("Starting server...")
		if checkTLSCert(c.TLSCertFile) && checkTLSKey(c.TLSKeyFile) {
			fmt.Printf("Listening on https://%s:%s\n\n", c.ip, c.port)
			log.Fatal(http.ListenAndServeTLS(":" + c.port, c.TLSKeyFile, c.TLSCertFile, nil))
		} else {
			log.Fatal("Could not start server")
		}
	}
	return true
}

func checkTLSCert(certFile string) bool {
	if len(certFile) > 0 {
		if _, err := os.Stat(certFile); !os.IsNotExist(err) {
			return true
		}
		log.Fatal("Certificate file not found - " + certFile)
		return false
	}
	log.Fatal("A certificate file must be specified!")
	return false
}

func checkTLSKey(keyFile string) bool {
	if len(keyFile) > 0 {
		if _, err := os.Stat(keyFile); !os.IsNotExist(err) {
			return true
		}
		log.Fatal("Certificate key file not found - " + keyFile)
		return false
	}
	log.Fatal("A certificate key file must be specified!")
	return false
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "API V1 Running")
}

func receiveMetrics(w http.ResponseWriter, r *http.Request) {
	message := models.Message{200,"Metric saved successfully"}

	js, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
func receiveAPM(w http.ResponseWriter, r *http.Request) {
	m := models.Metric{101010101,"cpu_load","45",time.Now()}

	js, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

}