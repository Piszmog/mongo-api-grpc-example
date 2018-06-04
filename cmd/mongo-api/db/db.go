package db

import (
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "go.uber.org/zap"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/util"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/model"
)

type IRepository interface {
    FindAll() ([]model.Movie, error)
    FindById(id string) (model.Movie, error)
    Insert(movie model.Movie) error
    Delete(id string) error
    Update(id string, movie model.Movie) error
}

type Connection struct {
    Server     string
    Database   string
    db         *mgo.Database
    collection mgo.Collection
}

const (
    COLLECTION = "movies"
)

var logger zap.SugaredLogger

// Creates the logger for the package
func init() {
    zapLogger := util.CreateLogger()
    defer zapLogger.Sync()
    logger = *zapLogger.Sugar()
}

// Connect to the Mongdb Database using the server name and the database name
func (m *Connection) Connect(server string, database string) error {
    session, err := mgo.Dial(m.Server)
    if err != nil {
        return err
    }
    m.Server = server
    m.Database = database
    m.db = session.DB(m.Database)
    m.collection = *m.db.C(COLLECTION)
    return nil
}

// Connect to the Mongodb Database using a URI connection string
func (m *Connection) ConnectWithURL(url string) {
    info, parseErr := mgo.ParseURL(url)
    if parseErr != nil {
        logger.Fatalf("failed to parse URL, %v", parseErr)
    }
    session, err := mgo.DialWithInfo(info)
    if err != nil {
        logger.Fatalf("failed to connect to the db, %v", err)
    }
    m.Database = info.Database
    m.db = session.DB(m.Database)
    m.collection = *m.db.C(COLLECTION)
}

// Finds all movies
func (m *Connection) FindAll() ([]model.Movie, error) {
    var movies []model.Movie
    err := m.collection.Find(bson.M{}).All(&movies)
    return movies, err
}

// Finds the movie matching the provided id
func (m *Connection) FindById(id string) (model.Movie, error) {
    var movie model.Movie
    query := m.collection.FindId(id)
    if query == nil {
        return movie, nil
    }
    err := query.One(&movie)
    return movie, err
}

// Inserts the provided movie
func (m *Connection) Insert(movie model.Movie) error {
    err := m.collection.Insert(&movie)
    return err
}

// Delete the movie matching the id
func (m *Connection) Delete(id string) error {
    err := m.collection.RemoveId(id)
    return err
}

// Updates the movie matching the id with the provided body
func (m *Connection) Update(id string, movie model.Movie) error {
    err := m.collection.UpdateId(id, &movie)
    return err
}
