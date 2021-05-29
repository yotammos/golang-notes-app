package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"notes-app/model"
	"notes-app/service"
	"strconv"
)
import "net/http"
import _ "github.com/hashicorp/go-memdb"

func hello() (int, string) {
	return http.StatusOK, "Hello world!"
}

func notFound() (int, string) {
	return http.StatusNotFound, "resource not found"
}

func main() {
	var notesService service.NotesService
	notesService = service.CreateService()
	m := martini.Classic()

	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		// Use indentations. &pretty=1
		pretty, _ := strconv.ParseBool(r.FormValue("pretty"))
		// Use null instead of empty object for json &null=1
		null, _ := strconv.ParseBool(r.FormValue("null"))
		// Some content negotiation
		switch r.Header.Get("Content-Type") {
			case "application/xml":
				c.MapTo(encoder.XmlEncoder{PrettyPrint: pretty}, (*encoder.Encoder)(nil))
				w.Header().Set("Content-Type", "application/xml; charset=utf-8")
			default:
				c.MapTo(encoder.JsonEncoder{PrettyPrint: pretty, PrintNull: null}, (*encoder.Encoder)(nil))
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
	})
	m.Get("/", hello)
	m.Get("/notes", func(enc encoder.Encoder) (int, []byte) {
		res, noteIds := notesService.ListNotes()
		if res != nil {
			return res.StatusCode, encoder.Must(enc.Encode(model.Message{
				Content: res.Message,
			}))
		}

		if noteIds == nil || len(noteIds) == 0 {
			return http.StatusNotFound, encoder.Must(enc.Encode(model.Message{
				Content: "resource not found",
			}))
		}

		ids, marshalingError := json.Marshal(noteIds)
		if marshalingError != nil {
			return http.StatusBadRequest, encoder.Must(enc.Encode(model.Message{
				Content: "failed building json",
			}))
		}
		return http.StatusOK, ids
	})
	m.Get("/notes/:id", func(c martini.Context, params martini.Params) (int, string) {
		res, note := notesService.DescribeNote(params["id"])
		if res != nil {
			return res.StatusCode, res.Message
		}
		jsonNote, marshalingError := json.Marshal(note)
		if marshalingError != nil {
			return http.StatusBadRequest, "failed building json"
		}
		return http.StatusOK, string(encoder.Must(jsonNote, nil))
	})
	m.Put("/notes", func(req *http.Request) (int, string) {
		var request model.CreateNoteRequest
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			return http.StatusBadRequest, err.Error()
		}
		res := notesService.CreateNote(request)
		return res.StatusCode, res.Message
	})
	m.NotFound(notFound)
	m.Run()
}
