package dao

import (
	"notes-app/model"
)
import "notes-app/util"

type NotesDao interface {
	CreateNote(request model.CreateNoteRequest) *util.Response
	DescribeNote(id string) (*util.Response, *model.Note)
	ListNotes() (*util.Response, []string)
}
