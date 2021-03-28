package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"notes-app/model"
	"notes-app/service"
)
import "net/http"
import _ "github.com/hashicorp/go-memdb"

func hello() (int, string) {
	return 200, "Hello world!"
}

func notFound() (int, string) {
	return 404, "resource not found"
}

func main() {
	var notesService service.NotesService
	notesService = service.CreateService()
	m := martini.Classic()
	m.Get("/", hello)
	m.Get("/notes", func() (int, string) {
		res, noteIds := notesService.ListNotes()
		if res != nil {
			return res.StatusCode, res.Message
		}

		if noteIds == nil || len(noteIds) == 0 {
			return 404, "resource not found"
		}

		ids, marshalingError := json.Marshal(noteIds)
		if marshalingError != nil {
			return 500, "failed building json"
		}
		return 200, string(ids)
	})
	m.Get("/notes/:id", func(params martini.Params) (int, string) {
		res, note := notesService.DescribeNote(params["id"])
		if res != nil {
			return res.StatusCode, res.Message
		}
		jsonNote, marshalingError := json.Marshal(note)
		if marshalingError != nil {
			return 500, "failed building json"
		}
		return 200, string(jsonNote)
	})
	m.Put("/notes", func(c martini.Context, req *http.Request) (int, string) {
		var request model.CreateNoteRequest
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			return 400, "invalid request"
		}
		res := notesService.CreateNote(request)
		return res.StatusCode, res.Message
	})
	m.NotFound(notFound)
	m.Run()
}
