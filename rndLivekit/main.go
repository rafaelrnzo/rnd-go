package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/livekit/protocol/auth"
)

// loadEnv loads environment variables
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  No .env file found, using system environment.")
	}
}

// generateToken creates a LiveKit JWT for joining a room
func generateToken(apiKey, apiSecret, roomName, identity string) (string, error) {
	at := auth.NewAccessToken(apiKey, apiSecret)

	canPub := true
	canSub := true

	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPub,
		CanSubscribe: &canSub,
	}

	at.SetVideoGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return "", err
	}

	return token, nil
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")

	if room == "" || identity == "" {
		http.Error(w, "room and identity are required", http.StatusBadRequest)
		return
	}

	token, err := generateToken(apiKey, apiSecret, room, identity)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"identity": identity,
		"room":     room,
		"token":    token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	loadEnv()

	http.HandleFunc("/token", tokenHandler)

	port := "8080"
	fmt.Printf("🚀 LiveKit Auth API running on http://localhost:%s/token?room=myroom&identity=user1\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
