/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package web

import (
	"easywebhook/domain"
	"easywebhook/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func login(c *gin.Context) {
	loginDTO := domain.LoginDTO{}
	err := c.ShouldBindBodyWith(&loginDTO, binding.JSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	if loginDTO.Account == "" || loginDTO.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Account and Password cannot empty"})
		return
	}

	// 查看是否是超级管理员
	if loginDTO.Account == os.Getenv("SUPER_ADMIN_ACCOUNT") && loginDTO.Password == os.Getenv("SUPER_ADMIN_PASSWORD") {
		token := LoginToken(-1)
		c.SetCookie("session", token, 3600, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var user domain.User
	if err := domain.Db.Where("account = ?", loginDTO.Account).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Account or password error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDTO.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Account or password error"})
		return
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	if err := domain.Db.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
	}

	token := LoginToken(user.ID)
	c.SetCookie("session", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

/**
 * 添加脚本
 * @Author gongquanlin
 * @Description
 * @Date 11:11 2024/3/23
 * @Param
 * @return
 **/
func addScript(c *gin.Context) {
	var tokenScript domain.TokenScript
	if err := c.BindJSON(&tokenScript); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	if tokenScript.Token == "" {
		tokenScript.Token = utils.GenerateRandomString(10)
	}

	tokenScript.UserId = c.GetInt("userId")

	// 插入数据到数据库
	if err := domain.Db.Create(&tokenScript).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add script"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script added successfully"})
}

/**
 * 获取脚本
 * @Author gongquanlin
 * @Description
 * @Date 11:11 2024/3/23
 * @Param
 * @return
 **/
func getScript(c *gin.Context) {
	var tokenScript domain.TokenScript
	id := c.Param("id")

	// 查询数据库，根据ID获取对应的脚本信息
	if err := domain.Db.First(&tokenScript, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Script not found"})
		return
	}

	c.JSON(http.StatusOK, tokenScript)
}

/**
 * 更新脚本
 * @Author gongquanlin
 * @Description
 * @Date 11:11 2024/3/23
 * @Param
 * @return
 **/
func updateScript(c *gin.Context) {
	var tokenScript domain.TokenScript
	id := c.Param("id")

	if err := domain.Db.First(&tokenScript, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Script not found"})
		return
	}

	if err := c.BindJSON(&tokenScript); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
		return
	}

	// 更新数据库中对应ID的脚本信息
	if err := domain.Db.Save(&tokenScript).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update script"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script updated successfully"})
}

/**
 * 删除脚本
 * @Author gongquanlin
 * @Description
 * @Date 11:11 2024/3/23
 * @Param
 * @return
 **/
func deleteScript(c *gin.Context) {
	var tokenScript domain.TokenScript
	id := c.Param("id")

	// 查询数据库，根据ID获取对应的脚本信息
	if err := domain.Db.First(&tokenScript, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Script not found"})
		return
	}

	// 删除数据库中对应ID的脚本信息
	if err := domain.Db.Delete(&tokenScript).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete script"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Script deleted successfully"})
}

/**
 * 处理webhook
 * @Author gongquanlin
 * @Description
 * @Date 11:09 2024/3/23
 * @Param
 * @return
 **/
func handleWebhook(c *gin.Context) {
	// 获取请求的Header中的指定信息，比如token
	token, ok := c.GetQuery("Token")
	if !ok || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	// 查询数据库，获取对应Token的脚本
	var tokenScript domain.TokenScript
	if err := domain.Db.Where("token = ?", token).First(&tokenScript).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 执行获取到的脚本
	cmd := exec.Command("sh", "-c", tokenScript.Script)
	err := cmd.Run()

	// 检查是否执行成功，如果执行失败则返回错误
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute script"})
		return
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "Script executed successfully"})
}
