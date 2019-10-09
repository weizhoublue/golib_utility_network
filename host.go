package golib_utility_network
import (
    "fmt"
    "net"
    //"strings"
)


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





