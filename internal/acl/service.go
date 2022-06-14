package acl

import (
	"strings"

	"github.com/casbin/casbin/v2"

	"github.com/meBazil/hr-mvp/internal/config"
)

type Service struct {
	enforcer casbin.IEnforcer
}

func NewAclService(cfg config.ACL) (*Service, error) {
	enforcer, err := casbin.NewEnforcer(cfg.ConfigPath, cfg.PolicyPath)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadModel(); err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Service{enforcer: enforcer}, nil
}

func (s *Service) Enforce(roles []string, path string, method string) bool {
	for _, role := range roles {
		if result, _ := s.enforcer.Enforce(strings.ToLower(role), path, method); result {
			return true
		}
	}

	return false
}
