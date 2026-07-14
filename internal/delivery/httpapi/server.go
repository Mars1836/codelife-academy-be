package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	authdomain "codelife-study-be/internal/domain/auth"
	domain "codelife-study-be/internal/domain/document"
	progressdomain "codelife-study-be/internal/domain/progress"
	authusecase "codelife-study-be/internal/usecase/auth"
	usecase "codelife-study-be/internal/usecase/document"
	progressusecase "codelife-study-be/internal/usecase/progress"
)

type Pinger interface{ Ping(context.Context) error }

type Server struct {
	documents       *usecase.Service
	auth            *authusecase.Service
	progress        *progressusecase.Service
	postgres, redis Pinger
	logger          *slog.Logger
	maxBodyBytes    int64
}

func New(documents *usecase.Service, auth *authusecase.Service, progress *progressusecase.Service, postgres, redis Pinger, logger *slog.Logger, maxBodyBytes int64) http.Handler {
	if maxBodyBytes <= 0 {
		maxBodyBytes = 1 << 20
	}
	s := &Server{documents: documents, auth: auth, progress: progress, postgres: postgres, redis: redis, logger: logger, maxBodyBytes: maxBodyBytes}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /readyz", s.ready)
	mux.HandleFunc("GET /api/v1/documents", s.listDocuments)
	mux.HandleFunc("GET /api/v1/documents/{slug}", s.getDocument)
	mux.HandleFunc("POST /api/v1/auth/register", s.register)
	mux.HandleFunc("POST /api/v1/auth/verify-email", s.verifyEmail)
	mux.HandleFunc("POST /api/v1/auth/login", s.login)
	mux.HandleFunc("GET /api/v1/auth/me", s.me)
	mux.HandleFunc("GET /api/v1/progress", s.listProgress)
	mux.HandleFunc("PUT /api/v1/progress/{slug}", s.updateProgress)
	return requestID(recoverer(logger, logging(logger, mux)))
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if s.postgres != nil && s.postgres.Ping(ctx) != nil {
		writeError(w, http.StatusServiceUnavailable, "postgres is unavailable")
		return
	}
	if s.redis != nil && s.redis.Ping(ctx) != nil {
		writeError(w, http.StatusServiceUnavailable, "redis is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := s.documents.List(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": documents})
}

func (s *Server) getDocument(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.authenticatedUser(w, r); !ok {
		return
	}
	document, err := s.documents.Get(r.Context(), r.PathValue("slug"))
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": document})
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if !s.authAvailable(w) {
		return
	}
	var input authusecase.RegisterInput
	if !s.readJSON(w, r, &input) {
		return
	}
	if err := s.auth.Register(r.Context(), input); errors.Is(err, authdomain.ErrEmailAlreadyExists) {
		writeError(w, http.StatusConflict, "email already exists")
		return
	} else if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": map[string]string{"message": "otp sent"}})
}

func (s *Server) verifyEmail(w http.ResponseWriter, r *http.Request) {
	if !s.authAvailable(w) {
		return
	}
	var input authusecase.VerifyEmailInput
	if !s.readJSON(w, r, &input) {
		return
	}
	user, err := s.auth.VerifyEmail(r.Context(), input)
	if errors.Is(err, authdomain.ErrInvalidOTP) {
		writeError(w, http.StatusBadRequest, "invalid or expired otp")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": user})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.authAvailable(w) {
		return
	}
	var input authusecase.LoginInput
	if !s.readJSON(w, r, &input) {
		return
	}
	session, err := s.auth.Login(r.Context(), input)
	if errors.Is(err, authdomain.ErrEmailNotVerified) {
		writeError(w, http.StatusForbidden, "email is not verified")
		return
	}
	if errors.Is(err, authdomain.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": session})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	if !s.authAvailable(w) {
		return
	}
	token := bearerToken(r.Header.Get("Authorization"))
	user, err := s.auth.Me(r.Context(), token)
	if errors.Is(err, authdomain.ErrInvalidCredentials) || errors.Is(err, authdomain.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": user})
}

func (s *Server) listProgress(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok || !s.progressAvailable(w) {
		return
	}
	items, err := s.progress.List(r.Context(), user.ID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (s *Server) updateProgress(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok || !s.progressAvailable(w) {
		return
	}
	var input progressusecase.UpdateInput
	if !s.readJSON(w, r, &input) {
		return
	}
	item, err := s.progress.Update(r.Context(), user.ID, r.PathValue("slug"), input)
	if errors.Is(err, progressdomain.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, progressdomain.ErrDocumentNotFound) {
		writeError(w, http.StatusNotFound, "document not found")
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
}

func (s *Server) authenticatedUser(w http.ResponseWriter, r *http.Request) (authdomain.User, bool) {
	if !s.authAvailable(w) {
		return authdomain.User{}, false
	}
	user, err := s.auth.Me(r.Context(), bearerToken(r.Header.Get("Authorization")))
	if errors.Is(err, authdomain.ErrInvalidCredentials) || errors.Is(err, authdomain.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return authdomain.User{}, false
	}
	if err != nil {
		s.internalError(w, err)
		return authdomain.User{}, false
	}
	return user, true
}

func (s *Server) authAvailable(w http.ResponseWriter) bool {
	if s.auth != nil {
		return true
	}
	writeError(w, http.StatusServiceUnavailable, "auth requires postgres")
	return false
}

func (s *Server) progressAvailable(w http.ResponseWriter) bool {
	if s.progress != nil {
		return true
	}
	writeError(w, http.StatusServiceUnavailable, "learning progress requires postgres")
	return false
}

func (s *Server) readJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodyBytes)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return false
	}
	return true
}

func (s *Server) internalError(w http.ResponseWriter, err error) {
	s.logger.Error("request failed", "error", err)
	writeError(w, 500, "internal server error")
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"message": message}})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func bearerToken(value string) string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(strings.ToLower(value), "bearer ") {
		return ""
	}
	return strings.TrimSpace(value[7:])
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = time.Now().UTC().Format("20060102150405.000000000")
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func logging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	})
}
func recoverer(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if value := recover(); value != nil {
				log.Error("panic recovered", "value", value)
				writeError(w, 500, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var _ = strings.Builder{}
