package main

import (
	"github.com/ca-gip/dploy/cmd"
)

func main() {

	//var TagsRe = regexp.MustCompile("([\\w-.\\/]+)([,]|)")
	//
	//for _, tag := range TagsRe.FindAllStringSubmatch("multipass-create-v_ms,os-metadata,os-provisionning-dns-add", -1) {
	//	fmt.Println(tag)
	//}

	cmd.Execute()
}
