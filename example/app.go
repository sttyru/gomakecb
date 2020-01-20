package main

import (
    "fmt"
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
    fmt.Println("Test application");
    os.Exit(0)
}
