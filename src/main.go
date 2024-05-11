package main

import (
	"log"
	"syscall"
	"unsafe"
	"os"
	"strings"
	"github.com/lxn/win"
)

var revert = false;

func main() {
	listener = &Listener{}
	go startListenerMessageLoop()
	list, err := GetAllWindows()
	if err != nil {
		log.Fatal(err)
	}
	revert = shouldReverse();
	for _, h := range list {
		hwnd := win.HWND(h)
		applyBasicTheme(hwnd);
	}
	if(revert) {
		os.Exit(0);
	}
	select {}
}

func applyBasicTheme(hwnd win.HWND){
	var policy DWMNCRENDERINGPOLICY;
	var ok = false;
	if(revert){
		policy = DWMNCRP_USEWINDOWSTYLE;
		ok = true;
	} else if (getDWMactive(hwnd)){
		policy = DWMNCRP_DISABLED;
		ok = true;
	}
	if(ok){
		var policyParameter = policy
		DwmSetWindowAttribute(hwnd, DWMWA_NCRENDERING_POLICY, &policyParameter, 8) // unsafe.Sizeof(int(0))
		policyParameter = policy
		DwmSetWindowAttribute(hwnd, DWMWA_FORCE_ICONIC_REPRESENTATION, &policyParameter, 8) // unsafe.Sizeof(int(0))
	}
}

func shouldReverse() bool {
	args := os.Args;
	for _, v := range args {
		if(strings.EqualFold(v, "-reverse")){
			return true;
		}
	}
	return false;
}

func GetAllWindows() ([]syscall.Handle, error) {
	var hwndList []syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, _ uintptr) uintptr {
		hwndList = append(hwndList, h)
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	return hwndList, nil
}

func getDWMactive(hWnd win.HWND) bool {
	var isNCRenderingEnabled = DWMNCRP_USEWINDOWSTYLE
	if DwmGetWindowAttribute(hWnd, DWMWA_NCRENDERING_ENABLED, &isNCRenderingEnabled, uint32(unsafe.Sizeof(isNCRenderingEnabled))) != 0 { // unsafe.Sizeof(int(0))
		log.Println("getDWMactive Err", hWnd)
	}
	return isNCRenderingEnabled == 1
}
