package golib_utility_network

import(
    utility "github.com/weizhouBlue/golib_utility_network"
    "fmt"
    "testing"
    
)


//---------------

func Test_ipv4route(t *testing.T){



    if entrys , err:= utility.GetIpv4RouteAllEntryFromAllTable() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range entrys {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " entry = %v   \n" ,   k  )

        }
    }



}



func Test_ipv6route(t *testing.T){

    if entrys , err:= utility.GetIpv6RouteAllEntryFromMainTable() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range entrys {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " entry = %v   \n" ,   k  )

        }
    }




}



func Test_ipv4entry(t *testing.T){

    if gw , viaInterface , detail , err:= utility.GetIpv4RouteDefaultFromMainTable() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " entry = %v , %v , %v   \n" ,   gw , viaInterface , detail  )
    }


    if entry , err:= utility.CalculateIpv4RouteByDst("1.1.1.1") ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " entry = %v   \n" ,   entry  )
    }




    //=== ip r a 5.0.0.0/8 dev dce-ovs table main ===
    // table:=utility.CONST_RouteTable_MAIN
    // dstNet:="5.0.0.0/8"
    // viaHost:="" 
    // viaInterface:="dce-ovs"

    //=== ip r a 5.0.0.1/32 dev dce-ovs table 100 ===
    // table:=100
    // dstNet:="5.0.0.1/32"
    // viaHost:="" 
    // viaInterface:="dce-ovs"

    //=== ip r a 6.0.0.0/8 via 172.16.0.211 dev dce-ovs table 100 ===
    // table:=100
    // dstNet:="6.0.0.0/8"
    // viaHost:="172.16.0.211" 
    // viaInterface:="dce-ovs"

    //=== ip r a default via 172.16.0.211 dev dce-ovs table 101 ===
    table:=101
    dstNet:=""
    viaHost:="172.16.0.211" 
    viaInterface:="dce-ovs"

    if  err:= utility.CreateIPv4RouteEntry(  table , dstNet , viaHost , viaInterface ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " create ipv4 route \n"   )
    }

    if  err:= utility.DelIPv4RouteEntry(  table , dstNet , viaHost , viaInterface ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " del ipv4 route \n"   )
    }



}




func Test_ipv6entry(t *testing.T){

    if gw , viaInterface , detail , err:= utility.GetIpv6RouteDefaultFromMainTable() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " entry = %v , %v , %v   \n" ,   gw , viaInterface , detail  )
    }

    if entry , err:= utility.CalculateIpv6RouteByDst("2000::1") ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " entry = %v   \n" ,   entry  )
    }


    //=== ip r a fdde::/64 dev dce-ovs table main ===
    table:=utility.CONST_RouteTable_MAIN
    dstNet:="fdde::/64"
    viaHost:="" 
    viaInterface:="dce-ovs"

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

    if  err:= utility.CreateIPv6RouteEntry(  table , dstNet , viaHost , viaInterface ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " create ipv6 route \n"   )
    }

    if  err:= utility.DelIPv6RouteEntry(  table , dstNet , viaHost , viaInterface ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " del ipv6 route \n"   )
    }

}

