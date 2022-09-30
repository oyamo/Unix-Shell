package src

type HistoryQueue struct {
	// The history queue
	Queue       []string
	CurrenIndex int
}

type histQueue interface {
	AddToQueue(cmd string)
	GetNext() string
	GetPrevious() string
	GetLast() string
	GetFirst() string
}

func (h *HistoryQueue) AddToQueue(cmd string) {
	h.Queue = append(h.Queue, cmd)
	h.CurrenIndex = len(h.Queue) - 1
}

func (h *HistoryQueue) GetNext() string {
	if h.CurrenIndex == len(h.Queue) {
		return ""
	}

	h.CurrenIndex++
	return h.Queue[h.CurrenIndex]
}

func (h *HistoryQueue) GetPrevious() string {
	if h.CurrenIndex == 0 {
		return ""
	}

	h.CurrenIndex--
	return h.Queue[h.CurrenIndex]
}

func (h *HistoryQueue) GetLast() string {
	return h.Queue[len(h.Queue)-1]
}

func (h *HistoryQueue) GetFirst() string {
	return h.Queue[0]
}
