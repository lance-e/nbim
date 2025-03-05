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
	resp, err := rpc.GetLogicExtClient().RegisterDevice(GetMetaData(c), &pb.RegisterDeviceReq{
		Type:          3, //TODO:暂时只支持pc端
		Brand:         c.GetHeader("Brand"),
		Model:         c.GetHeader("Model"),
		SystemVersion: c.GetHeader("System-Version"),
		SdkVersion:    c.GetHeader("Sdk-Version"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "successful",
		"data": resp.DeviceId,
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
	}

	resp, err := rpc.GetLogicExtClient().SignIn(GetMetaData(c), &pb.SignInReq{
		Username: username,
		Password: password,
		DeviceId: int64(deviceId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}

	token, err := GenerateToken(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "jwt error",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"is_new":  resp.IsNew,
		"user_id": resp.UserId,
		"token":   token,
	})
}

func GetUser(c *gin.Context) {
	userId := c.PostForm("user_id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
	}
	resp, err := rpc.GetLogicExtClient().GetUser(GetMetaData(c), &pb.GetUserReq{
		UserId: int64(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.User)
}

func UpdateUser(c *gin.Context) {
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
	}
	_, err = rpc.GetLogicExtClient().UpdateUser(GetMetaData(c), &pb.UpdateUserReq{
		Nickname:  nickname,
		Sex:       int32(sexnum),
		AvatarUrl: avatarUrl,
		Extra:     extra,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}
func SearchUser(c *gin.Context) {
	key := c.PostForm("key")

	resp, err := rpc.GetLogicExtClient().SearchUser(GetMetaData(c), &pb.SearchUserReq{
		Key: key,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.Users)
}

func AddFriend(c *gin.Context) {
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
	}
	_, err = rpc.GetLogicExtClient().AddFriend(GetMetaData(c), &pb.AddFriendReq{
		FriendId:    int64(id),
		Remarks:     remarks,
		Description: description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func AgreeFriend(c *gin.Context) {
	friend_id := c.PostForm("friend_id")
	remarks := c.PostForm("remarks")
	id, err := strconv.Atoi(friend_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "friend_id is wrong",
			"data": nil,
		})
	}
	_, err = rpc.GetLogicExtClient().AgreeFriend(GetMetaData(c), &pb.AgreeFriendReq{
		FriendId: int64(id),
		Remarks:  remarks,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func SetFriend(c *gin.Context) {
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
	}
	resp, err := rpc.GetLogicExtClient().SetFriend(GetMetaData(c), &pb.SetFriendReq{
		FriendId: int64(id),
		Remarks:  remarks,
		Extra:    extra,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func GetAllFriends(c *gin.Context) {
	resp, err := rpc.GetLogicExtClient().GetAllFriends(GetMetaData(c), &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.Friends)
}

type tempCreateGroupAllParam struct {
	req      pb.CreateGroupReq
	callerId int64 `json:caller_id`
	deviceId int64 `json:device_id`
}

func CreateGroup(c *gin.Context) {
	req := tempCreateGroupAllParam{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "creategroup request bind json failed",
			"data": nil,
		})
	}
	resp, err := rpc.GetLogicExtClient().CreateGroup(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", strconv.FormatInt(req.callerId, 10),
		"device_id", strconv.FormatInt(req.deviceId, 10),
	)), &req.req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func UpdateGroup(c *gin.Context) {
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
	}
	_, err = rpc.GetLogicExtClient().UpdateGroup(GetMetaData(c), &pb.UpdateGroupReq{
		GroupId:      int64(id),
		Name:         name,
		AvatarUrl:    avatarUrl,
		Extra:        extra,
		Introduction: introduction,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func GetGroup(c *gin.Context) {
	groupId := c.PostForm("group_id")
	id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
	}
	resp, err := rpc.GetLogicExtClient().GetGroup(GetMetaData(c), &pb.GetGroupReq{
		GroupId: int64(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.Group)
}

func GetAllGroup(c *gin.Context) {
	resp, err := rpc.GetLogicExtClient().GetAllGroup(GetMetaData(c), &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.Groups)
}

type tempAddGroupMemberAllParam struct {
	req      pb.AddGroupMemberReq
	callerId int64 `json:caller_id`
	deviceId int64 `json:device_id`
}

func AddGroupMember(c *gin.Context) {
	req := tempAddGroupMemberAllParam{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "add group member request bind json failed",
			"data": nil,
		})
	}
	resp, err := rpc.GetLogicExtClient().AddGroupMember(metadata.NewOutgoingContext(context.TODO(), metadata.Pairs(
		"user_id", strconv.FormatInt(req.callerId, 10),
		"device_id", strconv.FormatInt(req.deviceId, 10),
	)), &req.req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.UserIds)
}

func UpdateGroupMember(c *gin.Context) {
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
	}
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
	}
	member_type, err := strconv.Atoi(memberType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "member_type is wrong",
			"data": nil,
		})
	}
	_, err = rpc.GetLogicExtClient().UpdateGroupMember(GetMetaData(c), &pb.UpdateGroupMemberReq{
		GroupId:    int64(group_id),
		UserId:     int64(user_id),
		MemberType: pb.MemberType(member_type),
		Remarks:    remarks,
		Extra:      extra,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func DeleteGroupMember(c *gin.Context) {
	groupId := c.PostForm("group_id")
	userId := c.PostForm("user_id")
	group_id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
	}
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
	}

	_, err = rpc.GetLogicExtClient().DeleteGroupMember(GetMetaData(c), &pb.DeleteGroupMemberReq{
		GroupId: int64(group_id),
		UserId:  int64(user_id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}

func GetGroupMember(c *gin.Context) {
	groupId := c.PostForm("group_id")
	userId := c.PostForm("user_id")
	group_id, err := strconv.Atoi(groupId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "group_id is wrong",
			"data": nil,
		})
	}
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "user_id is wrong",
			"data": nil,
		})
	}

	resp, err := rpc.GetLogicExtClient().GetGroupMember(GetMetaData(c), &pb.GetGroupMemberReq{
		GroupId: int64(group_id),
		UserId:  int64(user_id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "rpc server wrong ",
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, resp.Members)
}
