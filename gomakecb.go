package main

import (
	// "bytes"
	// "bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"regexp"
	"syscall"
	"time"
	"github.com/mattn/go-shellwords"
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
	Os   string
}

// Command
type Command struct {
	Cmd     string        `json:"cmd"`
	Env     string        `json:"env"`
	Args    string        `json:"args"`
	Timeout time.Duration `json:"timeout"`
}

// Exec command result
type Exec_Command_Output struct {
	Error        error
	Output       []byte
	StdOutput    []byte
	StdErrOutput []byte
	Proc         *os.Process
}

func main() {
	var exec_path string = ""
	var err error
	var cmd_params string
	var cmd Command

	oa_flag := flag.String("osarch", "linux/amd64", "K/V list with GOOS/GOARCH.")
	bt_flag := flag.String("t", "", "Build tool: 'make' | 'go'.")
	bm_flag := flag.String("m", "build", "Build mode (e.g. 'build' or 'clear').")
	mf_flag := flag.String("f", "", "Path to Makefile (only if -t 'make').")
	p_flag := flag.String("p", "", "Another parameters for 'make'/'go' which should be passed.")
	env_flag := flag.String("e", "", "Environment variables.")
	tm_flag := flag.String("timeout", "24h", "Maximum timeout execution of 'make'/'go'.")
	dbg_flag := flag.Bool("d", false, "Debug output.")
	sim_flag := flag.Bool("s", false, "Perform a simulate mode.")
	flag.Parse()
	// Buildtool
	switch *bt_flag {
	// Make
	case "make":
		exec_path, err = exec.LookPath("make")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if *mf_flag == "" {
			log.Println("Required path to Makefile (see '-f' switch)")
			os.Exit(1)
		}
		err = file_exist(string(*mf_flag))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		cmd_params = fmt.Sprintf("-f %s", *mf_flag)
	// go build
	case "go":
		exec_path, err = exec.LookPath("go")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	default:
		log.Println("Unknown build mode.")
		os.Exit(1)
	}
	// Define build mode flag
	if *bm_flag == "" {
		log.Println("Required build mode (see '-m' switch)")
		os.Exit(1)
	}
	// Define timeout (if present)
	tm_exec, err := time.ParseDuration(*tm_flag)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Get ARCH/OS
	archos, err := get_arch_os(*oa_flag)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// If simulateus mode...
	if *sim_flag {
		fmt.Printf("(!) Running in simulate mode\n")
	}
	if *dbg_flag {
		for _, v := range archos {
			fmt.Printf("GOOS: %s, GOARCH: %s\n", v.Os, v.Arch)
		}
	}
	switch *bt_flag {
	case "make":
		for c, v := range archos {
			var cmd_args string
			if *p_flag != "" {
				// Add regexp for cases when used variables
				rg_goos, _ := regexp.Compile(`.GOOS`)
				rg_goarch, _ := regexp.Compile(`.GOARCH`)
				rep_goos := rg_goos.ReplaceAll([]byte(*p_flag), []byte(v.Os))
				rep_goarch := rg_goarch.ReplaceAll(rep_goos, []byte(v.Arch))
				prs_args := rep_goarch
				cmd_args = fmt.Sprintf("%s %s GOARCH=%s GOOS=%s %s", *bm_flag, cmd_params, v.Arch, v.Os, prs_args)
			} else {
				cmd_args = fmt.Sprintf("%s %s GOARCH=%s GOOS=%s", *bm_flag, cmd_params, v.Arch, v.Os)
			}
			cmd = Command{Cmd: exec_path, Env: *env_flag, Args: cmd_args, Timeout: tm_exec}
			if !*sim_flag {
				cmd_output, err := exec_command(cmd, *dbg_flag)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				fmt.Printf("\n%s", string(cmd_output.Output))
			} else {
				fmt.Printf("%v Cmd: %s %s %s\n", c, cmd.Env, cmd.Cmd, cmd.Args)
			}
		}
	case "go":
		for c, v := range archos {
			var cmd_args string
			envs := fmt.Sprintf("GOOS=%s GOARCH=%s %s", v.Os, v.Arch, *env_flag)
			// Add regexp for cases when used variables
			rg_goos, _ := regexp.Compile(`.GOOS`)
			rg_goarch, _ := regexp.Compile(`.GOARCH`)
			rep_goos := rg_goos.ReplaceAll([]byte(*p_flag), []byte(v.Os))
			rep_goarch := rg_goarch.ReplaceAll(rep_goos, []byte(v.Arch))
			prs_args := rep_goarch
			if *p_flag != "" {
				cmd_args = fmt.Sprintf("%s %s", *bm_flag, string(prs_args))
			} else {
				log.Println("Additional parameters are required (see '-p' switch).");
				os.Exit(1)
			}
			cmd = Command{Cmd: exec_path, Env: envs, Args: cmd_args, Timeout: tm_exec}
			if !*sim_flag {
				cmd_output, err := exec_command(cmd, *dbg_flag)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				fmt.Printf("\n%s", string(cmd_output.Output))
			} else {
				fmt.Printf("%v Cmd: %s %s %s\n", c, cmd.Env, cmd.Cmd, cmd.Args)
			}
		}
	}
	cmd = Command{Cmd: exec_path, Args: cmd_params, Timeout: tm_exec}
	os.Exit(0)
	stub(exec_path, archos)
}

