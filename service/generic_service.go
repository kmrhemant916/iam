package service

import (
	"github.com/kmrhemant916/iam/repositories"
)

type GenericService[T any] interface {
	Create(entity *T) (error)
	FindOne(entity *T, sqlQuery string, sqlQueryParams ...interface{}) (*T, error)
	FindMany(entities *[]T, sqlQuery string, sqlQueryParams ...interface{}) (error)
}

type genericService[T any] struct {
	genericRepository repositories.GenericRepository[T]
}

func NewGenericService[T any](genericRepository repositories.GenericRepository[T]) *genericService[T] {
	return &genericService[T]{
		genericRepository,
	}
}

func (s *genericService[T]) Create(entity *T) (error) {
	err := s.genericRepository.Create(entity)
	return err
}

func (s *genericService[T]) FindOne(entity *T, sqlQuery string, sqlQueryParams ...interface{}) (*T, error) {
	entity, err := s.genericRepository.FindOne(entity, sqlQuery, sqlQueryParams...)
	return entity, err
}

func (s *genericService[T]) FindMany(entities *[]T, sqlQuery string, sqlQueryParams ...interface{}) (error) {
	err := s.genericRepository.FindMany(entities, sqlQuery, sqlQueryParams...)
	return err
}