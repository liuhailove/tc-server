package dependencydescriptor

import (
	"errors"
	"io"
)

// BitStreamReader 比特流Reader
type BitStreamReader struct {
	buf           []byte // 缓冲大小
	pos           int    // 为止
	remainingBits int    // 剩余bit
}

// RemainingBits 剩余的bit数
func (b *BitStreamReader) RemainingBits() int {
	return b.remainingBits
}

// ReadBits 从比特流中读取“bits”。 `bits` 必须在 [0, 64] 范围内。
// 返回 [0, 2^bits - 1] 范围内的无符号整数。
// 失败时将 `BitstreamReader` 设置为失败状态并返回 0。
func (b *BitStreamReader) ReadBits(bits int) (uint64, error) {
	if bits < 0 || bits > 64 {
		return 0, errors.New("invalid number of bits, expected 0-64")
	}

	if b.remainingBits < bits {
		b.remainingBits -= bits
		return 0, io.EOF
	}

	remainingBitsInFirstByte := b.remainingBits % 8
	b.remainingBits -= bits
	if bits < remainingBitsInFirstByte {
		// 读取的位数少于当前字节中剩余的位数
		// 返回该字节中需要的部分。
		offset := remainingBitsInFirstByte - bits
		return uint64((b.buf[b.pos] >> offset) & ((1 << bits) - 1)), nil
	}

	var result uint64
	if remainingBitsInFirstByte > 0 {
		// 读取当前字节中剩余的所有位并消耗该字节。
		bits -= remainingBitsInFirstByte
		mask := byte((1 << remainingBitsInFirstByte) - 1)
		result = uint64(b.buf[b.pos]&mask) << bits
		b.pos++
	}

	// 读取尽可能多的完整字节。
	for bits >= 8 {
		bits -= 8
		result += uint64(b.buf[b.pos]) << bits
		b.pos++
	}

	// 剩下要读取的内容都小于一个字节，所以只获取需要的内容
	// 位并将它们移至最低位
	if bits > 0 {
		result |= uint64(b.buf[b.pos] >> (8 - bits))
	}
	return result, nil
}
