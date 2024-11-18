package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hoyci/ecom/types"
)

type mockUserStore struct {
	existingUser *types.User
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if m.existingUser != nil && m.existingUser.Email == email {
		return m.existingUser, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}

func TestUser(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "john",
			LastName:  "doe",
			Email:     "invalid",
			Password:  "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "john",
			LastName:  "doe",
			Email:     "valid@email.com",
			Password:  "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should fail if the user already exists", func(t *testing.T) {
		userStore := &mockUserStore{
			existingUser: &types.User{
				FirstName: "john",
				LastName:  "doe",
				Email:     "johndoe@email.com",
				Password:  "hashedpassword",
			},
		}

		handler := NewHandler(userStore)

		payload := types.RegisterUserPayload{
			FirstName: "john",
			LastName:  "doe",
			Email:     "johndoe@email.com",
			Password:  "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
