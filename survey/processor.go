package survey

import (
	"fmt"
	"time"

	"github.com/Koltsov1998/MiniPM/messenger"
	"github.com/Koltsov1998/MiniPM/task"
	"github.com/Koltsov1998/MiniPM/user"

	"github.com/sirupsen/logrus"
)

type SurveyProcessor[U user.User, T task.Task] struct {
	// dependencies
	taskRepository    task.ITaskRepository[U, T]
	userRepository    user.IUserRepository[U]
	messengerProvider messenger.IMessengerProvider
}

func NewSurveyProcessor[U user.User, T task.Task](
	taskRepository task.ITaskRepository[U, T],
	userRepository user.IUserRepository[U],
	messengerProvider messenger.IMessengerProvider,
) *SurveyProcessor[U, T] {
	return &SurveyProcessor[U, T]{
		taskRepository:    taskRepository,
		userRepository:    userRepository,
		messengerProvider: messengerProvider,
	}
}

func (s *SurveyProcessor[U, T]) DoSurveyForUser(user U) error {
	tasks, err := s.taskRepository.GetAllInProgressForUser(user)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		go func() {
			chatMessage := s.formatSurveyMessage(t)
			responseChan, err := s.messengerProvider.SendMessage(user.GetId(), chatMessage)
			if err != nil {
				logrus.Errorf("Error sending message: %v", err)
				return
			}
			tChan := time.After(20 * time.Hour)

			defer close(responseChan)

			select {
			case response := <-responseChan:
				logrus.Infof("Got survey response from user: %s", response)
				err = s.taskRepository.WriteTaskReport(t, user, response)
				if err != nil {
					logrus.Errorf("Error writing t report: %v", err)
				}
			case <-tChan:
				logrus.Warningf("Timeout waiting for response")
				err := s.messengerProvider.SendMessageWithoutResponse(user.GetId(), "Sorry, I didn't get your response")
				if err != nil {
					logrus.Errorf("Error sending message: %v", err)
				}
			}
		}()
	}

	return nil
}

func (s *SurveyProcessor[U, T]) formatSurveyMessage(task task.Task) string {
	return fmt.Sprintf("How is your progress on task: %s?", task.GetTitle())
}
