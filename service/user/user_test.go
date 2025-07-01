package user

import (
	Model "ThreeLayerArch/models"
	"errors"
	"testing"
)

type mockStore struct{}

func (m mockStore) Create(name string) (int, error) {
	return 0, nil
}

func (m mockStore) GetByID(id int) (*Model.User, error) {
	switch id {
	case 1:
		return &Model.User{UserID: 1, Name: "Ram"}, nil
	case 2:
		return &Model.User{UserID: 2, Name: "Shyam"}, nil
	default:
		return nil, errors.New("not found")
	}
}

func (m mockStore) ViewUsers() ([]Model.User, error) {
	return nil, nil
}

func TestGetByID(t *testing.T) {
	svc := Service{Store: mockStore{}}

	res, err := svc.GetUserByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.UserID != 1 || res.Name != "Ram" {
		t.Errorf("expected ID 1 and Name Ram, got ID: %d, Name: %s", res.UserID, res.Name)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := Service{Store: mockStore{}}

	res, err := svc.GetUserByID(999)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if res != nil {
		t.Errorf("expected nil user, got: %+v", res)
	}
}
