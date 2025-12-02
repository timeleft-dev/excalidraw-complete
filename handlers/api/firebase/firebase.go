package firebase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
)

type (
	BatchGetRequest struct {
		Documents []string `json:"documents"`
	}
	BatchGetEmptyResponse struct {
		Missing  string `json:"missing"`
		ReadTime string `json:"readTime"`
	}

	FoundInfoResponse struct {
		Name       string      `json:"name"`
		Fields     interface{} `json:"fields"`
		CreateTime string      `json:"createTime"`
		UpdateTime string      `json:"updateTime"`
	}
	BatchGetExistsResponse struct {
		Found    FoundInfoResponse `json:"found"`
		ReadTime string            `json:"readTime"`
	}

	UpdateRequest struct {
		Name   string      `json:"name"`
		Fields interface{} `json:"fields"`
	}
	WriteRequest struct {
		Update UpdateRequest `json:"update"`
	}
	BatchCommitRequest struct {
		Writes []WriteRequest `json:"writes"`
	}

	WriteResult struct {
		UpdateTime string `json:"updateTime"`
	}
	BatchCommitResponse struct {
		WriteResults []WriteResult `json:"writeResults"`
		CommitTime   string        `json:"commitTime"`
	}
)

var (
	db   *sql.DB
	once sync.Once
)

func getDB() *sql.DB {
	once.Do(func() {
		var err error
		dataSourceName := os.Getenv("DATA_SOURCE_NAME")
		if dataSourceName == "" {
			dataSourceName = ":memory:"
		}
		db, err = sql.Open("sqlite3", dataSourceName)
		if err != nil {
			panic(fmt.Sprintf("failed to open database: %v", err))
		}
		// Create table for Firebase documents
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS firebase_documents (
			name TEXT PRIMARY KEY,
			fields TEXT NOT NULL,
			create_time TEXT NOT NULL,
			update_time TEXT NOT NULL
		)`)
		if err != nil {
			panic(fmt.Sprintf("failed to create firebase_documents table: %v", err))
		}
	})
	return db
}

func (body *BatchGetRequest) Bind(r *http.Request) (err error) {
	return nil
}
func (body *BatchCommitRequest) Bind(r *http.Request) (err error) {
	return nil
}
func HandleBatchCommit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "project_id")
		databaseId := chi.URLParam(r, "database_id")
		_ = projectId
		_ = databaseId

		data := &BatchCommitRequest{}
		// Seems like requests is text/plain but content is json ...
		if err := render.DecodeJSON(r.Body, data); err != nil {
			fmt.Println(err)
			render.Status(r, http.StatusBadRequest)
			return
		}

		// Serialize fields to JSON
		fieldsJSON, err := json.Marshal(data.Writes[0].Update.Fields)
		if err != nil {
			fmt.Printf("failed to marshal fields: %v\n", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		now := time.Now().Format(time.RFC3339)
		documentName := data.Writes[0].Update.Name

		// Insert or update the document
		db := getDB()
		_, err = db.Exec(`INSERT INTO firebase_documents (name, fields, create_time, update_time)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(name) DO UPDATE SET
				fields = excluded.fields,
				update_time = excluded.update_time`,
			documentName, string(fieldsJSON), now, now)
		if err != nil {
			fmt.Printf("failed to save document: %v\n", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, BatchCommitResponse{
			CommitTime: now,
			WriteResults: []WriteResult{
				WriteResult{UpdateTime: now},
			},
		})
		render.Status(r, http.StatusOK)
		return

	}
}

func HandleBatchGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		projectId := chi.URLParam(r, "project_id")
		databaseId := chi.URLParam(r, "database_id")
		fmt.Printf("Got %v and %v\n", projectId, databaseId)
		data := &BatchGetRequest{}

		// Seems like requests is text/plain but content is json ...
		if err := render.DecodeJSON(r.Body, data); err != nil {
			fmt.Println(err)
			render.Status(r, http.StatusBadRequest)
			return
		}
		key := data.Documents[0]
		fmt.Printf("Got key %v \n", key)

		// Query the database
		db := getDB()
		var fieldsJSON string
		var createTime, updateTime string
		err := db.QueryRow(`SELECT fields, create_time, update_time FROM firebase_documents WHERE name = ?`, key).
			Scan(&fieldsJSON, &createTime, &updateTime)

		if err == sql.ErrNoRows {
			fmt.Println("missing key")
			render.JSON(w, r, []BatchGetEmptyResponse{BatchGetEmptyResponse{
				Missing:  key,
				ReadTime: time.Now().Format(time.RFC3339),
			}})
			render.Status(r, http.StatusOK)
			return
		} else if err != nil {
			fmt.Printf("database error: %v\n", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		// Deserialize fields from JSON
		var fields interface{}
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			fmt.Printf("failed to unmarshal fields: %v\n", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}

		fmt.Println("existing key")
		render.JSON(w, r, []BatchGetExistsResponse{BatchGetExistsResponse{
			Found: FoundInfoResponse{
				Name:       key,
				Fields:     fields,
				CreateTime: createTime,
				UpdateTime: updateTime,
			},
			ReadTime: time.Now().Format(time.RFC3339),
		}})
		render.Status(r, http.StatusOK)
		return

	}
}