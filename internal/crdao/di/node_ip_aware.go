package di

type NodeIPAware struct {
	ip string
}

func (t *NodeIPAware) SetNodeIP(ip string) {
	t.ip = ip
}

func (t *NodeIPAware) GetNodeIP() string {
	return t.ip
}
