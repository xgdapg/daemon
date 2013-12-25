## Usage
```go
package main

import "github.com/xgdapg/daemon"

func init() {
	daemon.Exec(daemon.Daemon) // send the process to the background
	daemon.Exec(daemon.Monitor) // keep the process running
	daemon.Exec(daemon.Daemon | daemon.Monitor) // both of the above
}
```