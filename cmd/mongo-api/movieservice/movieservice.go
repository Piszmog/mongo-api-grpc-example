package movieservice

import (
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/db"
    "golang.org/x/net/context"
    "github.com/Piszmog/mongo-api-grpc-example/cmd/mongo-api/model"
    "github.com/google/uuid"
    "github.com/micro/go-micro/server"
)

type MovieRepositoryService struct {
    repository db.IRepository
}

func RegisterService(server server.Server, repository db.IRepository) {
    RegisterMovieServiceHandler(server, &MovieRepositoryService{repository})
}

func (movieRepositoryService *MovieRepositoryService) CreateMovie(context context.Context, movie *Movie, movieId *MovieId) error {
    newMovieId := uuid.New().String()
    newMovie := model.Movie{Id: newMovieId, Name: movie.Name, Description: movie.Description}
    movieRepositoryService.repository.Insert(newMovie)
    movieId.Id = newMovieId
    return nil
}

func (movieRepositoryService *MovieRepositoryService) GetMovies(context context.Context, empty *Empty, stream MovieService_GetMoviesStream) error {
    movies, err := movieRepositoryService.repository.FindAll()
    if err != nil {
        return err
    }
    for _, movie := range movies {
        if err := stream.Send(&Movie{Id: movie.Id, Name: movie.Name, Description: movie.Description}); err != nil {
            return err
        }
    }
    return nil
}

func (movieRepositoryService *MovieRepositoryService) GetMovie(context context.Context, id *MovieId, movie *Movie) error {
    dbMovie, err := movieRepositoryService.repository.FindById(id.Id)
    if err != nil {
        return err
    }
    movie.Id = dbMovie.Id
    movie.Name = dbMovie.Name
    movie.Description = dbMovie.Description
    return nil
}

func (movieRepositoryService *MovieRepositoryService) UpdateMovie(context context.Context, requestMovie *Movie, responseMovie *Movie) error {
    updatedMovie := model.Movie{Id: requestMovie.Id, Name: requestMovie.Name, Description: requestMovie.Description}
    err := movieRepositoryService.repository.Update(requestMovie.Id, updatedMovie)
    if err != nil {
        return err
    }
    responseMovie = requestMovie
    return nil
}

func (movieRepositoryService *MovieRepositoryService) DeleteMovie(context context.Context, id *MovieId, empty *Empty) error {
    err := movieRepositoryService.repository.Delete(id.Id)
    if err != nil {
        return err
    }
    return nil
}
