package logichttp

import (
	"context"
	"nbim/pkg/protocol/pb"
	"nbim/pkg/rpc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func RegisterDevice(c *gin.Context) {
	resp, err := rpc.GetLogicExtClient().RegisterDevice(context.TODO(), &pb.RegisterDeviceReq{
		Type:          3, //TODO:暂时只支持pc端
		Brand:         c.GetHeader("Brand"),
		Model:         c.GetHeader("Model"),
		SystemVersion: c.GetHeader("System-Version"),
		SdkVersion:    c.GetHeader("Sdk-Version"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"msg":       "successful",
		"device_id": resp.DeviceId,
	})
}

func SignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	deviceIdstr := c.PostForm("device_id")
	deviceId, err := strconv.Atoi(deviceIdstr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "device_id is wrong",
			"data": nil,
		})
		return
	}

	resp, err := rpc.GetLogicExtClient().SignIn(context.TODO(), &pb.SignInReq{
		Username: username,
		Password: password,
		DeviceId: int64(deviceId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	token, err := GenerateToken(strconv.Itoa(int(resp.GetUserId())), username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "jwt error",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"is_new":  resp.GetIsNew(),
		"user_id": resp.GetUserId(),
		"token":   token,
	})
}

func GetUser(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	userId := c.PostForm("user_id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
		return
	}
	resp, err := rpc.GetLogicExtClient().GetUser(ctx, &pb.GetUserReq{
		UserId: int64(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.User)
}

func UpdateUser(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	nickname := c.PostForm("nickname")
	sex := c.PostForm("sex")
	avatarUrl := c.PostForm("avatar_url")
	extra := c.PostForm("extra")
	sexnum, err := strconv.Atoi(sex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "sex is wrong",
			"data": nil,
		})
		return
	}
	_, err = rpc.GetLogicExtClient().UpdateUser(ctx, &pb.UpdateUserReq{
		Nickname:  nickname,
		Sex:       int32(sexnum),
		AvatarUrl: avatarUrl,
		Extra:     extra,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}
func SearchUser(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	key := c.PostForm("key")

	resp, err := rpc.GetLogicExtClient().SearchUser(ctx, &pb.SearchUserReq{
		Key: key,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.Users)
}

func AddFriend(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	friend_id := c.PostForm("friend_id")
	remarks := c.PostForm("remarks")
	description := c.PostForm("description")
	id, err := strconv.Atoi(friend_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "friend_id is wrong",
			"data": nil,
		})
		return
	}
	_, err = rpc.GetLogicExtClient().AddFriend(ctx, &pb.AddFriendReq{
		FriendId:    int64(id),
		Remarks:     remarks,
		Description: description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func ViewAddFriend(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	resp, err := rpc.GetLogicExtClient().ViewAddFriend(ctx, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.Friends)
}

func AgreeFriend(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	friend_id := c.PostForm("friend_id")
	remarks := c.PostForm("remarks")
	id, err := strconv.Atoi(friend_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "friend_id is wrong",
			"data": nil,
		})
		return
	}
	_, err = rpc.GetLogicExtClient().AgreeFriend(ctx, &pb.AgreeFriendReq{
		FriendId: int64(id),
		Remarks:  remarks,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func SetFriend(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	friend_id := c.PostForm("friend_id")
	remarks := c.PostForm("remarks")
	extra := c.PostForm("extra")
	id, err := strconv.Atoi(friend_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "friend_id is wrong",
			"data": nil,
		})
		return
	}
	resp, err := rpc.GetLogicExtClient().SetFriend(ctx, &pb.SetFriendReq{
		FriendId: int64(id),
		Remarks:  remarks,
		Extra:    extra,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetAllFriends(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	resp, err := rpc.GetLogicExtClient().GetAllFriends(ctx, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.Friends)
}

type tempCreateGroupAllParam struct {
	Name         string  `json:"name,omitempty"`         //名称
	AvatarUrl    string  `json:"avatar_url,omitempty"`   //头像
	Introduction string  `json:"introduction,omitempty"` //介绍
	Extra        string  `json:"extra,omitempty"`        //附加字段
	MemberIds    []int64 `json:"member_ids,omitempty"`   //群组成员id列表
	CallerId     int64   `json:"caller_id"`              //元数据
	DeviceId     int64   `json:"device_id"`              //元数据
}

func CreateGroup(c *gin.Context) {
	req := tempCreateGroupAllParam{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "creategroup request bind json failed",
			"data": nil,
		})
		return
	}
	//token auth
	id, ok := c.Get("caller_id")
	if !ok || (ok && id != strconv.Itoa(int(req.CallerId))) {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}
	resp, err := rpc.GetLogicExtClient().CreateGroup(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", strconv.FormatInt(req.CallerId, 10),
		"device_id", strconv.FormatInt(req.DeviceId, 10),
	)), &pb.CreateGroupReq{
		Name:         req.Name,
		AvatarUrl:    req.AvatarUrl,
		Introduction: req.Introduction,
		Extra:        req.Extra,
		MemberIds:    req.MemberIds,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func UpdateGroup(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	groupId := c.PostForm("group_id")
	avatarUrl := c.PostForm("avatar_url")
	name := c.PostForm("name")
	introduction := c.PostForm("introduction")
	extra := c.PostForm("extra")
	id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
		return
	}
	_, err = rpc.GetLogicExtClient().UpdateGroup(ctx, &pb.UpdateGroupReq{
		GroupId:      int64(id),
		Name:         name,
		AvatarUrl:    avatarUrl,
		Extra:        extra,
		Introduction: introduction,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func GetGroup(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	groupId := c.PostForm("group_id")
	id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
		return
	}
	resp, err := rpc.GetLogicExtClient().GetGroup(ctx, &pb.GetGroupReq{
		GroupId: int64(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.Group)
}

func GetAllGroup(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	resp, err := rpc.GetLogicExtClient().GetAllGroup(ctx, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.Groups)
}

type tempAddGroupMemberAllParam struct {
	GroupId  int64   `json:"group_id"` //群组id
	UserIds  []int64 `json:"user_ids"` //用户id列表
	CallerId int64   `json:"caller_id"`
	DeviceId int64   `json:"device_id"`
}

func AddGroupMember(c *gin.Context) {
	req := tempAddGroupMemberAllParam{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "add group member request bind json failed",
			"data": nil,
		})
		return
	}
	//token auth
	id, ok := c.Get("caller_id")
	if !ok || (ok && id != strconv.Itoa(int(req.CallerId))) {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}
	resp, err := rpc.GetLogicExtClient().AddGroupMember(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", strconv.FormatInt(req.CallerId, 10),
		"device_id", strconv.FormatInt(req.DeviceId, 10),
	)), &pb.AddGroupMemberReq{
		GroupId: req.GroupId,
		UserIds: req.UserIds,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func UpdateGroupMember(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	groupId := c.PostForm("group_id")
	userId := c.PostForm("user_id")
	memberType := c.PostForm("member_type")
	remarks := c.PostForm("remarks")
	extra := c.PostForm("extra")
	group_id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
		return
	}
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
		return
	}
	member_type, err := strconv.Atoi(memberType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "member_type is wrong",
			"data": nil,
		})
		return
	}
	_, err = rpc.GetLogicExtClient().UpdateGroupMember(ctx, &pb.UpdateGroupMemberReq{
		GroupId:    int64(group_id),
		UserId:     int64(user_id),
		MemberType: pb.MemberType(member_type),
		Remarks:    remarks,
		Extra:      extra,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func DeleteGroupMember(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	groupId := c.PostForm("group_id")
	userId := c.PostForm("user_id")
	group_id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
		return
	}
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
		return
	}

	_, err = rpc.GetLogicExtClient().DeleteGroupMember(ctx, &pb.DeleteGroupMemberReq{
		GroupId: int64(group_id),
		UserId:  int64(user_id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func GetGroupMember(c *gin.Context) {
	ctx := GetMetaData(c)
	if ctx == nil {
		c.JSON(401, gin.H{
			"code": 401,
			"msg":  "token metadata is wrong ",
		})
		return
	}

	groupId := c.PostForm("group_id")
	userId := c.PostForm("user_id")
	group_id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
		return
	}
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
		return
	}

	resp, err := rpc.GetLogicExtClient().GetGroupMember(ctx, &pb.GetGroupMemberReq{
		GroupId: int64(group_id),
		UserId:  int64(user_id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, resp.Members)
}
