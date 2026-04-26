package main

import (
	"net/http"
	"protobuf/pb"

	// Replace "your_module_path/pd" with the actual module path where your generated pd package is located

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func getUserProto(c *gin.Context) {
	// 这里我们假设已经定义了一个 UserProto 的 protobuf 消息类型，并且已经生成了对应的 Go 代码
	user := &pb.User{
		Id:       1,
		Username: "test",
		Email:    "<EMAIL>",
		Age:      18,
		Active:   true,
		Tags:     []string{"golang", "protobuf"},
		Metadata: map[string]string{
			"role": "admin",
		},
	}
	c.Header("content-type", "application/x-protobuf")
	data, err := proto.Marshal(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/x-protobuf", data)
}
func createUserProto(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var req pb.CreateUserRequest
	err = proto.Unmarshal(data, &req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid protobuf data: " + err.Error()})
		return
	}
	user := &pb.User{
		Id:       1,
		Username: req.Username,
		Email:    req.Email,
		Age:      req.Age,
		Active:   true,
	}
	resp := &pb.CreateUserResponse{
		User:    user,
		Success: true,
		Message: "User create successfully",
	}
	data, _ = proto.Marshal(resp)
	c.Data(http.StatusOK, "application/x-protobuf", data)
}
func main() {
	r := gin.Default()
	r.GET("/getuser", getUserProto)
	r.GET("/createuser", createUserProto)
	r.Run()
}
