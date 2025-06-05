package squeue

// GetHeadItem deprecated
func (q *Queue) GetHeadItem() (any, error) {
	return q.HeadItem()
}

// GetTailItem deprecated
func (q *Queue) GetTailItem() (any, error) {
	return q.TailItem()
}

// IsEmpty deprecated
func (q *Queue) IsEmpty() bool {
	return q.Empty()
}

// CopyItems deprecated
func (q *Queue) CopyItems() []any {
	return q.CopyList()
}
