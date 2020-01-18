package main

import (
    // "bytes"
    // "bufio"
    "flag"
    "errors"
    "log"
    "os"
    "os/exec"
    "strings"
    "fmt"
)

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

type ArchosMap struct {
    Arch string
    Osv  string
}

func main(){
    var exec_path string = ""
    var err error
    var mf_path string = ""
    _ = mf_path
    oa_flag := flag.String("oa", "linux/amd64", "K/V list GOOS/GOARCH");
    bm_flag := flag.String("mode", "", "Build mode: 'make' | 'go'");
    mf_flag := flag.String("mf", "", "Path to Makefile");
    flag.Parse();
    switch *bm_flag {
    case "make":
        exec_path, err = exec.LookPath("make")
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }
        err = file_exist(string(*mf_flag))
        if err != nil {
            log.Println(err)
            os.Exit(1)
            }

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
    _ = exec_path
    _ =  mf_path
    /*
    fmt.Printf("Path to %s file: %s\n", *bm_flag, exec_path)
    fmt.Printf("Makefile path: %s\n", string(*mf_flag))
    */
    archs := strings.Split(*oa_flag, ",")
    archos   := ArchosMap{}
    archos_a := []ArchosMap{}
    _ = archos
    for _, e := range os_archs {
        for _, ao := range archs {
            if(ao == e){
                arch := fmt.Sprintf("%s", trim_string_after_s(e, "/"))
                osv := fmt.Sprintf("%s", trim_string_before_s(e, "/"))
                archos := ArchosMap{ Arch: arch, Osv: osv }
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
    os.Exit(0);
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
                return errors.New(fmt.Sprintf("'%s' has not found", filepath))
        }
        if finfo.IsDir() {
                return errors.New(fmt.Sprintf("'%s' is the directory (not a file).", filepath))
                } else {
    }
    return nil
}


