package server

import (
	"crud-example-go/src/dto"
	"crud-example-go/src/intreface/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Server struct {
	port           string
	userRepository repository.UserRepository
}

func New(port string, userRepository repository.UserRepository) *Server {
	return &Server{
		port:           port,
		userRepository: userRepository}
}

func (s *Server) Start() error {
	router := gin.Default()

	router.GET("/users", getUsers(s))
	routerGroup := router.Group("/user")
	{
		routerGroup.GET("/:id", getUser(s))
		routerGroup.POST("", saveUser(s))
		routerGroup.DELETE("/:id", deleteUser(s))
	}

	if err := router.Run(":" + s.port); err != nil {
		return err
	}
	return nil
}

func getUsers(s *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		if users, err := s.userRepository.GetUsers(); err != nil {
			log.Error(err)
			internalServerError(c, err)
		} else {
			ok(c, users)
		}
	}
}

func getUser(s *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err)
			badRequest(c)
			return
		}

		if user, err := s.userRepository.GetUser(id); err != nil {
			log.Error(err)
			internalServerError(c, err)
		} else {
			if user == nil {
				c.String(http.StatusNotFound, "user not found")
			} else {
				ok(c, user)
			}
		}
	}
}

func saveUser(s *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		var user dto.User
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Error(err)
			badRequest(c)
			return
		}

		if user, err := s.userRepository.SaveUser(&user); err != nil {
			log.Error(err)
			internalServerError(c, err)
		} else {
			ok(c, user)
		}
	}
}

func deleteUser(s *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error(err)
			badRequest(c)
			return
		}

		if err := s.userRepository.DeleteUser(id); err != nil {
			log.Error(err)
			internalServerError(c, err)
		} else {
			ok(c, nil)
		}
	}
}

func badRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "bad request")
}

func internalServerError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func ok(c *gin.Context, responseBody interface{}) {
	c.JSON(http.StatusOK, responseBody)
}
