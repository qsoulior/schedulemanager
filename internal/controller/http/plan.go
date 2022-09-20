package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/qsoulior/schedulemanager/internal/usecase"
)

type PlanController struct {
	service usecase.Plan
}

func NewPlanController(service usecase.Plan) *PlanController {
	return &PlanController{service}
}

func (p *PlanController) GetSchedules(c *fiber.Ctx) error {
	group := c.Query("group")

	if c.Query("latest") == "" {
		if schedules, err := p.service.GetSchedules(context.Background(), group); err != nil {
			return fiber.ErrInternalServerError
		} else {
			return c.JSON(schedules)
		}
	}

	if schedule, err := p.service.GetLatestSchedule(context.Background(), group); err != nil {
		return fiber.ErrInternalServerError
	} else {
		return c.JSON(schedule)
	}
}

func (p *PlanController) GetPlansInfo(c *fiber.Ctx) error {
	info, err := p.service.GetPlansInfo(context.Background())
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(info)
}
