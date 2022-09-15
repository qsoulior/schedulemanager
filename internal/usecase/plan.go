package usecase

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/qsoulior/schedulemanager/internal/entity"
	"github.com/qsoulior/schedulemanager/internal/repository"
	"github.com/qsoulior/schedulemanager/pkg/moodle"
	"github.com/qsoulior/scheduleparser"
)

type Plan interface {
	AddSchedules(ctx context.Context, sourceId int) error
	GetSchedules(ctx context.Context, group string) ([]entity.Schedule, error)
	GetLatestSchedule(ctx context.Context, group string) (*entity.Schedule, error)
	GetPlansInfo(ctx context.Context) ([]entity.PlanInfo, error)
}

type PlanService struct {
	db  repository.PlanDatabase
	web repository.PlanWeb
}

func NewPlanService(db repository.PlanDatabase, web repository.PlanWeb) Plan {
	return &PlanService{db, web}
}

func (service *PlanService) parseFiles(files []entity.File) ([]entity.Plan, error) {
	planCh := make(chan *entity.Plan)
	defer close(planCh)
	errCh := make(chan error)
	defer close(errCh)

	for _, file := range files {
		go func(file entity.File) {
			fileData, err := scheduleparser.ParseBytes(file.Data, file.Modified)
			if err != nil {
				errCh <- err
				return
			}

			var events []entity.Event
			if err := json.Unmarshal(fileData, &events); err != nil {
				errCh <- err
				return
			}
			group := strings.TrimSuffix(file.Name, ".pdf")
			schedules := []entity.Schedule{{Modified: file.Modified, Events: events}}

			planCh <- &entity.Plan{Group: group, Active: true, Schedules: schedules}
		}(file)
	}

	plans := make([]entity.Plan, 0)
	for i := 0; i < len(files); i++ {
		select {
		case plan := <-planCh:
			plans = append(plans, *plan)
		case err := <-errCh:
			return nil, err
		}
	}
	return plans, nil
}

func (service *PlanService) AddSchedules(ctx context.Context, sourceId int) error {
	filesInfo, err := service.web.GetFilesInfo(sourceId)
	if err != nil {
		return err
	}
	log.Printf("New files info received: %d\n", len(filesInfo))

	plansInfo, err := service.db.GetPlansInfo(ctx)
	if err != nil {
		return err
	}
	log.Printf("Plans info received: %d\n", len(plansInfo))

	plansInfoMap := make(map[string]time.Time, len(plansInfo))
	for _, planInfo := range plansInfo {
		plansInfoMap[planInfo.Group] = planInfo.Modified
	}

	newFilesInfo := make([]moodle.Content, 0)
	for _, fileInfo := range filesInfo {
		fileName := strings.TrimSuffix(fileInfo.FileName, ".pdf")
		if modified, ok := plansInfoMap[fileName]; (ok && fileInfo.TimeModified > modified.Unix()) || !ok {
			newFilesInfo = append(newFilesInfo, fileInfo)
		}

		delete(plansInfoMap, fileName)
	}

	for fileName := range plansInfoMap {
		err := service.db.DeactivatePlan(ctx, fileName)
		if err != nil {
			return err
		}
	}
	log.Printf("Old plans deactivated: %d\n", len(plansInfoMap))

	newFiles, err := service.web.GetFiles(newFilesInfo)
	if err != nil {
		return err
	}

	log.Printf("New files received: %d\n", len(newFiles))

	if len(newFiles) > 0 {
		plans, err := service.parseFiles(newFiles)
		if err != nil {
			return err
		}
		for _, plan := range plans {
			if err := service.db.AddSchedules(ctx, plan.Group, plan.Schedules...); err != nil {
				return err
			}
		}
		log.Printf("New plans saved: %d\n", len(plans))
	}

	return nil
}
func (service *PlanService) GetSchedules(ctx context.Context, group string) ([]entity.Schedule, error) {
	schedules, err := service.db.GetSchedules(ctx, group)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (service *PlanService) GetLatestSchedule(ctx context.Context, group string) (*entity.Schedule, error) {
	schedule, err := service.db.GetLatestSchedule(ctx, group)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (service *PlanService) GetPlansInfo(ctx context.Context) ([]entity.PlanInfo, error) {
	plansInfo, err := service.db.GetPlansInfo(ctx)
	if err != nil {
		return nil, err
	}
	return plansInfo, nil
}
