package MiniPM

import (
	"MiniPm/messenger"
	"MiniPm/survey"
	"MiniPm/task"
	"MiniPm/user"
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Worker struct {

	// dependencies
	userRepository user.IUserRepository
	taskRepository task.ITaskRepository

	// inner state
	lastNotificationDay int
	surveyProcessor     *survey.SurveyProcessor
}

func StartWorker(
	ctx context.Context,
	userRepository user.IUserRepository,
	taskRepository task.ITaskRepository,
	messageProvider messenger.IMessengerProvider,
	cfg *Config,
) {
	surveyProcessor := survey.NewSurveyProcessor(taskRepository, userRepository, messageProvider)

	w := &Worker{
		userRepository:  userRepository,
		taskRepository:  taskRepository,
		surveyProcessor: surveyProcessor,
	}

	spec := "0 9 * * *"
	if cfg != nil && cfg.DefaultSchedule != "" {
		spec = cfg.DefaultSchedule
	}
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, _ := cronParser.Parse(spec)

	notificationTicker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-notificationTicker.C:
				if w.isNotificationTime(schedule) {
					w.doSurvey()
				}
			case <-ctx.Done():
				notificationTicker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) isNotificationTime(schedule cron.Schedule) bool {
	now := time.Now()

	next := schedule.Next(now)

	if next.Sub(now) <= 1*time.Minute {

		currentDay := now.Day()
		if currentDay != w.lastNotificationDay {
			w.lastNotificationDay = currentDay
			return true
		}
	}

	return false
}

func (w *Worker) doSurvey() {
	users, err := w.userRepository.GetAll()
	if err != nil {
		logrus.Errorf("Error getting users: %v", err)
		return
	}

	for _, u := range users {
		err := w.surveyProcessor.DoSurveyForUser(u.ID)
		if err != nil {
			logrus.Errorf("Error surveying tasks: %v", err)
		}
	}
}