// Check for file is exist
func file_exist(filepath string) (err error) {
	finfo, _ := os.Stat(filepath)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("'%s' not found.", filepath))
	}
	if finfo.IsDir() {
		return errors.New(fmt.Sprintf("'%s' is the directory (not a file).", filepath))
	} else {
	}
	return nil
}

// Trim string AFTER a symbol ("linux/amd64" -> "linux")
func trim_string_after_s(s string, x string) (r string) {
	if idx := strings.Index(s, x); idx != -1 {
		return s[:idx]
	}
	return s
}

// Trim string BEFORE a symbol ("linux/amd64" -> "amd64")
func trim_string_before_s(s string, x string) (r string) {
	if idx := strings.LastIndex(s, x); idx != -1 {
		return s[idx+1:]
	}
	return s
}

// Parse ARCH/OS
func get_arch_os(s string) (ret []GoosArch, err error) {
	archs := strings.Split(s, ",")
	for _, e := range os_archs {
		for _, ao := range archs {
			if ao == e {
				arch := fmt.Sprintf("%s", trim_string_before_s(e, "/"))
				osv := fmt.Sprintf("%s", trim_string_after_s(e, "/"))
				archos := GoosArch{Arch: arch, Os: osv}
				ret = append(ret, archos)
			}
		}
	}
	if len(ret) == 0 {
		return nil, errors.New("No valid pair OS/ARCH were found.")
	}
	return ret, nil
}

// Exec command
func exec_command(cmd Command, dbg bool) (res Exec_Command_Output, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
	defer cancel()
	var eco Exec_Command_Output
	args, _ := shellwords.Parse(cmd.Args)
	envs, _ := shellwords.Parse(cmd.Env)
	// args := strings.Split(cmd.Args, " ")
	// envs := strings.Split(cmd.Env, " ")
	_ = envs
	comm := exec.CommandContext(ctx, cmd.Cmd, args...)
	env := os.Environ()
	// 'make' is a VERY sensitive for an empty environment variables
	for _, e := range envs {
		if(e != ""){
			env = append(env, e)
		}
	}
	comm.Env = env
	// if debug enabled
	if dbg {
		fmt.Printf("Commandline: %s\n\n", cmd.Cmd)
		fmt.Printf("Arguments: %s\n\n", args)
		fmt.Printf("Environment: %s\n\n", env)
	}
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
		Proc:         comm.Process}
	pgid, err := syscall.Getpgid(comm.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
	}
	comm.Wait()
	return eco, nil
}

// stub
func stub(i ...interface{}) {
}
