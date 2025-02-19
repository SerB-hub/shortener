package main

import (
	"github.com/SerB-hub/shortener/internal/app/hasher"
	"github.com/SerB-hub/shortener/internal/app/middlewares"
	"github.com/SerB-hub/shortener/internal/app/repo"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ShortenerHandler struct {
	repository   repo.Repository
	hasher       hasher.Hasher
	shortUrlBase string
}

func (sh *ShortenerHandler) PostSrcUrl(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentTypeHeader := r.Header.Get("Content-Type"); !strings.HasPrefix(contentTypeHeader, "text/plain") {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		log.Printf("Error while read body: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyStr := string(body)
	body = nil

	if bodyStr == "" ||
		!strings.HasPrefix(bodyStr, "http") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	hashKey, err := sh.hasher.Hash(bodyStr)

	if err != nil {
		log.Printf(
			"Error while creating hash key: %v",
			err,
		)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sh.repository.SaveSrcUrlByHashKey(
		hashKey,
		bodyStr,
	)

	if err != nil {
		log.Printf(
			"Error while saving URL %v by hash key %v: %v",
			bodyStr,
			hashKey,
			err,
		)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	shortUrl := sh.buildShortUrl(hashKey)

	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Content-Length", strconv.Itoa(len(shortUrl)))
	rw.WriteHeader(201)
	_, _ = rw.Write([]byte(shortUrl))
}

func (sh *ShortenerHandler) GetShortUrlBySrcUrl(rw http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").(map[string]string)
	id := params["id"]

	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf(`get "id" param: %v`, id)

	srcUrl, err := sh.repository.GetSrcUrlByHashKey(id)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Location", srcUrl)
	rw.WriteHeader(307)
}

func (sh *ShortenerHandler) buildShortUrl(route string) string {
	sb := strings.Builder{}
	sb.WriteString(sh.shortUrlBase)
	sb.WriteString("/")
	sb.WriteString(route)

	return sb.String()
}

func main() {
	repository := repo.NewMemoryRepo()
	hasherService := &hasher.MD5Base62Hasher{}
	shortenerHandler := ShortenerHandler{
		repository:   repository,
		hasher:       hasherService,
		shortUrlBase: "http://localhost:8080",
	}

	postSrcUrl := http.HandlerFunc(shortenerHandler.PostSrcUrl)
	getShortUrlBySrcUrl := http.HandlerFunc(shortenerHandler.GetShortUrlBySrcUrl)

	router := middlewares.NewRouter(
		map[string]*http.HandlerFunc{
			"/":     &postSrcUrl,
			"/{id}": &getShortUrlBySrcUrl,
		},
	)

	mux := http.NewServeMux()
	mux.Handle("/", router.ProcessRequest(nil))

	err := http.ListenAndServe("localhost:8080", mux)

	if err != nil {
		panic(err)
	}
}
