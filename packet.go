package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// implements the conditional size described in the packet docs
func getPluginPart2DataLen(v uint8) int {
	if v == 0 {
		return 0
	}

	if MAX_AUTH_PLUGIN_LEN > (v - 8) {
		return MAX_AUTH_PLUGIN_LEN
	}

	return int(v - 8)
}

// implements reading the packet header for a MySQL connection
// https://dev.mysql.com/doc/internals/en/mysql-packet.html#idm45663064450400
func ReadPacketSize(r *bufio.Reader) (int, error) {
	header := make([]byte, PACKET_HEADER_SIZE)
	_, err := io.ReadFull(r, header)

	if err != nil {
		return -1, err
	}

	n := ToFixedLengthInt3(header)

	if n <= 0 {
		return -1, errors.New("invalid packet size")
	}

	return n, nil
}

// reads the contents of a single packet into a byte slice
// https://dev.mysql.com/doc/internals/en/mysql-packet.html
func ReadPacket(r *bufio.Reader) ([]byte, int, error) {
	n, err := ReadPacketSize(r)
	if err != nil {
		return nil, -1, err
	}

	p := make([]byte, n)

	n, err = r.Read(p)
	if err != nil {
		return nil, -1, err
	}

	return p, n, nil
}

// reads various fields from a HandshakeV10 packet
// https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake
func ParseDescriptor(p []byte, n int) (*ServerDescriptor, error) {
	var d ServerDescriptor
	r := bufio.NewReader(bytes.NewReader(p))

	// ingest connection id
	u, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	d.ProtocolVersion = uint8(u)

	// ingest server version
	str, err := r.ReadString(0)
	if err != nil {
		return nil, err
	}
	d.ServerVersion = str

	// ingest connection id
	b := make([]byte, CONNECTION_ID_SIZE)
	if _, err = io.ReadFull(r, b); err != nil {
		return nil, err
	}
	d.ConnectionId = binary.LittleEndian.Uint32(b)

	// skip unused values
	if _, err = r.Discard(AUTH_PLUGIN_DATA_PART_1_SIZE + FILLER_1_SIZE); err != nil {
		return nil, err
	}

	// ingest lower capability flag bytes
	capabilities_lo := make([]byte, CAPABILITY_FLAGS_1_SIZE)
	if _, err = io.ReadFull(r, capabilities_lo); err != nil {
		return nil, err
	}

	// skip unused value
	if _, err = r.Discard(CHARACTER_SET_SIZE); err != nil {
		return nil, err
	}

	// ingest lower capability flag bytes
	b = make([]byte, STATUS_FLAGS_SIZE)
	if _, err = io.ReadFull(r, b); err != nil {
		return nil, err
	}
	d.Status = binary.LittleEndian.Uint16(b)

	// ingest upper capability flag bytes
	capabilities_hi := make([]byte, CAPABILITY_FLAGS_2_SIZE)
	if _, err = io.ReadFull(r, capabilities_hi); err != nil {
		return nil, err
	}
	d.Capabilities = binary.LittleEndian.Uint32(append(capabilities_lo, capabilities_hi...))

	// ingest auth plugin data len and skip that many bytes as well as the reserved padding
	u, err = r.ReadByte()
	if err != nil {
		return nil, err
	}

	_, err = r.Discard(getPluginPart2DataLen(uint8(u)) + AUTH_PLUGIN_DATA_RESERVED_PADDING_SIZE)
	if err != nil {
		return nil, err
	}

	// only if the plugin auth capability is supported
	if u != 0 {
		// ingest auth plugin name
		str, err = r.ReadString(0)
		if err != nil {
			return nil, err
		}
		d.AuthenticationPlugin = str
	}

	return &d, nil
}
