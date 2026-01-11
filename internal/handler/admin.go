package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"fiozap/internal/database/repository"
	"fiozap/internal/model"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

// ListUsers godoc
// @Summary List users or get single user
// @Description Get all users or a specific user by ID
// @Tags Admin
// @Produce json
// @Param id path string false "User ID"
// @Success 200 {object} model.Response
// @Failure 404 {object} model.Response
// @Security AdminAuth
// @Router /admin/users [get]
// @Router /admin/users/{id} [get]
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id != "" {
		user, err := h.userRepo.GetByID(id)
		if err != nil {
			model.RespondNotFound(w, errors.New("user not found"))
			return
		}
		model.RespondOK(w, user)
		return
	}

	users, err := h.userRepo.GetAll()
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, users)
}

// AddUser godoc
// @Summary Create new user
// @Description Create a new API user
// @Tags Admin
// @Accept json
// @Produce json
// @Param user body model.UserCreateRequest true "User data"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Security AdminAuth
// @Router /admin/users [post]
func (h *AdminHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	if req.Name == "" {
		model.RespondBadRequest(w, errors.New("name is required"))
		return
	}

	if req.Token == "" {
		model.RespondBadRequest(w, errors.New("token is required"))
		return
	}

	user, err := h.userRepo.Create(&req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondCreated(w, map[string]string{"id": user.ID})
}

// EditUser godoc
// @Summary Update user
// @Description Update an existing user
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body model.UserUpdateRequest true "User data"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Security AdminAuth
// @Router /admin/users/{id} [put]
func (h *AdminHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req model.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		model.RespondBadRequest(w, errors.New("invalid payload"))
		return
	}

	user, err := h.userRepo.Update(id, &req)
	if err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Admin
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.Response
// @Failure 500 {object} model.Response
// @Security AdminAuth
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.userRepo.Delete(id); err != nil {
		model.RespondInternalError(w, err)
		return
	}

	model.RespondOK(w, map[string]string{"details": "User deleted successfully"})
}
