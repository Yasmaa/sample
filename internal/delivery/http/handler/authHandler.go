package handler

import (
	"api/internal/delivery/http/handler/requests"
	"api/internal/delivery/http/validator"
	"api/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)

	Confirm(c *gin.Context)
	ResetLink(c *gin.Context)
	ResetPassword(c *gin.Context)
	VerifyToken(c *gin.Context)

	VerifyTF(c *gin.Context)
}

type authHandler struct {
	AuthService usecase.AuthService
	Validator   validation.CustomValidator
}

func NewAuthHandler(uc usecase.AuthService, v validation.CustomValidator) AuthHandler {
	return &authHandler{AuthService: uc, Validator: v}
}

// register
// @Summary Register a user
// @Description create a user
// @Accept  json
// @Produce  json
// @Router /api/create [post]
func (uh *authHandler) Register(c *gin.Context) {

	data := &requests.RegisterRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rs, err := uh.AuthService.Register(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		rs.Password = ""
		c.JSON(http.StatusOK, gin.H{
			"data": rs,
		})
	}

}

// login
func (uh *authHandler) Login(c *gin.Context) {

	data := &requests.LoginRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Validator.Validate(data); err != nil {

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, tfa, err := uh.AuthService.Login(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tfa {
		c.IndentedJSON(http.StatusTemporaryRedirect, gin.H{})
		return
	}

	c.SetCookie("session_cookie", token, 24*3*60*60, "/", "localhost", false, true)

	c.IndentedJSON(http.StatusOK, gin.H{})
}

func (uh *authHandler) Logout(c *gin.Context) {

	c.SetCookie("session_cookie", "", 0, "", "", true, true)
	c.IndentedJSON(http.StatusOK, gin.H{})

}

func (uh *authHandler) Confirm(c *gin.Context) {
	token := c.Param("token")
	b := uh.AuthService.Confirm(token)
	c.JSON(http.StatusOK, gin.H{
		"data": b,
	})
}

func (uh *authHandler) ResetLink(c *gin.Context) {

	data := &requests.ResetLinkRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.AuthService.ResetLink(data)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": ""})

}

func (uh *authHandler) ResetPassword(c *gin.Context) {

	data := &requests.ResetPasswordRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uh.AuthService.VerifyToken(data.Token)

	if err != nil {

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	if data.Password != data.Confirm {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match !"})
		return
	}

	err = uh.AuthService.ResetPassword(data)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": ""})

}

func (uh *authHandler) VerifyToken(c *gin.Context) {

	resetToken, _ := c.GetQuery("token")
	err := uh.AuthService.VerifyToken(resetToken)
	if err != nil {

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.IndentedJSON(http.StatusOK, gin.H{})

}

func (uh *authHandler) VerifyTF(c *gin.Context) {

	data := &requests.TwoFaToken{}
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uh.Validator.Validate(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, tfa, err := uh.AuthService.VerifyTF(data)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !tfa {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("session_cookie", token, 24*3*60*60, "/", "localhost", false, true)

	c.IndentedJSON(http.StatusOK, gin.H{})

}
