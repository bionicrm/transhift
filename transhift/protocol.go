package transhift

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

type Serializable interface {
    Serialize() []byte

    Deserialize([]byte)
}

type ProtoMetaInfo struct {
    passwordHash []byte
    fileName     string
    fileSize     uint64
    fileHash     []byte
}

func (m *ProtoMetaInfo) Serialize() []byte {
    var buffer bytes.Buffer

    // passwordHash
    buffer.Write(m.passwordHash)
    buffer.WriteRune('\n')
    // fileName
    buffer.WriteString(m.fileName)
    buffer.WriteRune('\n')
    // fileSize
    fileSizeBuffer := make([]byte, 8)
    binary.BigEndian.PutUint64(fileSizeBuffer, m.fileSize)
    buffer.Write(fileSizeBuffer)
    buffer.WriteRune('\n')
    // fileHash
    buffer.Write(m.fileHash)
    buffer.WriteRune('\n')

    return buffer.Bytes()
}

func (m *ProtoMetaInfo) Deserialize(b []byte) {
    buffer := bytes.NewBuffer(b)

    // passwordHash
    m.passwordHash, _ = buffer.ReadBytes('\n')
    m.passwordHash = m.passwordHash[:len(m.passwordHash) - 1] // trim leading \n
    // fileName
    m.fileName, _ = buffer.ReadString('\n')
    m.fileName = m.fileName[:len(m.fileName) - 1] // trim leading \n
    // fileSize
    fileSize, _ := buffer.ReadBytes('\n')
    fileSize = fileSize[:len(fileSize) - 1] // trim leading \n
    m.fileSize = binary.BigEndian.Uint64(fileSize)
    // fileHash
    m.fileHash, _ = buffer.ReadBytes('\n')
    m.fileHash = m.fileHash[:len(m.fileHash) - 1] // trim leading \n
}

func (m *ProtoMetaInfo) String() string {
    return fmt.Sprintf("{passwordHash=%x, fileName=%s, fileSize=%d, fileHash=%x}",
        m.passwordHash, m.fileName, m.fileSize, m.fileHash)
}

type ProtoChunkInfo struct {
    close bool
    data  []byte
}

func (c *ProtoChunkInfo) Serialize() []byte {
    var buffer bytes.Buffer

    // close
    buffer.WriteByte(Btobyte(c.close))
    // data
    buffer.Write(c.data)

    return buffer.Bytes()
}

func (c *ProtoChunkInfo) Deserialize(b []byte) {
    // close
    c.close = Btobool(b[0])
    // data
    c.data = b[1:]
}

func (c *ProtoChunkInfo) String() string {
    return fmt.Sprintf("{close=%t, data=(len)%d}", c.close, len(c.data))
}

func Btobyte(b bool) byte {
    if b {
        return 0x01
    }
    return 0x00
}

func Btobool(b byte) bool {
    if b == 0x00 {
        return false
    }
    return true
}
