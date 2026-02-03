package main

import (
	"runtime/debug"
	"tuck/cmd"
)

// Override variables with when building:
// $ go build -ldflags -o tuck \
//     "-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			var vcsRevision string
			var vcsModified bool
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					vcsRevision = setting.Value
				case "vcs.modified":
					vcsModified = setting.Value == "true"
				case "vcs.time":
					date = setting.Value
				}
			}
			if vcsRevision != "" {
				commit = vcsRevision
				if vcsModified {
					commit += "-dirty"
				}
			}
		}
	}

	cmd.SetVersion(version, commit, date)
	cmd.Execute()
}
