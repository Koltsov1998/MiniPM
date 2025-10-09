package survey

import (
	"fmt"
	"github.com/Koltsov1998/MiniPM/messenger"
	"github.com/Koltsov1998/MiniPM/task"
	"github.com/Koltsov1998/MiniPM/user"
	"time"

	"github.com/sirupsen/logrus"
)

type SurveyProcessor struct {
	// dependencies
	taskRepository    task.ITaskRepository
	userRepository    user.IUserRepository
	messengerProvider messenger.IMessengerProvider
}

func NewSurveyProcessor(
	taskRepository task.ITaskRepository,
	userRepository user.IUserRepository,
	messengerProvider messenger.IMessengerProvider,
) *SurveyProcessor {
	return &SurveyProcessor{
		taskRepository:    taskRepository,
		userRepository:    userRepository,
		messengerProvider: messengerProvider,
	}
}

func (s *SurveyProcessor) DoSurveyForUser(userId user.UserId) error {
	tasks, err := s.taskRepository.GetAllInProgressForUser(userId)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		go func() {
			chatMessage := s.formatSurveyMessage(t)
			responseChan, err := s.messengerProvider.SendMessage(userId, chatMessage)
			if err != nil {
				logrus.Errorf("Error sending message: %v", err)
				return
			}
			tChan := time.After(20 * time.Hour)

			defer close(responseChan)

			select {
			case response := <-responseChan:
				logrus.Infof("Got survey response from user: %s", response)
				err = s.taskRepository.WriteTaskReport(userId, response)
				if err != nil {
					logrus.Errorf("Error writing t report: %v", err)
				}
			case <-tChan:
				logrus.Warningf("Timeout waiting for response")
				err := s.messengerProvider.SendMessageWithoutResponse(userId, "Sorry, I didn't get your response")
				if err != nil {
					logrus.Errorf("Error sending message: %v", err)
				}
			}
		}()
	}

	return nil
}

func (s *SurveyProcessor) formatSurveyMessage(task task.Task) string {
	return fmt.Sprintf("How is your progress on task: %s?", task.Title)
}
