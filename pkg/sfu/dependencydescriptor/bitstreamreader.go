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

func NewBitStreamReader(buf []byte) *BitStreamReader {
	return &BitStreamReader{buf: buf, remainingBits: len(buf) * 8}
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

// ReadBool 从buf中读取bool值
func (b *BitStreamReader) ReadBool() (bool, error) {
	val, err := b.ReadBits(1)
	return val != 0, err
}

func (b *BitStreamReader) Ok() bool {
	return b.remainingBits >= 0
}

// Invalidate 使无效
func (b *BitStreamReader) Invalidate() {
	b.remainingBits = -1
}

//ReadNonSymmetric 读取范围[0，`num_values`-1]中的值。
//这种编码类似于ReadBits（val、Ceil（Log2（num_values）），
//但是减少了当对两个值范围的非幂进行编码时产生的浪费
//非对称值编码为：
// 1) n = bit_width(num_values)
// 2) k = (1 << n) - num_values
//范围[0，k-1]中的值v以（n-1）位编码。
//范围[k，num_values-1]中的值v被编码为n位的（v+k）。
//https://aomediacodec.github.io/av1-spec/#nsn
func (b *BitStreamReader) ReadNonSymmetric(numValues uint32) (uint32, error) {
	if numValues >= (uint32(1) << 31) {
		return 0, errors.New("invalid number of values, expected 0-2^31")
	}

	width := bitwidth(numValues)
	numMinBitsValue := (uint32(1) << width) - numValues

	val, err := b.ReadBits(width - 1)
	if err != nil {
		return 0, err
	}
	if val < uint64(numMinBitsValue) {
		return uint32(val), nil
	}
	bit, err := b.ReadBits(1)
	if err != nil {
		return 0, err
	}
	return uint32((val << 1) + bit - uint64(numMinBitsValue)), nil
}

func (b *BitStreamReader) BytesRead() int {
	if b.remainingBits%8 > 0 {
		return b.pos + 1
	}
	return b.pos
}

// bitwidth bit宽度
func bitwidth(n uint32) int {
	var w int
	for n != 0 {
		n >>= 1
		w++
	}
	return w
}
