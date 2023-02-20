package main

import (
	"context"
	"log"

	v1 "github.com/turao/topics/api/v1"
	userRepository "github.com/turao/topics/users/repository/user"
	userService "github.com/turao/topics/users/service/user"

	fileRepository "github.com/turao/topics/movies/repository/file"
	movieRepository "github.com/turao/topics/movies/repository/movie"
	fileService "github.com/turao/topics/movies/service/file"
	movieService "github.com/turao/topics/movies/service/movie"
)

func main() {
	testMovieService()
}

func testMovieService() {
	movieRepo, err := movieRepository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	fileRepo, err := fileRepository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	movieSvc, err := movieService.NewService(movieRepo, fileRepo)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	res, err := movieSvc.RegisterMovie(ctx, v1.RegisterMovieRequest{
		Title:   "John Wick",
		URI:     "uri-example",
		Tenancy: "tenancy/test",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = movieSvc.DeleteMovie(ctx, v1.DeleteMovieRequest{
		ID: res.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	movieinfo, err := movieSvc.GetMovie(ctx, v1.GetMovieRequest{
		ID: res.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(movieinfo)

	_, err = movieSvc.DownloadMovie(ctx, v1.DownloadMovieRequest{
		ID: res.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	movieinfos, err := movieSvc.ListMovies(ctx, v1.ListMoviesRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(movieinfos)

	fileSvc, err := fileService.NewService(fileRepo)
	if err != nil {
		log.Fatal(err)
	}

	filesByMovie, err := fileSvc.ListFilesByMovie(ctx, v1.ListFilesByMovieRequest{
		MovieID: res.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(filesByMovie)

}

func testUserService() {
	userRepo, err := userRepository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	userSvc, err := userService.NewService(userRepo)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	res, err := userSvc.RegisterUser(ctx, v1.RegisteUserRequest{
		Email:     "john@doe.com",
		FirstName: "john",
		LastName:  "doe",
		Tenancy:   "tenancy/test",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userSvc.DeleteUser(ctx, v1.DeleteUserRequest{
		ID: res.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	userinfo, err := userSvc.GetUserInfo(ctx, v1.GetUserInfoRequest{
		ID: res.ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("userinfo", userinfo)
}
