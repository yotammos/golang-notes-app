package dao

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"notes-app/model"
	"notes-app/util"
	"time"
)

type InMemoryDao struct {
	db *memdb.MemDB
}

func CreateDao() *InMemoryDao {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"note": {
				Name: "note",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return &InMemoryDao{
		db: db,
	}
}

func (dao InMemoryDao) StartDb() error {
	dao.db = CreateDao().db
	return nil
}

func (dao InMemoryDao) ListNotes() (*util.Response, []string) {
	transaction := dao.db.Txn(false)
	iter, err := transaction.Get("note", "id")
	if err != nil {
		res := util.Response{
			StatusCode: 500,
			Message:    "failed pulling from db",
		}
		return &res, nil
	}

	var noteIds []string
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		note := raw.(model.Note)
		noteIds = append(noteIds, note.Id)
	}
	return nil, noteIds
}

func (dao InMemoryDao) DescribeNote(id string) (*util.Response, *model.Note) {
	transaction := dao.db.Txn(false)
	raw, err := transaction.First("note", "id", id)
	if err != nil {
		res := util.Response{
			StatusCode: 500,
			Message:    "failed pulling from db",
		}
		return &res, nil
	}
	if raw == nil {
		res := util.Response{
			StatusCode: 404,
			Message:    "note not found",
		}
		return &res, nil
	}
	note := raw.(model.Note)
	return nil, &note
}

func (dao InMemoryDao) CreateNote(request model.CreateNoteRequest) *util.Response {
	uuidGenerator, uuidError := uuid.NewUUID()
	if uuidError != nil {
		return &util.Response{
			StatusCode: 500,
			Message: "failed inserting into db",
		}
	}

	note := model.Note{
		Id: uuidGenerator.String(),
		Message: request.Message,
		CreatedTimestamp: time.Now(),
	}
	transaction := dao.db.Txn(true)
	err := transaction.Insert("note", note)
	transaction.Commit()
	if err != nil {
		return &util.Response{
			StatusCode: 500,
			Message: "failed inserting into db",
		}
	}
	return &util.Response{
		StatusCode: 200,
		Message: "created a note!",
	}
}

