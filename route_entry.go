package golib_utility_network
import (
    "fmt"
    "github.com/vishvananda/netlink"

    //"github.com/vishvananda/netns"
    //"os"
    //"runtime"
    "net"
    "golang.org/x/sys/unix"
    "strconv"
    "strings"
	//"github.com/vishvananda/netlink/nl"
)

// attention ! this should be used on linux os 


/*
doc  https://godoc.org/github.com/vishvananda/netlink
	https://godoc.org/github.com/vishvananda/netns

github https://github.com/vishvananda/netlink

refere example : 
https://github.com/vishvananda/netlink/blob/master/route_test.go
https://github.com/vishvananda/netlink/blob/master/netns_test.go

*/


/*
func GetIpv4RouteAllEntryFromSpecifiedTable( tableNum int )([]netlink.Route , error )
func GetIpv4RouteAllEntryFromAllTable( ) ([]netlink.Route , error )
func GetIpv4RouteAllEntryFromMainTable( ) ([]netlink.Route , error )
func GetIpv4RouteDefaultFromMainTable( ) ( gw , viaInterface string ,  detailRoute netlink.Route , e error )

func GetIpv6RouteAllEntryFromSpecifiedTable( tableNum int )([]netlink.Route , error )
func GetIpv6RouteAllEntryFromAllTable( ) ([]netlink.Route , error )
func GetIpv6RouteAllEntryFromMainTable( ) ([]netlink.Route , error )
func GetIpv6RouteDefaultFromMainTable( ) ( gw , viaInterface string ,  detailRoute netlink.Route , e error )

func CalculateIpv4RouteByDst( dst string ) ( []netlink.Route , error )
func CalculateIpv6RouteByDst( dst string ) ( []netlink.Route , error )

func CreateIPv4RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error )
func DelIPv4RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) 

func CreateIPv6RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) 
func DelIPv6RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) 

*/


//------------------------------

var (
	CONST_RouteTable_UNSPEC=unix.RT_TABLE_UNSPEC //0
	CONST_RouteTable_COMPAT=unix.RT_TABLE_COMPAT  //0xfc
	CONST_RouteTable_MAIN=unix.RT_TABLE_MAIN  //0xfe
	CONST_RouteTable_DEFAULT=unix.RT_TABLE_DEFAULT //0xfd
	CONST_RouteTable_LOCAL=unix.RT_TABLE_LOCAL //0xff
	CONST_RouteTable_MAX=unix.RT_TABLE_MAX //0xffffffff
)

//---------------- get ipv4 route --------------
/*
Route
https://godoc.org/github.com/vishvananda/netlink#Route
type Route struct {
    LinkIndex  int  //下跳网卡的index
    ILinkIndex int
    Scope      Scope
    Dst        *net.IPNet  //目的IP
    Src        net.IP   //数据包使用的本地源IP
    Gw         net.IP
    MultiPath  []*NexthopInfo
    Protocol   int
    Priority   int
    Table      int
    Type       int
    Tos        int
    Flags      int
    MPLSDst    *int
    NewDst     Destination
    Encap      Encap
    MTU        int
    AdvMSS     int
    Hoplimit   int
}

*/
func GetIpv4RouteAllEntryFromSpecifiedTable( tableNum int )([]netlink.Route , error ){

	routeFilter := &netlink.Route{
		Table: tableNum ,
	}

	routes, err := netlink.RouteListFiltered(  netlink.FAMILY_V4  , routeFilter , netlink.RT_FILTER_TABLE )
	if err != nil {
		return nil , err
	}
	return routes , nil

}


func GetIpv4RouteAllEntryFromAllTable( ) ([]netlink.Route , error ){

	return GetIpv4RouteAllEntryFromSpecifiedTable(unix.RT_TABLE_UNSPEC)

}

func GetIpv4RouteAllEntryFromMainTable( ) ([]netlink.Route , error ){

	return GetIpv4RouteAllEntryFromSpecifiedTable(unix.RT_TABLE_MAIN)

}


func GetIpv4RouteDefaultFromMainTable( ) ( gw , viaInterface string ,  detailRoute netlink.Route , e error ){

	if routeList, err:=GetIpv4RouteAllEntryFromMainTable(); err!=nil{
		return "" , "" , netlink.Route{} , err
	}else{
		min:=-1
		position:=-1
		for m , k:=range routeList {
			if k.Gw!=nil && k.Dst==nil  {
				//default route
				//check priority
				if min<0 {
					min=k.Priority
					position=m
				}else if k.Priority<min {
					min=k.Priority
					position=m
				}
			}
		}
		if position>=0{
			gw=routeList[position].Gw.String()
			index:=routeList[position].LinkIndex

			if linkList , err := netlink.LinkList() ; err!=nil {
				return "" , "" , netlink.Route{} , fmt.Errorf("failed to find the interface, info=%v " ,err  )
			}else{
				for _ , k :=range linkList {
					if k.Attrs().Index==index{
						viaInterface=k.Attrs().Name
						break
					}
				}
			}
			return gw , viaInterface ,  routeList[position] , nil 
		}else{
			return "" , "" , netlink.Route{} , fmt.Errorf("no default gw")
		}
		
	}

}



