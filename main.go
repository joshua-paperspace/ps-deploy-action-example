package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

var _load sync.Mutex
var done chan struct{}

func startLoad() {
	_load.Lock()
	defer _load.Unlock()
	if done != nil {
		log.Info("Load is already being generated")
	}
	done = make(chan struct{}, 1)
	log.Info("Generating load...")
	for i := 0; i < runtime.NumCPU(); i++ {
		i := i
		go func() {
			log.Infof("Staring load generator on cpu %d", i)
			expensive := 0
			ticker := time.NewTicker(10 * time.Second)
			for {
				select {
				case <-ticker.C:
					log.Infof("Another 10 seconds for load generation have on cpu %d", i)
				case <-done:
					log.Infof("Shutting down load generator on cpu %d", i)
					return
				default:
					expensive++
				}
			}
		}()
	}
}

func stopLoad() {
	_load.Lock()
	defer _load.Unlock()
	log.Info("Shutting down load generation...")
	for i := 0; i < runtime.NumCPU(); i++ {
		done <- struct{}{}
	}
	done = nil
}

func ls(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to read directory %s: %w", path, err)
	}
	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

type Error struct {
	Message string `json:"message"`
}

func main() {
	workspacePath := "/workspace"
	modelPath := os.Getenv("MODEL_BASE_PATH")
	if modelPath == "" {
		modelPath = "/model"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ls", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Receieved request to /ls")
		resp := map[string][]string{}
		enc := json.NewEncoder(w)

		var err error
		resp[modelPath], err = ls(modelPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			if err := enc.Encode(Error{err.Error()}); err != nil {
				log.Fatal(err.Error())
			}
			return
		}
		resp[workspacePath], err = ls(workspacePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			if err := enc.Encode(Error{err.Error()}); err != nil {
				log.Fatal(err.Error())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := enc.Encode(resp); err != nil {
			log.Fatal(err.Error())
		}
	})
	mux.HandleFunc("/load/start", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Receieved request to /load/start, starting load...")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
		startLoad()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"load": "started"}`)
	})
	mux.HandleFunc("/load/stop", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Receieved request to /load/stop, stopping load...")
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
		stopLoad()
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"load": "stopped"}`)
	})
	mux.HandleFunc("/args", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Received request to /args")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w)
	})
	mux.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Receieved request to /env")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(os.Environ()); err != nil {
			log.Fatal(err.Error())
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Receieved request to /! :)")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		if err := enc.Encode(os.Args); err != nil {
			log.Fatal(err.Error())
		}
	})

	log.Info("Starting on :8888")
	if err := http.ListenAndServe(":8888", mux); err != nil {
		log.Fatal(err.Error())
	}
}