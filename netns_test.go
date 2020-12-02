package golib_utility_network

import(
    utility "github.com/weizhouBlue/golib_utility_network"
    "fmt"
    "testing"
    "github.com/vishvananda/netns"
    "time"
    //"path"
    "github.com/vishvananda/netlink"

)


//---------------

func Test1(t *testing.T){

	// show all netns of pid

	utility.EnableLog=true

	pathList , _ := utility.GetAllPidNetnsPath()

	handlelist , _:=utility.GetAllNetnsHandleByPath(pathList)
	utility.CloseAllNetnsHandle(handlelist)

}





func Test3(t *testing.T){
	// enter to an netns and show ip

	utility.EnableLog=true


	currentNs , _ := utility.GetCurrentNetnsHandle()


	pathList , _ := utility.GetAllPidNetnsPath()
	handlelist , _:=utility.GetAllNetnsHandleByPath(pathList)
	for _ , handle := range handlelist {
		if e:=utility.GotoNetns(handle) ; e!=nil{
			fmt.Printf("failed to enter to ns : %v \n" , handle.UniqueId() )
		}else{
			if ip4List, _ ,e:=utility.GetAllInterfaceUnicastAddrByName( ) ; e==nil {
				fmt.Printf("in Ns %v , ip: %v \n" , handle.UniqueId() ,ip4List  )
			}
		}
	}

	utility.CloseAllNetnsHandle(handlelist)
	utility.CloseAllNetnsHandle( []netns.NsHandle{currentNs}  )

}





func Test2(t *testing.T){

	// show all netns by "ip netns"
	utility.EnableLog=true

	nameList := utility.GetNamedNetnsName()
	if len(nameList)>0{
		for _ , name := range nameList{
			fmt.Printf("find named netns : %v \n" , name)
		}

		handlelist , _:=utility.GetAllNetnsHandleByNamedNs(nameList)
		for _ , handle := range handlelist{
			fmt.Printf(" netns : %v \n" , handle )
		}

		utility.CloseAllNetnsHandle(handlelist)
	}


}



func Test5(t *testing.T){
	// create named netns  
	utility.EnableLog=true

	netnsName:="test"
	utility.DeteleNamedNetns(  netnsName   )

	handle , e := utility.NewNamedNetns(netnsName)
	if e!=nil{
		fmt.Printf("failed to create ns: %v \n" , e)
		t.FailNow()
	}
	fmt.Printf("succeeded to create named netns : %v \n" , netnsName)


	// you can input " ip netns " to see the netns
	time.Sleep(10*time.Second)
	utility.CloseAllNetnsHandle( []netns.NsHandle{ handle }  )

	utility.DeteleNamedNetns(  netnsName   )

	
}




func Test6(t *testing.T){
	// watch ip change event   

	utility.EnableLog=true

	handle , _ := utility.GetCurrentNetnsHandle()

	ch:=make ( chan netlink.AddrUpdate )
	done:=make( chan struct{} )

	if e:=utility.WatchIPChangeEventByNetNs(handle ,  ch , done )  ; e!=nil {
		fmt.Printf("failed to watch ip event")
		utility.CloseAllNetnsHandle( []netns.NsHandle{ handle }  )
		t.FailNow()
	}
	
    select{
        case data , ok := <-ch :
			if ok {   
				// https://godoc.org/github.com/vishvananda/netlink#AddrUpdate
				fmt.Println ("got event:" , data ) 
			}else{
				fmt.Println ("failed "  )
			}
        case <-time.After(20*time.Second):
        	fmt.Println ("failed  to get event" )
    }

    close(done)

    utility.CloseAllNetnsHandle( []netns.NsHandle{ handle }  )

}





