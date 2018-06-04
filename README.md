
A POC linux distribution with a user-space comprising a single Go binary (a trivial web server)

`make.sh` runs a vagrant machine that creates the ISO. The Go exe is built into an initramfs and runs as init (process 1). It configures an ethernet port with a hard-wired driver (e1000) and IP address, then listens on port 80.