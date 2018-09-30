package entity

func init() {
	Cache = &EntityCache{
		Entities:   make(map[string]*Entity),
		Controller: make(chan *Entity),
	}

	go Start()
}
