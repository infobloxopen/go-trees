#include "textflag.h"

TEXT ·Compare(SB),NOSPLIT,$0-40
	MOVQ	aOff+0(FP), SI
	MOVQ	aLen+8(FP), CX
	MOVQ	bOff+16(FP), DI
	MOVQ	bLen+24(FP), AX
	LEAQ	ret+32(FP), R9

	XORQ	CX, AX
	JNZ	ret

loop:

	SUBQ	$8, CX
	JB	ret
	MOVQ	(SI), AX
	MOVQ	(DI), BX
	ADDQ	$8, SI
	ADDQ	$8, DI
	SUBQ	BX, AX
	JZ	loop

	MOVQ	AX, (R9)
	RET

ret:
	XORQ	AX, AX
	MOVQ    AX, (R9)
	RET
