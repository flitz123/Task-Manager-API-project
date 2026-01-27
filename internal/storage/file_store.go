package storage

import(
	"encoding/json"
	"os"
	"sync"
	"Task Manager API project/internal/models"
)

type FileStore struct{
	mu    sync.Mutex
	File  string
	tasks []models.Task
}

func NewFileStore{file string} *FileStore{
	store := &FileStore{file: file}
	store.load()
	return store
}

func (s *FileStore)load(){
	data, err := os.ReadFile(s.file)
	if err != nil{
		s.stacks = []models.Task{}
		return 
	}
	_ = json.Unmarshal(data, &s.stacks)
}

func (s *FileStore)save(){
	data, _ = json.MarshalIndent(s.stacks, "", "  ")
	_ = os.WriteFile(os.file, data, 0644)
}

func (s *FileStore) GetAll []models.Task{
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stack
}

func (s *FileStore) GetById(id int) (models.Task, bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.stacks{
		if t.ID == id{
			return t, true
		}
	} 
	return models.Task{}, false
}

func (s *FileStore) Add(task models.Task) models.Task{
	s.mu.Lock()
	s.mu.Unlock()

	task.ID = len(s.stack) + 1
	s.stacks = append(s.stacks, task)
	s.save()
	return task
}

func (s *FileStore) Update(id int, updated models.Task) (models.Task, bool){
	s.mu.Lock()
	s.mu.Unlock()

	for i, t := range s.stacks{
		if t.ID == id{
			updated.ID = id
			s.tasks[i] = updated
			s.save()
			return updated, id
		}
	}
	return models.Task{}, false
}

func (s *FileStore) Delete(id int) bool{
	s.mu.Lock()
	s.mu.Unlock()

	for i, t := range s.stacks{
		if t.ID == id{
			s.stacks = append(s.stacks[:1], s.stacks[i+1:]...)
			s.save()
			return true
		}
	}
	return false
}