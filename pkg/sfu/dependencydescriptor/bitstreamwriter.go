package dependencydescriptor

import (
	"errors"
	"fmt"
)

// BitStreamWriter 字节流写入
type BitStreamWriter struct {
	buf       []byte
	pos       int
	bitOffset int // 当前byte的bit偏移量
}

func NewBitStreamWriter(buf []byte) *BitStreamWriter {
	return &BitStreamWriter{buf: buf}
}

func (w *BitStreamWriter) RemainingBits() int {
	return (len(w.buf)-w.pos)*8 - w.bitOffset
}

func (w *BitStreamWriter) WriteBits(val uint64, bitCount int) error {
	if bitCount > w.RemainingBits() {
		return errors.New("insufficient space")
	}

	totalBits := bitCount

	// push bits to the highest bits of uint64
	val <<= 64 - bitCount

	buf := w.buf[w.pos:]

	// 第一个字节比较特殊；要写入的位偏移量可能会使我们
	// 在字节的中间，并且要写入的总位数可能需要
	// 保存字节末尾的位。
	remainingBitsInCurrentByte := 8 - w.bitOffset
	bitsInFirstByte := bitCount
	if bitsInFirstByte > remainingBitsInCurrentByte {
		bitsInFirstByte = remainingBitsInCurrentByte
	}

	buf[0] = w.writePartialByte(uint8(val>>56), bitsInFirstByte, buf[0], w.bitOffset)

	if bitCount <= remainingBitsInCurrentByte {
		//没有可写的位
		return w.consumeBits(totalBits)
	}

	// 写入其余的位
	val <<= bitsInFirstByte
	buf = buf[1:]
	bitCount -= bitsInFirstByte
	for bitCount >= 8 {
		buf[0] = uint8(val >> 56)
		buf = buf[1:]
		val <<= 8
		bitCount -= 8
	}

	//写入最后的位
	if bitCount > 0 {
		buf[0] = w.writePartialByte(uint8(val>>56), bitCount, buf[0], 0)
	}
	return w.consumeBits(totalBits)
}

func (w *BitStreamWriter) consumeBits(bitCount int) error {
	if bitCount > w.RemainingBits() {
		return errors.New("insufficient space")
	}

	w.pos += (w.bitOffset + bitCount) / 8
	w.bitOffset = (w.bitOffset + bitCount) % 8

	return nil
}

func (w *BitStreamWriter) writePartialByte(source uint8, sourceBitCount int, target uint8, targetBitOffset int) uint8 {
	// if !(targetBitOffset < 8 &&  sourceBitCount <= (8-targetBitOffset)) {
	// 	return fmt.Errorf("invalid argument, source %d, sourceBitCount %d, target %d, targetBitOffset %d", source, sourceBitCount, target, targetBitOffset)
	// }

	// 为要重写的位生成掩码，将源位移位到最高位，然后定位到目标位偏移
	mask := uint8(0xff<<(8-sourceBitCount)) >> uint8(targetBitOffset)
	// 清除目标位并写入源位
	return (target &^ mask) | (source >> targetBitOffset)
}

func (w *BitStreamWriter) WriteNonSymmetric(val, numValues uint32) error {
	if !(val < numValues && numValues <= 1<<31) {
		return fmt.Errorf("invalid argument, val %d, numValues %d", val, numValues)
	}
	if numValues == 1 {
		return nil
	}

	countBits := bitwidth(numValues)
	numMinBitsValues := (uint32(1) << countBits) - numValues
	if val < numMinBitsValues {
		return w.WriteBits(uint64(val), countBits-1)
	} else {
		return w.WriteBits(uint64(val+numMinBitsValues), countBits)
	}
}

func SizeNonSymmetricBits(val, numValues uint32) int {
	countBits := bitwidth(numValues)
	numMinBitsValues := (uint32(1) << countBits) - numValues
	if val < numMinBitsValues {
		return countBits - 1
	} else {
		return countBits
	}
}
