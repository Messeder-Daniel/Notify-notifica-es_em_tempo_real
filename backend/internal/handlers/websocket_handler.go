package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ws "github.com/gorilla/websocket"
	internalwebsocket "github.com/messederdaniel/real-time-notifications/backend/internal/websocket"
)

type WebSocketHandler struct {
	hub       *internalwebsocket.Hub
	jwtSecret string
	upgrader  ws.Upgrader
}

func NewWebSocketHandler(hub *internalwebsocket.Hub, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{
		hub:       hub,
		jwtSecret: jwtSecret,
		upgrader: ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (handler *WebSocketHandler) Connect(ctx *gin.Context) {
	tokenString := ctx.Query("token")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token query parameter is required"})
		return
	}

	userID, userEmail, err := handler.validateToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := handler.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := internalwebsocket.NewClient(handler.hub, conn, userID, userEmail)

	handler.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (handler *WebSocketHandler) validateToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(handler.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", "", fmt.Errorf("invalid token subject")
	}

	expiresAt, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(expiresAt) {
		return "", "", fmt.Errorf("expired token")
	}

	userEmail, _ := claims["email"].(string)

	return userID, userEmail, nil
}
