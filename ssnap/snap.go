package ssnap

type Snap struct {
	path string
	data any
}

func New(path string, data interface{}) *Snap {
	return &Snap{path: path, data: data}
}

func (s *Snap) Path() string {
	return s.path
}

func (s *Snap) Data() any {
	return s.data
}

func (s *Snap) Save() error {
	return Save(s.path, s.data)
}

func (s *Snap) Load() error {
	return Load(s.path, s.data)
}
