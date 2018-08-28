package sdntable64

import (
	"bufio"
	"encoding/binary"
	"io/ioutil"
	"os"
)

func (d domains) merge() ([]uint32, int, int, string, error) {
	dst, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, 0, 0, "", err
	}
	defer dst.Close()

	w := bufio.NewWriter(dst)
	defer w.Flush()

	m := len(d.keys)/len(d.idx) - 1
	block := blocks[m]

	src, err := os.Open(*d.path)
	if err != nil {
		return nil, 0, 0, "", err
	}
	defer src.Close()

	r := bufio.NewReaderSize(src, block.b)

	memKeys := d.keys
	memValues := d.values
	memDIdx := d.dIdx
	i := 0

	srcKeys := make([]int64, m*block.n)
	srcValues := make([]uint64, block.n)
	sKeys := srcKeys
	j := 0

	dstKeys := make([]int64, m*block.n)
	dstValues := make([]uint64, block.n)
	dKeys := dstKeys
	k := 0

	var c uint32

	n := d.blocks
	rem := d.rem
	if n > 0 {
		if err := binary.Read(r, binary.LittleEndian, srcKeys); err != nil {
			return nil, 0, 0, dst.Name(), err
		}
		sKeys = srcKeys

		if err := binary.Read(r, binary.LittleEndian, srcValues); err != nil {
			return nil, 0, 0, dst.Name(), err
		}

		n--
	} else if rem > 0 {
		if err := binary.Read(r, binary.LittleEndian, srcKeys); err != nil {
			return nil, 0, 0, dst.Name(), err
		}
		sKeys = srcKeys[:m*rem]

		if err := binary.Read(r, binary.LittleEndian, srcValues); err != nil {
			return nil, 0, 0, dst.Name(), err
		}

		rem = 0
	}

	blks := 0

	for len(memKeys) > 0 || len(sKeys) > 0 {
		if len(memKeys) > 0 && len(sKeys) > 0 {
			res := compare(memKeys[:m], sKeys[:m])
			if res < 0 {
				copy(dKeys, memKeys[:m])
				memKeys = memKeys[m:]

				dstValues[k] = memValues[i]
				memDIdx[i] = c
				i++
			} else if res > 0 {
				copy(dKeys, sKeys[:m])
				sKeys = sKeys[m:]

				dstValues[k] = srcValues[j]
				j++
			} else {
				copy(dKeys, memKeys[:m])
				memKeys = memKeys[m:]
				sKeys = sKeys[m:]

				dstValues[k] = memValues[i]
				memDIdx[i] = c
				i++
				j++
			}
		} else if len(memKeys) > 0 {
			copy(dKeys, memKeys[:m])
			memKeys = memKeys[m:]

			dstValues[k] = memValues[i]
			memDIdx[i] = c
			i++
		} else {
			copy(dKeys, sKeys[:m])
			sKeys = sKeys[m:]

			dstValues[k] = srcValues[j]
			j++
		}

		if len(sKeys) <= 0 {
			if n > 0 {
				if err := binary.Read(r, binary.LittleEndian, srcKeys); err != nil {
					return nil, 0, 0, dst.Name(), err
				}
				sKeys = srcKeys

				if err := binary.Read(r, binary.LittleEndian, srcValues); err != nil {
					return nil, 0, 0, dst.Name(), err
				}

				n--
			} else if rem > 0 {
				if err := binary.Read(r, binary.LittleEndian, srcKeys); err != nil {
					return nil, 0, 0, dst.Name(), err
				}
				sKeys = srcKeys[:m*rem]

				if err := binary.Read(r, binary.LittleEndian, srcValues); err != nil {
					return nil, 0, 0, dst.Name(), err
				}

				rem = 0
			}
		}

		c++

		dKeys = dKeys[m:]
		k++

		if k >= block.n {
			if err := binary.Write(w, binary.LittleEndian, dstKeys); err != nil {
				return nil, 0, 0, dst.Name(), err
			}

			if err := binary.Write(w, binary.LittleEndian, dstValues); err != nil {
				return nil, 0, 0, dst.Name(), err
			}

			dKeys = dstKeys
			k = 0

			blks++
		}
	}

	if k > 0 {
		if err := binary.Write(w, binary.LittleEndian, dstKeys[:m*k]); err != nil {
			return nil, 0, 0, dst.Name(), err
		}

		if err := binary.Write(w, binary.LittleEndian, dstValues[:k]); err != nil {
			return nil, 0, 0, dst.Name(), err
		}
	}

	return memDIdx, blks, k, dst.Name(), nil
}

