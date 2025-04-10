package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

var (
	cyan    = color.New(color.FgCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
)

type Tunnel struct {
	Listener   net.Listener
	Status     string
	LocalPort  string
	RemoteDest string
	StartTime  time.Time
}

type TunnelManager struct {
	SSHClient  *ssh.Client
	Tunnels    map[string]*Tunnel
	SSHConfig  *ssh.ClientConfig
	SSHAddress string
}

func (tm *TunnelManager) CreateTunnel(localPort, remoteDest string) error {
	if _, exists := tm.Tunnels[localPort]; exists {
		return fmt.Errorf("端口 %s 已存在", localPort)
	}

	listener, err := net.Listen("tcp", "localhost:"+localPort)
	if err != nil {
		return fmt.Errorf("端口监听失败: %v", err)
	}

	tunnel := &Tunnel{
		Listener:   listener,
		Status:     "运行中",
		LocalPort:  localPort,
		RemoteDest: remoteDest,
		StartTime:  time.Now(),
	}

	tm.Tunnels[localPort] = tunnel

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					tunnel.Status = "错误: " + err.Error()
				}
				return
			}
			go handleConnection(tm.SSHClient, conn, remoteDest)
		}
	}()

	return nil
}

func (tm *TunnelManager) StopTunnel(port string) {
	if tunnel, exists := tm.Tunnels[port]; exists {
		tunnel.Listener.Close()
		delete(tm.Tunnels, port)
	}
}

func (tm *TunnelManager) RestartTunnel(port string) error {
	tunnel, exists := tm.Tunnels[port]
	if !exists {
		return fmt.Errorf("隧道不存在")
	}

	tm.StopTunnel(port)
	return tm.CreateTunnel(tunnel.LocalPort, tunnel.RemoteDest)
}

func handleConnection(client *ssh.Client, localConn net.Conn, remoteDest string) {
	defer localConn.Close()

	remoteConn, err := client.Dial("tcp", remoteDest)
	if err != nil {
		showError(fmt.Sprintf("远程连接失败: %v", err))
		return
	}
	defer remoteConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		io.Copy(remoteConn, localConn)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(localConn, remoteConn)
		done <- struct{}{}
	}()

	<-done
	<-done
}

func showBanner() {
	color.Cyan(`
███████╗███████╗██╗  ██╗
██╔════╝██╔════╝██║  ██║
███████╗███████╗███████║
╚════██║╚════██║██╔══██║
███████║███████║██║  ██║
╚══════╝╚══════╝╚═╝  ╚═╝
DBLens for SSH Tunnel v1.0 `)
	color.Magenta("SSH隧道管理工具 v1.0 - 支持多隧道")
}

func showStatus(tm *TunnelManager) {
	color.Cyan("\n当前隧道状态:")
	if len(tm.Tunnels) == 0 {
		color.Yellow("没有活动的隧道")
		return
	}

	for port, tunnel := range tm.Tunnels {
		color.Yellow("本地端口: %s", cyan(port))
		color.Yellow("远程目标: %s", cyan(tunnel.RemoteDest))
		color.Yellow("运行状态: %s", green(tunnel.Status))
		color.Yellow("运行时长: %s", cyan(time.Since(tunnel.StartTime).Round(time.Second)))
		fmt.Println(strings.Repeat("-", 40))
	}
}

func showError(msg string) {
	color.Red("✗ 错误: %s", msg)
}

func showSuccess(msg string) {
	color.Green("✓ %s", msg)
}

func mainMenu(tm *TunnelManager) {
	reader := bufio.NewReader(os.Stdin)

	for {
		showStatus(tm)
		color.Cyan("\n请选择操作:")
		color.Green("1. 添加新隧道")
		color.Red("2. 删除隧道")
		color.Yellow("3. 重启隧道")
		color.Magenta("4. 删除所有隧道")
		color.Blue("5. 退出程序")

		fmt.Printf("%s ", cyan("输入选项:"))
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			addTunnelMenu(tm, reader)
		case "2":
			removeTunnelMenu(tm, reader)
		case "3":
			restartTunnelMenu(tm, reader)
		case "4":
			removeAllTunnels(tm)
		case "5":
			cleanupAndExit(tm)
		default:
			showError("无效选项")
		}

		fmt.Print("\033[H\033[2J")
		showBanner()
	}
}

func addTunnelMenu(tm *TunnelManager, reader *bufio.Reader) {
	fmt.Printf("%s ", yellow("本地端口:"))
	localPort, _ := reader.ReadString('\n')
	localPort = strings.TrimSpace(localPort)

	fmt.Printf("%s ", yellow("远程目标 (host:port):"))
	remoteDest, _ := reader.ReadString('\n')
	remoteDest = strings.TrimSpace(remoteDest)

	if err := tm.CreateTunnel(localPort, remoteDest); err != nil {
		showError(err.Error())
	} else {
		showSuccess(fmt.Sprintf("隧道 %s->%s 创建成功", localPort, remoteDest))
	}
}

func removeTunnelMenu(tm *TunnelManager, reader *bufio.Reader) {
	fmt.Printf("%s ", yellow("输入要删除的本地端口:"))
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)

	if _, exists := tm.Tunnels[port]; !exists {
		showError("隧道不存在")
		return
	}

	tm.StopTunnel(port)
	showSuccess("隧道已删除")
}

func restartTunnelMenu(tm *TunnelManager, reader *bufio.Reader) {
	fmt.Printf("%s ", yellow("输入要重启的本地端口:"))
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)

	if err := tm.RestartTunnel(port); err != nil {
		showError(err.Error())
	} else {
		showSuccess("隧道重启成功")
	}
}

func removeAllTunnels(tm *TunnelManager) {
	for port := range tm.Tunnels {
		tm.StopTunnel(port)
	}
	showSuccess("所有隧道已删除")
}

func cleanupAndExit(tm *TunnelManager) {
	removeAllTunnels(tm)
	if tm.SSHClient != nil {
		tm.SSHClient.Close()
	}
	color.Magenta("\n感谢使用，再见！")
	os.Exit(0)
}

func main() {
	showBanner()
	reader := bufio.NewReader(os.Stdin)

	color.Cyan("\n请输入SSH服务器配置:")
	fmt.Printf("%s ", yellow("SSH地址 (host:port):"))
	sshAddress, _ := reader.ReadString('\n')
	sshAddress = strings.TrimSpace(sshAddress)

	fmt.Printf("%s ", yellow("SSH用户名:"))
	sshUser, _ := reader.ReadString('\n')
	sshUser = strings.TrimSpace(sshUser)

	fmt.Printf("%s ", yellow("SSH密码:"))
	sshPassword, _ := reader.ReadString('\n')
	sshPassword = strings.TrimSpace(sshPassword)

	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", sshAddress, config)
	if err != nil {
		showError(fmt.Sprintf("SSH连接失败: %v", err))
		os.Exit(1)
	}

	tm := &TunnelManager{
		SSHClient:  client,
		Tunnels:    make(map[string]*Tunnel),
		SSHConfig:  config,
		SSHAddress: sshAddress,
	}

	fmt.Print("\033[H\033[2J")
	showBanner()
	mainMenu(tm)
}