package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

const StreamKey = "myfreekey"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("name")
		if key != StreamKey {
			log.Println("Unauthorized attempt:", key)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Denied"))
			return
		}
		log.Println("Authorized stream key:", key)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve HLS files
	r.PathPrefix("/live/").Handler(http.StripPrefix("/live/", http.FileServer(http.Dir("/tmp"))))

	// Web page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<body>
				<h1>Live Stream</h1>
				<video width="720" height="480" controls autoplay>
					<source src="/live/mystream.m3u8" type="application/x-mpegURL">
				</video>
			</body>
			</html>
		`))
	})

	log.Println("Go server running at :8081")
	http.ListenAndServe(":8081", r)
}
