package main

import (
	"encoding/json"
	"net/http"
	"snackable/cache"
	"snackable/ext/snackable"
	"snackable/handler"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("server startup")

	snackableAPI := snackable.NewClient("http://interview-api.snackable.ai")
	fileCache := cache.NewInMemoryFileCache(nil)
	fileHandler := handler.NewFileHandler(snackableAPI, fileCache)

	router := httprouter.New()
	router.GET("/protected/file/:fileid", authorize(fileHandlerIntegration(fileHandler)))
	router.GET("/file/:fileid", fileHandlerIntegration(fileHandler))

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal().Err(err).Msg("failde to initialize http router")
	}
}

func fileHandlerIntegration(fileHandler handler.FileHandler) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fileID := p.ByName("fileid")
		resp, err := fileHandler.Handle(fileID)
		if err != nil && err == handler.ErrFileNotFound {
			log.Debug().Str("file_id", fileID).Msg("unknown id requested")
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if err != nil && err == handler.ErrFileNotFinished {
			log.Debug().Str("file_id", fileID).Msg("unfinished file")
			http.Error(rw, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(resp)
	}
}

func authorize(handler httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		hdr := r.Header.Get("Authorization")
		if hdr == "" {
			log.Debug().Msg("unauthorized request with empty Authorization header")
			http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// TODO: implement service that check Authorization token or integration 3rd service
		handler(rw, r, p)
	}
}
