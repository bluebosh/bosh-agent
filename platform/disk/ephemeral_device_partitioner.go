package disk

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"code.cloudfoundry.org/clock"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshretry "github.com/cloudfoundry/bosh-utils/retrystrategy"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type EphemeralDevicePartitioner struct {
	partedPartitioner PartedPartitioner
	deviceUtil        Util
	logger            boshlog.Logger

	logTag       string
	cmdRunner    boshsys.CmdRunner
	fs           boshsys.FileSystem
	timeService  clock.Clock
	settingsPath string
}

type Settings struct {
	AgentID string `json:"agent_id"`
}

func NewEphemeralDevicePartitioner(
	partedPartitioner PartedPartitioner,
	deviceUtil Util,
	logger boshlog.Logger,
	cmdRunner boshsys.CmdRunner,
	fs boshsys.FileSystem,
	timeService clock.Clock,
	settingsPath string,
) *EphemeralDevicePartitioner {
	return &EphemeralDevicePartitioner{
		partedPartitioner: partedPartitioner,
		deviceUtil:        deviceUtil,
		logger:            logger,
		logTag:            "EphemeralDevicePartitioner",
		cmdRunner:         cmdRunner,
		fs:                fs,
		timeService:       timeService,
		settingsPath:      settingsPath,
	}
}

func (p *EphemeralDevicePartitioner) Partition(devicePath string, partitions []Partition) error {
	agentID, err := p.getAgentID()
	if err != nil {
		return bosherr.WrapError(err, "Getting agent ID")
	}

	p.partedPartitioner.setPartitionNamePrefix(agentID)

	existingPartitions, deviceFullSizeInBytes, err := p.partedPartitioner.getPartitions(devicePath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Getting existing partitions of `%s'", devicePath)
	}

	if p.partedPartitioner.partitionsMatch(existingPartitions, partitions, deviceFullSizeInBytes) && p.namesMatch(existingPartitions, agentID) {
		p.logger.Debug(p.logTag, "Existing partitions match desired partitions")
		return nil
	}

	err = p.removePartitions(existingPartitions, devicePath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Removing existing partitions of `%s'", devicePath)
	}

	return p.partedPartitioner.Partition(devicePath, partitions)
}

func (p *EphemeralDevicePartitioner) GetDeviceSizeInBytes(devicePath string) (uint64, error) {
	return p.partedPartitioner.GetDeviceSizeInBytes(devicePath)
}

func (p EphemeralDevicePartitioner) getAgentID() (string, error) {
	opts := boshsys.ReadOpts{Quiet: true}
	existingSettingsJSON, readError := p.fs.ReadFileWithOpts(p.settingsPath, opts)
	if readError != nil {
		return "", bosherr.WrapErrorf(readError, "Failed reading settings from file %s", p.settingsPath)
	}

	settings := Settings{}
	err := json.Unmarshal(existingSettingsJSON, &settings)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed unmarshalling settings from file %s", p.settingsPath)
	}

	return settings.AgentID, nil
}

func (p EphemeralDevicePartitioner) namesMatch(existingPartitions []existingPartition, agentID string) bool {
	for _, existingPartition := range existingPartitions {
		if strings.HasPrefix(existingPartition.Name, agentID) {
			return true
		}
	}

	return false
}

func (p EphemeralDevicePartitioner) removePartitions(partitions []existingPartition, devicePath string) error {
	partitionPaths, err := p.getPartitionPaths(devicePath)
	if err != nil {
		return bosherr.WrapErrorf(err, "Getting partition paths of disk `%s'", devicePath)
	}

	p.logger.Debug(p.logTag, "Erasing old partition paths")
	for _, partitionPath := range partitionPaths {
		partitionRetryable := boshretry.NewRetryable(func() (bool, error) {
			_, _, _, err := p.cmdRunner.RunCommand(
				"wipefs",
				"-a",
				partitionPath,
			)
			if err != nil {
				return true, bosherr.WrapError(err, fmt.Sprintf("Erasing partition path `%s' ", partitionPath))
			}

			p.logger.Info(p.logTag, "Successfully erased partition path `%s'", partitionPath)
			return false, nil
		})

		partitionRetryStrategy := NewPartitionStrategy(partitionRetryable, p.timeService, p.logger)
		err := partitionRetryStrategy.Try()

		if err != nil {
			return bosherr.WrapErrorf(err, "Erasing partition `%s' paths", devicePath)
		}
	}

	p.logger.Debug(p.logTag, "Removing old partitions")
	for _, partition := range partitions {
		partitionRetryable := boshretry.NewRetryable(func() (bool, error) {
			_, _, _, err := p.cmdRunner.RunCommand(
				"parted",
				devicePath,
				"rm",
				strconv.Itoa(partition.Index),
			)
			if err != nil {
				return true, bosherr.WrapError(err, "Removing partition using parted")
			}

			p.logger.Info(p.logTag, "Successfully removed partition %s from %s", partition.Name, devicePath)
			return false, nil
		})

		partitionRetryStrategy := NewPartitionStrategy(partitionRetryable, p.timeService, p.logger)
		err := partitionRetryStrategy.Try()

		if err != nil {
			return bosherr.WrapErrorf(err, "Removing partitions of disk `%s'", devicePath)
		}
	}
	return nil
}

func (p EphemeralDevicePartitioner) getPartitionPaths(devicePath string) ([]string, error) {
	stdout, _, _, err := p.cmdRunner.RunCommand("blkid")
	if err != nil {
		return []string{}, err
	}

	pathRegExp := devicePath + "[0-9]+"
	re := regexp.MustCompile(pathRegExp)
	match := re.FindAllString(stdout, -1)

	if nil == match {
		return []string{}, nil
	}

	return match, nil
}
