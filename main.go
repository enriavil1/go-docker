//six levels of isolation
/*
Unix Timesharing System
Process IDs
Mounts
Networks
User IDs
InterProcess Communication

*/

// docker run cmd
//go run main.go cmd args...

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)


func main() {
	switch os.Args[1]{
		case "run":
			run()
		case "child":
			child()
		default:
			panic("Invalid command")

	}
}


// runs the command given by the stdin so essentially the run after the main.go
// go run main.go [run](cmd) [args...]
func run(){
	fmt.Printf("Running %v as %d", os.Args[3:], os.Getegid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// create isolated process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Cloning is what creates the process(container) in which we would be running our command.
		// CLONE_NEWUTS will allow to have our own hostname inside our container by creating a new unix timesharing system.
		// CLONE_NEWPID assigns pids to only process inside the new namspace.
		// CLONE_NEWNS new namespace for mount.
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,

		// prevents sharing of new name space with the host
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child(){
	fmt.Printf("Running %v as %d", os.Args[3:], os.Getegid())

	// sets hostname for newly created namespace
	must(syscall.SetHostname([]byte(os.Args[2])))

	// sets container default directory
	must(syscall.Chroot("/."))

	// mounting proc dir to set the process running inside the container
	must(syscall.Mount("proc", "proc", "proc", 0, ""))


	cmd := exec.Command(os.Args[3], os.Args[4:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

}


func must(err error){
	if err != nil{
		panic(err)
	}
}