func CalculateIpv4RouteByDst( dst string ) ( []netlink.Route , error ){

	if CheckIPv4Format(dst)==false{
		return nil , fmt.Errorf("dst %v is not ipv4 address " , dst) 
	}

	return netlink.RouteGet( net.ParseIP(dst) ) 

}


//---------------- get ipv6 route --------------
/*
Route
https://godoc.org/github.com/vishvananda/netlink#Route
*/
func GetIpv6RouteAllEntryFromSpecifiedTable( tableNum int )([]netlink.Route , error ){

	routeFilter := &netlink.Route{
		Table: tableNum ,
	}

	routes, err := netlink.RouteListFiltered(  netlink.FAMILY_V6  , routeFilter , netlink.RT_FILTER_TABLE )
	if err != nil {
		return nil , err
	}
	return routes , nil

}


func GetIpv6RouteAllEntryFromAllTable( ) ([]netlink.Route , error ){

	return GetIpv6RouteAllEntryFromSpecifiedTable(unix.RT_TABLE_UNSPEC)

}

func GetIpv6RouteAllEntryFromMainTable( ) ([]netlink.Route , error ){

	return GetIpv6RouteAllEntryFromSpecifiedTable(unix.RT_TABLE_MAIN)

}


func GetIpv6RouteDefaultFromMainTable( ) ( gw , viaInterface string ,  detailRoute netlink.Route , e error ){

	if routeList, err:=GetIpv6RouteAllEntryFromMainTable(); err!=nil{
		return "" , "" , netlink.Route{} , err
	}else{
		min:=-1
		position:=-1
		for m , k:=range routeList {
			if k.Gw!=nil && k.Dst==nil  {
				//default route
				//check priority
				if min<0 {
					min=k.Priority
					position=m
				}else if k.Priority<min {
					min=k.Priority
					position=m
				}
			}
		}
		if position>=0{
			gw=routeList[position].Gw.String()
			if name , err := GetInterfaceNameByIndex( routeList[position].LinkIndex ) ; err!=nil {
				return "" , "" , netlink.Route{} , fmt.Errorf("failed to find the interface, info=%v " ,err  )
			}else{
				viaInterface=name
			}
			return gw , viaInterface ,  routeList[position] , nil 
		}else{
			return "" , "" , netlink.Route{} , fmt.Errorf("no default gw")
		}
		
	}

}



func CalculateIpv6RouteByDst( dst string ) ( []netlink.Route , error ){

	if CheckIPv6Format(dst)==false{
		return nil , fmt.Errorf("dst %v is not ipv6 address " , dst) 
	}

	return netlink.RouteGet( net.ParseIP(dst) ) 

}



//---------------- add and del ipv4 route --------------


func generateIPv4Route( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( netlink.Route , error ) {

	ToDstNet:=net.IPNet{}
	route := netlink.Route{  
		Table: tableNum,
	}

	// check input 
	if tableNum<=CONST_RouteTable_UNSPEC && tableNum>=CONST_RouteTable_MAX {
		return route , fmt.Errorf("tableNum %v outof range " , tableNum)
	}
	if len(dstNet)==0 && len(viaHost)==0 {
		return route , fmt.Errorf("miss dstNet and viaHost" )
	}


	if len(dstNet)!=0 {
		if CheckIPv4FormatWithMask(dstNet)==false{
			return route , fmt.Errorf("dstNet %v is not ipv4 subnet" , dstNet)
		}
		v := strings.Split(dstNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToDstNet=net.IPNet{IP: net.ParseIP( strings.Split(dstNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 32)}
		}else{
			return route , fmt.Errorf("failed to get mask from dstNet %v " )
		}		
		route.Dst=&ToDstNet
	}

	if len(viaHost)!=0 {
		if CheckIPv4Format(viaHost)==false{
			return route , fmt.Errorf("viaHost %v is not ipv4 address" , viaHost)
		}
		route.Gw=net.ParseIP(viaHost)
	}

	if len(viaInterface)!=0{
		if link , e:=GetInterfaceByName( viaInterface ); e!=nil {
			return route , fmt.Errorf("no interface %v " , viaInterface)
		}else{
			route.LinkIndex=link.Attrs().Index
		}
	}

	return route  , nil
}



/*
usage:
    //=== ip r a 5.0.0.0/8 dev dce-ovs table main ===
     table:=utility.CONST_RouteTable_MAIN
     dstNet:="5.0.0.0/8"
     viaHost:="" 
     viaInterface:="dce-ovs"

    //=== ip r a 5.0.0.1/32 dev dce-ovs table 100 ===
     table:=100
     dstNet:="5.0.0.1/32"
     viaHost:="" 
     viaInterface:="dce-ovs"

    //=== ip r a 6.0.0.0/8 via 172.16.0.211 dev dce-ovs table 100 ===
     table:=100
     dstNet:="6.0.0.0/8"
     viaHost:="172.16.0.211" 
     viaInterface:="dce-ovs"

    //=== ip r a default via 172.16.0.211 dev dce-ovs table 101 ===
    table:=101
    dstNet:=""
    viaHost:="172.16.0.211" 
    viaInterface:="dce-ovs"
*/

func CreateIPv4RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) {

	var route netlink.Route

	route , e:=generateIPv4Route( tableNum     , dstNet , viaHost , viaInterface  )
	if e!=nil{
		return e
	}

	//------------

	if err := netlink.RouteAdd(&route); err != nil {
		return err
	}
	return nil

}


func DelIPv4RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) {

	var route netlink.Route

	route , e:=generateIPv4Route( tableNum     , dstNet , viaHost , viaInterface  )
	if e!=nil{
		return e
	}

	//------------

	if err := netlink.RouteDel(&route); err != nil {
		return err
	}
	return nil

}




