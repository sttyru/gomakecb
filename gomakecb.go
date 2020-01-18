package main

import (
    // "bytes"
    // "bufio"
    "context"
    "io/ioutil"
    "syscall"
    "flag"
    "errors"
    "log"
    "os"
    "os/exec"
    "strings"
    "fmt"
    "time"
)

// K/V list of GOOS/GOARCH
var os_archs = []string{
    "darwin/386",
    "darwin/amd64",
    "dragonfly/amd64",
    "freebsd/386",
    "freebsd/amd64",
    "freebsd/arm",
    "linux/386",
    "linux/amd64",
    "linux/arm",
    "linux/arm64",
    "linux/ppc64",
    "linux/ppc64le",
    "linux/mips",
    "linux/mipsle",
    "linux/mips64",
    "linux/mips64le",
    "linux/s390x",
    "nacl/386",
    "nacl/amd64p32",
    "nacl/arm",
    "netbsd/386",
    "netbsd/amd64",
    "netbsd/arm",
    "openbsd/386",
    "openbsd/amd64",
    "openbsd/arm",
    "plan9/386",
    "plan9/amd64",
    "plan9/arm",
    "solaris/amd64",
    "windows/386",
    "windows/amd64",
}

// GOOS/GOARCH
type GoosArch struct {
    Arch string
    Osv  string
}

// Command
type Command struct {
	Cmd      string        `json:"cmd"`
	Args     string        `json:"args"`
	Timeout  time.Duration `json:"timeout"`
}

// Exec command result
type Exec_Command_Output struct {
	Error     error
	Output       []byte
	StdOutput    []byte
	StdErrOutput []byte
	Proc     *os.Process
}

func main(){
    var exec_path string = ""
    var err error
    var mf_path string = ""
    var archos GoosArch
    var archos_a []GoosArch
    var cmd Command

    oa_flag := flag.String("osarch", "linux/amd64", "K/V list with GOOS/GOARCH");
    bm_flag := flag.String("mode", "", "Build mode: 'make' | 'go'");
    mf_flag := flag.String("f", "", "Path to Makefile");
    flag.Parse();
    switch *bm_flag {
    case "make":
        exec_path, err = exec.LookPath("make")
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
	if(*mf_flag == ""){
		log.Println("Required path to Makefile (see '-f' switch)");
		os.Exit(1)
	}
        err = file_exist(string(*mf_flag))
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
	mf_path = *mf_flag
    case "go":
        exec_path, err = exec.LookPath("go")
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
    default:
        log.Println("Unknown build mode.");
        os.Exit(1)
    }
    archs := strings.Split(*oa_flag, ",")
    for _, e := range os_archs {
        for _, ao := range archs {
            if(ao == e){
                arch := fmt.Sprintf("%s", trim_string_after_s(e, "/"))
                osv := fmt.Sprintf("%s", trim_string_before_s(e, "/"))
                archos := GoosArch{ Arch: arch, Osv: osv }
                archos_a = append(archos_a, archos)
            }
        }
    }
    if(len(archos_a) == 0){
        log.Println("No valid pair OS/ARCH were found.");
        os.Exit(1);
    }
    for _, v := range archos_a{
        fmt.Printf("GOOS: %s, GOARCH: %s\n", v.Arch, v.Osv)
    }
    def_args := fmt.Sprintf("-f %s", mf_path);
    cmd = Command{Cmd: exec_path, Args: def_args, Timeout: 5000000000000 }
    cmd_output, err := exec_command(cmd)
    if err != nil {
	    log.Println(err);
	    os.Exit(1);
    }
    fmt.Printf("\n%s", string(cmd_output.Output))
    os.Exit(0);
    stub(exec_path, mf_path, archos)
}

// Trim string AFTER a symbol ("linux/amd64" -> "linux")
func trim_string_after_s(s string, x string)(r string) {
	if idx :=strings.Index(s, x); idx != -1 {
		return s[:idx]
	}
	return s
}

// Trim string BEFORE a symbol ("linux/amd64" -> "amd64")
func trim_string_before_s(s string, x string)(r string) {
	if idx := strings.LastIndex(s, x); idx != -1 {
		return s[idx+1:]
	}
	return s
}

// Check for file is exist
func file_exist(filepath string)(err error) {
        finfo, _ := os.Stat(filepath);
        if _, err := os.Stat(filepath); os.IsNotExist(err) {
                // return err
                return errors.New(fmt.Sprintf("'%s' not found.", filepath))
        }
        if finfo.IsDir() {
                return errors.New(fmt.Sprintf("'%s' is the directory (not a file).", filepath))
                } else {
    }
    return nil
}

// Exec command
func exec_command(cmd Command) (res Exec_Command_Output, err error) {
	/* DEBUG */
	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout*time.Second)
	defer cancel()
	var eco Exec_Command_Output
	args := strings.Split(cmd.Args, " ")
	comm := exec.CommandContext(ctx, cmd.Cmd, args...)
	env := os.Environ()
	comm.Env = env
	comm.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	comm_out, err := comm.StdoutPipe()
	comm_stderr, err := comm.StderrPipe()
	err = comm.Start()
	if err != nil {
		return eco, err
	}
	comm_out_std, _ := ioutil.ReadAll(comm_out)
	comm_std_err, _ := ioutil.ReadAll(comm_stderr)
	var output []byte
	output = append(comm_out_std, comm_std_err...)
	eco = Exec_Command_Output{Error: err,
		Output:       output,
		StdOutput:    comm_out_std,
		StdErrOutput: comm_std_err,
		Proc:     comm.Process,}
	pgid, err := syscall.Getpgid(comm.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
	}
	comm.Wait()
	return eco, nil
}

// stub
func stub(i...interface{})(){
}
