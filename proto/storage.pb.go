// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: storage.proto

package proto

import (
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

type StorageList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files []string `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *StorageList) Reset() {
	*x = StorageList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorageList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageList) ProtoMessage() {}

func (x *StorageList) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageList.ProtoReflect.Descriptor instead.
func (*StorageList) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{0}
}

func (x *StorageList) GetFiles() []string {
	if x != nil {
		return x.Files
	}
	return nil
}

type StorageStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deleted bool `protobuf:"varint,1,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *StorageStatus) Reset() {
	*x = StorageStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorageStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageStatus) ProtoMessage() {}

func (x *StorageStatus) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageStatus.ProtoReflect.Descriptor instead.
func (*StorageStatus) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{1}
}

func (x *StorageStatus) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type StorageRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bucket   string `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *StorageRef) Reset() {
	*x = StorageRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorageRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageRef) ProtoMessage() {}

func (x *StorageRef) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageRef.ProtoReflect.Descriptor instead.
func (*StorageRef) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{2}
}

func (x *StorageRef) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *StorageRef) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

type StorageParcel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size    int32  `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	Offset  int64  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *StorageParcel) Reset() {
	*x = StorageParcel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorageParcel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageParcel) ProtoMessage() {}

func (x *StorageParcel) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageParcel.ProtoReflect.Descriptor instead.
func (*StorageParcel) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{3}
}

func (x *StorageParcel) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *StorageParcel) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *StorageParcel) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_storage_proto protoreflect.FileDescriptor

var file_storage_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x29, 0x0a,
	0x0d, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x40, 0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x55, 0x0a, 0x0d, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x63, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x32, 0xe9, 0x01, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a,
	0x04, 0x50, 0x75, 0x6c, 0x6c, 0x12, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x66, 0x1a, 0x16, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x50, 0x61, 0x72, 0x63,
	0x65, 0x6c, 0x22, 0x00, 0x30, 0x01, 0x12, 0x37, 0x0a, 0x04, 0x50, 0x75, 0x73, 0x68, 0x12, 0x16,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x50, 0x61, 0x72, 0x63, 0x65, 0x6c, 0x1a, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x66, 0x22, 0x00, 0x28, 0x01, 0x12,
	0x37, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x66, 0x1a, 0x16,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x66, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x42, 0x09, 0x5a,
	0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storage_proto_rawDescOnce sync.Once
	file_storage_proto_rawDescData = file_storage_proto_rawDesc
)

func file_storage_proto_rawDescGZIP() []byte {
	file_storage_proto_rawDescOnce.Do(func() {
		file_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_storage_proto_rawDescData)
	})
	return file_storage_proto_rawDescData
}

var file_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_storage_proto_goTypes = []interface{}{
	(*StorageList)(nil),   // 0: service.StorageList
	(*StorageStatus)(nil), // 1: service.StorageStatus
	(*StorageRef)(nil),    // 2: service.StorageRef
	(*StorageParcel)(nil), // 3: service.StorageParcel
}
var file_storage_proto_depIdxs = []int32{
	2, // 0: service.Storage.Pull:input_type -> service.StorageRef
	3, // 1: service.Storage.Push:input_type -> service.StorageParcel
	2, // 2: service.Storage.Delete:input_type -> service.StorageRef
	2, // 3: service.Storage.List:input_type -> service.StorageRef
	3, // 4: service.Storage.Pull:output_type -> service.StorageParcel
	2, // 5: service.Storage.Push:output_type -> service.StorageRef
	1, // 6: service.Storage.Delete:output_type -> service.StorageStatus
	0, // 7: service.Storage.List:output_type -> service.StorageList
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_storage_proto_init() }
func file_storage_proto_init() {
	if File_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorageList); i {
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
		file_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorageStatus); i {
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
		file_storage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorageRef); i {
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
		file_storage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorageParcel); i {
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
			RawDescriptor: file_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storage_proto_goTypes,
		DependencyIndexes: file_storage_proto_depIdxs,
		MessageInfos:      file_storage_proto_msgTypes,
	}.Build()
	File_storage_proto = out.File
	file_storage_proto_rawDesc = nil
	file_storage_proto_goTypes = nil
	file_storage_proto_depIdxs = nil
}
