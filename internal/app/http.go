package app

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// NewHTTPHandler 返回带健康检查与 demo 生产者的根处理器。
// pubSub 为 nil 时 /demo/publish 返回 503。
func NewHTTPHandler(pubSub *gochannel.GoChannel) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /demo/publish", func(w http.ResponseWriter, r *http.Request) {
		if pubSub == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "pubsub unavailable"})
			return
		}
		payload, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "read body"})
			return
		}
		if len(payload) == 0 {
			payload = []byte(`{"hello":"watermill"}`)
		}
		msg := message.NewMessage(watermill.NewUUID(), payload)
		if err := pubSub.Publish(demoTopic, msg); err != nil {
			slog.ErrorContext(r.Context(), "publish demo message", "error", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "publish failed"})
			return
		}
		slog.InfoContext(r.Context(), "demo message published", "uuid", msg.UUID, "topic", demoTopic)
		writeJSON(w, http.StatusAccepted, map[string]string{
			"status": "published",
			"uuid":   msg.UUID,
			"topic":  demoTopic,
		})
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
