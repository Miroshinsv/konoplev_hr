package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	"github.com/meBazil/hr-mvp/internal/rest/conv"
	"github.com/meBazil/hr-mvp/internal/rest/forms/profile"
	"github.com/meBazil/hr-mvp/internal/rest/middleware"
	"github.com/meBazil/hr-mvp/internal/rest/response"
)

type Profile struct {
	profileService profileService
}

func NewProfileHandler(profileService profileService) *Profile {
	return &Profile{
		profileService: profileService,
	}
}

func (h *Profile) EnrichRoutes(baseRouter *mux.Router) {
	profileRoutes := baseRouter.PathPrefix("/me").Subrouter()
	profileRoutes.HandleFunc("", h.meAction).Methods(http.MethodGet).Name("profile")
	profileRoutes.HandleFunc("", h.updateAction).Methods(http.MethodPut).Name("profile_update")

	profileAdminRoutes := baseRouter.PathPrefix("/profile").Subrouter()
	profileAdminRoutes.HandleFunc("/{id}/roles", h.updateRoles).Methods(http.MethodPut).Name("profile_admin_roles")
}

func (h *Profile) meAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)
	prof, err := h.profileService.Get(user.ID)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	result := conv.ConvertProfile(prof)

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Profile) updateAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)
	prof, err := h.profileService.Get(user.ID)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	form, verr := profile.NewUpdateForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)
		return
	}

	uform := form.(*profile.UpdateForm)
	if uform.Name != nil {
		prof.Name = *uform.Name
	}

	if uform.MiddleName != nil {
		prof.MiddleName = *uform.MiddleName
	}

	if uform.SureName != nil {
		prof.SureName = *uform.SureName
	}

	err = h.profileService.Update(prof)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	result := conv.ConvertProfile(prof)

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Profile) updateRoles(w http.ResponseWriter, r *http.Request) {
	targetID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		HandleError(response.NewNotFoundError(), w)
		return
	}

	prof, err := h.profileService.Get(uint(targetID))
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	form, verr := profile.NewUpdateRoleForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)
		return
	}

	prof.Roles = form.(*profile.UpdateRoleForm).Roles

	err = h.profileService.Update(prof)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	result := conv.ConvertProfile(prof)

	_ = json.NewEncoder(w).Encode(result)
}
