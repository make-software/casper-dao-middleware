package di

import (
	"casper-dao-middleware/internal/persistence"
)

type EntityManagerAware struct {
	entityManager persistence.EntityManager
}

func (a *EntityManagerAware) SetEntityManager(manager persistence.EntityManager) {
	a.entityManager = manager
}

func (a *EntityManagerAware) GetEntityManager() persistence.EntityManager {
	return a.entityManager
}
