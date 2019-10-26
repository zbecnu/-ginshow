package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main()  {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh,syscall.SIGINT,syscall.SIGTERM,syscall.SIGQUIT)
	go func() {
		for s:= range signalCh{
			switch s {
			case syscall.SIGINT:
				fmt.Println("hahaha")
				os.Exit(0)
			default:
				fmt.Println("kother")
			}
		}
	}()

	sum:=0


	for{

		sum++
		time.Sleep(time.Second)
	}

	cmds  := []*exec.Cmd{
		exec.Command("ps","aux"),
		exec.Command("grep","signal"),
	}

	output,err := runCmds(cmds)
	if err != nil{

		fmt.Printf("Command Execution Error:",err)
		return
	}

	pids,err := getPids(output)

	if err != nil{
		return
	}

	for _,pid  := range pids{
		proc ,err:=  os.FindProcess(pid)
		if err != nil{
			return
		}

		sig:= syscall.SIGQUIT
		err =proc.Signal(sig)
		if err != nil{
			return
		}
	}


}

func runCmds(cmds []*exec.Cmd)([]string,error){
	if cmds != nil || len(cmds) ==0{
		return nil,errors.New("error")
	}
	var err error
	first := true

	var output []byte

	for _,cmd := range  cmds{
		fmt.Printf("Run command:%v\n",getCmdPlaintext(cmd))

		if first{
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}

		var stdoutBuf bytes.Buffer

		cmd.Stdout = &stdoutBuf

		if err = cmd.Start();err != nil{
			return nil,getError(err,cmd)
		}

		if err = cmd.Wait(); err != nil{
			return nil ,getError(err,cmd)
		}
		output = stdoutBuf.Bytes()
		fmt.Printf("output:\n%s\n",string(output))
		if first{
			first =  false
		}
	}

	var lines  []string

	var outputBuf bytes.Buffer

	outputBuf.Write(output)

	for{
		line, err := outputBuf.ReadBytes('\n')
		if err != nil{
			if err ==io.EOF{
				break
			}else{
				return nil,nil
			}
		}

		lines =  append(lines,string(line))
	}

	return lines,nil

}


func getCmdPlaintext(cmd *exec.Cmd) string{
	var buff bytes.Buffer
	buff.WriteString(cmd.Path)

	for _, arg:= range  cmd.Args[1:]{
		 buff.WriteRune(' ')
		 buff.WriteString(arg)
	}
	return buff.String()

}

func getError(err error, cmd *exec.Cmd,extrainfo ...string)error{
	var errMsg string
	if cmd != nil{
		errMsg = fmt.Sprintf("%s%s%v",err, cmd.Path,cmd.Args)

	}else{
		errMsg = fmt.Sprintf("%s",err)

	}
	if len(extrainfo)>0{
		errMsg =  fmt.Sprintf("%s%v",errMsg,extrainfo)
	}
	return errors.New(errMsg)
}


func getPids(strs []string)([]int ,error){
	var pids []int
	for _,str := range  strs{
		pid,err := strconv.Atoi(strings.TrimSpace(str))

		if err != nil{
			return nil,err
		}

		pids =  append(pids,pid)

	}

	return pids,nil
}


func main02()  {
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		time.Sleep(5*time.Second)

		context.String(http.StatusOK,"welcome to my website!")
	})

	srv := &http.Server{
		Addr:":8080",
		Handler:router,
	}

	go func() {
		if err:= srv.ListenAndServe();err != nil && err != http.ErrServerClosed{
			log.Fatalf("listen :%s \n" ,err)
		}
	}()

//	优雅地关闭服务器
	quit := make(chan os.Signal)

	signal.Notify(quit,os.Interrupt)

	<-quit
	log.Println("shut down server ing...")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)

	defer cancel()

	if err:= srv.Shutdown(ctx); err!=nil{
		log.Fatal("Server shutdown:",err)
	}

	log.Println("Server exiting")

}