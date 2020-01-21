// +build linux windows
// +build ignore

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/CrowdSurge/banner"
	"github.com/mattn/go-shellwords"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// K/V list of GOOS/GOARCH
var os_archs = []string{
	"aix/ppc64",
	"android/386",
	"android/amd64",
	"android/arm",
	"android/arm64",
	"darwin/386",
	"darwin/amd64",
	"darwin/arm",
	"darwin/arm64",
	"dragonfly/amd64",
	"freebsd/386",
	"freebsd/amd64",
	"freebsd/arm",
	"js/wasm",
	"linux/386",
	"linux/amd64",
	"linux/arm",
	"linux/arm64",
	"linux/mips",
	"linux/mips64",
	"linux/mips64le",
	"linux/mipsle",
	"linux/ppc64",
	"linux/ppc64le",
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
	"windows/arm",
}

// version
var (
	version   string
	branch    string
	buildnum  string
	builddate string
	buildtime string
)

// GOOS/GOARCH
type GoosArch struct {
	Arch string
	Os   string
}

// Command
type Command struct {
	Cmd     string        `json:"cmd"`
	Env     string        `json:"env"`
	Env_ow  bool          `json:"env_ow"`
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

// Usage
var Usage = func() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Printf("\nExamples:\n")
	fmt.Printf(`%s -t "make" -f="Makefile" -m="build" -osarch="linux/amd64,windows/amd64" %s`, os.Args[0], "\n")
	fmt.Printf(`%s -t "go" -m="build" -osarch="linux/amd64,windows/amd64" -p="-o bin/\$GOOS/\$GOARCH/app -v app.go" %s`, os.Args[0], "\n")
}

type void struct{}

func main() {
	var exec_path string = ""
	var err error
	var cmd_params string
	var cmd Command
	flag.Usage = Usage
	oa_flag := flag.String("osarch", "linux/amd64", "Set GOOS/GOARCH. Use 'all' for build for all OS/ARCHs.")
	ls_flag := flag.Bool("list", false, "Print the list of supported GOOS/GOARCH.")
	bt_flag := flag.String("t", "", "Build tool: 'make' | 'go'.")
	bm_flag := flag.String("m", "build", "Build mode (e.g. 'build' or 'clear').")
	mf_flag := flag.String("f", "", "Path to Makefile (only if -t 'make').")
	p_flag := flag.String("p", "", "Another parameters for 'make'/'go' which should be passed.")
	env_flag := flag.String("e", "", "Environment variables.")
	env_oflag := flag.Bool("eow", false, "Overwriting of environment variables.")
	tm_flag := flag.String("timeout", "24h", "Maximum timeout execution of 'make'/'go'.")
	dbg_flag := flag.Bool("d", false, "Debug output.")
	sim_flag := flag.Bool("s", false, "Perform a simulate mode.")
	ver_flag := flag.Bool("v", false, "Show version.")
	flag.Parse()
	// If empty
	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	// If 'ls_flag'
	if *ls_flag {
		fmt.Printf("- List of supported GOOS/GOARCH:\n")
		for _, oa := range os_archs {
			fmt.Printf("%s\n", oa)
		}
		fmt.Printf("- Total: %v\n", len(os_archs))
		os.Exit(0)
	}
	// Show version.
	if *ver_flag {
		banner.Print("gomakecb")
		fmt.Printf("\nGo Make Crossplatform Builder.\n\n")
		fmt.Printf("Version: %s, branch: %s, build number: %s, build date: %s, build time: %s\n\n", version, branch, buildnum, builddate, buildtime)
		os.Exit(0)
	}
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
			cmd = Command{Cmd: exec_path, Env: *env_flag, Env_ow: *env_oflag, Args: cmd_args, Timeout: tm_exec}
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
				log.Println("Additional parameters are required (see '-p' switch).")
				os.Exit(1)
			}
			cmd = Command{Cmd: exec_path, Env: envs, Env_ow: *env_oflag, Args: cmd_args, Timeout: tm_exec}
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
	if s == "all" {
		for _, e := range os_archs {
			arch := fmt.Sprintf("%s", trim_string_before_s(e, "/"))
			osv := fmt.Sprintf("%s", trim_string_after_s(e, "/"))
			archos := GoosArch{Arch: arch, Os: osv}
			ret = append(ret, archos)
		}
		return ret, nil
	}
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
	w_os_archs := missing(os_archs, archs)
	if len(w_os_archs) != 0 {
		for _, w := range w_os_archs {
			log.Printf("Pair GOOS/GOARCH '%s' is an invalid and will be ignored. \n", w)
		}
	}
	// If ret is an empty...
	if len(ret) == 0 {
		return nil, errors.New("No valid pair GOOS/GOARCH were found.")
	}
	return ret, nil
}

// Missing compares two slices and returns slice of differences
// copied from: https://stackoverflow.com/questions/57446978/comparing-two-slices-for-missing-element
func missing(a, b []string) []string {
	// create map with length of the 'a' slice
	ma := make(map[string]void, len(a))
	diffs := []string{}
	// Convert first slice to map with empty struct (0 bytes)
	for _, ka := range a {
		ma[ka] = void{}
	}
	// find missing values in a
	for _, kb := range b {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}

// Exec command
func exec_command(cmd Command, dbg bool) (res Exec_Command_Output, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.Timeout)
	defer cancel()
	var env []string
	var eco Exec_Command_Output
	args, _ := shellwords.Parse(cmd.Args)
	envs, _ := shellwords.Parse(cmd.Env)
	// args := strings.Split(cmd.Args, " ")
	// envs := strings.Split(cmd.Env, " ")
	_ = envs
	comm := exec.CommandContext(ctx, cmd.Cmd, args...)
	if !cmd.Env_ow {
		env = os.Environ()
	}
	// 'make' is a VERY sensitive for an empty environment variables
	for _, e := range envs {
		if e != "" {
			env = append(env, e)
		}
	}
	comm.Env = env
	// if debug enabled
	if dbg {
		fmt.Printf("Commandline: %s\n\n", cmd.Cmd)
		fmt.Printf("Arguments: %s\n\n", args)
		fmt.Printf("Env variables: %s\n\n", env)
	}
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
	comm.Wait()
	return eco, nil
}

// stub
func stub(i ...interface{}) {
}
