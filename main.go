package main

import ctr "go-demo/containerd"

func main() {
	ctr.CheckImageExist("docker.io/library/redis:5.0.9", "arm64")
}
