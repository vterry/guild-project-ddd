package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// responseWriter aprimorado para capturar status e corpo da resposta.
type responseWriter struct {
	http.ResponseWriter
	status       int
	responseBody *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // Status padrão
		responseBody:   new(bytes.Buffer),
	}
}

// WriteHeader armazena o código de status.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captura o corpo da resposta e o escreve na resposta original.
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.responseBody.Write(b) // Captura o corpo
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware aprimorado.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// --- Leitura do Corpo da Requisição ---
		var requestBodyBytes []byte
		if r.Body != nil {
			requestBodyBytes, _ = io.ReadAll(r.Body)
		}
		// Restaura o corpo para que o próximo handler possa lê-lo.
		r.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

		// Cria nosso responseWriter customizado para capturar status e corpo.
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Loga as informações do request e response usando slog.Group para melhor estrutura.
		slog.Info(
			"requisição processada",
			slog.Group("http",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Int("status", rw.status),
				slog.Duration("duration", duration),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_body", string(requestBodyBytes)),
				slog.String("response_body", rw.responseBody.String()),
			),
		)
	})
}
