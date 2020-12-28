package golib_utility_network

import(
    utility "github.com/weizhouBlue/golib_utility_network"
    "fmt"
    "testing"
)


func Test_ipformat(t *testing.T){

	var ip string 

	//ipv4 or ipv6
    ip="1.0.0.0"
    fmt.Printf( "%s is ip ? %t \n" , ip , utility.CheckIPFormat( ip )   )
    
    ip="fc80::0"
    fmt.Printf( "%s is ip ? %t \n" , ip , utility.CheckIPFormat( ip )   )


    //
    m:="172.19.0.0/16"
    fmt.Printf( "%s is ip with mask ? %t \n" , m , utility.CheckIPv4FormatWithMask( m )   )

    m="fc00::/64"
    fmt.Printf( "%s is ip with mask ? %t \n" , m , utility.CheckIPv6FormatWithMask( m )   )


	//ipv4
    ip="1.0.0.0"
    fmt.Printf( "%s is ipv4 ? %t \n" , ip , utility.CheckIPv4Format( ip )   )
    
    ip="fc80::0"
    fmt.Printf( "%s is ipv4 ? %t \n" , ip , utility.CheckIPv4Format( ip )   )

    ip="test"
    fmt.Printf( "%s is ipv4 ? %t \n" , ip , utility.CheckIPv4Format( ip )   )


    //ipv6
    ip="fc80::0"
    fmt.Printf( "%s is ipv6 ? %t \n" , ip , utility.CheckIPv6Format( ip )   )
    
    ip="1.0.0.0"
    fmt.Printf( "%s is ipv6 ? %t \n" , ip , utility.CheckIPv6Format( ip )   )
    

    //
    m="fc00::/64"
    fmt.Printf( "%s is ip with mask ? %t \n" , m , utility.CheckIPv6v4FormatWithMask( m )   )
    m="172.19.0.0/16"
    fmt.Printf( "%s is ip with mask ? %t \n" , m , utility.CheckIPv6v4FormatWithMask( m )   )

    //
    m="172.19.0.0/16"
    fmt.Printf( "%s is subnet with mask ? %t \n" , m , utility.CheckIPv4SubnetWithMask( m )   )
    m="fd00::0/64"
    fmt.Printf( "%s is subnet with mask ? %t \n" , m , utility.CheckIPv6SubnetWithMask( m )   )





}



func Test_maskip(t *testing.T){

	var ip  string 
	//var mask int

	//ipv4 or ipv6
    ip="1.1.1.0"
    if result , err:=utility.MaskIPv4( ip , "34" ) ; err !=nil {
    	fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
    	fmt.Printf( "good ,  %v \n" ,    result )    	
    }


	//ipv4 or ipv6
    ip="fc00:0:0:1::"
    if result , err:=utility.MaskIPv6( ip , "16" ) ; err !=nil {
    	fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
    	fmt.Printf( "good ,  %v \n" ,    result )    	
    }
    


    //ipv4 or ipv6
    ip1:="fc00::"
    ip2:="fc00:0::0"
    re:=utility.EqualIPv6v4( ip1 , ip2 )  
    
    fmt.Printf( "  %v == %v ? %v \n" , ip1 , ip2 ,   re )

   


}


func Test_samesubnet(t *testing.T){

    ip1:="1.1.1.0"
    ip2:="1.1.0.0"
    length1:="24"
    //length1:=16
    if result , err:=utility.CheckSameIPv4Subnet( ip1, ip2 , length1 ) ; err !=nil {
    	fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
    	fmt.Printf( "result =  %v \n" ,    result )    	
    }



    ip3:="fc00:0:0:1::"
    ip4:="fc00:0:0:2::"
    length2:="0"
    //length1:=16
    if result , err:=utility.CheckSameIPv6Subnet( ip3, ip4 , length2 ) ; err !=nil {
    	fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
    	fmt.Printf( "result =  %v \n" ,    result )    	
    }


    sub1:="1.1.0.0/16"
    sub2:="1.0.0.0/8"
    if result , err:=utility.CheckIPv4SubnetOverlay( sub1, sub2  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( "CheckIPv4SubnetOverlay result =  %v \n" ,    result )     
    }


    sub1="fc00:0:0:1::/64"
    sub2="fc00::/16"
    if result , err:=utility.CheckIPv6SubnetOverlay( sub1, sub2  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( "CheckIPv6SubnetOverlay result =  %v \n" ,    result )     
    }


}


