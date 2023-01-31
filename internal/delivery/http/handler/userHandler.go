package handler

import (
	"api/internal/delivery/http/handler/requests"
	"api/internal/delivery/http/validator"
	"api/internal/domain"
	"api/internal/usecase"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
)

type UserHandler interface {
	GetAllUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	UpdateCurrentUser(c *gin.Context)
	UpdateCurrentUserPassword(c *gin.Context)
	GetCurrentUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	DeleteUsers(c *gin.Context)
	GetQr(c *gin.Context)
	VerifyTF(c *gin.Context)
	DisableTF(c *gin.Context)
}

type userHandler struct {
	UserService usecase.UserService
	Validator   validation.CustomValidator
}

func NewUserHandler(uc usecase.UserService, v validation.CustomValidator) UserHandler {
	return &userHandler{UserService: uc, Validator: v}
}

func (uh *userHandler) UpdateUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	data := &requests.UpdateUserInfoRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{}
	copier.Copy(user, &data)
	u, err := uh.UserService.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": u,
		})
	}
}

func (uh *userHandler) UpdateCurrentUser(c *gin.Context) {

	data := &requests.UpdateUserInfoRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("session_cookie")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

	}

	user := &domain.User{}
	copier.Copy(user, &data)

	
	r, err := uh.UserService.UpdateCurrentUser(token, user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": r})

}

func (uh *userHandler) UpdateCurrentUserPassword(c *gin.Context) {

	data := &requests.UpdateUserPasswordRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("session_cookie")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

	}

	r, err := uh.UserService.UpdateCurrentUserPassword(token, data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": r})

}

func (uh *userHandler) GetCurrentUser(c *gin.Context) {

	token, err := c.Cookie("session_cookie")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

	}

	r, err := uh.UserService.GetCurrentUser(token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, r)

}

func (uh *userHandler) DeleteUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	_, err := uh.UserService.DeleteUser(id)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "user deleted",
		})
	}

}

func (uh *userHandler) GetAllUsers(c *gin.Context) {
	u, _ := uh.UserService.GetAllUsers()
	c.JSON(http.StatusOK, gin.H{
		"users": u,
	})
}

func (uh *userHandler) DeleteUsers(c *gin.Context) {

	data := &requests.MultiID{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := uh.UserService.DeleteUsers(&data.Ids)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{})
	}

}

func (uh *userHandler) GetQr(c *gin.Context) {

	token, err := c.Cookie("session_cookie")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

	}

	s, r, err := uh.UserService.GetQr(token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"secret": s, "link": r})

}

func (uh *userHandler) VerifyTF(c *gin.Context) {

	data := &requests.TwoFaToken{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.Cookie("session_cookie")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

	}

	r, err := uh.UserService.VerifyTF(token, data.Token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !r {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)

}

func (uh *userHandler) DisableTF(c *gin.Context) {

	token, err := c.Cookie("session_cookie")

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

	}

	err = uh.UserService.DisableTF(token)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{})

}
