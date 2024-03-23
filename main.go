/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package main

import (
	"easywebhook/domain"
	"easywebhook/web"
)

func init() {
	domain.InitDB()
}

func main() {
	web.InitRouter()
}
