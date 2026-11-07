// Package presentation adapts auth use cases to HTTP.
package presentation

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	authapp "boilerplate/api/internal/auth/application"
	"boilerplate/api/internal/auth/domain"
	"boilerplate/api/internal/platform/session"
	"boilerplate/api/internal/shared/httpx"
)

// Handler serves /auth endpoints.
type Handler struct {
	svc      *authapp.Service
	sessions *session.Manager
}

// NewHandler builds a Handler.
func NewHandler(svc *authapp.Service, sessions *session.Manager) *Handler {
	return &Handler{svc: svc, sessions: sessions}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login verifies credentials and sets the HttpOnly session cookie.
func (h *Handler) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return httpx.Error(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	user, err := h.svc.Login(c.Context(), authapp.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if handled, herr := httpx.FieldError(c, err); handled {
			return herr
		}
		if errors.Is(err, authapp.ErrInvalidCredentials) {
			return httpx.Error(c, fiber.StatusUnauthorized, "invalid username or password")
		}
		return err
	}
	h.sessions.SetCookie(c, h.sessions.Issue(user.ID))
	return c.JSON(user)
}

// Me returns the authenticated user profile.
func (h *Handler) Me(c fiber.Ctx) error {
	user, err := h.svc.Me(c.Context(), session.UserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return httpx.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}
		return err
	}
	return c.JSON(user)
}

// Logout clears the session cookie.
func (h *Handler) Logout(c fiber.Ctx) error {
	h.sessions.ClearCookie(c)
	return c.JSON(fiber.Map{"ok": true})
}

// CSRF exists so SPAs can obtain a CSRF cookie with a safe GET.
// The csrf middleware sets the cookie; the body is just an ack.
func (h *Handler) CSRF(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"ok": true})
}
