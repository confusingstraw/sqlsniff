package main

const PACKET_HEADER_SIZE = 4
const CONNECTION_ID_SIZE = 4
const AUTH_PLUGIN_DATA_PART_1_SIZE = 8
const FILLER_1_SIZE = 1
const CAPABILITY_FLAGS_1_SIZE = 2
const CHARACTER_SET_SIZE = 1
const STATUS_FLAGS_SIZE = 2
const CAPABILITY_FLAGS_2_SIZE = 2
const AUTH_PLUGIN_DATA_LEN_SIZE = 1
const AUTH_PLUGIN_DATA_RESERVED_PADDING_SIZE = 10
const MAX_AUTH_PLUGIN_LEN = 13

type ServerDescriptor struct {
	AuthenticationPlugin string // authentication plugin name, if specified
	Capabilities         uint32 // capability flags, if specified
	ConnectionId         uint32 // connection id, if specifed X
	ProtocolVersion      uint8  // version of the MySQL wire protocol used in the connection
	ServerVersion        string // version of MySQL running on the server X
	Status               uint16 // status flags, if specified
}
