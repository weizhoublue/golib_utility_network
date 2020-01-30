package golib_utility_network
import (
    "fmt"
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"
    "os"
    "runtime"
    "net"
)

// attention ! this should be used on linux os 


/*
doc  https://godoc.org/github.com/vishvananda/netlink
	https://godoc.org/github.com/vishvananda/netns

github https://github.com/vishvananda/netlink

refere example : 
https://github.com/vishvananda/netlink/blob/master/link_test.go
https://github.com/vishvananda/netlink/blob/master/netns_test.go
*/



/*
func CheckHostInterfaceByName(name string) bool 

func ListHostAllInterfaces() ( []string, error )
func ListHostUpInterfaces() ( []string, error )

func GetHostInterfaceMac(name string) ( string , error )
func GetHostInterfaceMtu(name string) ( int , error ) 
func CheckHostInterfaceUp(name string) ( bool , error ) 

func GetInterfaceUnicastAddrByName( name string ) ( ipv4_list , ipv6_list []string , err error ) 

func GetAllInterfaceUnicastAddrByName( ) ( list_ipv4_all , list_ipv6_all []string , err error ) 
func GetAllUpInterfaceUnicastAddrByName( ) ( list_ipv4_all , list_ipv6_all []string , err error ) 

//----------------------------------------

func GetAllInterface( )( []netlink.Link , error )
func GetAllInterfaceInNetns( pid int )( []netlink.Link , error )

func GetInterfaceByName( name string )( netlink.Link , error )
func GetInterfaceByNameInNetns( pid int , name string )( netlink.Link , error )
func GetInterfaceNameByIndex( index int )( string , error )
func GetInterfaceNameByIndexInNetns( pid int , index int )( string , error )

func SetInterfaceUp( interfaceName string ) error 
func SetInterfaceDown( interfaceName string ) error 

func CreateInterfaceVeth( name , vethName string ) error 
func CreateInterfaceVlan( parentName , name  string  , vlanId int ) error 

func DelInterfaceByName( name string ) error 
func DelInterfaceByNameInNetns( pid int , name string ) error 

func SetInterfaceToNetns( pid int , name string ) error 
func SetInterfaceNewName(  oldName , newName string ) error 
func SetInterfaceNewNameInNetns( pid int ,  oldName , newName string ) error 
func SetInterfaceMtu( name string  , mtu int ) error 
func SetInterfaceMtuInNetns( pid int , name string  , mtu int ) error 
func GetInterfaceMtu( name string   ) ( int , error) 
func GetInterfaceMtuInNetns( pid int , name string   ) ( int , error) 
func SetInterfaceMac( name string  , mac string ) error 
func GetInterfaceMac( name string   ) ( string , error) 

*/
//=================================

func checkRootPrivilege() int {

	if os.Getuid() != 0 {
		return -1
	}
	return 0
}


//================================

/*
netlink.Link  https://godoc.org/github.com/vishvananda/netlink#Link
type Link interface { 
    Attrs() *LinkAttrs  // https://godoc.org/github.com/vishvananda/netlink#LinkAttrs
    Type() string  // device , bridge , veth , ipip , ...... 
}
*/
func GetAllInterface( )( []netlink.Link , error ){
	return netlink.LinkList()
}

func GetAllInterfaceInNetns( pid int )( []netlink.Link , error ){
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return nil , e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return nil , e
		}
		return GetAllInterface()

	}
}



func GetInterfaceByName( name string )( netlink.Link , error ){
	return netlink.LinkByName(name)
}


func GetInterfaceByNameInNetns( pid int , name string )( netlink.Link , error ){
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return nil , e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return nil , e
		}
		return GetInterfaceByName(name)
	}
}


func GetInterfaceNameByIndex( index int )( string , error ){
	if link , err:=netlink.LinkByIndex(index) ; err!=nil{
		return "" , err
	}else{
		return link.Attrs().Name , nil
	}
}

func GetInterfaceNameByIndexInNetns( pid int , index int )( string , error ){
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return "" , e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return "" , e
		}
		return GetInterfaceNameByIndex(index)
	}
}



//================================

func SetInterfaceUp( interfaceName string ) error {

	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}
	return nil
}



func SetInterfaceDown( interfaceName string ) error {

	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetDown(link); err != nil {
		return err
	}
	return nil
}




//================================


