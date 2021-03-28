package service

import (
	"notes-app/dao"
	"notes-app/model"
	"notes-app/util"
)

type INotesService interface {
	ListNotes() (*util.Response, []string)
	DescribeNote(id string) (*util.Response, *model.Note)
	CreateNote(request model.CreateNoteRequest) *util.Response
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

func (service NotesService) ListNotes() (*util.Response, []string) {
	return service.database.ListNotes()
}

func (service NotesService) DescribeNote(id string) (*util.Response, *model.Note) {
	return service.database.DescribeNote(id)
}

func (service NotesService) CreateNote(request model.CreateNoteRequest) *util.Response {
	return service.database.CreateNote(request)
}
