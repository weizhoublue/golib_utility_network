package golib_utility_network
import (
    "fmt"
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"
    "os"
    "runtime"
    //"net"
    "path/filepath"
    "strings"
    "strconv"
    "path"

)

// attention ! this should be used on linux os 


/*
doc  	
		https://godoc.org/github.com/vishvananda/netns
		https://godoc.org/github.com/vishvananda/netlink

github 
		https://github.com/vishvananda/netlink

refere example : 
https://github.com/vishvananda/netlink/blob/master/netns_test.go
*/




/*

func GetCurrentNetnsHandle() (netns.NsHandle, error)

#for named netns
func GetNamedNetnsName()  []string 
func GetAllNetnsHandleByNamedNs( nameList []string ) ([]netns.NsHandle , error )
func CloseAllNetnsHandle( list []netns.NsHandle )

#for named netns
func NewNamedNetns( name string ) ( netns.NsHandle, error) 
func DeteleNamedNetns( fpath string ) error 
func CloseAllNetnsHandle( list []netns.NsHandle )


# for pid netns
func GetAllPidNetnsPath(  ) ([]string , error )
func GetAllNetnsHandleByPath( pathList []string ) ([]netns.NsHandle , error ) 
func CloseAllNetnsHandle( list []netns.NsHandle )

func GotoNetns( ns  netns.NsHandle  ) error 


func WatchIPChangeEventByNetNs( ns netns.NsHandle, ch chan<- netlink.AddrUpdate, done <-chan struct{}) error 


*/
//=================================


var (
	// named net ns , like ip netns ls 
	ConstNamedNetnsPath="/run/netns"
    EnableLog=false
)



func toString(a interface{}) string {
	if v, p := a.(int); p {
		return strconv.Itoa(v)
	}
	if v, p := a.(int16); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(int32); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(uint); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(float32); p {
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	}
	if v, p := a.(float64); p {
		return strconv.FormatFloat(v, 'f', -1, 32)
	}
	return ""
}

func existFile(filePath string) bool {
	if info, err := os.Stat(filePath); err == nil {
		if !info.IsDir() {
			return true
		}
	}
	return false
}


func getFileName( fpath string) string {
	b := strings.LastIndex(fpath, "/")
	if b >= 0 {
		return fpath[b+1:]
	} else {
		return fpath
	}
}

func log( format string, v ...interface{}) {
	prefix:=""
	funcName, filepath, line, ok := runtime.Caller(1)
	if ok {
		file := getFileName(filepath)
		funcname := getFileName(runtime.FuncForPC(funcName).Name())
		prefix += "[" + file + " " + funcname + " " + toString(line) + "]     "
	}
	fmt.Printf(prefix+format, v...)

}



//===================================================


func IsDir( fpath string) bool {  
    s, err := os.Stat(fpath)  
    if err != nil {  
        return false  
    }  
    return s.IsDir()  
}

// return the absolute dir path of file list (not include subDirecotry )
func GetDirFileList( dirPath string ) ([]string , error ) {

	// https://godoc.org/path/filepath#Glob
	filepathNames,e := filepath.Glob( filepath.Join( dirPath ,"*") )
	if e != nil {
		return nil , e
	}
 
 	list:=[]string{}

	for _ , v := range filepathNames {
		if IsDir( v )==false {
			list=append(list,v )
		}
	}
	return list , nil 
}

// func GetDirFileListUnderPwd( ) ([]string , error) {
// 	pwd,_ := os.Getwd()
// 	return GetDirFileList(pwd)
// }






//===================================================

//get all named nens file path under ConstNamedNetnsPath
func GetNamedNetnsName()  []string {

	NamedFileList , _ := GetDirFileList( ConstNamedNetnsPath )

	nsNameList:=[]string{}
	for _ , fpath := range NamedFileList {
		nsNameList=append(nsNameList , path.Base(fpath) )
	}
	return nsNameList

}




func GetAllPidNetnsPath(  ) ([]string , error ) {

	// https://godoc.org/path/filepath#Glob
	filepathNames,e := filepath.Glob( "/proc/*/ns/net" )
	if e != nil {
		return nil , e
	}
 
 	list:=[]string{}

	for _ , v := range filepathNames {
		if IsDir( v )==false {
			list=append(list,v )
		}
	}
	return list , nil 

}




//===================================================

func GetAllNetnsHandleByPath( pathList []string ) ([]netns.NsHandle , error ) {

	if len(pathList)==0 {
		return nil , fmt.Errorf("empty pathList")
	}


	list:=[]netns.NsHandle{}
	OUTER:
	for _ , nspath := range pathList {
		handle , e := netns.GetFromPath(nspath)
		if e==nil &&  handle!=0 {
			// compare with existed netns
			for _, n :=range list {
				if n.Equal(handle) == true {
					handle.Close()
					continue OUTER
				}
			}
			list=append(list,handle)
			log("find netns : %v \n" , handle.UniqueId() )
		}
	}

	return list , nil

}


func GetAllNetnsHandleByNamedNs( nameList []string ) ([]netns.NsHandle , error ) {
	pathList:=[]string{}
	for _ , name :=range nameList {
		pathList=append( pathList , path.Join(ConstNamedNetnsPath,name) )
	}
	return GetAllNetnsHandleByPath(pathList)

}






func GetCurrentNetnsHandle() (netns.NsHandle, error){
	return netns.Get() 
}


func CloseAllNetnsHandle( list []netns.NsHandle ){
	if list==nil || len(list)==0 {
		return
	}

	for _ , handle := range list {
		log("close netns : %v \n" , handle.UniqueId() )
		handle.Close()
	}
	return
}


//===================================================


func GotoNetns( ns  netns.NsHandle  ) error {
	return netns.Set(ns)
}

//===================================================



func NewNamedNetns( name string ) ( netns.NsHandle, error) {
	// https://github.com/vishvananda/netns/blob/master/netns_linux.go#L51
	// create a netns and bind to the file
	return netns.NewNamed(name)
}



// ns for ip netns
func DeteleNamedNetns( name string ) error {
	return netns.DeleteNamed( name)


}



//============

// event data struct : https://godoc.org/github.com/vishvananda/netlink#AddrUpdate
func WatchIPChangeEventByNetNs( ns netns.NsHandle, ch chan<- netlink.AddrUpdate, done <-chan struct{}) error {
	return netlink.AddrSubscribeAt(ns , ch , done ) 
}

