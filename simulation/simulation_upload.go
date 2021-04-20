package simulation

import (
	pb "com.github.patrickz98.omnet/proto"
	"com.github.patrickz98.omnet/simple"
	"com.github.patrickz98.omnet/storage"
)

func Upload(config Config) (ref *pb.StorageRef, err error) {

	logger.Println("zipping", config.Path)

	buf, err := simple.TarGz(config.Path, config.Id)
	if err != nil {
		return
	}

	logger.Println("uploading", config.Id)

	ref, err = storage.Upload(&buf, storage.FileMeta{
		Bucket:   config.Id,
		Filename: "source.tar.gz",
	})
	if err != nil {
		return
	}

	return
}
