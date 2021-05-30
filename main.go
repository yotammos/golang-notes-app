package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"notes-app/model"
	"notes-app/service"
)
import "net/http"
import _ "github.com/hashicorp/go-memdb"

func buildGinResponse(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, obj)
}

func listNotes(notesService service.NotesService) (code int, obj interface{}) {
	err, noteIds := notesService.ListNotes()
	if err != nil {
		return err.StatusCode, gin.H{
			"message": err.Message,
		}
	}
	if noteIds == nil || len(noteIds) == 0 {
		return http.StatusNotFound, gin.H{
			"message": "resource not found",
		}
	}
	return http.StatusOK, noteIds
}

func createNote(c *gin.Context, notesService service.NotesService) (code int, obj interface{}) {
	data, err := c.GetRawData()
	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message": "failed getting request body",
		}
	}

	var request model.CreateNoteRequest
	err = json.Unmarshal(data, &request)
	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		}
	}

	response := notesService.CreateNote(request)
	return response.StatusCode, response.Message
}

func describeNote(c *gin.Context, notesService service.NotesService) (code int, obj interface{}) {
	id := c.Param("id")
	errorResponse, note := notesService.DescribeNote(id)
	if errorResponse != nil {
		return errorResponse.StatusCode, gin.H{
			"message": errorResponse.Message,
		}
	}

	return http.StatusOK, note
}

func startGin(notesService service.NotesService) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		buildGinResponse(c, http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/notes", func(c *gin.Context) {
		code, data := listNotes(notesService)
		buildGinResponse(c, code, data)
	})
	r.PUT("/notes", func(c *gin.Context) {
		code, data := createNote(c, notesService)
		buildGinResponse(c, code, data)
	})
	r.GET("/notes/:id", func(c *gin.Context) {
		code, data := describeNote(c, notesService)
		buildGinResponse(c, code, data)
	})
	r.Run()
}

func main() {
	var notesService service.NotesService
	notesService = service.CreateService()
	startGin(notesService)
}
