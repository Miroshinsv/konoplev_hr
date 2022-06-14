package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/meBazil/hr-mvp/internal/config"
	ptypes "github.com/meBazil/hr-mvp/internal/profile/types"
	helper "github.com/meBazil/hr-mvp/internal/rest/helpers"
)

const (
	expiresAccessTTL  = 24 * 60 * 60 * time.Second
	expiresRefreshTTL = 3 * 24 * 60 * 60 * time.Second
)

var (
	ErrDuplicateEmail  = errors.New("email already registered")
	ErrDuplicateMobile = errors.New("mobile number already registered")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInActiveProfile    = errors.New("profile is blocked")
)

type profileService interface {
	RegisterUser(email string, mobile string, pwd string) (*ptypes.Profile, error)
	GetByEmail(email string) (*ptypes.Profile, error)
	GetByMobile(mobile string) (*ptypes.Profile, error)
	IsPasswordCorrect(prof *ptypes.Profile, pwd string) bool
	Get(ID uint) (*ptypes.Profile, error)
	ChangePassword(profile *ptypes.Profile, pwd string) error
}

type Service struct {
	profileService profileService
	config         config.JWT

	JWTIssuer    *Issuer
	JWTValidator *Validator
}

func NewService(profileService profileService, cfg config.JWT) *Service {
	return &Service{
		profileService: profileService,
		config:         cfg,
		JWTIssuer:      NewIssuer(cfg),
		JWTValidator:   NewValidator(cfg),
	}
}

func (s *Service) AuthByToken(token string) (*ptypes.Profile, error) {
	claims, err := s.JWTValidator.ValidateJWT(token)
	if err != nil {
		return nil, helper.ErrInvalidToken
	}

	if claims.Type != AccessToken {
		return nil, helper.ErrInvalidTokenType
	}

	if claims.Version != s.config.TokenVersion {
		return nil, helper.ErrInvalidTokenVersion
	}

	prof, err := s.profileService.Get(claims.UserID)
	if err != nil {
		return nil, helper.ErrInvalidTokenType
	}

	if !prof.IsActive {
		return nil, ErrInActiveProfile
	}

	return prof, nil
}

func (s *Service) AuthByEmail(email string, pwd string) (*ptypes.Profile, error) {
	prof, err := s.profileService.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !s.profileService.IsPasswordCorrect(prof, pwd) {
		return nil, ErrInvalidCredentials
	}

	if !prof.IsActive {
		return nil, ErrInActiveProfile
	}

	return prof, nil
}

func (s *Service) AuthByMobile(mobile string, pwd string) (*ptypes.Profile, error) {
	prof, err := s.profileService.GetByMobile(mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !s.profileService.IsPasswordCorrect(prof, pwd) {
		return nil, ErrInvalidCredentials
	}

	if !prof.IsActive {
		return nil, ErrInActiveProfile
	}

	return prof, nil
}

func (s *Service) Register(email string, mobile string, pwd string) (*ptypes.Profile, error) {
	_, err := s.profileService.GetByEmail(email)
	if err == nil {
		return nil, ErrDuplicateEmail
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_, err = s.profileService.GetByMobile(email)
	if err == nil {
		return nil, ErrDuplicateMobile
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return s.profileService.RegisterUser(email, mobile, pwd)
}

func (s *Service) GenerateAccessToken(prof *ptypes.Profile) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   UserSubject,
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(expiresAccessTTL).Unix(),
			Issuer:    IssuerName,
		},
		UserID:  prof.ID,
		Version: s.config.TokenVersion,
		Type:    AccessToken,
	}

	return s.JWTIssuer.Issue(claims)
}

func (s *Service) GenerateRefreshToken(prof *ptypes.Profile) (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   UserSubject,
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(expiresRefreshTTL).Unix(),
			Issuer:    IssuerName,
		},
		UserID:  prof.ID,
		Version: s.config.TokenVersion,
		Type:    RefreshToken,
	}

	return s.JWTIssuer.Issue(claims)
}

func (s *Service) GetValidator() *Validator {
	return s.JWTValidator
}

func (s *Service) ChangePassword(prof *ptypes.Profile, pwd, pwdNew string) error {
	if !s.profileService.IsPasswordCorrect(prof, pwd) {
		return ErrInvalidCredentials
	}

	if !prof.IsActive {
		return ErrInActiveProfile
	}

	return s.profileService.ChangePassword(prof, pwdNew)
}

func (s *Service) SetPassword(prof *ptypes.Profile, pwdNew string) error {
	if !prof.IsActive {
		return ErrInActiveProfile
	}

	return s.profileService.ChangePassword(prof, pwdNew)
}

func (s *Service) IsPasswordCorrect(prof *ptypes.Profile, pwd string) bool {
	return s.profileService.IsPasswordCorrect(prof, pwd)
}

func (s *Service) HasRole(prof *ptypes.Profile, role ptypes.Role) bool {
	for _, v := range prof.Roles {
		if ptypes.Role(v) == role {
			return true
		}
	}

	return false
}
