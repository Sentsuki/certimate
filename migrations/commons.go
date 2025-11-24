package migrations

type mWorkflowGraphWalker struct {
	visitors []mWorkflowNodeVisitor
}

type mWorkflowNodeVisitor func(node *mWorkflowNode) (bool, error)

type mWorkflowNode struct {
	Id     string           `json:"id"`
	Type   string           `json:"type"`
	Data   map[string]any   `json:"data"`
	Blocks []*mWorkflowNode `json:"blocks,omitempty,omitzero"`
}

func (w *mWorkflowGraphWalker) Define(visitor mWorkflowNodeVisitor) {
	if w.visitors == nil {
		w.visitors = make([]mWorkflowNodeVisitor, 0)
	}
	w.visitors = append(w.visitors, visitor)
}

func (w *mWorkflowGraphWalker) Visit(nodes []*mWorkflowNode) (bool, error) {
	changed := false

	if w.visitors == nil {
		return changed, nil
	}

	for _, node := range nodes {
		for _, visitor := range w.visitors {
			nodeChanged, err := visitor(node)
			if err != nil {
				return changed, err
			}
			if nodeChanged {
				changed = true
			}

			if len(node.Blocks) > 0 {
				blocksChanged, err := w.Visit(node.Blocks)
				if err != nil {
					return changed, err
				}
				if blocksChanged {
					changed = true
				}
			}
		}
	}

	return changed, nil
}
