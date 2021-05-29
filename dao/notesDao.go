package dao

import (
	"notes-app/model"
)

type NotesDao interface {
	CreateNote(request model.CreateNoteRequest) *model.Response
	DescribeNote(id string) (*model.Response, *model.Note)
	ListNotes() (*model.Response, []string)
}
