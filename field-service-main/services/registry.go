package services

import (
	"field-service/common/storage"
	"field-service/repositories"
	fieldService "field-service/services/field"
	fieldScheduleService "field-service/services/fieldschedule"
	timeService "field-service/services/time"
)

type Registry struct {
	repository    repositories.IRepositoryRegistry
	storageClient storage.IClient
}

type IServiceRegistry interface {
	GetField() fieldService.IFieldService
	GetFieldSchedule() fieldScheduleService.IFieldScheduleService
	GetTime() timeService.ITimeService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry, storageClient storage.IClient) IServiceRegistry {
	return &Registry{
		repository:    repository,
		storageClient: storageClient,
	}
}

func (r *Registry) GetField() fieldService.IFieldService {
	return fieldService.NewFieldService(r.repository, r.storageClient)
}

func (r *Registry) GetFieldSchedule() fieldScheduleService.IFieldScheduleService {
	return fieldScheduleService.NewFieldScheduleService(r.repository)
}

func (r *Registry) GetTime() timeService.ITimeService {
	return timeService.NewTimeService(r.repository)
}
