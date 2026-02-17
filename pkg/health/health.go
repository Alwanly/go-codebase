package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response represents the health check response
type Response struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// Handler handles health check requests
type Handler struct {
	serviceName    string
	serviceVersion string
}

// NewHandler creates a new health check handler
func NewHandler(serviceName, serviceVersion string) *Handler {
	return &Handler{
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
	}
}

// Check godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /health [get]
func (h *Handler) Check(c *fiber.Ctx) error {
	return c.JSON(Response{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   h.serviceName,
		Version:   h.serviceVersion,
	})
}

// Readiness godoc
// @Summary Readiness check
// @Description Check if the service is ready to accept requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /ready [get]
func (h *Handler) Readiness(c *fiber.Ctx) error {
	return c.JSON(Response{
		Status:    "ready",
		Timestamp: time.Now(),
		Service:   h.serviceName,
		Version:   h.serviceVersion,
	})
}

// Liveness godoc
// @Summary Liveness check
// @Description Check if the service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /live [get]
func (h *Handler) Liveness(c *fiber.Ctx) error {
	return c.JSON(Response{
		Status:    "alive",
		Timestamp: time.Now(),
		Service:   h.serviceName,
		Version:   h.serviceVersion,
	})
}
