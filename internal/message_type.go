package internal

type ClientMessageType uint16
type ServerMessageType uint16

// WARNING! This values must be kept in sync with client ones.

// Client message types.
const (
	ClientResponse       ClientMessageType = 0
	ClientIntroduction   ClientMessageType = 1
	FullSyncRequest      ClientMessageType = 2
	HostSyncRequest      ClientMessageType = 3
	HostTextUpdate       ClientMessageType = 4
	SyncThisHost         ClientMessageType = 5
	ClientMessageTypeMax ClientMessageType = SyncThisHost
)

// Server message types.
const (
	ServerResponse       ServerMessageType = 256
	ServerIntroduction   ServerMessageType = 257
	HostConnected        ServerMessageType = 258
	HostDisconnected     ServerMessageType = 259
	TextUpdate           ServerMessageType = 260
	HostSynced           ServerMessageType = 261
	ServerMessageTypeMax ServerMessageType = HostSynced
)
