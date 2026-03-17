package storage

import(
	"encoding/json"
	"os"
	"sync"
	"Task_Manager_API_project/internal/models"
)

type FileStore struct{
	mu    sync.Mutex
	File  string
	tasks []models.Task
}

func NewFileStore(file string) *FileStore{
	store := &FileStore{File: file}
	store.load()
	return store
}

func (s *FileStore)load(){
	data, err := os.ReadFile(s.File)
	if err != nil{
		s.tasks = []models.Task{}
		return 
	}
	_ = json.Unmarshal(data, &s.tasks)
}

func (s *FileStore)save(){
	data, _ := json.MarshalIndent(s.tasks, "", "  ")
	_ = os.WriteFile(s.File, data, 0644)
}

func (s *FileStore) GetAll() []models.Task{
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks
}

func (s *FileStore) GetByID(id int) (models.Task, bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.tasks{
		if t.ID == id{
			return t, true
		}
	} 
	return models.Task{}, false
}

func (s *FileStore) Add(task models.Task) models.Task{
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = len(s.tasks) + 1
	s.tasks = append(s.tasks, task)
	s.save()
	return task
}

func (s *FileStore) Update(id int, updated models.Task) (models.Task, bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks{
		if t.ID == id{
			updated.ID = id
			s.tasks[i] = updated
			s.save()
			return updated, true
		}
	}
	return models.Task{}, false
}

func (s *FileStore) Delete(id int) bool{
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks{
		if t.ID == id{
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			s.save()
			return true
		}
	}
	return false
}