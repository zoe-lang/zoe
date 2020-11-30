package main

func init() {
	handlers["shutdown"] = HandleShutdown
	handlers["exit"] = HandleExit
}

func HandleShutdown(req *LspRequest) error {
	// FIXME should probably do some cleanup...
	req.Conn.receivedShutdown = true
	req.Reply(nil)
	return nil
}

func HandleExit(req *LspRequest) error {
	// if req.Conn.receivedShutdown {
	// 	os.Exit(0)
	// }
	// os.Exit(1)
	return nil
}
