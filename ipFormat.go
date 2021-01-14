package golib_utility_network
import (
	"fmt"
	"net"
	"strings"
	"strconv"
	"github.com/mdlayher/ndp"
)
// https://godoc.org/net


/*
func CheckIPFormat( ip string ) bool 
func CheckIPv4Format( ip string ) bool
func CheckIPv6Format( ip string ) bool 

func CheckIPv4FormatWithMask( ip string ) (bool ) 
func CheckIPv6FormatWithMask( ip string ) (bool ) 
func CheckIPv6v4FormatWithMask( ip string ) (bool ) 

func CheckIPv4SubnetWithMask( ip string ) bool
func CheckIPv6SubnetWithMask( ip string ) bool

func EqualIPv6v4( ip1 , ip2 string ) bool 


func MaskIPv4( ip , mask string ) (string,error) {
func MaskIPv6( ip , mask string ) (string , error) 

func CheckSameIPv4Subnet( subnet1 , subnet2 , mask string ) ( bool , error) {
func CheckSameIPv6Subnet( subnet1 , subnet2  , mask string ) ( bool , error)

func CheckIPv4SubnetOverlay( subnet1 , subnet2 string ) ( overlay bool , err error ) 
func CheckIPv6SubnetOverlay( subnet1 , subnet2 string ) ( overlay bool , err error ) 

func CheckIPTypeUnicast( ip string ) (bool , error) 
func CheckIPTypeLoopback( ip string ) (bool , error) 
func CheckIPTypeUnspecified( ip string ) (bool , error) 
func CheckIPTypeLinkLocalUnicast( ip string ) (bool , error) 

func ConvertIPv6FullFormat() string

*/


//==============================
//1.1.1.1
//fc00::
func CheckIPFormat( ip string ) bool {
	result := net.ParseIP(ip)
	if result==nil {
		return false
	}
	return true
}

//1.1.1.1
func CheckIPv4Format( ip string ) bool {
	result := net.ParseIP(ip)
	if result==nil {
		return false
	}
	if result.To4()==nil {
		return false
	}
	return true
}

// 1.1.0.0/16
// 1.1.0.1/16
func CheckIPv4FormatWithMask( ip string ) (bool ) {
	re:=strings.Split(ip, "/" )
	if len(re)!=2{
		return false 
	}

	if s, err := strconv.ParseInt(re[1], 10, 64); err == nil {
	    if s>32 || s<1 {
	    	return false
	    } 
	}else{
		return false 
	}


	if CheckIPv4Format(re[0]) {
		return true
	}else{
		return false
	}

}


//fc00::
func CheckIPv6Format( ip string ) bool {
	result := net.ParseIP(ip)
	if result==nil {
		return false
	}
	if result.To4()==nil {
		return true
	}
	return false
}


// fc00::1/64
// fc00::0/64
func CheckIPv6FormatWithMask( ip string ) (bool ) {
	re:=strings.Split(ip, "/" )
	if len(re)!=2{
		return false 
	}

	if s, err := strconv.ParseInt(re[1], 10, 64); err == nil {
	    if s>128 || s<1 {
	    	return false
	    } 
	}else{
		return false 
	}

	if CheckIPv6Format(re[0]) {
		return true
	}else{
		return false
	}

}


func CheckIPv6v4FormatWithMask( ip string ) (bool ) {
	if CheckIPv6FormatWithMask(ip) {
		return true
	}
	if CheckIPv4FormatWithMask(ip) {
		return true
	}
	return false
}





// 1.1.0.0/16
func CheckIPv4SubnetWithMask( ip string ) bool {
	if CheckIPv4FormatWithMask(ip)==false{
		return false
	}

	re:=strings.Split(ip, "/" )
	if len(re)!=2{
		return false 
	}

	subnet1 , e2:= MaskIPv4(re[0], re[1] ) 
	if e2!=nil {
		return false
	}

	if EqualIPv6v4(subnet1 , re[0] ) ==false {
		return false
	}
	return true
	
	
}