//---------------- add ipv6 route --------------

func generateIPv6Route( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( netlink.Route , error ) {

	ToDstNet:=net.IPNet{}
	route := netlink.Route{  
		Table: tableNum,
	}

	// check input 
	if tableNum<=CONST_RouteTable_UNSPEC && tableNum>=CONST_RouteTable_MAX {
		return route , fmt.Errorf("tableNum %v outof range " , tableNum)
	}
	if len(dstNet)==0 && len(viaHost)==0 {
		return route , fmt.Errorf("miss dstNet and viaHost" )
	}


	if len(dstNet)!=0 {
		if CheckIPv6FormatWithMask(dstNet)==false{
			return route , fmt.Errorf("dstNet %v is not ipv6 subnet" , dstNet)
		}
		v := strings.Split(dstNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToDstNet=net.IPNet{IP: net.ParseIP( strings.Split(dstNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 128)}
		}else{
			return route , fmt.Errorf("failed to get mask from dstNet %v " )
		}		
		route.Dst=&ToDstNet
	}

	if len(viaHost)!=0 {
		if CheckIPv6Format(viaHost)==false{
			return route , fmt.Errorf("viaHost %v is not ipv6 address" , viaHost)
		}
		route.Gw=net.ParseIP(viaHost)
	}

	if len(viaInterface)!=0{
		if link , e:=GetInterfaceByName( viaInterface ); e!=nil {
			return route , fmt.Errorf("no interface %v " , viaInterface)
		}else{
			route.LinkIndex=link.Attrs().Index
		}
	}

	return route  , nil
}


/*
usage:
    //=== ip r a fdde::/64 dev dce-ovs table main ===
    // table:=utility.CONST_RouteTable_MAIN
    // dstNet:="fdde::/64"
    // viaHost:="" 
    // viaInterface:="dce-ovs"

    //=== ip -6 r a fdde::22/128 dev dce-ovs table 101 ===
    // table:=101
    // dstNet:="fdde::22/128"
    // viaHost:="" 
    // viaInterface:="dce-ovs"

    //=== ip -6 r a fddd::/64 via fc02::11 dev dce-ovs table 101 ===
    // table:=101
    // dstNet:="fddd::/64"
    // viaHost:="fc02::11" 
    // viaInterface:="dce-ovs"

    //=== ip -6 r a default via fc02::11 dev dce-ovs table 101 ===
    // table:=101
    // dstNet:=""
    // viaHost:="fc02::11" 
    // viaInterface:="dce-ovs"
*/


func CreateIPv6RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) {

	var route netlink.Route

	route , e:=generateIPv6Route( tableNum     , dstNet , viaHost , viaInterface  )
	if e!=nil{
		return e
	}

	//------------

	if err := netlink.RouteAdd(&route); err != nil {
		return err
	}
	return nil

}


func DelIPv6RouteEntry( tableNum  int  , dstNet , viaHost , viaInterface string  ) ( error ) {

	var route netlink.Route

	route , e:=generateIPv6Route( tableNum     , dstNet , viaHost , viaInterface  )
	if e!=nil{
		return e
	}

	//------------

	if err := netlink.RouteDel(&route); err != nil {
		return err
	}
	return nil

}








