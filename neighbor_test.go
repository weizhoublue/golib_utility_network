package golib_utility_network

import(
    utility "github.com/weizhouBlue/golib_utility_network"
    "fmt"
    "testing"

    "github.com/vishvananda/netlink"
    "net"
)


//---------------

func Test_nei(t *testing.T){

    //utility.EnableLog=true

    mac1,_:=net.ParseMAC("66:60:f9:89:6b:9e")
    neighFilter:=netlink.Neigh{
        IP: net.ParseIP("172.28.28.4") ,
        HardwareAddr: mac1 ,
    }
    if entrys , err:= utility.GetNeighByFilter(neighFilter) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range entrys {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " entry = %v   \n" ,   k  )

        }
    }

            fmt.Printf( " ------------   \n\n"  )



    if entrys , err:= utility.GetNeighAll() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range entrys {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " entry = %v   \n" ,   k  )

        }
    }

    if entrys , err:= utility.GetNeighAllByIpv4() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range entrys {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " entry = %v   \n" ,   k  )

        }
    }



    if entrys , err:= utility.GetNeighAllByIpv6() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range entrys {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " entry = %v   \n" ,   k  )

        }
    }

}

func Test_nei2(t *testing.T){

    //utility.EnableLog=true


    ip:="172.18.0.165"
    if mac , viaInterface  , state, detail   , err:= utility.GetNeighByIp(ip) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " %v %v %v %v   \n" ,   mac , viaInterface , state , detail  )
    }

    ip="fe80::250:56ff:feb4:b4ec"
    if mac , viaInterface  , state , detail   , err:= utility.GetNeighByIp(ip) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " %v %v %v %v   \n" ,   mac , viaInterface , state , detail  )
    }


    mac:="00:50:56:b4:9b:6a"
    if ip , viaInterface  , state , detail   , err:= utility.GetNeighByMac(mac) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " %v %v %v %v   \n" ,   ip , viaInterface , state , detail  )
    }

}



func Test_neiAdd(t *testing.T){

    ip:="10.6.185.199"
    viaInterface:="dce-ext"
    mac:="00:22:33:44:55:11"
    if err:= utility.AddPermanentNeigh( ip , mac ,viaInterface  ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " create nei   \n"  )
    }


    if mac , viaInterface  , state, detail   , err:= utility.GetNeighByIp(ip) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " %v %v %v %v   \n" ,   mac , viaInterface , state , detail  )
    }

    if err:= utility.DelPermanentNeigh( ip , mac ,viaInterface  ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " del nei   \n"  )
    }


}




func Test_neidel(t *testing.T){


    utility.EnableLog=true


    ip:="10.6.185.192"
    viaInterface:="dce-ext"
    mac:="00:22:33:44:55:12"
    if err:= utility.AddPermanentNeigh( ip , mac ,viaInterface  ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " create nei   \n"  )
    }


    mac1,_:=net.ParseMAC(mac)
    link , _ :=utility.GetInterfaceByName( viaInterface )
    linkAttr:=link.Attrs()
    neighFilter:=netlink.Neigh{
        IP: net.ParseIP(ip) ,
        HardwareAddr: mac1 ,
        LinkIndex: linkAttr.Index ,
    }
    if err:= utility.DelNeighsByfilter(neighFilter) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " del entry  \n" )

    }


}
