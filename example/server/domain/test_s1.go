package domain

import "errors"

type TestS1 struct {
	Id      [16]byte
	Emmiter [16]byte
}

func SerializeTestS1(t TestS1) []byte {
	data := []byte{}
	data = append(data, t.Id[:]...)
	data = append(data, t.Emmiter[:]...)
	return data
}

func DeserializeTestS1(b []byte) (TestS1, error) {
	size := len(b)
	if size != 32 {
		return TestS1{}, errors.New("invalid test_s1 size for deserialization")
	}

	id := b[:16]
	emitter := b[16:]

	return TestS1{[16]byte(id), [16]byte(emitter)}, nil
}
