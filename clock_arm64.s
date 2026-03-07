#include "textflag.h"

// func counterValue() uint64
TEXT ·counterValue(SB),NOSPLIT,$0-8
	MRS CNTVCT_EL0, R0
	MOVD R0, ret+0(FP)
	RET

// func counterFreq() uint64
TEXT ·counterFreq(SB),NOSPLIT,$0-8
	MRS CNTFRQ_EL0, R0
	MOVD R0, ret+0(FP)
	RET
