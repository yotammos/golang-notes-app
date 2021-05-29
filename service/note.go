package service

import (
	"notes-app/dao"
	"notes-app/model"
)

type INotesService interface {
	ListNotes() (*model.Response, []string)
	DescribeNote(id string) (*model.Response, *model.Note)
	CreateNote(request model.CreateNoteRequest) *model.Response
}

type NotesService struct {
	database dao.NotesDao
}

func CreateService() NotesService {
	var database dao.NotesDao = dao.CreateDao()
	return NotesService{
		database: database,
	}
}

func (service NotesService) ListNotes() (*model.Response, []string) {
	return service.database.ListNotes()
}

func (service NotesService) DescribeNote(id string) (*model.Response, *model.Note) {
	return service.database.DescribeNote(id)
}

func (service NotesService) CreateNote(request model.CreateNoteRequest) *model.Response {
	return service.database.CreateNote(request)
}
