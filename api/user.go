package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/trenchesdeveloper/fingo-backend/db/sqlc"
	"net/http"
)

type User struct {
	server *Server
}

func (u *User) router(server *Server) {
	u.server = server
	serverGroup := server.router.Group("/users", AuthenticatedMiddleware())

	serverGroup.GET("", u.listUsers)
	serverGroup.GET("/me", u.getCurrentUser)

}

type UserParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (u *User) listUsers(c *gin.Context) {
	users, err := u.server.queries.ListUsers(c, db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response []UserResponse

	for _, user := range users {
		response = append(response, ToUserResponse(user))
	}

	c.JSON(http.StatusOK, response)
}

func ToUserResponse(user db.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func (u *User) getCurrentUser(c *gin.Context) {
	userID, exist := c.Get("user_id")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
		return
	}

	convertedUserID, ok := userID.(int64)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	user, err := u.server.queries.GetUserByID(c, convertedUserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ToUserResponse(user))
}
