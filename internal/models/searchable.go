package models

func (n Note) GetID() int64       { return n.ID }
func (n Note) GetContent() string { return n.Title + " " + n.Body }
func (n Note) GetType() string    { return "note" }

func (t Todo) GetID() int64       { return t.ID }
func (t Todo) GetContent() string { return t.Title + " " + t.Description }
func (t Todo) GetType() string    { return "todo" }
