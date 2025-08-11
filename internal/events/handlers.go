package events

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) List(c *fiber.Ctx) error {
	// Optional filters: location, from, to (RFC3339 or YYYY-MM-DD)
	var filters ListFilters
	if loc := c.Query("location"); loc != "" {
		filters.Location = loc
	}
	if from := c.Query("from"); from != "" {
		if t, err := parseTime(from); err == nil {
			filters.From = &t
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid from"})
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := parseTime(to); err == nil {
			// Set end of day if only date
			if len(to) == len("2006-01-02") {
				tt := t.Add(24 * time.Hour).Add(-time.Nanosecond)
				filters.To = &tt
			} else {
				filters.To = &t
			}
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid to"})
		}
	}

	var (
		list []Event
		err  error
	)
	if filters.Location != "" || filters.From != nil || filters.To != nil {
		list, err = h.repo.ListFiltered(filters)
	} else {
		list, err = h.repo.List()
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list events"})
	}
	return c.JSON(list)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	ev, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(ev)
}

type upsertEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartsAt    string `json:"starts_at"`
	Price       int64  `json:"price"`
	Capacity    int    `json:"capacity"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req upsertEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	startsAt, err := parseTime(req.StartsAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid starts_at"})
	}
	ev := &Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartsAt:    startsAt,
		Price:       req.Price,
		Capacity:    req.Capacity,
	}
	if err := h.repo.Create(ev); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create"})
	}
	return c.Status(fiber.StatusCreated).JSON(ev)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	ev, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	var req upsertEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if req.Title != "" {
		ev.Title = req.Title
	}
	if req.Description != "" {
		ev.Description = req.Description
	}
	if req.Location != "" {
		ev.Location = req.Location
	}
	if req.StartsAt != "" {
		t, err := parseTime(req.StartsAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid starts_at"})
		}
		ev.StartsAt = t
	}
	if req.Price != 0 {
		ev.Price = req.Price
	}
	if req.Capacity != 0 {
		ev.Capacity = req.Capacity
	}
	if err := h.repo.Update(ev); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update"})
	}
	return c.JSON(ev)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
