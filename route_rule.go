package golib_utility_network
import (
    "fmt"
    "github.com/vishvananda/netlink"
    "github.com/vishvananda/netns"
    //"os"
    "runtime"
    "net"
    "golang.org/x/sys/unix"
    "strconv"
    "strings"
)

// attention ! this should be used on linux os 


/*
doc  https://godoc.org/github.com/vishvananda/netlink
	https://godoc.org/github.com/vishvananda/netns

github https://github.com/vishvananda/netlink

refere example : 
https://github.com/vishvananda/netlink/blob/master/rule_test.go
https://github.com/vishvananda/netlink/blob/master/route_test.go
https://github.com/vishvananda/netlink/blob/master/netns_test.go

*/

/*
func GetIPv4RouteRule() ( []netlink.Rule , error ) 
func GetIPv4RouteRuleInNetns( pid int ) ( []netlink.Rule , error ) 
func GetIPv6RouteRule() ( []netlink.Rule , error ) 
func GetIPv6RouteRuleInNetns( pid int ) ( []netlink.Rule , error )

func CreateIPv4RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error )
func CreateIPv4RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) 
func CreateIPv6RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) 
func CreateIPv6RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) 

func DelIPv4RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string ) ( error )
func DelIPv4RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) 
func DelIPv6RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string  ) ( error )
func DelIPv6RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error )

*/

//-------------------------------

/*
Rule
https://godoc.org/github.com/vishvananda/netlink#Rule
*/
func GetIPv4RouteRule() ( []netlink.Rule , error ) {

	rules, err := netlink.RuleList(unix.AF_INET)
	if err != nil {
		return nil , err
	}
	return rules , nil 
}


func GetIPv4RouteRuleInNetns( pid int ) ( []netlink.Rule , error ) {

	// when thread operate netns , goroutine should lock the thread in case that the operation goroutine switch to other thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return nil , e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return nil , e
		}
		return GetIPv4RouteRule( )
	}
}


func GetIPv6RouteRule() ( []netlink.Rule , error ) {

	rules, err := netlink.RuleList(unix.AF_INET6)
	if err != nil {
		return nil , err
	}
	return rules , nil 
}


func GetIPv6RouteRuleInNetns( pid int ) ( []netlink.Rule , error ) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return nil , e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return nil , e
		}
		return GetIPv6RouteRule( )
	}
}

//-------------------------------
/*
Rule
https://godoc.org/github.com/vishvananda/netlink#Rule
*/
func CreateIPv4RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) {

	ToSrcNet:=net.IPNet{}
	ToDstNet:=net.IPNet{}

	// check input 
	if tableNum<=CONST_RouteTable_UNSPEC && tableNum>=CONST_RouteTable_MAX {
		return fmt.Errorf("tableNum %v outof range " , tableNum)
	}

	if tablePriority<0 && tablePriority>32767 {
		return fmt.Errorf("tablePriority %v outof range 0-32767 " , tablePriority)
	}

	if len(srcNet)!=0 {
		if CheckIPv4FormatWithMask(srcNet)==false{
			return fmt.Errorf("srcNet %v is not ipv4 subnet" , srcNet)
		}
		v := strings.Split(srcNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToSrcNet=net.IPNet{IP: net.ParseIP( strings.Split(srcNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 32)}
		}else{
			return fmt.Errorf("failed to get mask from srcNet %v " )
		}		
	}

	if len(dstNet)!=0 {
		if CheckIPv4FormatWithMask(dstNet)==false{
			return fmt.Errorf("dstNet %v is not ipv4 subnet" , dstNet)
		}
		v := strings.Split(dstNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToDstNet=net.IPNet{IP: net.ParseIP( strings.Split(dstNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 32)}
		}else{
			return fmt.Errorf("failed to get mask from dstNet %v " )
		}		
	}

	if len(InIf)!=0 {
		if _ , e := GetInterfaceByName(InIf);e!=nil {
			return fmt.Errorf("interface %v does not exit" , InIf)
		}
	}

	if len(OutIf)!=0 {
		if _ , e := GetInterfaceByName(OutIf);e!=nil {
			return fmt.Errorf("interface %v does not exit" , OutIf)
		}
	}

	//------------

	rulesBegin, err := netlink.RuleList(unix.AF_INET)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}

	rule := netlink.NewRule()
	rule.Table = tableNum
	rule.Family = unix.AF_INET
	rule.Src = &ToSrcNet
	rule.Dst = &ToDstNet
	rule.Priority = tablePriority
	rule.OifName = OutIf
	rule.IifName = InIf
	rule.Invert = logicalNot

	if err := netlink.RuleAdd(rule); err != nil {
		return err
	}

	rules, err := netlink.RuleList(unix.AF_INET)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}	
	if len(rules) != len(rulesBegin)+1 {
		return fmt.Errorf("Rule not added properly" )
	}

	return nil
}


