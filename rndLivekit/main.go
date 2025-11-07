package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/livekit/protocol/auth"
	livekit "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

// --- Config & helpers ---

type Config struct {
	APIKey    string
	APISecret string
	ServerURL string // boleh http(s)://; kita akan normalkan jadi ws(s):// untuk client
	Port      string
}

func loadConfig() Config {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		APIKey:    os.Getenv("LIVEKIT_API_KEY"),
		APISecret: os.Getenv("LIVEKIT_API_SECRET"),
		ServerURL: os.Getenv("LIVEKIT_SERVER_URL"),
		Port:      port,
	}
}

// konversi http(s) -> ws(s). kalau sudah ws(s) tetap.
func normalizeWS(u string) string {
	if u == "" {
		return u
	}
	trimmed := strings.TrimRight(u, "/")
	if strings.HasPrefix(trimmed, "http://") {
		return "ws://" + strings.TrimPrefix(trimmed, "http://")
	}
	if strings.HasPrefix(trimmed, "https://") {
		return "wss://" + strings.TrimPrefix(trimmed, "https://")
	}
	return trimmed
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// --- Token ---

func generateToken(apiKey, apiSecret, roomName, identity string, ttl time.Duration) (string, error) {
	at := auth.NewAccessToken(apiKey, apiSecret)

	canPub := true
	canSub := true
	canData := true

	grant := &auth.VideoGrant{
		RoomJoin:       true,
		Room:           roomName,
		CanPublish:     &canPub,
		CanSubscribe:   &canSub,
		CanPublishData: &canData,
	}

	at.SetVideoGrant(grant).
		SetIdentity(identity).
		SetValidFor(ttl)

	return at.ToJWT()
}

// --- Room Service client factory ---

func newRoomService(cfg Config) *lksdk.RoomServiceClient {
	return lksdk.NewRoomServiceClient(cfg.ServerURL, cfg.APIKey, cfg.APISecret)
}

// --- HTTP Handlers ---

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// GET /token?room=..&identity=..
// POST /token  { "room": "..", "identity": ".." }
func tokenHandler(cfg Config) http.HandlerFunc {
	type reqBody struct {
		Room     string `json:"room"`
		Identity string `json:"identity"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var room, identity string
		switch r.Method {
		case http.MethodGet:
			room = r.URL.Query().Get("room")
			identity = r.URL.Query().Get("identity")
		case http.MethodPost:
			var body reqBody
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid json payload", http.StatusBadRequest)
				return
			}
			room = body.Room
			identity = body.Identity
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if room == "" || identity == "" {
			http.Error(w, "room and identity are required", http.StatusBadRequest)
			return
		}

		jwt, err := generateToken(cfg.APIKey, cfg.APISecret, room, identity, 2*time.Hour)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to generate token: %v", err), http.StatusInternalServerError)
			return
		}

		resp := map[string]string{
			"identity": identity,
			"room":     room,
			"token":    jwt,
			// kirim host untuk JS client; wajib ws(s)://
			"host": normalizeWS(cfg.ServerURL),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// POST /rooms {name, emptyTimeout?, maxParticipants?, metadata?}
func createRoomHandler(cfg Config) http.HandlerFunc {
	type createReq struct {
		Name            string `json:"name"`
		EmptyTimeout    int32  `json:"emptyTimeout,omitempty"`
		MaxParticipants uint32 `json:"maxParticipants,omitempty"`
		Metadata        string `json:"metadata,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body createReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
			http.Error(w, "invalid payload (need name)", http.StatusBadRequest)
			return
		}
		svc := newRoomService(cfg)
		room, err := svc.CreateRoom(r.Context(), &livekit.CreateRoomRequest{
			Name:            body.Name,
			EmptyTimeout:    uint32(body.EmptyTimeout),
			MaxParticipants: body.MaxParticipants,
			Metadata:        body.Metadata,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("create room error: %v", err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(room)
	}
}

// GET /rooms
func listRoomsHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		svc := newRoomService(cfg)
		resp, err := svc.ListRooms(r.Context(), &livekit.ListRoomsRequest{})
		if err != nil {
			http.Error(w, fmt.Sprintf("list rooms error: %v", err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp.Rooms)
	}
}

// DELETE /rooms?name=room1
func deleteRoomHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		svc := newRoomService(cfg)
		if _, err := svc.DeleteRoom(r.Context(), &livekit.DeleteRoomRequest{Room: name}); err != nil {
			http.Error(w, fmt.Sprintf("delete room error: %v", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// GET /participants?room=room1
func listParticipantsHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		room := r.URL.Query().Get("room")
		if room == "" {
			http.Error(w, "room is required", http.StatusBadRequest)
			return
		}
		svc := newRoomService(cfg)
		resp, err := svc.ListParticipants(r.Context(), &livekit.ListParticipantsRequest{Room: room})
		if err != nil {
			http.Error(w, fmt.Sprintf("list participants error: %v", err), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp.Participants)
	}
}

// DELETE /participants?room=room1&identity=user1
func removeParticipantHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		room := r.URL.Query().Get("room")
		identity := r.URL.Query().Get("identity")
		if room == "" || identity == "" {
			http.Error(w, "room and identity are required", http.StatusBadRequest)
			return
		}
		svc := newRoomService(cfg)
		_, err := svc.RemoveParticipant(r.Context(), &livekit.RoomParticipantIdentity{
			Room:     room,
			Identity: identity,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("remove participant error: %v", err), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	cfg := loadConfig()
	if cfg.APIKey == "" || cfg.APISecret == "" || cfg.ServerURL == "" {
		log.Fatal("LIVEKIT_API_KEY / LIVEKIT_API_SECRET / LIVEKIT_SERVER_URL are required")
	}

	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/healthz", healthHandler)

	// token
	mux.HandleFunc("/token", tokenHandler(cfg))

	// rooms
	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createRoomHandler(cfg)(w, r)
		case http.MethodGet:
			listRoomsHandler(cfg)(w, r)
		case http.MethodDelete:
			deleteRoomHandler(cfg)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// participants
	mux.HandleFunc("/participants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listParticipantsHandler(cfg)(w, r)
		case http.MethodDelete:
			removeParticipantHandler(cfg)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":" + cfg.Port
	fmt.Printf("✅ LiveKit backend running on http://localhost%s\n", addr)
	fmt.Printf("   - GET/POST /token?room=ROOM&identity=USER\n")
	fmt.Printf("   - POST     /rooms {name, emptyTimeout?, maxParticipants?, metadata?}\n")
	fmt.Printf("   - GET      /rooms\n")
	fmt.Printf("   - DELETE   /rooms?name=ROOM\n")
	fmt.Printf("   - GET      /participants?room=ROOM\n")
	fmt.Printf("   - DELETE   /participants?room=ROOM&identity=USER\n")

	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}
