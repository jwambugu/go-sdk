// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.4
// source: payment_state.proto

package com_elarian_hera_proto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PaymentState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerNumbers     []*CustomerNumber          `protobuf:"bytes,1,rep,name=customer_numbers,json=customerNumbers,proto3" json:"customer_numbers,omitempty"`
	ChannelNumbers      []*PaymentChannelNumber    `protobuf:"bytes,2,rep,name=channel_numbers,json=channelNumbers,proto3" json:"channel_numbers,omitempty"`
	TransactionLog      []*PaymentTransaction      `protobuf:"bytes,3,rep,name=transaction_log,json=transactionLog,proto3" json:"transaction_log,omitempty"`
	PendingTransactions []*PaymentTransaction      `protobuf:"bytes,4,rep,name=pending_transactions,json=pendingTransactions,proto3" json:"pending_transactions,omitempty"`
	Wallets             map[string]*PaymentBalance `protobuf:"bytes,5,rep,name=wallets,proto3" json:"wallets,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *PaymentState) Reset() {
	*x = PaymentState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payment_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentState) ProtoMessage() {}

func (x *PaymentState) ProtoReflect() protoreflect.Message {
	mi := &file_payment_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentState.ProtoReflect.Descriptor instead.
func (*PaymentState) Descriptor() ([]byte, []int) {
	return file_payment_state_proto_rawDescGZIP(), []int{0}
}

func (x *PaymentState) GetCustomerNumbers() []*CustomerNumber {
	if x != nil {
		return x.CustomerNumbers
	}
	return nil
}

func (x *PaymentState) GetChannelNumbers() []*PaymentChannelNumber {
	if x != nil {
		return x.ChannelNumbers
	}
	return nil
}

func (x *PaymentState) GetTransactionLog() []*PaymentTransaction {
	if x != nil {
		return x.TransactionLog
	}
	return nil
}

func (x *PaymentState) GetPendingTransactions() []*PaymentTransaction {
	if x != nil {
		return x.PendingTransactions
	}
	return nil
}

func (x *PaymentState) GetWallets() map[string]*PaymentBalance {
	if x != nil {
		return x.Wallets
	}
	return nil
}

var File_payment_state_proto protoreflect.FileDescriptor

var file_payment_state_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x61, 0x72, 0x69,
	0x61, 0x6e, 0x2e, 0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x13, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x04, 0x0a, 0x0c, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x51, 0x0a, 0x10, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x2e,
	0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0f, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x55, 0x0a, 0x0f, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x61, 0x72, 0x69, 0x61,
	0x6e, 0x2e, 0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x52, 0x0e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x12, 0x53, 0x0a, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6c, 0x6f, 0x67, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x65, 0x6c, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x2e, 0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x67, 0x12, 0x5d, 0x0a, 0x14, 0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x61, 0x72, 0x69,
	0x61, 0x6e, 0x2e, 0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x13, 0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x4b, 0x0a, 0x07, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x61,
	0x72, 0x69, 0x61, 0x6e, 0x2e, 0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x57, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x73, 0x1a, 0x62, 0x0a, 0x0c, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x3c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6c, 0x61, 0x72, 0x69, 0x61,
	0x6e, 0x2e, 0x68, 0x65, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_payment_state_proto_rawDescOnce sync.Once
	file_payment_state_proto_rawDescData = file_payment_state_proto_rawDesc
)

func file_payment_state_proto_rawDescGZIP() []byte {
	file_payment_state_proto_rawDescOnce.Do(func() {
		file_payment_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_payment_state_proto_rawDescData)
	})
	return file_payment_state_proto_rawDescData
}

var file_payment_state_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_payment_state_proto_goTypes = []interface{}{
	(*PaymentState)(nil),         // 0: com.elarian.hera.proto.PaymentState
	nil,                          // 1: com.elarian.hera.proto.PaymentState.WalletsEntry
	(*CustomerNumber)(nil),       // 2: com.elarian.hera.proto.CustomerNumber
	(*PaymentChannelNumber)(nil), // 3: com.elarian.hera.proto.PaymentChannelNumber
	(*PaymentTransaction)(nil),   // 4: com.elarian.hera.proto.PaymentTransaction
	(*PaymentBalance)(nil),       // 5: com.elarian.hera.proto.PaymentBalance
}
var file_payment_state_proto_depIdxs = []int32{
	2, // 0: com.elarian.hera.proto.PaymentState.customer_numbers:type_name -> com.elarian.hera.proto.CustomerNumber
	3, // 1: com.elarian.hera.proto.PaymentState.channel_numbers:type_name -> com.elarian.hera.proto.PaymentChannelNumber
	4, // 2: com.elarian.hera.proto.PaymentState.transaction_log:type_name -> com.elarian.hera.proto.PaymentTransaction
	4, // 3: com.elarian.hera.proto.PaymentState.pending_transactions:type_name -> com.elarian.hera.proto.PaymentTransaction
	1, // 4: com.elarian.hera.proto.PaymentState.wallets:type_name -> com.elarian.hera.proto.PaymentState.WalletsEntry
	5, // 5: com.elarian.hera.proto.PaymentState.WalletsEntry.value:type_name -> com.elarian.hera.proto.PaymentBalance
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_payment_state_proto_init() }
func file_payment_state_proto_init() {
	if File_payment_state_proto != nil {
		return
	}
	file_common_model_proto_init()
	file_payment_model_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_payment_state_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_payment_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_payment_state_proto_goTypes,
		DependencyIndexes: file_payment_state_proto_depIdxs,
		MessageInfos:      file_payment_state_proto_msgTypes,
	}.Build()
	File_payment_state_proto = out.File
	file_payment_state_proto_rawDesc = nil
	file_payment_state_proto_goTypes = nil
	file_payment_state_proto_depIdxs = nil
}