package utils

type Queue []interface{} 

func NewQueue() *Queue {
	return new(Queue)
}

func (q *Queue) Encqueue(a interface{}) {
	*q = append(*q, a)
}

func (q *Queue) Dequeue() (interface{}, bool) {
	if len(*q) == 0 {
		return nil, false
	}
	value := (*q)[0]
	*q = (*q)[1:]
	return value, true
}