func (d domains) writeAll() (int, int, string, error) {
	f, err := ioutil.TempFile("", "*")
	if err != nil {
		return 0, 0, "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	m := len(d.keys)/len(d.idx) - 1
	block := blocks[m]

	keys := d.keys
	values := d.values

	blks := 0
	count := block.n
	n := len(d.idx)
	for n > count {
		if err := binary.Write(w, binary.LittleEndian, keys[:m*count]); err != nil {
			return 0, 0, f.Name(), err
		}
		keys = keys[m*count:]

		if err := binary.Write(w, binary.LittleEndian, values[:count]); err != nil {
			return 0, 0, f.Name(), err
		}
		values = values[count:]

		blks++
		n -= count
	}

	if n > 0 {
		if err := binary.Write(w, binary.LittleEndian, keys); err != nil {
			return 0, 0, f.Name(), err
		}

		if err := binary.Write(w, binary.LittleEndian, values); err != nil {
			return 0, 0, f.Name(), err
		}
	}

	return blks, n, f.Name(), nil
}

type dBlock struct {
	n int
	b int
}

var blocks = []dBlock{
	{
		n: 5120,
		b: 102400,
	},
	{
		n: 3664,
		b: 102592,
	},
	{
		n: 2848,
		b: 102528,
	},
	{
		n: 2336,
		b: 102784,
	},
	{
		n: 1984,
		b: 103168,
	},
	{
		n: 1712,
		b: 102720,
	},
	{
		n: 1520,
		b: 103360,
	},
	{
		n: 1360,
		b: 103360,
	},
	{
		n: 1232,
		b: 103488,
	},
	{
		n: 1120,
		b: 103040,
	},
	{
		n: 1024,
		b: 102400,
	},
	{
		n: 960,
		b: 103680,
	},
	{
		n: 896,
		b: 103936,
	},
	{
		n: 832,
		b: 103168,
	},
	{
		n: 784,
		b: 103488,
	},
	{
		n: 736,
		b: 103040,
	},
	{
		n: 704,
		b: 104192,
	},
	{
		n: 672,
		b: 104832,
	},
	{
		n: 640,
		b: 104960,
	},
	{
		n: 608,
		b: 104576,
	},
	{
		n: 576,
		b: 103680,
	},
	{
		n: 560,
		b: 105280,
	},
	{
		n: 528,
		b: 103488,
	},
	{
		n: 512,
		b: 104448,
	},
	{
		n: 496,
		b: 105152,
	},
	{
		n: 480,
		b: 105600,
	},
	{
		n: 464,
		b: 105792,
	},
	{
		n: 448,
		b: 105728,
	},
	{
		n: 432,
		b: 105408,
	},
	{
		n: 416,
		b: 104832,
	},
	{
		n: 400,
		b: 104000,
	},
	{
		n: 384,
		b: 102912,
	},
	{
		n: 384,
		b: 105984,
	},
	{
		n: 368,
		b: 104512,
	},
	{
		n: 352,
		b: 102784,
	},
	{
		n: 352,
		b: 105600,
	},
	{
		n: 336,
		b: 103488,
	},
	{
		n: 336,
		b: 106176,
	},
	{
		n: 320,
		b: 103680,
	},
	{
		n: 320,
		b: 106240,
	},
	{
		n: 304,
		b: 103360,
	},
	{
		n: 304,
		b: 105792,
	},
	{
		n: 288,
		b: 102528,
	},
	{
		n: 288,
		b: 104832,
	},
	{
		n: 288,
		b: 107136,
	},
	{
		n: 272,
		b: 103360,
	},
	{
		n: 272,
		b: 105536,
	},
	{
		n: 272,
		b: 107712,
	},
	{
		n: 256,
		b: 103424,
	},
	{
		n: 256,
		b: 105472,
	},
	{
		n: 256,
		b: 107520,
	},
	{
		n: 240,
		b: 102720,
	},
	{
		n: 240,
		b: 104640,
	},
	{
		n: 240,
		b: 106560,
	},
	{
		n: 240,
		b: 108480,
	},
	{
		n: 224,
		b: 103040,
	},
	{
		n: 224,
		b: 104832,
	},
	{
		n: 224,
		b: 106624,
	},
	{
		n: 224,
		b: 108416,
	},
	{
		n: 224,
		b: 110208,
	},
	{
		n: 208,
		b: 104000,
	},
	{
		n: 208,
		b: 105664,
	},
	{
		n: 208,
		b: 107328,
	},
	{
		n: 208,
		b: 108992,
	},
	{
		n: 208,
		b: 110656,
	},
	{
		n: 192,
		b: 103680,
	},
	{
		n: 192,
		b: 105216,
	},
	{
		n: 192,
		b: 106752,
	},
	{
		n: 192,
		b: 108288,
	},
	{
		n: 192,
		b: 109824,
	},
	{
		n: 192,
		b: 111360,
	},
	{
		n: 176,
		b: 103488,
	},
	{
		n: 176,
		b: 104896,
	},
	{
		n: 176,
		b: 106304,
	},
	{
		n: 176,
		b: 107712,
	},
	{
		n: 176,
		b: 109120,
	},
	{
		n: 176,
		b: 110528,
	},
	{
		n: 176,
		b: 111936,
	},
	{
		n: 160,
		b: 103040,
	},
	{
		n: 160,
		b: 104320,
	},
	{
		n: 160,
		b: 105600,
	},
	{
		n: 160,
		b: 106880,
	},
	{
		n: 160,
		b: 108160,
	},
	{
		n: 160,
		b: 109440,
	},
	{
		n: 160,
		b: 110720,
	},
	{
		n: 160,
		b: 112000,
	},
	{
		n: 160,
		b: 113280,
	},
	{
		n: 144,
		b: 103104,
	},
	{
		n: 144,
		b: 104256,
	},
	{
		n: 144,
		b: 105408,
	},
	{
		n: 144,
		b: 106560,
	},
	{
		n: 144,
		b: 107712,
	},
	{
		n: 144,
		b: 108864,
	},
	{
		n: 144,
		b: 110016,
	},
	{
		n: 144,
		b: 111168,
	},
	{
		n: 144,
		b: 112320,
	},
	{
		n: 144,
		b: 113472,
	},
	{
		n: 144,
		b: 114624,
	},
	{
		n: 128,
		b: 102912,
	},
	{
		n: 128,
		b: 103936,
	},
	{
		n: 128,
		b: 104960,
	},
	{
		n: 128,
		b: 105984,
	},
	{
		n: 128,
		b: 107008,
	},
	{
		n: 128,
		b: 108032,
	},
	{
		n: 128,
		b: 109056,
	},
	{
		n: 128,
		b: 110080,
	},
	{
		n: 128,
		b: 111104,
	},
	{
		n: 128,
		b: 112128,
	},
	{
		n: 128,
		b: 113152,
	},
	{
		n: 128,
		b: 114176,
	},
	{
		n: 128,
		b: 115200,
	},
	{
		n: 128,
		b: 116224,
	},
	{
		n: 112,
		b: 102592,
	},
	{
		n: 112,
		b: 103488,
	},
	{
		n: 112,
		b: 104384,
	},
	{
		n: 112,
		b: 105280,
	},
	{
		n: 112,
		b: 106176,
	},
	{
		n: 112,
		b: 107072,
	},
	{
		n: 112,
		b: 107968,
	},
	{
		n: 112,
		b: 108864,
	},
	{
		n: 112,
		b: 109760,
	},
	{
		n: 112,
		b: 110656,
	},
	{
		n: 112,
		b: 111552,
	},
	{
		n: 112,
		b: 112448,
	},
	{
		n: 112,
		b: 113344,
	},
	{
		n: 112,
		b: 114240,
	},
	{
		n: 112,
		b: 115136,
	},
}
