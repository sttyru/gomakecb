package main

import (
    "fmt"
    "flag"
    "os"
)

var (
    version   string
    branch    string
    buildnum  string
    builddate string
    buildtime string
)

func main(){
    e := flag.Bool("v", false, "Show version")
    flag.Parse();
    if(*e){
        fmt.Printf("%s %s %s %s %s\n", version, branch, buildnum, builddate, buildtime)
        os.Exit(0)
    }
    fmt.Println("Test application");
    os.Exit(0)
}
