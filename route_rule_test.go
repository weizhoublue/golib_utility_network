package golib_utility_network

import(
    utility "github.com/weizhouBlue/golib_utility_network"
    "fmt"
    "testing"
    
)


//---------------

func Test_ipv4rule(t *testing.T){

    if ruleList , err:= utility.GetIPv4RouteRule() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range ruleList {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " rule = %v   \n" ,   k  )

        }
    }


    if ruleList , err:= utility.GetIPv4RouteRuleInNetns( 10214 ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range ruleList {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " rule = %v   \n" ,   k  )

        }
    }


    tableNum:=100 
    tablePriority:=3000  
    logicalNot:=false 
    srcNet:="172.10.0.0/16"  
    dstNet:="172.11.0.0/16"  
    InIf:="dce-ext"
    OutIf:="dce-ovs"
    // tableNum:=100 
    // tablePriority:=3000  
    // logicalNot:=false 
    // srcNet:="172.10.0.0/16"  
    // dstNet:=""  
    // InIf:=""
    // OutIf:=""
    if err:=utility.CreateIPv4RouteRule( tableNum , tablePriority  , logicalNot , srcNet , dstNet  , InIf , OutIf  ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " create rule  \n"  )
    }

    // we can just use the tableNum and tablePriority to delete the rule
    // tableNum=100
    // tablePriority=3000  
    // logicalNot=false 
    // srcNet=""  
    // dstNet=""  
    if err:=utility.DelIPv4RouteRule( tableNum , tablePriority  , logicalNot , srcNet , dstNet   ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " del rule  \n"  )
    }


}





func Test_ipv6rule(t *testing.T){

    if ruleList , err:= utility.GetIPv6RouteRule() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range ruleList {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " get rule = %v   \n" ,   k  )

        }
    }

    if ruleList , err:= utility.GetIPv6RouteRuleInNetns( 10214 ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        for _ , k := range ruleList {
            fmt.Printf( " ------------   \n"  )
            // https://godoc.org/github.com/vishvananda/netlink#Rule
            fmt.Printf( " get netns rule = %v   \n" ,   k  )

        }
    }


    tableNum:=100 
    tablePriority:=3008
    logicalNot:=false 
    srcNet:="fd01::/64"
    dstNet:="fd02::/64"
    InIf:="dce-ext"
    OutIf:="dce-ovs"
    // tableNum:=100 
    // tablePriority:=3000  
    // logicalNot:=false 
    // srcNet:="fd01::/64"  
    // dstNet:=""  
    // InIf:=""
    //OutIf:=""
    if err:=utility.CreateIPv6RouteRule( tableNum , tablePriority  , logicalNot , srcNet , dstNet  , InIf , OutIf  ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " create rule  \n"  )
    }

    if err:=utility.DelIPv6RouteRule( tableNum , tablePriority  , logicalNot , srcNet , dstNet   ) ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( " del rule  \n"  )
    }

}



//---------------




