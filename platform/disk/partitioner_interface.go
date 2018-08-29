package disk

import "fmt"

type PartitionType string

const (
	PartitionTypeSwap    PartitionType = "swap"
	PartitionTypeLinux   PartitionType = "linux"
	PartitionTypeEmpty   PartitionType = "empty"
	PartitionTypeUnknown PartitionType = "unknown"
	PartitionTypeGPT     PartitionType = "gpt"
)

type Partition struct {
	SizeInBytes uint64
	Type        PartitionType
}

type Partitioner interface {
	Partition(devicePath string, partitions []Partition) (err error)
	GetDeviceSizeInBytes(devicePath string) (size uint64, err error)
}

type PartedPartitioner interface {
	Partitioner
	partitionsMatch(existingPartitions []existingPartition, desiredPartitions []Partition, deviceSizeInBytes uint64) bool
	getPartitions(devicePath string) (partitions []existingPartition, deviceFullSizeInBytes uint64, err error)
	setPartitionNamePrefix(partitionNamePrefix string)
}

func (p Partition) String() string {
	return fmt.Sprintf("[Type: %s, SizeInBytes: %d]", p.Type, p.SizeInBytes)
}
