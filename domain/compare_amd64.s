#include "textflag.h"

TEXT ·Compare(SB),NOSPLIT,$0-40
	MOVQ	aOff+0(FP), SI
	MOVQ	aLen+8(FP), CX
	MOVQ	bOff+16(FP), DI
	LEAQ	ret+32(FP), R9

	XORQ	AX, AX

loop:
	SUBQ	$8, CX
	JB	tail7
	MOVQ	(SI), AX
	MOVQ	(DI), BX
	ADDQ	$8, SI
	ADDQ	$8, DI
	SUBQ	BX, AX
	JZ	loop

	MOVQ	AX, (R9)
	RET

tail7:
	CMPQ	CX, $4
	JB	tail3
	MOVLQZX	(SI), AX
	MOVLQZX	(DI), DX
	SUBQ	DX, AX
	JNZ	ret
	ADDQ	$4, SI
	ADDQ	$4, DI
	SUBQ	$4, CX

tail3:
	CMPQ	CX, $2
	JB	tail1
	MOVWQZX	(SI), AX
	MOVWQZX	(DI), DX
	SUBQ	DX, AX
	JNZ	ret
	ADDQ	$2, SI
	ADDQ	$2, DI
	SUBQ	$2, CX

tail1:
	CMPQ    CX, $1
	JB	ret
	MOVBQZX (SI), AX
	MOVBQZX (DI), DX
	SUBQ	DX, AX

ret:
	MOVQ    AX, (R9)
	RET
