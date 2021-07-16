package solver

type nodeListElement struct {
	node     Node
	config   string
	artefact string
	step     int
}

type NodeList struct {
	list []nodeListElement
}

func (n *NodeList) AddNode(config string, artefact string, node Node) {
	idx := 0

	latest := n.getLatestElement(config, artefact)
	if latest != nil {
		idx = latest.step + 1
	}

	n.list = append(n.list, nodeListElement{
		node:     node,
		config:   config,
		artefact: artefact,
		step:     idx,
	})
}

func (n *NodeList) GetList() []Node {
	nodes := []Node{}

	for _, node := range n.list {
		nodes = append(nodes, node.node)
	}

	return nodes
}

func (n *NodeList) getLatestElement(config string, artefact string) *nodeListElement {
	idx := -1
	var latest *nodeListElement

	for _, node := range n.list {
		if node.config == config && node.artefact == artefact && node.step > idx {
			idx = node.step
			latest = &node
		}
	}

	return latest
}

func (n *NodeList) GetLatest(config string, artefact string) *Node {
	latest := n.getLatestElement(config, artefact)

	if latest == nil {
		return nil
	}

	return &latest.node
}
