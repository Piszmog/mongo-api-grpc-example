package main

import (
    "log"
    "os"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/movieservice"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/db"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/model"
    "encoding/json"
    "github.com/micro/go-micro"
    "fmt"
    "github.com/micro/go-plugins/registry/eureka"
    "github.com/micro/go-micro/registry"
)

const (
    CfServices      = "VCAP_SERVICES"
    DefaultDatabase = "test"
    DefaultServer   = "localhost"
)

// Runs the application
func main() {
    connection := db.Connection{}
    err := connectToDB(&connection)
    if err != nil {
        log.Fatalf("failed to connect to the database, %v", err)
    }
    // Create a new service. Optionally include some options here.
    service := micro.NewService(
        // This name must match the package name given in your protobuf definition
        micro.Name("movieservice"),
        micro.Version("latest"),
        micro.Registry(eureka.NewRegistry(registry.Addrs("http://localhost:8761/eureka"))),
    )
    //required to read env variables or commandline args
    service.Init()
    movieservice.RegisterService(service.Server(), &connection)
    if err := service.Run(); err != nil {
        fmt.Println(err)
    }
}

func connectToDB(connection *db.Connection) error {
    cfServices := os.Getenv(CfServices)
    if len(cfServices) == 0 {
        connection.Connect(DefaultServer, DefaultDatabase)
    } else {
        var env model.CloudFoundryEnvironment
        err := json.Unmarshal([]byte(cfServices), &env)
        if err != nil {
            return err
        }
        connection.ConnectWithURL(env.Mlab[0].Credentials.Uri)
    }
    return nil
}
