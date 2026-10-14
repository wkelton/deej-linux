package deej

import (
	"fmt"
	"net"

	"github.com/jfreymuth/pulse/proto"
	"go.uber.org/zap"
)

type paSessionFinder struct {
	logger        *zap.SugaredLogger
	sessionLogger *zap.SugaredLogger

	client *proto.Client
	conn   net.Conn

	config *CanonicalConfig
}

func newSessionFinder(logger *zap.SugaredLogger, config *CanonicalConfig) (SessionFinder, error) {
	client, conn, err := proto.Connect("")
	if err != nil {
		logger.Warnw("Failed to establish PulseAudio connection", "error", err)
		return nil, fmt.Errorf("establish PulseAudio connection: %w", err)
	}

	request := proto.SetClientName{
		Props: proto.PropList{
			"application.name": proto.PropListString("deej"),
		},
	}
	reply := proto.SetClientNameReply{}

	if err := client.Request(&request, &reply); err != nil {
		return nil, err
	}

	sf := &paSessionFinder{
		logger:        logger.Named("session_finder"),
		sessionLogger: logger.Named("sessions"),
		client:        client,
		conn:          conn,
		config:        config,
	}

	sf.logger.Debug("Created PA session finder instance")

	return sf, nil
}

func (sf *paSessionFinder) GetAllSessions() ([]Session, error) {
	sessions := []Session{}

	// get the master sink session
	masterSink, err := sf.getMasterSinkSession()
	if err == nil {
		sessions = append(sessions, masterSink)
	} else {
		sf.logger.Warnw("Failed to get master audio sink session", "error", err)
	}

	// get the master source session
	masterSource, err := sf.getMasterSourceSession()
	if err == nil {
		sessions = append(sessions, masterSource)
	} else {
		sf.logger.Warnw("Failed to get master audio source session", "error", err)
	}

	// enumerate sink inputs and add sessions along the way
	if err := sf.enumerateAndAddSessions(&sessions); err != nil {
		sf.logger.Warnw("Failed to enumerate audio sessions", "error", err)
		return nil, fmt.Errorf("enumerate audio sessions: %w", err)
	}

	return sessions, nil
}

func (sf *paSessionFinder) Release() error {
	if err := sf.conn.Close(); err != nil {
		sf.logger.Warnw("Failed to close PulseAudio connection", "error", err)
		return fmt.Errorf("close PulseAudio connection: %w", err)
	}

	sf.logger.Debug("Released PA session finder instance")

	return nil
}

func (sf *paSessionFinder) getMasterSinkSession() (Session, error) {
	request := proto.GetSinkInfo{
		SinkIndex: proto.Undefined,
	}
	reply := proto.GetSinkInfoReply{}

	if err := sf.client.Request(&request, &reply); err != nil {
		sf.logger.Warnw("Failed to get master sink info", "error", err)
		return nil, fmt.Errorf("get master sink info: %w", err)
	}

	// create the master sink session
	sink := newMasterSession(sf.sessionLogger, sf.client, reply.SinkIndex, reply.Channels, true)

	return sink, nil
}

func (sf *paSessionFinder) getMasterSourceSession() (Session, error) {
	request := proto.GetSourceInfo{
		SourceIndex: proto.Undefined,
	}
	reply := proto.GetSourceInfoReply{}

	if err := sf.client.Request(&request, &reply); err != nil {
		sf.logger.Warnw("Failed to get master source info", "error", err)
		return nil, fmt.Errorf("get master source info: %w", err)
	}

	// create the master source session
	source := newMasterSession(sf.sessionLogger, sf.client, reply.SourceIndex, reply.Channels, false)

	return source, nil
}

func (sf *paSessionFinder) dumpProperties(rawProperties proto.PropList, name string) {
	properties := make(map[string]string, len(rawProperties))
	for key, value := range rawProperties {
		properties[key] = value.String()
	}
	sf.logger.Infow(fmt.Sprintf("LINUX %v", name), "properties", properties)
}

func (sf *paSessionFinder) enumerateAndAddSessions(sessions *[]Session) error {
	request := proto.GetSinkInputInfoList{}
	reply := proto.GetSinkInputInfoListReply{}

	if err := sf.client.Request(&request, &reply); err != nil {
		sf.logger.Warnw("Failed to get sink input list", "error", err)
		return fmt.Errorf("get sink input list: %w", err)
	}

	for _, info := range reply {
		if sf.config.Linux.LogSessionProperties {
			sf.dumpProperties(info.Properties, "sink")
		}

		sessionNames := buildSessionNames(
			info.Properties,
			proto.PropList{},
		)

		if sf.config.Linux.UseExtraAppNameFields {
			// For apps that connect via the native PipeWire protocol (e.g. Spotify), neither
			// application property is present on the sink-input itself — they live on the
			// associated PulseAudio client object instead.
			clientReq := proto.GetClientInfo{ClientIndex: info.ClientIndex}
			clientReply := proto.GetClientInfoReply{}
			if err := sf.client.Request(&clientReq, &clientReply); err == nil {
				if sf.config.Linux.LogSessionProperties {
					sf.dumpProperties(clientReply.Properties, "client")
				}

				sessionNames = buildSessionNames(
					info.Properties,
					clientReply.Properties,
				)
			}
		}

		if sf.config.Linux.LogSessionNames {
			if sf.config.Linux.UseExtraAppNameFields {
				sf.logger.Infow("Session",
					"sinkAppProcessBinary", sessionNames.sinkAppProcessBinary,
					"sinkAppName", sessionNames.sinkAppName,
					"clientAppProcessBinary", sessionNames.clientAppProcessBinary,
					"clientAppName", sessionNames.clientAppName,
					"mediaName", sessionNames.mediaName,
					"objectId", sessionNames.objectId,
				)
			} else {
				sf.logger.Infow("Session",
					"key", sessionNames.key,
					"sinkAppProcessBinary", sessionNames.sinkAppProcessBinary,
					"sinkAppName", sessionNames.sinkAppName,
					"mediaName", sessionNames.mediaName,
					"objectId", sessionNames.objectId,
				)
			}
		}

		if sessionNames.key == "" {
			sf.logger.Warnw("Failed to get sink input's process name",
				"sinkInputIndex", info.SinkInputIndex)
			continue
		}

		// create the deej session object
		newSession := newPASession(sf.sessionLogger, sf.client, info.SinkInputIndex, info.Channels, sessionNames)

		// add it to our slice
		*sessions = append(*sessions, newSession)

	}

	return nil
}
