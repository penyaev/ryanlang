package object

type Storage struct {
	objects []Object
	has     map[string]uint16
}

func NewStorage() *Storage {
	return &Storage{
		has: map[string]uint16{},
	}
}
func (s *Storage) Add(obj Object) uint16 {
	if IsHashable(obj) {
		id, ok := s.has[obj.(Hashable).Hash()]
		if ok {
			return id
		}
		s.has[obj.(Hashable).Hash()] = uint16(len(s.objects))
	}
	s.objects = append(s.objects, obj)
	return uint16(len(s.objects) - 1)
}
func (s *Storage) Get(id uint16) (Object, bool) {
	if int(id) >= len(s.objects) {
		return nil, false
	}
	return s.objects[id], true
}
func (s *Storage) GetRef(id uint16) (*Object, bool) {
	if int(id) >= len(s.objects) {
		return nil, false
	}
	return &s.objects[id], true
}
func (s *Storage) Len() uint16 {
	return uint16(len(s.objects))
}
