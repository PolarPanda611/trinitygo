/**
 * @ Author: Daniel Tan
 * @ Date: 2020-08-04 14:57:51
 * @ LastEditTime: 2020-08-04 16:18:57
 * @ LastEditors: Daniel Tan
 * @ Description:
 * @ FilePath: /trinitygo/startup/startup.go
 * @
 */

package startup

import "fmt"

var (
	_startupDebugger     bool = false
	_startupDebuggerInfo []string
	_requestMapping      []string
)

func AppendRequestMapping(method, url string, handler ...string) {
	if len(handler) == 0 {
		_requestMapping = append(_requestMapping, fmt.Sprintf("request mapping : %-6s  %-30s ", method, url))
	} else {
		_requestMapping = append(_requestMapping, fmt.Sprintf("request mapping : %-6s  %-30s => %v", method, url, handler[0]))
	}

}
func AppendStartupDebuggerInfo(msg string) {
	_startupDebuggerInfo = append(_startupDebuggerInfo, msg)
}
func SetStartupDebugger(isLog bool) {
	_startupDebugger = isLog
}

func GetStartupDebugger() bool {
	return _startupDebugger
}

func GetStartupDebuggerInfo() []string {
	return _startupDebuggerInfo
}

func GetRequestMapping() []string {
	return _requestMapping
}
