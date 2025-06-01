package internal

type ClientMessageType uint16
type ServerMessageType uint16

// WARNING! This values must be kept in sync with client ones.

// Client message types.
const (
	ClientResponse       ClientMessageType = 0
	FullSyncRequest      ClientMessageType = 1
	HostSyncRequest      ClientMessageType = 2
	HostTextUpdate       ClientMessageType = 3
	SyncThisHost         ClientMessageType = 4
	ClientMessageTypeMax ClientMessageType = SyncThisHost
)

// Server message types.
const (
	ServerResponse       ServerMessageType = 256
	HostConnected        ServerMessageType = 257
	HostDisconnected     ServerMessageType = 258
	TextUpdate           ServerMessageType = 259
	HostSynced           ServerMessageType = 260
	ServerMessageTypeMax ServerMessageType = HostSynced
)
