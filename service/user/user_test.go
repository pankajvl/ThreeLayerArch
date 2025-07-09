package user

import (
	"ThreeLayerArch/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		mockReturnID  int
		mockReturnErr error
		expectResult  *models.User
		expectErr     error
		mockStoreCall bool
	}{
		{
			name:          "Empty name",
			input:         "",
			expectResult:  nil,
			expectErr:     errors.New("name cannot be empty"),
			mockStoreCall: false,
		},
		{
			name:          "Store error",
			input:         "Alice",
			mockReturnID:  0,
			mockReturnErr: errors.New("store error"),
			expectResult:  nil,
			expectErr:     errors.New("store error"),
			mockStoreCall: true,
		},
		{
			name:          "Success",
			input:         "Alice",
			mockReturnID:  1,
			mockReturnErr: nil,
			expectResult:  &models.User{UserID: 1, Name: "Alice"},
			expectErr:     nil,
			mockStoreCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := NewMockUserStore(ctrl)
			svc := &Service{Store: mockStore}
			ctx := &gofr.Context{}

			if tt.mockStoreCall {
				mockStore.EXPECT().
					Create(ctx, tt.input).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}

			result, err := svc.CreateUser(ctx, tt.input)
			assert.Equal(t, tt.expectResult, result)

			if tt.expectErr != nil {
				assert.EqualError(t, err, tt.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		inputID     int
		mockReturn  *models.User
		mockError   error
		expectUser  *models.User
		expectError error
	}{
		{
			name:        "User not found",
			inputID:     1,
			mockReturn:  nil,
			mockError:   nil,
			expectUser:  nil,
			expectError: nil,
		},
		{
			name:        "Store error",
			inputID:     2,
			mockReturn:  nil,
			mockError:   errors.New("db error"),
			expectUser:  nil,
			expectError: errors.New("db error"),
		},
		{
			name:        "Success",
			inputID:     3,
			mockReturn:  &models.User{UserID: 3, Name: "Bob"},
			mockError:   nil,
			expectUser:  &models.User{UserID: 3, Name: "Bob"},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := NewMockUserStore(ctrl)
			svc := &Service{Store: mockStore}
			ctx := &gofr.Context{}

			mockStore.EXPECT().
				GetByID(ctx, tt.inputID).
				Return(tt.mockReturn, tt.mockError)

			user, err := svc.GetUserByID(ctx, tt.inputID)
			assert.Equal(t, tt.expectUser, user)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_View_Users(t *testing.T) {
	tests := []struct {
		name        string
		mockReturn  []models.User
		mockError   error
		expectUsers []models.User
		expectError error
	}{
		{
			name:        "Store error",
			mockReturn:  nil,
			mockError:   errors.New("fetch error"),
			expectUsers: nil,
			expectError: errors.New("fetch error"),
		},
		{
			name: "Success",
			mockReturn: []models.User{
				{UserID: 1, Name: "Alice"},
				{UserID: 2, Name: "Bob"},
			},
			mockError: nil,
			expectUsers: []models.User{
				{UserID: 1, Name: "Alice"},
				{UserID: 2, Name: "Bob"},
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := NewMockUserStore(ctrl)
			svc := &Service{Store: mockStore}
			ctx := &gofr.Context{}

			mockStore.EXPECT().
				ViewUsers(ctx).
				Return(tt.mockReturn, tt.mockError)

			users, err := svc.View_Users(ctx)
			assert.Equal(t, tt.expectUsers, users)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