// fc00::0/64
func CheckIPv6SubnetWithMask( ip string ) bool {
	if CheckIPv6FormatWithMask(ip)==false{
		return false
	}

	re:=strings.Split(ip, "/" )
	if len(re)!=2{
		return false 
	}

	subnet1 , e2:= MaskIPv6(re[0],re[1] ) 
	if e2!=nil {
		return false
	}

	if EqualIPv6v4(subnet1 , re[0] ) ==false {
		return false
	}
	return true
	
}




//================================

// fd00::  fd00:0::0 -> true
func EqualIPv6v4( ip1 , ip2 string ) bool {

	r1 := net.ParseIP(ip1)
	if r1==nil {
		return false
	}

	r2 := net.ParseIP(ip2)
	if r2==nil {
		return false
	}

	if r1.Equal(r2) {
		return true
	}
	return false
}





//1.1.1.1  , 16 -> 1.1.1.0
func MaskIPv4( ip , mask string ) (string,error) {

	length , e := strconv.ParseInt(mask, 10, 32)
	if e != nil {
	    return "", fmt.Errorf("wrong mask %v , info=%v ", mask , e ) 
	}

	if length <0 || length > 32 {
		return "", fmt.Errorf("error subnet length=%v " , length )
	}


	if ! CheckIPv4Format(ip) {
		return "" , fmt.Errorf("error ip=%v " , ip ) 
	}
	//fmt.Println("ok")
	r:=net.ParseIP(ip)
	to:=r.Mask( net.CIDRMask( int(length) , 32)  )
	return to.String() , nil
}

//fc00:0:0:1::  , 16 ->  fc00::
func MaskIPv6( ip , mask string ) (string , error) {

	length , e := strconv.ParseInt(mask, 10, 32)
	if e != nil {
	    return "", fmt.Errorf("wrong mask %v , info=%v ", mask , e ) 
	}

	if length <0 || length > 128 {
		return "", fmt.Errorf("error subnet length=%v " , length )
	}

	if ! CheckIPv6Format(ip) {
		return "" , fmt.Errorf("error ip=%v " , ip ) 
	}
	//fmt.Println("ok")
	r:=net.ParseIP(ip)
	to:=r.Mask( net.CIDRMask( int(length) , 128)  )
	return to.String() , nil
}



//================================
//1.1.1.0 , 1.1.0.0 , 16 -> true
func CheckSameIPv4Subnet( subnet1 , subnet2 , mask string ) ( bool , error) {
	var err error
	var one , two string
	if one , err =MaskIPv4(subnet1 , mask) ; err!=nil{
		return false, err
	}
	if two , err =MaskIPv4(subnet2 , mask) ; err!=nil{
		return false, err
	}
	if one==two {
		return true , nil
	}else{
		return false, nil
	}
}

//fc00:0:0:1:: , fc00:0:0:2::  , 64 ->  false
func CheckSameIPv6Subnet( subnet1 , subnet2 , mask string ) ( bool , error) {

	var err error
	var one , two string
	if one , err =MaskIPv6(subnet1 , mask) ; err!=nil{
		return false, err
	}
	if two , err =MaskIPv6(subnet2 , mask) ; err!=nil{
		return false, err
	}
	if one==two {
		return true , nil
	}else{
		return false, nil
	}
}


//================================

func CheckIPv4SubnetOverlay( subnet1 , subnet2 string ) ( overlay bool , err error ) {

	if ! CheckIPv4FormatWithMask( subnet1 )  {
		return false , fmt.Errorf("error ipv4 subnet=%v " , subnet1 )
	}
	if ! CheckIPv4FormatWithMask( subnet2 )  {
		return false , fmt.Errorf("error ipv4 subnet=%v " , subnet2 )
	}

	ip1:=strings.Split( subnet1 , "/" )[0]
	len1:=strings.Split( subnet1 , "/" )[1]

	ip2:=strings.Split( subnet2 , "/" )[0]
	len2:=strings.Split( subnet2 , "/" )[1]

	re , er := CheckSameIPv4Subnet( ip1 , ip2 , len1 )
	if er!=nil {
		return false , er
	}
	if re{
		return true , nil
	}

	re , er = CheckSameIPv4Subnet( ip1 , ip2 , len2 )
	if er!=nil {
		return false , er
	}
	if re{
		return true , nil
	}

	return false , nil
}



