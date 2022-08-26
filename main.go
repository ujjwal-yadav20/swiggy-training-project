package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"sync"
	"time"
)


var statusMap = map[string]string{}
func checkStatus(str string,wg *sync.WaitGroup){
	defer wg.Done()
	resp,err := http.Get(str)
	if err != nil {
		statusMap[str]="DOWN"
		return	
		//log.Fatal(err)
	}
	status_code := resp.StatusCode
	if status_code==200{
		statusMap[str]="UP"
	}else{
		statusMap[str]="DOWN"
	}
	fmt.Println(status_code,":",str)

}
func mapWebsites(w http.ResponseWriter, r *http.Request) {

	typeRequest:=r.Method
	if typeRequest =="GET"{
		param := r.URL.Query()
		listWebsites,ok := param["name"]
		if ok{
			for _,url := range listWebsites{
				urlStatus,check := statusMap[url]
				if check{
					fmt.Fprintln(w,fmt.Sprintf("The Status of Website %+v is: %+v",string(url),string(urlStatus)))
				}else{
					fmt.Fprintln(w,fmt.Sprintf("The Website %+v is not present in the Map",url))

				}
				
				
			}

		}else{
			for key,val := range statusMap {
				fmt.Fprintln(w,"Status of Website",key ,"is :", val)
			}
		}

	


	}else{
		var arr []string
		err := json.NewDecoder(r.Body).Decode(&arr)
		if err != nil {
			fmt.Fprintf(w,fmt.Sprintf("Error:%+v",err))
			return
		}
		//fmt.Println(arr)
		currLen:=len(arr)
		for i:=0;i<currLen ;i++{
			fmt.Println(arr[i])
			val , ok := statusMap[arr[i]]
			if ok{
				statusMap[arr[i]]=val
			}else{
				statusMap[arr[i]]="CHECK"
			}
		}
		
	}
	
}

func runCheck() {
	//fmt.Println("Inside runcheck")
	var wg sync.WaitGroup
	currSize := len(statusMap)
	wg.Add(currSize)
	for key, _ := range statusMap {
		go checkStatus(key,&wg)
        //fmt.Println("Website:", key)
    }
	wg.Wait()
	//fmt.Println("Status of Websites is:")
	for key,val := range statusMap{
		fmt.Println(key,"->",val)
	}
}
func loopOver(){
	for{
		runCheck()
		time.Sleep(10 * time.Second)
	}

}
func main() {
	//http.HandleFunc("/", getRoot)
	//var wg sync.WaitGroup
	go loopOver()
	http.HandleFunc("/", mapWebsites)
	fmt.Println("Starting the Server")
	err := http.ListenAndServe(":8000", nil)
  	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
	
	// for{
	// 	runCheck()
	// 	time.Sleep(60 * time.Second)
	// }
}