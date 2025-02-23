package gerror

import (
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/any"

	spb "google.golang.org/genproto/googleapis/rpc/status"
)

const name = "im"
const TypeUrlStack = "type_url_stack"

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   []byte(stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func WrapRPCError(err error) error {
	if err == nil {
		return nil
	}
	e, _ := status.FromError(err)
	s := &spb.Status{
		Code:    int32(e.Code()),
		Message: e.Message(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   []byte((GetErrorStack(e) + "--grpc-- \n" + stack())),
			},
		},
	}
	return status.FromProto(s).Err()
}

func GetErrorStack(s *status.Status) string {
	pbs := s.Proto()
	for i := range pbs.Details {
		if pbs.Details[i].TypeUrl == TypeUrlStack {
			return string(pbs.Details[i].Value)
		}
	}
	return ""
}

// 获取堆栈信息
func stack() string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, name)
		if n != -1 {
			s := fmt.Sprintf(" %s:%d\n", file[n:], line)
			build.WriteString(s)
		}
	}
	return build.String()
}