func CreateInterfaceVeth( name , vethName string ) error {

	link := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:        name ,
		},
		PeerName:         vethName ,
	}

	if err := netlink.LinkAdd(link); err != nil {
		return err
	}

	//check what we created 
	base := link.Attrs()
	result, err := netlink.LinkByName(link.Name)
	if err != nil {
		return err
	}
	rBase := result.Attrs()
	if veth, ok := result.(*netlink.Veth); ok {
		original  := link
		if original.PeerName != "" {
			other, err := netlink.LinkByName(original.PeerName)
			if err != nil {
				return fmt.Errorf("Peer %s not created", veth.PeerName)
			}
			if _ , ok = other.(*netlink.Veth); !ok {
				return fmt.Errorf("Peer %s is incorrect type", veth.PeerName)
			}
		}
		
	} else {
		// recent kernels set the parent index for veths in the response
		if rBase.ParentIndex == 0 && base.ParentIndex != 0 {
			return fmt.Errorf("Created link doesn't have parent %d but it should", base.ParentIndex)
		} else if rBase.ParentIndex != 0 && base.ParentIndex == 0 {
			return fmt.Errorf("Created link has parent %d but it shouldn't", rBase.ParentIndex)
		} else if rBase.ParentIndex != 0 && base.ParentIndex != 0 {
			if rBase.ParentIndex != base.ParentIndex {
				return fmt.Errorf("Link.ParentIndex doesn't match %d != %d", rBase.ParentIndex, base.ParentIndex)
			}
		}
	}
	return nil

}




//================================

func CreateInterfaceVlan( parentName , name  string  , vlanId int ) error {

	parent, err := netlink.LinkByName(parentName)
	if err != nil {
		return err
	}

	//https://sourcegraph.com/-/godoc/refs?def=Vlan&pkg=github.com%2Fvishvananda%2Fnetlink&repo=github.com%2Fvishvananda%2Fnetlink
	link := &netlink.Vlan{
		netlink.LinkAttrs{ 
			Name: name , 
			ParentIndex: parent.Attrs().Index ,
		}, 
		vlanId ,
		netlink.VLAN_PROTOCOL_8021Q ,
	}

	if err = netlink.LinkAdd(link); err != nil {
		return err
	}

	//check
	result, err := netlink.LinkByName(link.Name)
	if err != nil {
		return err
	}
	other, ok := result.(*netlink.Vlan)
	if !ok {
		return fmt.Errorf("the type of created link is not a vlan")
	}
	if link.VlanId != other.VlanId {
		return fmt.Errorf("the VlanId of created link doesn't match")
	}

	return nil
}



//================================

func SetInterfaceToNetns( pid int , name string ) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()


	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	if e:=netlink.LinkSetNsPid(link , pid ) ; e!=nil {
		return e
	}

	// check
	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return e
		}
		if _, e:=netlink.LinkByName(name) ; e!=nil{
			return fmt.Errorf("link %v is not in netns with pid=%v" , name ,pid )
		}
	}
	return nil 
}

//================================


func DelInterfaceByName( name string ) error {

	var link netlink.Link
	var e error
	if link , e= netlink.LinkByName(name) ; e!=nil{
		return e
	}

	if e = netlink.LinkDel(link); e != nil {
		return e
	}

	if links, err := netlink.LinkList() ; err != nil {
		return err
	}else{
		for _, l := range links {
			if l.Attrs().Name == link.Attrs().Name {
				return fmt.Errorf("Link %v not removed properly" , name )
			}
		}			
	}
	return nil
}



func DelInterfaceByNameInNetns( pid int , name string ) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()
		if e:=netns.Set(f) ; e!=nil {
			return  e
		}

		return DelInterfaceByName(name)

	}


}



//================================

func SetInterfaceNewName(  oldName , newName string ) error {

	var link netlink.Link
	var e error
	if link , e = netlink.LinkByName(oldName) ; e!=nil{
		return e
	}
	if e = netlink.LinkSetName( link , newName) ; e!=nil{
		return e
	}
	return nil

}

func SetInterfaceNewNameInNetns( pid int ,  oldName , newName string ) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()
		if e:=netns.Set(f) ; e!=nil {
			return  e
		}

		return SetInterfaceNewName(oldName , newName )

	}

}


//================================


func SetInterfaceMtu( name string  , mtu int ) error {

	var link netlink.Link
	var e error
	if link , e = netlink.LinkByName(name) ; e!=nil{
		return e
	}
	if e = netlink.LinkSetMTU( link , mtu) ; e!=nil{
		return e
	}
	return nil

}

func SetInterfaceMtuInNetns( pid int , name string  , mtu int ) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()
		if e:=netns.Set(f) ; e!=nil {
			return  e
		}
		return SetInterfaceMtu(name , mtu )
	}


}


func GetInterfaceMtu( name string   ) ( int , error) {

	var link netlink.Link
	var e error
	if link , e = netlink.LinkByName(name) ; e!=nil{
		return 0 , e
	}
	return link.Attrs().MTU , nil 

}


func GetInterfaceMtuInNetns( pid int , name string   ) ( int , error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  0, e

	}else{
		defer f.Close()
		if e:=netns.Set(f) ; e!=nil {
			return  0, e
		}
		return GetInterfaceMtu(name  )
	}


}