func CreateIPv4RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return   e
		}
		return CreateIPv4RouteRule( tableNum , tablePriority   , logicalNot   , srcNet , dstNet   , InIf , OutIf    )
	}

}


func CreateIPv6RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) {

	ToSrcNet:=net.IPNet{}
	ToDstNet:=net.IPNet{}

	// check input 
	if tableNum<=CONST_RouteTable_UNSPEC && tableNum>=CONST_RouteTable_MAX {
		return fmt.Errorf("tableNum %v outof range" , tableNum)
	}

	if tablePriority<0 && tablePriority>32767 {
		return fmt.Errorf("tablePriority %v outof range 0-32767 " , tablePriority)
	}

	if len(srcNet)!=0 {
		if CheckIPv6FormatWithMask(srcNet)==false{
			return fmt.Errorf("srcNet %v is not ipv6 subnet" , srcNet)
		}
		v := strings.Split(srcNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToSrcNet=net.IPNet{IP: net.ParseIP( strings.Split(srcNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 128)}
		}else{
			return fmt.Errorf("failed to get mask from srcNet %v " )
		}		
	}

	if len(dstNet)!=0 {
		if CheckIPv6FormatWithMask(dstNet)==false{
			return fmt.Errorf("dstNet %v is not ipv6 subnet" , dstNet)
		}
		v := strings.Split(dstNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToDstNet=net.IPNet{IP: net.ParseIP( strings.Split(dstNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 128 )}
		}else{
			return fmt.Errorf("failed to get mask from dstNet %v " )
		}		
	}

	if len(InIf)!=0 {
		if _ , e := GetInterfaceByName(InIf);e!=nil {
			return fmt.Errorf("interface %v does not exit" , InIf)
		}
	}

	if len(OutIf)!=0 {
		if _ , e := GetInterfaceByName(OutIf);e!=nil {
			return fmt.Errorf("interface %v does not exit" , OutIf)
		}
	}

	//------------

	rulesBegin, err := netlink.RuleList(unix.AF_INET6)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}

	rule := netlink.NewRule()
	rule.Table = tableNum
	rule.Family = unix.AF_INET6
	rule.Src = &ToSrcNet
	rule.Dst = &ToDstNet
	rule.Priority = tablePriority
	rule.OifName = OutIf
	rule.IifName = InIf
	rule.Invert = logicalNot

	if err := netlink.RuleAdd(rule); err != nil {
		return err
	}

	rules, err := netlink.RuleList(unix.AF_INET6)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}	
	if len(rules) != len(rulesBegin)+1 {
		return fmt.Errorf("Rule not added properly" )
	}

	return nil
}


func CreateIPv6RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return   e
		}
		return CreateIPv6RouteRule( tableNum , tablePriority   , logicalNot   , srcNet , dstNet   , InIf , OutIf    )
	}

}


//-------------------------------

