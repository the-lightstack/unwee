package utils

import (
		"os"
		"fmt"
		"bufio"
		"io"
		"strings"
		"sync"
		"net/http"
)


//displays usage details
func DisplayHelp(){
banner := `                          
 _  _ _ ___ __ _____ ___ 
| || | ' \ V  V / -_) -_)
 \_,_|_||_\_/\_/\___\___| v1.1.0

usage: unwee [options...]

options: 
[-u] (url) takes in single URL as input and gives it's unshortened form.
[-o] (output) takes in file name as input and saves result in the file
[-f] (file) takes in file name containing list of shortened URLs
[-h] (help) prints help menu

Example:

- Single URL Input from stdin:

	> echo "http://shorturl/xyz" | go run main.go

- Input single url

	> go run main.go -u http://shorturl/xyz

- URL List file Input from stdin

	> cat url_list.txt | go run main.go 

- Input from file

	> go run main.go -f url_list.txt

- Setting output file to save results (default is saved at "output/<data_time>.txt")

	> go run main.go -o outputfile.txaer
`
fmt.Println(banner)
os.Exit(0)
}

//check error
func CheckErr(err error){
	if err!=nil{
		fmt.Println("[ERR] ",err)
		os.Exit(0)
	}
}

//return slice of data in the file
func GetFileData(filePath string)([]string){
	var output []string
	fl,err:=os.Open(filePath)
	if os.IsNotExist(err) {
		fmt.Println("[ERR] ",err)
		os.Exit(0)
	}
 	defer fl.Close()
	
	reader := bufio.NewReader(fl)
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, strings.Split(input, "\n")[0])
	}
	return output
}

//checks if the results are piped to another command
func isPiped() bool {
	info, _ := os.Stdin.Stat()
	if info.Mode()&os.ModeCharDevice == 0 {
		return true
	} 
		return false
}

//read input from stdin
func GetStdin() []string {
	var output []string
	if !isPiped() {
		return output 
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, strings.Split(input, "\n")[0])
	}

	return output
}

//writes slice of string line by line to a file
func WriteToFile(resList []string,filePath string){
	//if file doesn't exists then create one
	outFile, err := os.Create(filePath)
	CheckErr(err)
	writer:=bufio.NewWriter(outFile)
	for _,data := range resList{
		_,_=writer.WriteString(data+"\n")
	}
	writer.Flush()
	outFile.Close()
}

//return unshortened URL
func Start(url string, wg *sync.WaitGroup, resultList *[]string){
	defer wg.Done()
	res,err:=http.Head(url)
	CheckErr(err)
	fmtRes:=fmt.Sprintf("\n%d  %s  %s",res.StatusCode,url,res.Request.URL.String())
	*resultList=append(*resultList,fmtRes)
	fmt.Println(fmtRes)
}
