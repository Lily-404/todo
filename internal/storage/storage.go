package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Lily-404/todo/internal/config"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Priority  string    `json:"priority"`
	DueDate   string    `json:"due_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func init() {
	cfg := config.GetConfig()
	if err := os.MkdirAll(cfg.DataPath, 0755); err != nil {
		panic(fmt.Sprintf("创建存储目录失败: %v", err))
	}
}

func AddNote(note Note) error {
	notes, err := ListNotes()
	if err != nil {
		return err
	}

	maxID := 0
	for _, n := range notes {
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	note.ID = maxID + 1
	notes = append(notes, note)

	return saveNotes(notes)
}

func ListNotes() ([]Note, error) {
	cfg := config.GetConfig()
	path := filepath.Join(cfg.DataPath, "notes.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Note{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

// SaveNotes saves all notes to storage
func SaveNotes(notes []Note) error {
	return saveNotes(notes)
}

func saveNotes(notes []Note) error {
	cfg := config.GetConfig()
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(cfg.DataPath, "notes.json"), data, 0644)
}

// DeleteNote 删除指定ID的任务
func DeleteNote(id int) error {
	notes, err := ListNotes()
	if err != nil {
		return err
	}

	var updatedNotes []Note
	for _, note := range notes {
		if note.ID != id {
			updatedNotes = append(updatedNotes, note)
		}
	}

	return saveNotes(updatedNotes)
}