// we can just use the tableNum and tablePriority to delete the rule
func DelIPv4RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string ) ( error ) {

	ToSrcNet:=net.IPNet{}
	ToDstNet:=net.IPNet{}

	// check input 
	if tableNum<=CONST_RouteTable_UNSPEC && tableNum>=CONST_RouteTable_MAX {
		return fmt.Errorf("tableNum %v outof range  " , tableNum)
	}

	if tablePriority<0 && tablePriority>32767 {
		return fmt.Errorf("tablePriority %v outof range 0-32767 " , tablePriority)
	}

	if len(srcNet)!=0 {
		if CheckIPv4FormatWithMask(srcNet)==false{
			return fmt.Errorf("srcNet %v is not ipv4 subnet" , srcNet)
		}
		v := strings.Split(srcNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToSrcNet=net.IPNet{IP: net.ParseIP( strings.Split(srcNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 32)}
		}else{
			return fmt.Errorf("failed to get mask from srcNet %v " )
		}		
	}

	if len(dstNet)!=0 {
		if CheckIPv4FormatWithMask(dstNet)==false{
			return fmt.Errorf("dstNet %v is not ipv4 subnet" , dstNet)
		}
		v := strings.Split(dstNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToDstNet=net.IPNet{IP: net.ParseIP( strings.Split(dstNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 32)}
		}else{
			return fmt.Errorf("failed to get mask from dstNet %v " )
		}		
	}


	//------------

	rulesBegin, err := netlink.RuleList(unix.AF_INET)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}

	rule := netlink.NewRule()
	rule.Table = tableNum
	rule.Family = unix.AF_INET
	rule.Src = &ToSrcNet
	rule.Dst = &ToDstNet
	rule.Priority = tablePriority
	rule.OifName = ""
	rule.IifName = ""
	rule.Invert = logicalNot

	if err := netlink.RuleDel(rule); err != nil {
		return err
	}

	rules, err := netlink.RuleList(unix.AF_INET)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}	
	if len(rules) != len(rulesBegin)-1 {
		return fmt.Errorf("Rule not removed properly" )
	}

	return nil

}



func DelIPv4RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return   e
		}
		return DelIPv4RouteRule( tableNum , tablePriority   , logicalNot   , srcNet , dstNet       )
	}

}



// we can just use the tableNum and tablePriority to delete the rule
func DelIPv6RouteRule( tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string  ) ( error ) {

	ToSrcNet:=net.IPNet{}
	ToDstNet:=net.IPNet{}

	// check input 
	if tableNum<=CONST_RouteTable_UNSPEC && tableNum>=CONST_RouteTable_MAX {
		return fmt.Errorf("tableNum %v outof range " , tableNum)
	}

	if tablePriority<0 && tablePriority>32767 {
		return fmt.Errorf("tablePriority %v outof range 0-32767 " , tablePriority)
	}

	if len(srcNet)!=0 {
		if CheckIPv6FormatWithMask(srcNet)==false{
			return fmt.Errorf("srcNet %v is not ipv6 subnet" , srcNet)
		}
		v := strings.Split(srcNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToSrcNet=net.IPNet{IP: net.ParseIP( strings.Split(srcNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 128)}
		}else{
			return fmt.Errorf("failed to get mask from srcNet %v " )
		}		
	}

	if len(dstNet)!=0 {
		if CheckIPv6FormatWithMask(dstNet)==false{
			return fmt.Errorf("dstNet %v is not ipv6 subnet" , dstNet)
		}
		v := strings.Split(dstNet , "/" )[1]
		if s, err := strconv.ParseInt(v, 10, 64); err == nil {
		    ToDstNet=net.IPNet{IP: net.ParseIP( strings.Split(dstNet , "/" )[0] ), Mask: net.CIDRMask( int(s) , 128 )}
		}else{
			return fmt.Errorf("failed to get mask from dstNet %v " )
		}		
	}


	//------------

	rulesBegin, err := netlink.RuleList(unix.AF_INET6)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}

	rule := netlink.NewRule()
	rule.Table = tableNum
	rule.Family = unix.AF_INET6
	rule.Src = &ToSrcNet
	rule.Dst = &ToDstNet
	rule.Priority = tablePriority
	rule.OifName = ""
	rule.IifName = ""
	rule.Invert = logicalNot

	if err := netlink.RuleDel(rule); err != nil {
		return err
	}

	rules, err := netlink.RuleList(unix.AF_INET6)
	if err != nil {
		return fmt.Errorf("failed to get current rule " )
	}	
	if len(rules) != len(rulesBegin)-1 {
		return fmt.Errorf("Rule not removed properly" )
	}

	return nil

}


func DelIPv6RouteRuleInNetns( pid int , tableNum , tablePriority int , logicalNot bool , srcNet , dstNet string , InIf , OutIf string ) ( error ) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	origNs, _ := netns.Get()
	defer func() {
		err := netns.Set(origNs)
		if err != nil {
			panic("failed to restore network ns, bailing")
		}
	}()

	if f , e:=netns.GetFromPid(pid); e!=nil{
		return  e

	}else{
		defer f.Close()

		if e:=netns.Set(f) ; e!=nil {
			return   e
		}
		return DelIPv6RouteRule( tableNum , tablePriority   , logicalNot   , srcNet , dstNet      )
	}

}


