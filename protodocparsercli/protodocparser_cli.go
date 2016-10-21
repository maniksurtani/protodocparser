package main

import (
	cli "gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"github.com/maniksurtani/protodocparser"
	"os"
)

func main() {
	protofile := cli.Arg("protofile", "The proto file to fetch the metadata from").Required().String()
	filepath := cli.Flag("source", "Where the proto came from originally").String()
	url := cli.Flag("url", "The url for the proto file").String()
	sha := cli.Flag("sha", "The git sha for the proto file version").String()
	cli.Version("0.0.1")
	cli.Parse()
	protoFile, _ := os.Open(*protofile)
	pf := &protodocparser.ProtoFile{
		ProtoFileSource: protoFile,
		ProtoFilePath:   *filepath,
		Url:             *url,
		Sha:             *sha}
	output := protodocparser.ParseAsString([]*protodocparser.ProtoFile{pf})
	fmt.Println(output)
}