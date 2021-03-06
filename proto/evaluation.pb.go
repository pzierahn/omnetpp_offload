// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: evaluation.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Scenario struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scenario   string `protobuf:"bytes,1,opt,name=scenario,proto3" json:"scenario,omitempty"`
	Trail      string `protobuf:"bytes,2,opt,name=trail,proto3" json:"trail,omitempty"`
	Simulation string `protobuf:"bytes,3,opt,name=simulation,proto3" json:"simulation,omitempty"`
}

func (x *Scenario) Reset() {
	*x = Scenario{}
	if protoimpl.UnsafeEnabled {
		mi := &file_evaluation_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Scenario) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Scenario) ProtoMessage() {}

func (x *Scenario) ProtoReflect() protoreflect.Message {
	mi := &file_evaluation_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Scenario.ProtoReflect.Descriptor instead.
func (*Scenario) Descriptor() ([]byte, []int) {
	return file_evaluation_proto_rawDescGZIP(), []int{0}
}

func (x *Scenario) GetScenario() string {
	if x != nil {
		return x.Scenario
	}
	return ""
}

func (x *Scenario) GetTrail() string {
	if x != nil {
		return x.Trail
	}
	return ""
}

func (x *Scenario) GetSimulation() string {
	if x != nil {
		return x.Simulation
	}
	return ""
}

type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId     string `protobuf:"bytes,1,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
	Hostname     string `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Timesent     string `protobuf:"bytes,3,opt,name=timesent,proto3" json:"timesent,omitempty"`
	Timereceived string `protobuf:"bytes,4,opt,name=timereceived,proto3" json:"timereceived,omitempty"`
	Os           string `protobuf:"bytes,5,opt,name=os,proto3" json:"os,omitempty"`
	Arch         string `protobuf:"bytes,6,opt,name=arch,proto3" json:"arch,omitempty"`
	NumCPUs      uint32 `protobuf:"varint,7,opt,name=numCPUs,proto3" json:"numCPUs,omitempty"`
	NumJobs      uint32 `protobuf:"varint,8,opt,name=numJobs,proto3" json:"numJobs,omitempty"`
	Connect      string `protobuf:"bytes,9,opt,name=connect,proto3" json:"connect,omitempty"`
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_evaluation_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_evaluation_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_evaluation_proto_rawDescGZIP(), []int{1}
}

func (x *Device) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *Device) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Device) GetTimesent() string {
	if x != nil {
		return x.Timesent
	}
	return ""
}

func (x *Device) GetTimereceived() string {
	if x != nil {
		return x.Timereceived
	}
	return ""
}

func (x *Device) GetOs() string {
	if x != nil {
		return x.Os
	}
	return ""
}

func (x *Device) GetArch() string {
	if x != nil {
		return x.Arch
	}
	return ""
}

func (x *Device) GetNumCPUs() uint32 {
	if x != nil {
		return x.NumCPUs
	}
	return 0
}

func (x *Device) GetNumJobs() uint32 {
	if x != nil {
		return x.NumJobs
	}
	return 0
}

func (x *Device) GetConnect() string {
	if x != nil {
		return x.Connect
	}
	return ""
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId   string `protobuf:"bytes,1,opt,name=eventId,proto3" json:"eventId,omitempty"`
	DeviceId  string `protobuf:"bytes,2,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
	Timestamp string `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Activity  string `protobuf:"bytes,4,opt,name=activity,proto3" json:"activity,omitempty"`
	State     uint32 `protobuf:"varint,5,opt,name=state,proto3" json:"state,omitempty"`
	Config    string `protobuf:"bytes,6,opt,name=config,proto3" json:"config,omitempty"`
	RunNum    string `protobuf:"bytes,7,opt,name=runNum,proto3" json:"runNum,omitempty"`
	Error     string `protobuf:"bytes,8,opt,name=error,proto3" json:"error,omitempty"`
	ByteSize  uint64 `protobuf:"varint,9,opt,name=byteSize,proto3" json:"byteSize,omitempty"` // download & upload size
	Filename  string `protobuf:"bytes,10,opt,name=filename,proto3" json:"filename,omitempty"` // filename
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_evaluation_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_evaluation_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_evaluation_proto_rawDescGZIP(), []int{2}
}

