package storage

import (
	"Task_Manager_API_project/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileStore struct {
	mu    sync.Mutex
	File  string
	tasks []models.Task
}

func NewFileStore(file string) *FileStore {
	store := &FileStore{File: file}
	store.load()
	return store
}

func (s *FileStore) load() {
	_ = os.MkdirAll(filepath.Dir(s.File), 0755)
	data, err := os.ReadFile(s.File)
	if err != nil {
		s.tasks = []models.Task{}
		return
	}
	if err := json.Unmarshal(data, &s.tasks); err != nil || s.tasks == nil {
		s.tasks = []models.Task{}
	}
}

func (s *FileStore) save() error {
	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}
	_ = os.MkdirAll(filepath.Dir(s.File), 0755)
	temp, err := os.CreateTemp(filepath.Dir(s.File), "tasks-*.tmp")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	if _, err = temp.Write(data); err != nil {
		temp.Close()
		os.Remove(tempName)
		return err
	}
	if err = temp.Close(); err != nil {
		os.Remove(tempName)
		return err
	}
	if err = os.Rename(tempName, s.File); err != nil {
		os.Remove(tempName)
		return err
	}
	return nil
}

func (s *FileStore) GetAll() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]models.Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

func (s *FileStore) GetByID(id int) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return models.Task{}, false
}

func (s *FileStore) Add(task models.Task) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	maxID := 0
	for _, t := range s.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	task.ID = maxID + 1
	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now
	if task.Priority == "" {
		task.Priority = "medium"
	}
	s.tasks = append(s.tasks, task)
	_ = s.save()
	return task
}

func (s *FileStore) Update(id int, updated models.Task) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			updated.ID = id
			updated.CreatedAt = t.CreatedAt
			updated.UpdatedAt = time.Now().UTC()
			if updated.Priority == "" {
				updated.Priority = t.Priority
			}
			s.tasks[i] = updated
			s.save()
			return updated, true
		}
	}
	return models.Task{}, false
}

func (s *FileStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			s.save()
			return true
		}
	}
	return false
}
