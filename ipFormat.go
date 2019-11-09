package golib_utility_network
import (
	"fmt"
	"net"
	"strings"
	"strconv"
)
// https://godoc.org/net

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


// fc00::/64
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
 

//================================
//1.1.1.1  , 16 -> 1.1.1.0
func MaskIPv4( ip string, mask int ) (string,error) {
	if mask <0 || mask > 32 {
		return "", fmt.Errorf("error subnet length=%v " , mask )
	}

	if ! CheckIPv4Format(ip) {
		return "" , fmt.Errorf("error ip=%v " , ip ) 
	}
	fmt.Println("ok")
	r:=net.ParseIP(ip)
	to:=r.Mask( net.CIDRMask(mask, 32)  )
	return to.String() , nil
}

//fc00:0:0:1::  , 64 ->  fc00::
func MaskIPv6( ip string, mask int ) (string , error) {
	if mask <0 || mask > 128 {
		return "", fmt.Errorf("error subnet length=%v " , mask )
	}
	if ! CheckIPv6Format(ip) {
		return "" , fmt.Errorf("error ip=%v " , ip ) 
	}
	fmt.Println("ok")
	r:=net.ParseIP(ip)
	to:=r.Mask( net.CIDRMask(mask, 128)  )
	return to.String() , nil
}



//================================
//1.1.1.0 , 1.1.0.0 , 16 -> true
func CheckSameIPv4Subnet( subnet1 , subnet2 string , mask int) ( bool , error) {
	if mask <0 || mask >32 {
		return false , fmt.Errorf("error subnet length=%v " , mask )
	}
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
func CheckSameIPv6Subnet( subnet1 , subnet2 string , mask int) ( bool , error) {
	if mask <0 || mask >128  {
		return false , fmt.Errorf("error subnet length=%v " , mask )
	}
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
	mask1 , _ :=strconv.ParseInt(len1, 10, 64)

	ip2:=strings.Split( subnet2 , "/" )[0]
	len2:=strings.Split( subnet2 , "/" )[1]
	mask2,_:=strconv.ParseInt(len2, 10, 64)

	re , er := CheckSameIPv4Subnet( ip1 , ip2 , int(mask1) )
	if er!=nil {
		return false , er
	}
	if re{
		return true , nil
	}

	re , er = CheckSameIPv4Subnet( ip1 , ip2 , int(mask2) )
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
	mask1 , _ :=strconv.ParseInt(len1, 10, 64)

	ip2:=strings.Split( subnet2 , "/" )[0]
	len2:=strings.Split( subnet2 , "/" )[1]
	mask2,_:=strconv.ParseInt(len2, 10, 64)

	re , er := CheckSameIPv6Subnet( ip1 , ip2 , int(mask1) )
	if er!=nil {
		return false , er
	}
	if re{
		return true , nil
	}

	re , er = CheckSameIPv6Subnet( ip1 , ip2 , int(mask2) )
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
//fc00:: -> true
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

// fe80::1 -> true
func CheckIPTypeLinkLocalUnicast( ip string ) (bool , error) {
	result := net.ParseIP(ip)
	if result==nil {
		return false, fmt.Errorf("error ip=%v " , ip )
	}
	return result.IsLinkLocalUnicast() , nil
}