func (x *Event) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *Event) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *Event) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Event) GetActivity() string {
	if x != nil {
		return x.Activity
	}
	return ""
}

func (x *Event) GetState() uint32 {
	if x != nil {
		return x.State
	}
	return 0
}

func (x *Event) GetConfig() string {
	if x != nil {
		return x.Config
	}
	return ""
}

func (x *Event) GetRunNum() string {
	if x != nil {
		return x.RunNum
	}
	return ""
}

func (x *Event) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *Event) GetByteSize() uint64 {
	if x != nil {
		return x.ByteSize
	}
	return 0
}

func (x *Event) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

var File_evaluation_proto protoreflect.FileDescriptor

var file_evaluation_proto_rawDesc = []byte{
	0x0a, 0x10, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x08, 0x53, 0x63, 0x65, 0x6e,
	0x61, 0x72, 0x69, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x63, 0x65, 0x6e, 0x61, 0x72, 0x69, 0x6f,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x69, 0x6d, 0x75,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xf2, 0x01, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x69, 0x6d,
	0x65, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x63,
	0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x63, 0x68, 0x12, 0x18, 0x0a,
	0x07, 0x6e, 0x75, 0x6d, 0x43, 0x50, 0x55, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x6e, 0x75, 0x6d, 0x43, 0x50, 0x55, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x75, 0x6d, 0x4a, 0x6f,
	0x62, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6e, 0x75, 0x6d, 0x4a, 0x6f, 0x62,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x22, 0x8b, 0x02, 0x0a, 0x05,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x75, 0x6e, 0x4e, 0x75, 0x6d, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x75, 0x6e, 0x4e, 0x75, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x79, 0x74, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x79, 0x74, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0xda, 0x01, 0x0a, 0x0a, 0x45, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74,
	0x12, 0x0f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2d, 0x0a, 0x03, 0x4c, 0x6f, 0x67,
	0x12, 0x0e, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x12, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x63, 0x65, 0x6e,
	0x61, 0x72, 0x69, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x06,
	0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_evaluation_proto_rawDescOnce sync.Once
	file_evaluation_proto_rawDescData = file_evaluation_proto_rawDesc
)

func file_evaluation_proto_rawDescGZIP() []byte {
	file_evaluation_proto_rawDescOnce.Do(func() {
		file_evaluation_proto_rawDescData = protoimpl.X.CompressGZIP(file_evaluation_proto_rawDescData)
	})
	return file_evaluation_proto_rawDescData
}

var file_evaluation_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_evaluation_proto_goTypes = []interface{}{
	(*Scenario)(nil),      // 0: service.Scenario
	(*Device)(nil),        // 1: service.Device
	(*Event)(nil),         // 2: service.Event
	(*emptypb.Empty)(nil), // 3: google.protobuf.Empty
}
var file_evaluation_proto_depIdxs = []int32{
	1, // 0: service.Evaluation.Init:input_type -> service.Device
	2, // 1: service.Evaluation.Log:input_type -> service.Event
	0, // 2: service.Evaluation.Start:input_type -> service.Scenario
	3, // 3: service.Evaluation.Finish:input_type -> google.protobuf.Empty
	3, // 4: service.Evaluation.Init:output_type -> google.protobuf.Empty
	3, // 5: service.Evaluation.Log:output_type -> google.protobuf.Empty
	3, // 6: service.Evaluation.Start:output_type -> google.protobuf.Empty
	3, // 7: service.Evaluation.Finish:output_type -> google.protobuf.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_evaluation_proto_init() }
func file_evaluation_proto_init() {
	if File_evaluation_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_evaluation_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Scenario); i {
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
		file_evaluation_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
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
		file_evaluation_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
			RawDescriptor: file_evaluation_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_evaluation_proto_goTypes,
		DependencyIndexes: file_evaluation_proto_depIdxs,
		MessageInfos:      file_evaluation_proto_msgTypes,
	}.Build()
	File_evaluation_proto = out.File
	file_evaluation_proto_rawDesc = nil
	file_evaluation_proto_goTypes = nil
	file_evaluation_proto_depIdxs = nil
}
