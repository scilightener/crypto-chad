package app

import (
	"crypto-chad-lib/rsa"
	"fmt"
	"log/slog"
	"sync"
)

type App struct {
	logger *slog.Logger
	keys   *rsa.Keys
	users  map[string]*rsa.PublicKey
	mu     *sync.RWMutex
}

func NewApp(logger *slog.Logger) *App {
	logger.Info("generating app rsa key pair")
	keys := rsa.NewKeys()
	logger.Info("app rsa key pair is generated")
	logger.Info("app rsa public key",
		slog.String("e", keys.PublicKey.E.String()),
		slog.String("n", keys.PublicKey.N.String()))
	return &App{
		logger: logger,
		keys:   keys,
		users:  make(map[string]*rsa.PublicKey),
		mu:     &sync.RWMutex{},
	}
}

func (s *App) IssueCert(name string) (*rsa.Keys, error) {
	s.logger.Info("user issued an rsa key pair", slog.String("user name", name))
	if _, ok := s.users[name]; ok {
		return nil, fmt.Errorf("user %s already issued its certificate", name)
	}
	keys := rsa.NewKeys()
	s.logger.Info("rsa key pair generated successfully",
		slog.String("user name", name),
		slog.String("e", keys.PublicKey.E.String()),
		slog.String("n", keys.PublicKey.N.String()))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[name] = keys.PublicKey
	return keys, nil
}

func (s *App) ForgetUser(name string) bool {
	s.logger.Info("forgetting user", slog.String("user name", name))
	if _, ok := s.users[name]; ok {
		delete(s.users, name)
		return true
	}
	return false
}

func (s *App) RetrieveCert(name string) (*rsa.PublicKey, error) {
	s.logger.Info("retrieving certificate", slog.String("user name", name))
	s.mu.RLock()
	defer s.mu.RUnlock()
	if k, ok := s.users[name]; ok {
		return k, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (s *App) Keys() rsa.Keys {
	return *s.keys
}
