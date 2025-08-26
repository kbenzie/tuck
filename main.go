package main

import "tuck/cmd"

// Override variables with when building:
// $ go build -ldflags -o tuck \
//		"-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" \
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)
	cmd.Execute()
}