//================================

func SetInterfaceMac( name string  , mac string ) error {

	if hw , err:=net.ParseMAC(mac) ; err!=nil{
		return err
	}else{

		var link netlink.Link
		var e error
		if link , e = netlink.LinkByName(name) ; e!=nil{
			return e
		}
		if e = netlink.LinkSetHardwareAddr( link , hw ) ; e!=nil{
			return e
		}
		return nil		
	}

}


func GetInterfaceMac( name string   ) ( string , error) {

	var link netlink.Link
	var e error
	if link , e = netlink.LinkByName(name) ; e!=nil{
		return "" , e
	}
	return link.Attrs().HardwareAddr.String() , nil 

}



//======================
//===========================================

func CheckHostInterfaceByName(name string) bool {
    if _, err:= net.InterfaceByName(name) ; err!=nil{
        return false
    }
    return true
}

func ListHostAllInterfaces() ( []string, error ){
	list , err:=net.Interfaces()
	if err!=nil {
		return nil , err
	}
	name_list:=[]string {}
	for _ , v := range list {
		name_list=append(name_list , v.Name )
	}
	if len(name_list)==0{
		return nil, nil
	}
	return name_list , nil 
}

func ListHostUpInterfaces() ( []string, error ){
	list , err:=net.Interfaces()
	if err!=nil {
		return nil , err
	}
	name_list:=[]string {}
	for _ , v := range list {
		if ( v.Flags & 0x1 ) != 0 {
			name_list=append(name_list , v.Name )
		}
	}
	if len(name_list)==0{
		return nil, nil
	}
	return name_list , nil 
}


func GetHostInterfaceMac(name string) ( string , error ) {
    if result , err:= net.InterfaceByName(name) ; err!=nil{
		return "" , fmt.Errorf("no host interface with name=%v " , name )
    }else{
    	fmt.Println(result )
    	if len(result.HardwareAddr.String())==0 {
    		return "" , fmt.Errorf("failed to get the Mac " )
    	}
    	return result.HardwareAddr.String()  , nil
    }
}



func GetHostInterfaceMtu(name string) ( int , error ) {
    if result , err:= net.InterfaceByName(name) ; err!=nil{
		return 0 , fmt.Errorf("no host interface with name=%v " , name )
    }else{
    	return result.MTU  , nil
    }
}



func CheckHostInterfaceUp(name string) ( bool , error ) {
    if result , err:= net.InterfaceByName(name) ; err!=nil{
		return false , fmt.Errorf("no host interface with name=%v " , name )
    }else{
    	tmp:=result.Flags & 0x1
    	if tmp==0 {
    		return false, nil
    	}else{
    		return true, nil
    	}
    }
}



func GetInterfaceUnicastAddrByName( name string ) ( ipv4_list , ipv6_list []string , err error ) {
    result , err:= net.InterfaceByName(name)
    if err!=nil{
        return nil ,nil , fmt.Errorf("no host interface with name=%v " , name )
    }
    list , erro:=result.Addrs()
    if erro!=nil{
        return nil ,nil , erro
    }
    for _ , v := range list {
    		m:=v.String()
    		if CheckIPv4FormatWithMask( m ) {
    			ipv4_list=append(ipv4_list , m )
    		}else{
    			ipv6_list=append(ipv6_list , m )
    		}
    }
    return
}

func GetAllInterfaceUnicastAddrByName( ) ( list_ipv4_all , list_ipv6_all []string , err error ) {
	re , err :=ListHostAllInterfaces()
	if err!=nil{
		return nil , nil , err
	}

	for _ , name :=range re {
		ipv4_list , ipv6_list , err := GetInterfaceUnicastAddrByName(name)
		if err != nil {
			return nil , nil , err
		}
		if len(ipv4_list)>0{
			list_ipv4_all=append( list_ipv4_all , ipv4_list... )
		}
		if len(ipv6_list)>0{
			list_ipv6_all=append( list_ipv6_all , ipv6_list... )
		}
	}
	return 
}

func GetAllUpInterfaceUnicastAddrByName( ) ( list_ipv4_all , list_ipv6_all []string , err error ) {
	re , err :=ListHostUpInterfaces()
	if err!=nil{
		return nil , nil , err
	}

	for _ , name :=range re {
		ipv4_list , ipv6_list , err := GetInterfaceUnicastAddrByName(name)
		if err != nil {
			return nil , nil , err
		}
		if len(ipv4_list)>0{
			list_ipv4_all=append( list_ipv4_all , ipv4_list... )
		}
		if len(ipv6_list)>0{
			list_ipv6_all=append( list_ipv6_all , ipv6_list... )
		}
	}
	return 
}


