package build

import "fmt"

var (
	Version   string
	GitCommit string
	BuildTime string
)

func Print() {
	fmt.Println("Version: ", Version)
	fmt.Println("GitCommit: ", GitCommit)
	fmt.Println("BuildTime: ", BuildTime)
}
