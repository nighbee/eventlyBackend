package booking

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

type createBookingRequest struct {
	EventID    uint `json:"event_id"`
	SeatNumber int  `json:"seat_number"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := toUint(userIDVal)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req createBookingRequest
	if err := c.BodyParser(&req); err != nil || req.EventID == 0 || req.SeatNumber <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// Validate seat number is within event capacity
	if err := h.repo.ValidateSeat(req.EventID, req.SeatNumber); err != nil {
		if err.Error() == "event not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "event not found"})
		}
		if err.Error() == "seat number out of range" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "validation failed"})
	}

	b := &Booking{
		UserID:     userID,
		EventID:    req.EventID,
		SeatNumber: req.SeatNumber,
		Active:     true,
		Paid:       false,
		Status:     "reserved",
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	if err := h.repo.CreateIfSeatFree(c.Context(), b); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "seat already booked"})
	}
	return c.Status(fiber.StatusCreated).JSON(b)
}

func (h *Handler) ListMine(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := toUint(userIDVal)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	list, err := h.repo.ListByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list"})
	}
	return c.JSON(list)
}

func (h *Handler) Cancel(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := toUint(userIDVal)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.repo.Cancel(userID, uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func toUint(v interface{}) (uint, bool) {
	switch t := v.(type) {
	case int:
		if t < 0 {
			return 0, false
		}
		return uint(t), true
	case int64:
		if t < 0 {
			return 0, false
		}
		return uint(t), true
	case uint:
		return t, true
	case float64:
		if t < 0 {
			return 0, false
		}
		return uint(t), true
	case string:
		if t == "" {
			return 0, false
		}
		n, err := strconv.ParseUint(t, 10, 64)
		if err != nil {
			return 0, false
		}
		return uint(n), true
	default:
		return 0, false
	}
}
