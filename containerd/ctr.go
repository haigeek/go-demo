package ctr

import (
	"context"
	"fmt"
	"go-demo/utils"
	"log"
	"os/exec"
	"runtime"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// 定义一个包级别的 containerd 客户端对象
var client *containerd.Client
var k8sCtx context.Context
var containerdRuntime = "/run/k3s/containerd/containerd.sock"

func init() {
	// 在包初始化期间创建 containerd 客户端
	c, err := containerd.New("")
	// 创建一个容器上下文
	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	if err != nil {
		log.Fatalf("无法连接到 containerd 服务: %v", err)
	}
	client = c
	k8sCtx = ctx

}

// CheckImageExist 检测指定平台的镜像版本是否存在
func CheckImageExist(imageName, platform string) bool {
	//检测镜像是否携带版本
	imageVersion := utils.ExtractVersion(imageName)
	if imageVersion == "" {
		fmt.Printf("检测到镜像名%s未携带版本号，默认使用latest版本\n", imageName)
		imageName = imageName + ":latest"
	}
	//如果未传递platform 默认使用当前机器架构
	//定义架构信息
	imagePlatform := specs.Platform{}

	if platform == "" {
		platform = runtime.GOARCH
	}
	imagePlatform.Architecture = platform
	imagePlatform.OS = "linux"

	//检查镜像platfrom是否和当前环境一致,如果不一致，创建对应架构的client
	if platform != runtime.GOARCH {
		platforms.Any(imagePlatform)
		platformClient, err := containerd.New(containerdRuntime, containerd.WithDefaultPlatform(platforms.Any(imagePlatform)))
		if err != nil {
			log.Fatalf("无法连接到 containerd 服务: %v", err)
		}
		client = platformClient

	}
	// 查找镜像是否存在
	image, _ := client.GetImage(k8sCtx, imageName)

	// 如果镜像不存在则先拉取镜像，docker pull xxx
	if image == nil {
		fmt.Println("镜像[" + imageName + "]不存在，开始拉取...")
		return PullImage(imageName, platform)
	} else {
		//如果镜像存在，检测镜像平台
		//判断现有镜像和目标镜像是否相同
		if image.Platform().Match(imagePlatform) {
			fmt.Println("已经存在的镜像 " + imageName + "为" + platform + "平台镜像")
		} else {
			fmt.Println("开始拉取镜像 " + imageName + " " + platform + " 平台镜像,请确保存在" + platform + "版本镜像，否则会拉取默认amd64平台")
			return PullImage(imageName, platform)
		}
	}
	return true
}

func PullImage(imageName, platform string) bool {
	cmd := exec.Command("ctr", "i", "pull", "--platform", platform, imageName)
	_, err := utils.HandleCommandOutput(cmd)
	if err != nil {
		fmt.Printf("拉取镜像失败:%v\n", err)
		return false
	}
	return true
}
