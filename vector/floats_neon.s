//go:build !noasm && arm64
// AUTO-GENERATED BY GOAT -- DO NOT EDIT

TEXT ·vdot(SB), $0-32
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	MOVD ret+24(FP), R3
	WORD $0xa9bf7bfd    // stp	x29, x30, [sp,
	WORD $0x91000c48    // add	x8, x2,
	WORD $0xf100005f    // cmp	x2,
	WORD $0x9a82b108    // csel	x8, x8, x2, lt
	WORD $0x9342fd0a    // asr	x10, x8,
	WORD $0x927ef508    // and	x8, x8,
	WORD $0x7100055f    // cmp	w10,
	WORD $0xcb080048    // sub	x8, x2, x8
	WORD $0x910003fd    // mov	x29, sp
	WORD $0x540002ab    // b.lt	.LBB0_5
	WORD $0x3cc10400    // ldr	q0, [x0],
	WORD $0x3cc10421    // ldr	q1, [x1],
	WORD $0x71000549    // subs	w9, w10,
	WORD $0x6e21dc00    // fmul	v0.4s, v0.4s, v1.4s
	WORD $0x54000200    // b.eq	.LBB0_6
	WORD $0xb27d7beb    // mov	x11,
	WORD $0x8b0a096a    // add	x10, x11, x10, lsl
	WORD $0x927e7d4a    // and	x10, x10,
	WORD $0x9100114b    // add	x11, x10,
	WORD $0x8b0b080a    // add	x10, x0, x11, lsl
	WORD $0xaa0103ec    // mov	x12, x1

LBB0_3:
	WORD $0x3cc10401 // ldr	q1, [x0],
	WORD $0x3cc10582 // ldr	q2, [x12],
	WORD $0x71000529 // subs	w9, w9,
	WORD $0x6e22dc21 // fmul	v1.4s, v1.4s, v2.4s
	WORD $0x4e21d400 // fadd	v0.4s, v0.4s, v1.4s
	WORD $0x54ffff61 // b.ne	.LBB0_3
	WORD $0x8b0b0821 // add	x1, x1, x11, lsl
	WORD $0xaa0a03e0 // mov	x0, x10
	WORD $0x14000001 // b	.LBB0_6

LBB0_5:
LBB0_6:
	WORD $0x1e2703e1 // fmov	s1, wzr
	WORD $0x5e0c0402 // mov	s2, v0.s[1]
	WORD $0x5e140403 // mov	s3, v0.s[2]
	WORD $0x5e1c0404 // mov	s4, v0.s[3]
	WORD $0x1e212800 // fadd	s0, s0, s1
	WORD $0x1e202840 // fadd	s0, s2, s0
	WORD $0x1e202860 // fadd	s0, s3, s0
	WORD $0x1e202880 // fadd	s0, s4, s0
	WORD $0x7100011f // cmp	w8,
	WORD $0xbd000060 // str	s0, [x3]
	WORD $0x5400012d // b.le	.LBB0_9
	WORD $0x92407d08 // and	x8, x8,

LBB0_8:
	WORD $0xbc404401 // ldr	s1, [x0],
	WORD $0xbc404422 // ldr	s2, [x1],
	WORD $0xf1000508 // subs	x8, x8,
	WORD $0x1e220821 // fmul	s1, s1, s2
	WORD $0x1e212800 // fadd	s0, s0, s1
	WORD $0xbd000060 // str	s0, [x3]
	WORD $0x54ffff41 // b.ne	.LBB0_8

LBB0_9:
	WORD $0xa8c17bfd // ldp	x29, x30, [sp],
	WORD $0xd65f03c0 // ret
