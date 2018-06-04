package main

import (
    "github.com/micro/go-plugins/registry/eureka"
    "io/ioutil"
    "encoding/json"
    "log"
    "os"
    "golang.org/x/net/context"
    "io"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/movieservice"
    "github.com/micro/go-micro/cmd"
    "github.com/micro/go-micro/client"
    "github.com/micro/go-micro/registry"
)

const (
    //ADDRESS         = "mongo-api-grpc-demo.cfapps.io"
    ADDRESS         = "localhost:8080"
    DEFAULTFILENAME = "C:/Users/rande/go/src/github.com/Piszmog/mongo-api-grpc-example/cmd/cli/movie.json"
)

func parseFile(file string) (*movieservice.Movie, error) {
    var company *movieservice.Movie
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    json.Unmarshal(data, &company)
    return company, err
}

func main() {
    cmd.Init()
    movieClient := movieservice.NewMovieService("movieservice",
        client.NewClient(client.Registry(eureka.NewRegistry(registry.Addrs("http://localhost:8761/eureka")))))

    // Contact the server and print out its response.
    file := DEFAULTFILENAME
    if len(os.Args) > 1 {
        file = os.Args[1]
    }

    movieFile, err := parseFile(file)

    if err != nil {
        log.Fatalf("Could not parse file: %v", err)
    }

    r, err := movieClient.CreateMovie(context.Background(), movieFile)
    if err != nil {
        log.Fatalf("Could not create movie: %v", err)
    }
    log.Printf("Created: %v", r.Id)
    movie, err := movieClient.GetMovie(context.Background(), r)
    if err != nil {
        log.Fatalf("Could not get movie: %v", err)
    }
    log.Printf("Found movie %v", movie)
    movieClient.UpdateMovie(context.Background(), &movieservice.Movie{Id: movie.Id, Name: "A new movie", Description: "A new movie"})

    getAll, err := movieClient.GetMovies(context.Background(), &movieservice.Empty{})
    if err != nil {
        log.Fatalf("Could not list movies: %v", err)
    }
    for {
        movie, err := getAll.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("%v.getMovies(_) = _, %v", movieClient, err)
        }
        log.Printf("Movie: %v", movie)
        movieClient.DeleteMovie(context.Background(), &movieservice.MovieId{Id: movie.Id})
    }
}
