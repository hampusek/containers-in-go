package main

import (
	"fmt"
    "io/ioutil"
	"os"
	"os/exec"
    "path/filepath"
    "strconv"
    "syscall"
)

func main() {
    switch os.Args[1] {
    case "run":
        run()
    case "child":
        child()
    default:
        panic("bad command")
    }

}

func run() {
    fmt.Printf("running %v as %d\n", os.Args[2:], os.Getpid())

    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
        UidMappings: []syscall.SysProcIDMap{{
            ContainerID: 0,
            HostID: 1000,
            Size: 1,
        }},
    }

    must(cmd.Run())
}

func child() {
    fmt.Printf("running %v as %d\n", os.Args[2:], os.Getpid())

    //cg()

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    must(syscall.Sethostname([]byte("container")))
    must(syscall.Chroot("/mnt/home/a392673/ubuntufs"))
    must(syscall.Chdir("/"))
    must(syscall.Mount("proc", "proc", "proc", 0, ""))

    must(cmd.Run())

    must(syscall.Unmount("/proc", 0))
}

func cg() {
    cgroups := "/sys/fs/cgroup/"
    pids := filepath.Join(cgroups, "pids")
    os.Mkdir(filepath.Join(pids, "ek"), 0755)
    must(ioutil.WriteFile(filepath.Join(pids, "ek/pids.max"), []byte("20"), 0700))

    // Cleanup
    must(ioutil.WriteFile(filepath.Join(pids, "ek/notify_on_release"), []byte("1"), 0700))
    must(ioutil.WriteFile(filepath.Join(pids, "ek/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
