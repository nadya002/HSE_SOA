package protoPack

import (
	"HW1/cmd/server/pb"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func Serial_protobuf(m pb.ProtoMessage) ([]byte, error) {
	//fmt.Println(m)
	p, err := proto.Marshal(&m)
	if err != nil {
		fmt.Println("marshaling error: ", err)
		return nil, nil
	}
	//fmt.Println(p)
	return p, nil
}

func Deserial_protobuf(b []byte) (pb.ProtoMessage, error) {
	msg := pb.ProtoMessage{}
	proto.Unmarshal(b, &msg)
	return msg, nil
}
