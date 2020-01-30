package golib_utility_network
import (
    "fmt"
    "github.com/vishvananda/netlink"

    //"github.com/vishvananda/netns"
    //"os"
    //"runtime"
    "net"
    //"golang.org/x/sys/unix"
    //"strconv"
    //"strings"
	//"github.com/vishvananda/netlink/nl"
)

// attention ! this should be used on linux os 


/*
doc  https://godoc.org/github.com/vishvananda/netlink
	https://godoc.org/github.com/vishvananda/netns

github https://github.com/vishvananda/netlink

refere example : 
https://github.com/vishvananda/netlink/blob/master/neigh_test.go
https://github.com/vishvananda/netlink/blob/master/netns_test.go

*/



//------------------------

// https://github.com/vishvananda/netlink/blob/master/neigh_linux.go
var (
	CONST_NeighState_none=netlink.NUD_NONE  //0x00
	CONST_NeighState_incomplete=netlink.NUD_INCOMPLETE  //0x01
	CONST_NeighState_reachable=netlink.NUD_REACHABLE //0x02
	CONST_NeighState_stale=netlink.NUD_STALE //0x04
	CONST_NeighState_delay=netlink.NUD_DELAY //0x08
	CONST_NeighState_probe=netlink.NUD_PROBE //0x10
	CONST_NeighState_failed=netlink.NUD_FAILED  //0x20   
	CONST_NeighState_noarp=netlink.NUD_NOARP   //0x40   
	CONST_NeighState_permanent=netlink.NUD_PERMANENT //0x80 	
)

//=================================

/*
Neigh
https://godoc.org/github.com/vishvananda/netlink#Neigh
type Neigh struct {
    LinkIndex    int
    Family       int
    State        int
    Type         int
    Flags        int
    IP           net.IP
    HardwareAddr net.HardwareAddr
    LLIPAddr     net.IP //Used in the case of NHRP
    Vlan         int
    VNI          int
    MasterIndex  int
}
*/
func GetNeighAll() ([]netlink.Neigh , error ) {

	msg:=netlink.Ndmsg{
		Family: uint8(0) ,
	}

	return netlink.NeighListExecute( msg )
}


func GetNeighByIp( ip string )( mac string, viaInterface string , state int , detail netlink.Neigh , e error) {
	if CheckIPFormat(ip)==false{
		e=fmt.Errorf("%v is not an ip address" , ip)
		return
	}

	if neighList , err:=GetNeighAll() ; err!=nil {
		e=err
		return
	}else{
		for _ , v :=range neighList {
			if v.IP.String()==ip {
				if name , err:=GetInterfaceNameByIndex(v.LinkIndex) ; err!=nil{
					continue 
				}else{
					viaInterface=name
					if len(viaInterface)==0 {
						continue 
					}
				}
				mac=v.HardwareAddr.String()
				if len(mac)==0 {
					continue 
				}
				detail=v
				state=v.State
				e=nil
				return
			}
		}
		e=fmt.Errorf("no entry for ip %v" , ip)
		return 
	}

}


func GetNeighByMac( mac string )( ip string, viaInterface string , state int , detail netlink.Neigh , e error) {
	
	if _ , err:=net.ParseMAC(mac);  err!=nil {
		e=fmt.Errorf("%v is not an mac , info=%v " , mac , err)
		return
	}

	if neighList , err:=GetNeighAll() ; err!=nil {
		e=err
		return
	}else{
		for _ , v :=range neighList {
			if v.HardwareAddr.String()==mac {
				if name , err:=GetInterfaceNameByIndex(v.LinkIndex) ; err!=nil{
					continue 
				}else{
					viaInterface=name
					if len(viaInterface)==0 {
						continue 
					}
				}
				ip=v.IP.String()
				if len(ip)==0 {
					continue 
				}
				detail=v
				state=v.State
				e=nil
				return
			}
		}
		e=fmt.Errorf("no entry for mac %v" , mac)
		return 
	}

}





//=================================


func addNeigh( ip , mac , viaInterface string  , state int ) error {

	toIp :=net.ParseIP(ip)
	if toIp==nil {
		return fmt.Errorf( "%v is not ip  " , ip   )
	}

	toMac , e1:=net.ParseMAC(mac);  
	if e1 !=nil || toMac==nil {
		return fmt.Errorf( "%v is not mac , %v " , mac  , e1 )
	}

	link , e2:=GetInterfaceByName( viaInterface  )
	if e2!=nil {
		return fmt.Errorf( "%v is not interface , %v " , viaInterface , e2 )
	}

	entry := netlink.Neigh{
		LinkIndex: link.Attrs().Index  ,
		State:     state ,
		IP:        toIp ,
		HardwareAddr:  toMac ,
	}

	err := netlink.NeighAdd(&entry)
	if err != nil {
		return err
	}
	return nil
}

func delNeigh( ip , mac , viaInterface string  , state int ) error {

	toIp :=net.ParseIP(ip)
	if toIp==nil {
		return fmt.Errorf( "%v is not ip  " , ip   )
	}

	toMac , e1:=net.ParseMAC(mac);  
	if e1 !=nil || toMac==nil {
		return fmt.Errorf( "%v is not mac , %v " , mac  , e1 )
	}

	link , e2:=GetInterfaceByName( viaInterface  )
	if e2!=nil {
		return fmt.Errorf( "%v is not interface , %v " , viaInterface , e2 )
	}

	entry := netlink.Neigh{
		LinkIndex: link.Attrs().Index  ,
		State:     state ,
		IP:        toIp ,
		HardwareAddr:  toMac ,
	}

	err := netlink.NeighDel(&entry)
	if err != nil {
		return err
	}
	return nil
}


func AddPermanentNeigh( ip , mac , viaInterface string ) error {
	return addNeigh( ip , mac , viaInterface , CONST_NeighState_permanent )
}

func DelPermanentNeigh( ip , mac , viaInterface string ) error {
	return delNeigh( ip , mac , viaInterface , CONST_NeighState_permanent )
}


