package main

import (
	"fmt"
	"os/exec"
	"strings"
	"container/list"

	"time"
	"strconv"
)

type threadStack struct {
	stack string
	threadId string
	cpu string
	createTime time.Time
}
func main(){
	commond := "9397"
	cmd := exec.Command("jstack",commond)
	//cmd :=exec.Command("top")
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		return
	}

	//stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return
	}
	/*
	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line,_, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 || line == nil ||len(line) == 0 {
			break
		}
		fmt.Println(string(line))
	}
	*/

	strs := analyseJstack(output)
	tops := analyseTop(output)

	for e := strs.Front(); e != nil; e = e.Next(){
		fmt.Println(e.Value)
	}

	fmt.Println("Execute Command finished.\n")
}
func process(topList *list.List,stackList *list.List){
	for stack := topList.Front();stack != nil; stack = stack.Next(){
		for thread :=stackList.Front();thread != nil; thread = thread.Next(){
			var threadValue = thread.Value.(string)
			var stackValue= stack.Value.(threadStack)
			if  strings.Contains(threadValue,"nid=0x"+stackValue.threadId){
				stackValue.stack=threadValue
				break
			}
		}
	}

}
func analyseJstack(output []byte ) *list.List{
	var threadStr string
	threadStr = ""
	out := string(output)
	strs := strings.Split(out,"\n")
	opp :=list.New()
	for _,num := range strs{
		if strings.Contains(num,"nid=0x") && strings.Contains(num,"tid=0x"){
			if strings.Contains(threadStr,"nid=0x") && strings.Contains(threadStr,"tid=0x"){
				opp.PushBack(threadStr)
				//fmt.Println("111:",threadStr)
				threadStr = num + "\n";

			}else{
				threadStr=num+"\n";
			}
		}else {
			threadStr = threadStr + num  + "\n";
		}

	}
	opp.PushBack(threadStr)
	return opp
}

func analyseTop(output []byte ) *list.List{
	date :=time.Now()
	out := string(output)
	strs := strings.Split(out,"\n")
	opp :=list.New()
	for _,num := range strs{
		var stack *threadStack = new(threadStack)
		stack.createTime=date
		var cpus []string = strings.Split(num," ")
		if cpus == nil || len(cpus) ==0 {
			continue
		}

		cpu,_:=strconv.Atoi(cpus[0])
		stack.threadId=strconv.FormatInt(int64(cpu),16)
		stack.cpu=(cpus[1])
		opp.PushBack(stack)
	}

	return opp
}