func Test_ipType(t *testing.T){


    //ip1:="10.1.1.1"
    ip1:="fc00::1"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="2000::1"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="fec0::1"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="FC00::0"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="Fe80::1"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="00aa::1"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="ff00::1"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }

    ip1="0.0.0.0"
    if result , err:=utility.CheckIPTypeUnicast( ip1  ) ; err !=nil {
        fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
        fmt.Printf( " %v result =  %v \n" ,   ip1 ,  result )     
    }



    //------------------------------------------



    ip2:="127.0.0.1"
    //ip2:="::1"
    if result , err:=utility.CheckIPTypeLoopback( ip2  ) ; err !=nil {
    	fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
    	fmt.Printf( "result =  %v \n" ,    result )    	
    }

    ip3:="::"
    //ip3:="0.0.0.0"
    if result , err:=utility.CheckIPTypeUnspecified( ip3  ) ; err !=nil {
    	fmt.Printf( "failed ,  %v \n" ,    err )
    }else{
    	fmt.Printf( "result =  %v \n" ,    result )    	
    }


}



func Test_interface(t *testing.T){

    name:="ens192"
    if ok :=utility.CheckHostInterfaceByName( name  ) ; ok  {
    	fmt.Printf( "exist   %v \n" ,    name )
    }else{
    	fmt.Printf( "no %v \n" ,    name )    	
    }



    name="utun0"
    if mac , err :=utility.GetHostInterfaceMac( name  ) ; err==nil  {
    	fmt.Printf( "mac=   %v \n" ,    mac )
    }else{
    	fmt.Printf( "error= %v \n" ,    err )    	
    }




    name="utun2"
    if ok , err :=utility.CheckHostInterfaceUp( name  ) ; err==nil  {
    	fmt.Printf( "up? =   %v \n" ,    ok )
    }else{
    	fmt.Printf( "error= %v \n" ,    err )    	
    }


    name="utun2"
    if mtu , err :=utility.GetHostInterfaceMtu( name  ) ; err==nil  {
    	fmt.Printf( "mtu =   %v \n" ,  mtu )
    }else{
    	fmt.Printf( "error= %v \n" ,    err )    	
    }

    if int_list , err :=utility.ListHostAllInterfaces(  ) ; err==nil  {
    	fmt.Printf( "interface list =   %v \n" ,  int_list )
    }else{
    	fmt.Printf( "error= %v \n" ,    err )    	
    }


    if int_list , err :=utility.ListHostUpInterfaces(  ) ; err==nil  {
    	fmt.Printf( "up interface list =   %v \n" ,  int_list )
    }else{
    	fmt.Printf( "error= %v \n" ,    err )    	
    }



    if ipv4list , ipv6list , err :=utility.GetInterfaceUnicastAddrByName( "en0" ) ; err==nil  {
    	fmt.Printf( "ip list =  %v  , %v \n" ,    ipv4list , ipv6list )
    }else{
    	fmt.Printf( "error= %v   \n" ,    err  )    	
    }


    if ipv4list , ipv6list , err :=utility.GetAllInterfaceUnicastAddrByName(  ) ; err==nil  {
    	fmt.Printf( "ip list =  %v  , %v \n" ,    ipv4list , ipv6list )
    }else{
    	fmt.Printf( "error= %v   \n" ,    err  )    	
    }



    if ipv4list , ipv6list , err :=utility.GetAllUpInterfaceUnicastAddrByName(  ) ; err==nil  {
    	fmt.Printf( "ip list =  %v  , %v \n" ,    ipv4list , ipv6list )
    }else{
    	fmt.Printf( "error= %v   \n" ,    err  )    	
    }

}




func Test_type(t *testing.T){


    ip:="fe80::1"
    if r , e:=utility.CheckIPTypeLinkLocalUnicast( ip ) ; e==nil {
        fmt.Printf( "ip %v is local type ? %v \n" ,    ip , r )
    }else{
        fmt.Printf( "error= %v   \n" ,    e  )        
    }

    ip="169.254.1.1"
    if r , e:=utility.CheckIPTypeLinkLocalUnicast( ip ) ; e==nil {
        fmt.Printf( "ip %v is local type ? %v \n" ,    ip , r )
    }else{
        fmt.Printf( "error= %v   \n" ,    e  )        
    }


}


func Test_link(t *testing.T){

    if intList , err:= utility.GetAllInterface() ; err!=nil{
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{

        fmt.Printf( "all interface = %v   \n" ,   intList  )        

    }

    name:="lo0"
    if err:=utility.SetInterfaceUp(name) ;  err!=nil {
        fmt.Printf( "error= %v   \n" ,   err  )        
    }else{
        fmt.Printf( "up %v   \n" ,   name  )        
    }



}



func Test_ip6(t *testing.T){

    ip:="fd11:1::2"
    if new := utility.ConvertIPv6FullFormat(ip) ; len(new)>0 {
        fmt.Printf( "ip %v to %v   \n" ,  ip, new  )        
    }else{
        fmt.Printf( "failed %v   \n" ,   ip  )
    }

}


