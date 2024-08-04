package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/trenchesdeveloper/fingo-backend/db/sqlc"
	"github.com/trenchesdeveloper/fingo-backend/utils"
	"net/http"
)

type Auth struct {
	server *Server
}

func (a *Auth) router(server *Server) {
	a.server = server
	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/register", a.register)
	serverGroup.POST("/login", a.login)
}

type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (a *Auth) register(c *gin.Context) {
	var args UserParams

	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := utils.GenerateHashedPassword(args.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := a.server.queries.CreateUser(c, db.CreateUserParams{
		Email:          args.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "email already exists",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(user))
}

func (a *Auth) login(c *gin.Context) {
	var params LoginParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := a.server.queries.GetUserByEmail(c, params.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = utils.CompareHashedPassword(user.HashedPassword, params.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	token, err := tokenController.CreateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}
