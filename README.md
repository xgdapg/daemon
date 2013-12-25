## Usage
```go
import "github.com/xgdapg/daemon"

daemon.Daemon() // send the process to the background
daemon.Monitor() // keep the process running
daemon.DaemonAndMonitor() // both of the above
```
**Do *NOT* Use Monitor Before Daemon! EVER!**