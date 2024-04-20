package entity

type Consumer struct {
	Location
}

func (c Consumer) GetType() NodeType {
	return NodeTypeConsumer
}

func (c Consumer) GetId() int {
	return c.Id
}
