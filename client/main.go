package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mjudeikis/grpc-video-streaming/proto"
	"google.golang.org/grpc"
)

type Service struct {
	server *http.Server
	router *mux.Router
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

var client proto.StreamServiceClient

func main() {
	log.Print("start")

	serverURI := "localhost:4040"
	if os.Getenv("SERVER") != "" {
		serverURI = os.Getenv("SERVER")
	}

	// this is global....
	conn, err := grpc.Dial(serverURI, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	client = proto.NewStreamServiceClient(conn)

	log.Print("listening :8000")

	s := &Service{}

	s.router = setupRouter()

	s.router.Handle("/", http.FileServer(http.Dir("./public"))).Methods("GET")
	s.router.HandleFunc("/media/{mID:[0-9]+}/stream/", streamHandler).Methods("GET")
	s.router.HandleFunc("/media/{mID:[0-9]+}/stream/{segName:index[0-9]+.ts}", streamHandler).Methods("GET")

	s.server = &http.Server{
		Addr: ":8000",
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
		)(s),
	}

	err = s.server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}

func streamHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	segName, ok := vars["segName"]
	if !ok {
		serveHlsM3u8(response, request, "./public/media", "index.m3u8")
	} else {
		serveHlsTs(response, request, "./public/media", segName)
	}
}

func serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	mediaFile := fmt.Sprintf("%s/%s", mediaBase, m3u8Name)
	log.Println(mediaFile)
	req := &proto.Request{Filename: mediaFile}

	if response, err := client.GetFile(context.Background(), req); err == nil {
		fo, err := os.Create(mediaFile)
		if err != nil {
			log.Fatal("Failed to Create File")
		}
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		if _, err := fo.Write(response.GetContent()); err != nil {
			log.Fatal("Error in writing file: ", err.Error())
		}
		http.ServeFile(w, r, mediaFile)
		w.Header().Set("Content-Type", "application/x-mpegURL")
	} else {
		fmt.Println("Error getting response: ", err.Error())
	}
}

func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	mediaFile := fmt.Sprintf("%s/%s", mediaBase, segName)
	log.Println(mediaFile)
	req := &proto.Request{Filename: mediaFile}

	if response, err := client.GetFile(context.Background(), req); err == nil {
		fo, err := os.Create(mediaFile)
		if err != nil {
			log.Fatal("Failed to Create File")
		}
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		if _, err := fo.Write(response.GetContent()); err != nil {
			log.Fatal("Error in writing file: ", err.Error())
		}
		http.ServeFile(w, r, mediaFile)
		w.Header().Set("Content-Type", "video/MP2T")
	} else {
		fmt.Println("Error getting response: ", err.Error())
	}
}
