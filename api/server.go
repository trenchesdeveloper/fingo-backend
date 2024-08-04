package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/trenchesdeveloper/fingo-backend/db/sqlc"
	"github.com/trenchesdeveloper/fingo-backend/utils"
	"net/http"

	_ "github.com/lib/pq"
)

var tokenController *utils.JWTToken

type Server struct {
	queries *db.Queries
	router  *gin.Engine
	config  *utils.Config
}

func NewServer(envPath string) *Server {
	config, err := utils.LoadConfig(envPath)

	if err != nil {
		panic(fmt.Sprintf("cannot load config: %v", err))
	}

	conn, err := sql.Open(config.DBdriver, config.DB_source_live)

	if err != nil {
		panic(fmt.Sprintf("cannot connect to db: %v", err))
	}

	queries := db.New(conn)

	g := gin.Default()

	tokenController = utils.NewJWTToken(config)

	return &Server{
		queries: queries,
		router:  g,
		config:  config,
	}
}

func (s *Server) Start(port int) {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Fingo API",
		})
	})

	user := User{}
	auth := Auth{}
	user.router(s)
	auth.router(s)

	s.router.Run(fmt.Sprintf(":%d", port))
}
