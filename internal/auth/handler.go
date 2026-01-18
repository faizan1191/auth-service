package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/faizan1191/auth-service/internal/email"
	redisClient "github.com/faizan1191/auth-service/internal/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	repo  *Repository
	rdb   *redis.Client
	email *email.Sender
}

func NewHandler(repo *Repository, rdb *redis.Client, email *email.Sender) *Handler {
	return &Handler{repo: repo, rdb: rdb, email: email}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotRequest struct {
	Email string `json:"email"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest

	// Bind JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	// validate
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password con't be empty"})
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// store hashedPassword in DB
	user := &User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.repo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "signup success",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	// Bind JSON Input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and Password can't be empty"})
		return
	}

	// get user from db by email
	user, err := h.repo.GetUserByEmail(req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
	}

	refreshToken := GenerateRefreshToken()
	if err := redisClient.SetRefreshToken(h.rdb, refreshToken, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func GenerateRefreshToken() string {
	return uuid.NewString()
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid payload"})
		return
	}

	userID, err := redisClient.GetUserIDByRefreshToken(h.rdb, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find user"})
	}

	newAccessToken, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if err := redisClient.DeleteRefreshToken(h.rdb, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout unsuccessfull"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		// real error → log it
		log.Printf("forgot-password GetUserByEmail error: %v", err)

		// still return generic response
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a reset link has been sent",
		})
		return
	}

	if user == nil {
		// email does not exist → do nothing
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a reset link has been sent",
		})
		return
	}

	// create reset token
	token, err := GenerateResetToken()
	if err != nil {
		log.Printf("forgot-password  GenerateResetToken error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a reset link has been sent",
		})
		return
	}

	// store token in redis
	if err := redisClient.SetResetToken(h.rdb, token, user.ID); err != nil {
		log.Printf("forgot-password  SetResetToken error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a reset link has been sent",
		})
		return
	}

	// send email to user
	resetURL := fmt.Sprintf(
		"http://localhost:8080/auth/reset-password?token=%s",
		token,
	)

	if err = h.email.SendResetPassword(user.Email, resetURL); err != nil {
		log.Printf("forgot-password SendResetPassword error: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a reset link has been sent",
	})

}

func GenerateResetToken() (string, error) {
	b := make([]byte, 32) // 256-bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
