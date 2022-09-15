package repository

import (
	"context"

	"github.com/qsoulior/schedulemanager/internal/entity"
	"github.com/qsoulior/schedulemanager/pkg/moodle"
)

type PlanDatabase interface {
	AddSchedules(ctx context.Context, group string, schedules ...entity.Schedule) error
	GetSchedules(ctx context.Context, group string) ([]entity.Schedule, error)
	GetLatestSchedule(ctx context.Context, group string) (*entity.Schedule, error)
	GetPlansInfo(ctx context.Context) ([]entity.PlanInfo, error)
	DeactivatePlan(ctx context.Context, group string) error
}

type PlanWeb interface {
	GetFilesInfo(courseId int) ([]moodle.Content, error)
	GetFiles(filesInfo []moodle.Content) ([]entity.File, error)
}
