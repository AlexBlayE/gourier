package domain

import (
	"encoding/binary"
	"errors"
)

type TestS2 struct {
	TestS1
	Count uint32
}

func SerializeTestS2(t TestS2) []byte {
	data := []byte{}
	data = append(data, SerializeTestS1(t.TestS1)...)
	data = binary.BigEndian.AppendUint32(data, uint32(t.Count))
	return data
}

func DeserializeTestS2(b []byte) (TestS2, error) {
	size := len(b)
	if size != 32+4 {
		return TestS2{}, errors.New("invalid test_s1 size for deserialization")
	}

	t1, err := DeserializeTestS1(b[:33])
	if err != nil {
		return TestS2{}, err
	}

	count := binary.BigEndian.Uint32(b[33:])

	return TestS2{t1, count}, nil
}
