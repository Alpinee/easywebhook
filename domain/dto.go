/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package domain

type LoginDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
