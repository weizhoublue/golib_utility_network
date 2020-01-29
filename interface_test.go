package golib_utility_network

import(
    utility "github.com/weizhouBlue/golib_utility_network"
    "fmt"
    "testing"
    
)


func Test_info(t *testing.T){

    if intList , err:= utility.GetAllInterface() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range intList {
            fmt.Printf( " ------------   \n"  )

            fmt.Printf( " interface type = %v   \n" ,   k.Type()  )
            fmt.Printf( " interface name = %v   \n" ,   k.Attrs().Name  )        
            fmt.Printf( " interface info = %v   \n" ,   k.Attrs()  )        
            fmt.Printf( " ------------   \n"  )

        }
    }

    name:="lo"
    if k , err:=utility.GetInterfaceByName(name) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " interface type = %v   \n" ,   k.Type()  )
            fmt.Printf( " interface name = %v   \n" ,   k.Attrs().Name  )        
            fmt.Printf( " interface info = %v   \n" ,   k.Attrs()  )    
    }



}


func Test_NetnsInfo(t *testing.T){

    if intList , err:= utility.GetAllInterfaceInNetns(10214) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range intList {
            fmt.Printf( " ------------   \n"  )

            fmt.Printf( " interface type = %v   \n" ,   k.Type()  )
            fmt.Printf( " interface name = %v   \n" ,   k.Attrs().Name  )        
            fmt.Printf( " interface info = %v   \n" ,   k.Attrs()  )        
            fmt.Printf( " ------------   \n"  )

        }
    }

    name:="eth0"
    if k , err:=utility.GetInterfaceByNameInNetns( 10214 , name) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " interface type = %v   \n" ,   k.Type()  )
            fmt.Printf( " interface name = %v   \n" ,   k.Attrs().Name  )        
            fmt.Printf( " interface info = %v   \n" ,   k.Attrs()  )    
    }



}



func Test_VethAdddel(t *testing.T){

    if err:=utility.CreateInterfaceVeth("welan" , "veth-welan") ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " create veth   \n"  ) 
    }


    if err:=utility.SetInterfaceToNetns( 10214 , "veth-welan"  ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " set veth to netns  \n"  ) 
    }


    if err:=utility.DelInterfaceByNameInNetns(10214, "veth-welan" ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " del veth   \n"  ) 
    }


}


func Test_VlanAdddel(t *testing.T){

    if err:=utility.CreateInterfaceVlan("dce-ovs" , "vlan100" , 100 ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " create vlan   \n"  ) 
    }

    if err:=utility.SetInterfaceNewName(  "vlan100" , "newvlan100" ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " set name    \n"  ) 
    }


    if err:=utility.SetInterfaceMtu(   "newvlan100"  , 1300 ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " set mtu    \n"  ) 
    }

    if mtu , err:=utility.GetInterfaceMtu(   "newvlan100"  ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " get mtu %v   \n" , mtu ) 
    }


    if err:=utility.SetInterfaceMac(   "newvlan100"  , "00:11:22:33:44:59" ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " set mac    \n"  ) 
    }

    if mac , err:=utility.GetInterfaceMac(   "newvlan100"  ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " get mac %v   \n" , mac ) 
    }


    if err:=utility.DelInterfaceByName( "newvlan100" ) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
            fmt.Printf( " del vlan   \n"  ) 
    }



}




