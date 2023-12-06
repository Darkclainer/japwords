package anki

import "strconv"

type NoteID int64

func (id NoteID) String() string {
	if id == 0 {
		return ""
	}
	return strconv.FormatInt(int64(id), 10)
}
