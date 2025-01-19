package service

import (
	"net/http"
	"timeslot-app/models"
	"timeslot-app/repository"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(db *pgx.Conn) *UserService {
	service := new(UserService)
	service.userRepo = repository.NewUserRepo(db)
	return service
}

// ShowAccount godoc
// @Summary      Create a user
// @Description  Create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body   body   	models.UserCreateRequest   true "Create User request body"
// @Success      200  {object}  string "User created successfully"
// @Failure      400  {object}  string "Invalid request body"
// @Failure      500  {object}  string "Error creating user"
// @Router       /user [post]
func (ts *UserService) CreateUser(ctx *gin.Context) {
	// create time slot
	var userReq models.UserCreateRequest
	if err := ctx.BindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := models.User{}
	userID, err := uuid.NewV4()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating user ID"})
		return
	}

	user.ID = userID
	user.Name = userReq.Name
	// save the time slot for the user.

	err = ts.userRepo.Create(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
