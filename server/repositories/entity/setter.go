package entity

func (entity *Entity) Insert() error {
	_, err := entity.CreateEntity()

	if err != nil {
		return err
	}

	Cache.Controller <- entity

	return nil
}

func (entity *Entity) SetReceiveChannels(channels []int) error {
	entity.ReceiveChannels = channels

	_, err := entity.UpdateChannels()

	if err != nil {
		return err
	}

	Cache.Controller <- entity

	return nil
}

func (entity *Entity) SetSendChannels(channels []int) error {
	entity.SendChannels = channels

	_, err := entity.UpdateChannels()

	if err != nil {
		return err
	}

	Cache.Controller <- entity

	return nil
}
