package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	"github.com/meBazil/hr-mvp/internal/rest/conv"
	vforms "github.com/meBazil/hr-mvp/internal/rest/forms/vacancy"
	"github.com/meBazil/hr-mvp/internal/rest/middleware"
	"github.com/meBazil/hr-mvp/internal/rest/response"
	vtypes "github.com/meBazil/hr-mvp/internal/vacancy/types"
)

type Vacancy struct {
	authService    authService
	vacancyService vacancyService
}

func NewVacancyHandler(vs vacancyService, as authService) *Vacancy {
	return &Vacancy{
		vacancyService: vs,
		authService:    as,
	}
}

func (h *Vacancy) EnrichRoutes(baseRouter *mux.Router) {
	profileRoutes := baseRouter.PathPrefix("/vacancy").Subrouter()
	profileRoutes.HandleFunc("", h.createAction).Methods(http.MethodPost).Name("vacancy_create")
	profileRoutes.HandleFunc("/{id}", h.getAction).Methods(http.MethodGet).Name("vacancy_get")
	profileRoutes.HandleFunc("/{id}", h.updateAction).Methods(http.MethodPut).Name("vacancy_update")
	profileRoutes.HandleFunc("/{id}", h.deleteAction).Methods(http.MethodDelete).Name("vacancy_delete")
}

func (h *Vacancy) createAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)

	form, verr := vforms.NewCreateForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)
		return
	}

	f := form.(*vforms.CreateForm)

	v := &vtypes.Vacancy{
		Title:       f.Title,
		Address:     f.Address,
		Lat:         f.Lat,
		Long:        f.Long,
		Description: f.Description,
		ProfileID:   user.ID,
		Profile:     user,
	}

	result, err := h.vacancyService.Create(v)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	_ = json.NewEncoder(w).Encode(conv.ConvertVacancy(result))
}

func (h *Vacancy) updateAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)
	vacancyID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		HandleError(response.NewNotFoundError(), w)
		return
	}

	form, verr := vforms.NewUpdateForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)
		return
	}

	existedVacancy, err := h.vacancyService.Get(uint(vacancyID))
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	if existedVacancy.ProfileID != user.ID && h.authService.HasRole(user, ptypes.RoleAdmin) {
		HandleError(response.NewPermissionDeniedError(), w)
		return
	}

	f := form.(*vforms.UpdateForm)
	if f.Title != nil {
		existedVacancy.Title = *f.Title
	}

	if f.Address != nil {
		existedVacancy.Address = *f.Address
	}

	if f.Description != nil {
		existedVacancy.Description = *f.Description
	}

	if f.Lat != nil {
		existedVacancy.Lat = *f.Lat
	}

	if f.Long != nil {
		existedVacancy.Long = *f.Long
	}

	if err := h.vacancyService.Update(existedVacancy); err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	_ = json.NewEncoder(w).Encode(conv.ConvertVacancy(existedVacancy))
}

func (h *Vacancy) getAction(w http.ResponseWriter, r *http.Request) {
	vacancyID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		HandleError(response.NewNotFoundError(), w)
		return
	}

	vac, err := h.vacancyService.Get(uint(vacancyID))
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	_ = json.NewEncoder(w).Encode(conv.ConvertVacancy(vac))
}

func (h *Vacancy) deleteAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)
	vacancyID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		HandleError(response.NewNotFoundError(), w)
		return
	}

	vac, err := h.vacancyService.Get(uint(vacancyID))
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	if vac.ProfileID != user.ID && h.authService.HasRole(user, ptypes.RoleAdmin) {
		HandleError(response.NewPermissionDeniedError(), w)
		return
	}

	h.vacancyService.Delete(vac)

	w.WriteHeader(http.StatusOK)
}
