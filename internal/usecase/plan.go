package usecase

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/qsoulior/schedulemanager/internal/entity"
	"github.com/qsoulior/schedulemanager/internal/repository"
	"github.com/qsoulior/schedulemanager/pkg/moodle"
	"github.com/qsoulior/scheduleparser"
	"github.com/rs/zerolog"
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
	log *zerolog.Logger
}

func NewPlanService(db repository.PlanDatabase, web repository.PlanWeb, log *zerolog.Logger) Plan {
	return &PlanService{db, web, log}
}

func (s *PlanService) parseFiles(files []entity.File) ([]entity.Plan, error) {
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
			s.log.Error().Err(err).Msg("Failed to parse file")
			return nil, err
		}
	}
	return plans, nil
}

func (s *PlanService) AddSchedules(ctx context.Context, sourceId int) error {
	filesInfo, err := s.web.GetFilesInfo(sourceId)
	if err != nil {
		return err
	}
	s.log.Info().Str("count", strconv.Itoa(len(filesInfo))).Msg("New files info received")

	plansInfo, err := s.db.GetPlansInfo(ctx)
	if err != nil {
		return err
	}
	s.log.Info().Str("count", strconv.Itoa(len(plansInfo))).Msg("Plans info received")

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
		err := s.db.DeactivatePlan(ctx, fileName)
		if err != nil {
			return err
		}
	}
	s.log.Info().Str("count", strconv.Itoa(len(plansInfoMap))).Msg("Old plans deactivated")

	newFiles, err := s.web.GetFiles(newFilesInfo)
	if err != nil {
		return err
	}

	s.log.Info().Str("count", strconv.Itoa(len(newFiles))).Msg("New files received")

	if len(newFiles) > 0 {
		plans, err := s.parseFiles(newFiles)
		if err != nil {
			return err
		}
		for _, plan := range plans {
			if err := s.db.AddSchedules(ctx, plan.Group, plan.Schedules...); err != nil {
				return err
			}
		}
		s.log.Info().Str("count", strconv.Itoa(len(plans))).Msg("New plans saved")
	}

	return nil
}

func (s *PlanService) GetSchedules(ctx context.Context, group string) ([]entity.Schedule, error) {
	schedules, err := s.db.GetSchedules(ctx, group)
	if err != nil {
		s.log.Error().Err(err).Str("group", group).Msg("Failed to get schedules")
		return nil, err
	}
	return schedules, nil
}

func (s *PlanService) GetLatestSchedule(ctx context.Context, group string) (*entity.Schedule, error) {
	schedule, err := s.db.GetLatestSchedule(ctx, group)
	if err != nil {
		s.log.Error().Err(err).Str("group", group).Msg("Failed to get latest schedule")
		return nil, err
	}
	return schedule, nil
}

func (s *PlanService) GetPlansInfo(ctx context.Context) ([]entity.PlanInfo, error) {
	plansInfo, err := s.db.GetPlansInfo(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get plans info")
		return nil, err
	}
	return plansInfo, nil
}
