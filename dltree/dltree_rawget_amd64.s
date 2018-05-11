// +build amd64

#include "textflag.h"
TEXT ·treeRawGet(SB),NOSPLIT,$0-56
	MOVQ	t+0(FP), R8                      // +00: t -> R8
		                                 //   t+00: t.root

	TESTQ	R8, R8                           // if t(R8) == nil return
	JZ	ret

	MOVQ	(R8), R8                         // t.root(t(R8)+0) -> n(R8)
		                                 //   n+00: ptr(n.key)
		                                 //   n+08: len(n.key)
		                                 //   n+16: n.keySize
		                                 //   n+24: n.value1
		                                 //   n+32: n.value2
		                                 //   n+40: n.chld[0]
		                                 //   n+48: n.chld[1]
		                                 //   n+56: n.red

	MOVQ	keyOff+8(FP), R9                 // +08: ptr(key) -> R9
	MOVQ	keyLen+16(FP), R10               // +16: len(key) -> R10

	MOVQ	keySize+24(FP), DX               // +24: keySize -> DX

loop:
	TESTQ	R8, R8                           // if n == nil return
	JZ	ret

	MOVQ	16(R8), CX                       // n.keySize -> CX
	CMPQ	CX, DX                           // n.keySize(CX) - keySize(DX)?
	JE	cmp                              // compare labels if equal
	JL	right                            // go to right node if less
		                                 // otherwise ...
left:
	MOVQ	40(R8), R8                       // go to left node n.chld[0](n(R8)+40) -> n(R8)
	JMP	loop

right:
	MOVQ	48(R8), R8                       // go to right node n.chld[1](n(R8)+48) -> n(R8)
	JMP	loop

cmp:
	MOVQ	0(R8), SI                        // ptr(n.key(n(R8)+0)) -> SI
	MOVQ	R9, DI                           // ptr(key(R9)) -> DI
	MOVQ	R10, CX                          // len(key(R10)) -> CX

cmp_loop:
	CMPQ	CX, $8                           // for len(key)(CX) >= 8
	JB	found                            // go to return n.value, true if len(key)(CX) < 8 (labels match)

	MOVQ	(SI), AX                         // *((n.key + i)(SI)) -> AX
	MOVQ	(DI), BX                         // *((key + j)(DI)) -> BX
	ADDQ	$8, SI                           // i += 8
	ADDQ	$8, DI                           // j += 8
	SUBQ	$8, CX                           // len(key)(CX) -= 8
	SUBQ	BX, AX                           // n.key[i] - key[j]?
	JZ	cmp_loop                         // go to next 8 bytes if equal
	JL	right                            // go to right node if less
	JMP	left                             // otherwise go to left node

found:
	MOVQ	24(R8), AX                       // n.value1(n(R8)+24) -> AX
	MOVQ	AX, ret+32(FP)                   // n.value1(AX) -> ret interface{}1
	MOVQ	32(R8), AX                       // n.value2(n(R8)+32) -> AX
	MOVQ	AX, ret+40(FP)                   // n.value1(AX) -> ret interface{}2
	MOVB	$1, ret+48(FP)                   // true -> ret bool
	RET

ret:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)                   // nil -> ret interface{}
	MOVQ	AX, ret+40(FP)
	MOVB	AL, ret+48(FP)                   // false -> ret bool
	RET