func CheckIPv6SubnetOverlay( subnet1 , subnet2 string ) ( overlay bool , err error ) {

	if ! CheckIPv6FormatWithMask( subnet1 )  {
		return false , fmt.Errorf("error ipv6 subnet=%v " , subnet1 )
	}
	if ! CheckIPv6FormatWithMask( subnet2 )  {
		return false , fmt.Errorf("error ipv6 subnet=%v " , subnet2 )
	}

	ip1:=strings.Split( subnet1 , "/" )[0]
	len1:=strings.Split( subnet1 , "/" )[1]

	ip2:=strings.Split( subnet2 , "/" )[0]
	len2:=strings.Split( subnet2 , "/" )[1]

	re , er := CheckSameIPv6Subnet( ip1 , ip2 , len1 )
	if er!=nil {
		return false , er
	}
	if re{
		return true , nil
	}

	re , er = CheckSameIPv6Subnet( ip1 , ip2 , len2 )
	if er!=nil {
		return false , er
	}
	if re{
		return true , nil
	}

	return false , nil
}



//================================
//10.1.1.1 -> true
// fc00::1 result =  true
// 2000::1 result =  true
// fec0::1 result =  true
// FC00::1 result =  true
// Fe80::1 result =  false
// 00aa::1 result =  true
//  ff00::1 result =  false
func CheckIPTypeUnicast( ip string ) (bool , error) {
	result := net.ParseIP(ip)
	if result==nil {
		return false, fmt.Errorf("error ip=%v " , ip )
	}
	return result.IsGlobalUnicast() , nil
}

//127.0.0.1 -> true
//::1 -> true
func CheckIPTypeLoopback( ip string ) (bool , error) {
	result := net.ParseIP(ip)
	if result==nil {
		return false, fmt.Errorf("error ip=%v " , ip )
	}
	return result.IsLoopback() , nil
}

//0.0.0.0 -> true
//:: -> true
func CheckIPTypeUnspecified( ip string ) (bool , error) {
	result := net.ParseIP(ip)
	if result==nil {
		return false, fmt.Errorf("error ip=%v " , ip )
	}
	return result.IsUnspecified() , nil
}

// 169.254.x.x -> true
// fe80::1 -> true
func CheckIPTypeLinkLocalUnicast( ip string ) (bool , error) {
	result := net.ParseIP(ip)
	if result==nil {
		return false, fmt.Errorf("error ip=%v " , ip )
	}
	return result.IsLinkLocalUnicast() , nil
}



 // ip fd11:1::2 to fd1101000000000002 
func ConvertIPv6FullFormat( ip string ) string  {
	result := net.ParseIP(ip)
	if result==nil {
		return ""
	}
	if result.To4()!=nil {
		return ""
	}
	a:=[]byte( result.To16() )

	b:=""
	for _ , v:=range a {
		b+=fmt.Sprintf("%02x", v)
	}
	return b
}

// fc01::1111:2345 -> ff02::1:ff11:2345
func GetIPv6SolicitedNodeAddress( unicastIp string ) (  string , error ) {
	if CheckIPv6Format( unicastIp )==false {
		return "" , fmt.Errorf("invalid ipv6 address")
	}

	result := net.ParseIP(unicastIp)
	if result==nil {
		return "" , fmt.Errorf("failed to parse ipv6 address")
	}

	if s , e:= ndp.SolicitedNodeMulticast(result) ; e!=nil || len(s.String())==0 {
		return "" , fmt.Errorf("fail: %v" , e)
	}else{
		return s.String() , nil
	}

}
