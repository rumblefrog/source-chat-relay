package entity

func (entity *Entity) Insert() error {
	_, err := entity.CreateEntity()

	if err != nil {
		return err
	}

	WriteCache(entity)

	return nil
}

func (entity *Entity) Delete() error {
	_, err := entity.QDelete()

	if err != nil {
		return err
	}

	Cache.Lock()
	defer Cache.Unlock()

	delete(Cache.Entities, entity.ID)

	return nil
}

func (entity *Entity) SetReceiveChannels(channels []int) error {
	entity.ReceiveChannels = channels

	return entity.Propagate()
}

func (entity *Entity) SetSendChannels(channels []int) error {
	entity.SendChannels = channels

	return entity.Propagate()
}

func (entity *Entity) SetDisplayName(name string) error {
	entity.DisplayName = name

	return entity.Propagate()
}

func (entity *Entity) Propagate() error {
	_, err := entity.UpdateEntity()

	if err != nil {
		return err
	}

	WriteCache(entity)

	return nil
}
