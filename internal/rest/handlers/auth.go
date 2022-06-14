package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	authform "github.com/meBazil/hr-mvp/internal/rest/forms/auth"
	authmodels "github.com/meBazil/hr-mvp/internal/rest/models/auth"

	"github.com/meBazil/hr-mvp/internal/rest/middleware"
	"github.com/meBazil/hr-mvp/internal/rest/response"
)

type Auth struct {
	authService    authService
	profileService profileService
}

func NewAuthHandler(as authService, ps profileService) APIHandler {
	return &Auth{
		authService:    as,
		profileService: ps,
	}
}

func (h *Auth) EnrichRoutes(baseRouter *mux.Router) {
	authRoute := baseRouter.PathPrefix("/auth").Subrouter()
	authRoute.HandleFunc("/sign-up", h.signupAction).Methods(http.MethodPost).Name("public_signup")
	authRoute.HandleFunc("/sign-in", h.signinAction).Methods(http.MethodPost).Name("public_signin")
	authRoute.HandleFunc("/refresh", h.refreshAction).Methods(http.MethodPost).Name("public_refresh")

	authRoute.HandleFunc("/email", h.changeEmailAction).Methods(http.MethodPut).Name("change_email")
	authRoute.HandleFunc("/pwd", h.changePasswordAction).Methods(http.MethodPut).Name("change_password")
}

func (h *Auth) signupAction(w http.ResponseWriter, r *http.Request) {
	form, verr := authform.NewRegisterForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)

		return
	}

	prof, err := h.authService.Register(
		form.(*authform.RegisterForm).Email,
		form.(*authform.RegisterForm).Mobile,
		form.(*authform.RegisterForm).Password,
	)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	result, err := h.getTokenPair(prof)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Auth) signinAction(w http.ResponseWriter, r *http.Request) {
	form, verr := authform.NewSigninForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)

		return
	}

	prof, err := h.authService.AuthByMobile(
		form.(*authform.SigninForm).Mobile,
		form.(*authform.SigninForm).Password,
	)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	result, err := h.getTokenPair(prof)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Auth) refreshAction(w http.ResponseWriter, r *http.Request) {
	form, verr := authform.NewRefreshForm(h.authService.GetValidator()).ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)

		return
	}

	prof, err := h.profileService.Get(form.(*authform.RefreshForm).Claims.UserID)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	if !prof.IsActive {
		HandleError(response.NewPermissionDeniedError(), w)
		return
	}

	result, err := h.getTokenPair(prof)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *Auth) changeEmailAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)

	form, verr := authform.NewChangeEmailForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)

		return
	}

	prof, err := h.profileService.Get(user.ID)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	if !prof.IsActive {
		HandleError(response.NewPermissionDeniedError(), w)
		return
	}

	if !h.authService.IsPasswordCorrect(prof, form.(*authform.ChangeEmailForm).Password) {
		HandleError(response.NewPermissionDeniedError(), w)
		return
	}

	prof.Email = form.(*authform.ChangeEmailForm).EmailNew
	err = h.profileService.Update(prof)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Auth) changePasswordAction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.ContextUserField).(*ptypes.Profile)

	form, verr := authform.NewChangePasswordForm().ParseAndValidate(r)
	if verr != nil {
		HandleError(verr, w)

		return
	}

	prof, err := h.profileService.Get(user.ID)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	if !prof.IsActive {
		HandleError(response.NewPermissionDeniedError(), w)
		return
	}

	err = h.authService.ChangePassword(prof, form.(*authform.ChangePasswordForm).Password, form.(*authform.ChangePasswordForm).PasswordNew)
	if err != nil {
		HandleError(response.ResolveError(err), w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Auth) getTokenPair(prof *ptypes.Profile) (*authmodels.JWT, error) {
	access, err := h.authService.GenerateAccessToken(prof)
	if err != nil {
		return nil, err
	}

	refresh, err := h.authService.GenerateRefreshToken(prof)
	if err != nil {
		return nil, err
	}

	return &authmodels.JWT{
		Access:  access,
		Refresh: refresh,
	}, nil
}
