package infrastructure

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	boshplatform "github.com/cloudfoundry/bosh-agent/platform"
	boshsettings "github.com/cloudfoundry/bosh-agent/settings"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type InstanceMetadataSettingsSource struct {
	metadataHost    string
	metadataHeaders map[string]string
	settingsPath    string

	platform boshplatform.Platform
	logger   boshlog.Logger

	logTag          string
	metadataService DynamicMetadataService
}

func NewInstanceMetadataSettingsSource(
	metadataHost string,
	metadataHeaders map[string]string,
	settingsPath string,
	platform boshplatform.Platform,
	logger boshlog.Logger,
) *InstanceMetadataSettingsSource {
	logTag := "InstanceMetadataSettingsSource"
	return &InstanceMetadataSettingsSource{
		metadataHost:    metadataHost,
		metadataHeaders: metadataHeaders,
		settingsPath:    settingsPath,

		platform: platform,
		logger:   logger,

		logTag: logTag,
		// The HTTPMetadataService provides more functionality than we need (like custom DNS), so we
		// pass zero values to the New function and only use its GetValueAtPath method.
		metadataService: NewHTTPMetadataService(metadataHost, metadataHeaders, "", "", "", nil, platform, logger),
	}
}

func NewInstanceMetadataSettingsSourceWithoutRetryDelay(
	metadataHost string,
	metadataHeaders map[string]string,
	settingsPath string,
	platform boshplatform.Platform,
	logger boshlog.Logger,
) *InstanceMetadataSettingsSource {
	logTag := "InstanceMetadataSettingsSource"
	return &InstanceMetadataSettingsSource{
		metadataHost:    metadataHost,
		metadataHeaders: metadataHeaders,
		settingsPath:    settingsPath,

		platform: platform,
		logger:   logger,

		logTag: logTag,
		// The HTTPMetadataService provides more functionality than we need (like custom DNS), so we
		// pass zero values to the New function and only use its GetValueAtPath method.
		metadataService: NewHTTPMetadataServiceWithCustomRetryDelay(metadataHost, metadataHeaders, "", "", "", nil, platform, logger, 0*time.Second),
	}
}

func (s InstanceMetadataSettingsSource) PublicSSHKeyForUsername(string) (string, error) {
	return "", nil
}

func (s *InstanceMetadataSettingsSource) Settings() (boshsettings.Settings, error) {
	var settings boshsettings.Settings
	contents, err := s.metadataService.GetValueAtPath(s.settingsPath)
	if err != nil {
		return settings, bosherr.WrapError(err, fmt.Sprintf("Reading settings from instance metadata at path %q", s.settingsPath))
	}

	if strings.Contains(s.settingsPath, "SoftLayer_Resource_Metadata") {
		settingsBytesWithoutQuotes := strings.Replace(string(contents), `"`, ``, -1)
		decodedSettings, err := base64.RawURLEncoding.DecodeString(settingsBytesWithoutQuotes)
		if err != nil {
			return settings, bosherr.WrapError(err, "Decoding url encoded instance metadata settings")
		}

		err = json.Unmarshal([]byte(decodedSettings), &settings)
		if err != nil {
			return settings, bosherr.WrapErrorf(err, "Parsing instance metadata settings after decoding from %q", decodedSettings)
		}
	} else {
		err = json.Unmarshal([]byte(contents), &settings)
		if err != nil {
			return settings, bosherr.WrapErrorf(
				err, "Parsing instance metadata settings from %q", contents)
		}
	}

	return settings, nil
}
