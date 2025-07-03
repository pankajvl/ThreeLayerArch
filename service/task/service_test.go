package task

import (
	"ThreeLayerArch/models"
	"errors"
	"gofr.dev/pkg/gofr"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestAddTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockTaskStore(ctrl)
	service := New(mockStore)

	ctx := &gofr.Context{}

	t.Run("Success", func(t *testing.T) {
		mockStore.EXPECT().AddTask(ctx, "New Task").Return(true, nil)
		ok, err := service.Add_Task(ctx, "New Task")
		if !ok || err != nil {
			t.Errorf("Expected success, got err: %v", err)
		}
	})

	t.Run("Failure - Empty Task", func(t *testing.T) {
		ok, err := service.Add_Task(ctx, "")
		if ok || err == nil {
			t.Errorf("Expected error for empty task, got ok=%v err=%v", ok, err)
		}
	})
}

func TestViewTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockTaskStore(ctrl)
	service := New(mockStore)
	ctx := &gofr.Context{}
	t.Run("Success", func(t *testing.T) {
		mockStore.EXPECT().ViewTask(ctx).Return([]models.Tasks{{Tid: 1, Task: "Test", Completed: false}}, nil)
		tasks, err := service.View_Task(ctx)
		if err != nil || len(tasks) != 1 {
			t.Errorf("Expected 1 task, got %v, err: %v", len(tasks), err)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		mockStore.EXPECT().ViewTask(ctx).Return(nil, errors.New("db error"))
		_, err := service.View_Task(ctx)
		if err == nil {
			t.Errorf("Expected error from ViewTask")
		}
	})
}

func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockTaskStore(ctrl)
	service := New(mockStore)
	ctx := &gofr.Context{}
	t.Run("Found", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 1).Return(true)
		mockStore.EXPECT().GetByID(ctx, 1).Return(models.Tasks{Tid: 1, Task: "Test", Completed: false}, nil)

		taskResult, err := service.Get_By_ID(ctx, 1)
		if err != nil || taskResult.Tid != 1 {
			t.Errorf("Expected valid task, got: %v, err: %v", taskResult, err)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 99).Return(false)
		_, err := service.Get_By_ID(ctx, 99)
		if err == nil {
			t.Errorf("Expected error for missing ID")
		}
	})

	t.Run("GetByID Error", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 2).Return(true)
		mockStore.EXPECT().GetByID(ctx, 2).Return(models.Tasks{}, errors.New("db error"))

		_, err := service.Get_By_ID(ctx, 2)
		if err == nil {
			t.Errorf("Expected error from GetByID")
		}
	})
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockTaskStore(ctrl)
	service := New(mockStore)
	ctx := &gofr.Context{}
	t.Run("Success", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 1).Return(true)
		mockStore.EXPECT().UpdateTask(ctx, 1).Return(true, nil)
		ok, err := service.Update_Task(ctx, 1)
		if !ok || err != nil {
			t.Errorf("Expected success, got err: %v", err)
		}
	})

	t.Run("Failure - Update Error", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 99).Return(true)
		mockStore.EXPECT().UpdateTask(ctx, 99).Return(false, errors.New("update failed"))
		_, err := service.Update_Task(ctx, 99)
		if err == nil {
			t.Errorf("Expected error for invalid ID")
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 100).Return(false)
		_, err := service.Update_Task(ctx, 100)
		if err == nil {
			t.Errorf("Expected error when task does not exist")
		}
	})
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := NewMockTaskStore(ctrl)
	service := New(mockStore)
	ctx := &gofr.Context{}
	t.Run("Success", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 1).Return(true)
		mockStore.EXPECT().DeleteTask(ctx, 1).Return(true, nil)
		ok, err := service.Delete_Task(ctx, 1)
		if !ok || err != nil {
			t.Errorf("Expected success, got err: %v", err)
		}
	})

	t.Run("Failure - Delete Error", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 99).Return(true)
		mockStore.EXPECT().DeleteTask(ctx, 99).Return(false, errors.New("delete failed"))
		_, err := service.Delete_Task(ctx, 99)
		if err == nil {
			t.Errorf("Expected error for invalid ID")
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		mockStore.EXPECT().CheckIfExists(ctx, 100).Return(false)
		_, err := service.Delete_Task(ctx, 100)
		if err == nil {
			t.Errorf("Expected error when task does not exist")
		}
	})
}
