	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 26, 0	sdk_version 26, 2
	.globl	_ffnt                           ## -- Begin function ffnt
	.p2align	4, 0x90
_ffnt:                                  ## @ffnt
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$32, %rsp
	movq	%rdi, -8(%rbp)
	movq	-8(%rbp), %rdi
	callq	_strlen
	movslq	%eax, %rdi
	callq	_malloc
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rax
	movq	%rax, -24(%rbp)
	movw	$0, -26(%rbp)
LBB0_1:                                 ## =>This Inner Loop Header: Depth=1
	movq	-8(%rbp), %rax
	movsbl	(%rax), %eax
	cmpl	$0, %eax
	je	LBB0_8
## %bb.2:                               ##   in Loop: Header=BB0_1 Depth=1
	movq	-8(%rbp), %rax
	movsbl	(%rax), %eax
	cmpl	$32, %eax
	je	LBB0_4
## %bb.3:                               ##   in Loop: Header=BB0_1 Depth=1
	movq	-8(%rbp), %rax
	movsbl	(%rax), %eax
	movw	%ax, %cx
	movq	-16(%rbp), %rax
	movw	%cx, (%rax)
	movq	-16(%rbp), %rax
	addq	$2, %rax
	movq	%rax, -16(%rbp)
	jmp	LBB0_7
LBB0_4:                                 ##   in Loop: Header=BB0_1 Depth=1
	movswl	-26(%rbp), %eax
	cmpl	$0, %eax
	jne	LBB0_6
## %bb.5:                               ##   in Loop: Header=BB0_1 Depth=1
	movw	$1, -26(%rbp)
	movq	-16(%rbp), %rax
	movw	$46, (%rax)
	movq	-16(%rbp), %rax
	addq	$2, %rax
	movq	%rax, -16(%rbp)
LBB0_6:                                 ##   in Loop: Header=BB0_1 Depth=1
	jmp	LBB0_7
LBB0_7:                                 ##   in Loop: Header=BB0_1 Depth=1
	movq	-8(%rbp), %rax
	addq	$1, %rax
	movq	%rax, -8(%rbp)
	jmp	LBB0_1
LBB0_8:
	movq	-24(%rbp), %rax
	addq	$32, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_fntf                           ## -- Begin function fntf
	.p2align	4, 0x90
_fntf:                                  ## @fntf
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$48, %rsp
	movq	%rdi, -8(%rbp)
	movl	$16, %edi
	callq	_malloc
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rax
	movq	%rax, -24(%rbp)
	movw	$0, -26(%rbp)
	movl	$0, -32(%rbp)
LBB1_1:                                 ## =>This Loop Header: Depth=1
                                        ##     Child Loop BB1_5 Depth 2
	movq	-8(%rbp), %rax
	movsbl	(%rax), %eax
	cmpl	$0, %eax
	je	LBB1_12
## %bb.2:                               ##   in Loop: Header=BB1_1 Depth=1
	movq	-8(%rbp), %rax
	movsbl	(%rax), %eax
	cmpl	$46, %eax
	jne	LBB1_10
## %bb.3:                               ##   in Loop: Header=BB1_1 Depth=1
	movswl	-26(%rbp), %eax
	cmpl	$0, %eax
	jne	LBB1_9
## %bb.4:                               ##   in Loop: Header=BB1_1 Depth=1
	movw	$1, -26(%rbp)
	movl	$12, -36(%rbp)
	movl	-36(%rbp), %eax
	subl	-32(%rbp), %eax
	movl	%eax, -40(%rbp)
	movl	$0, -44(%rbp)
LBB1_5:                                 ##   Parent Loop BB1_1 Depth=1
                                        ## =>  This Inner Loop Header: Depth=2
	movl	-44(%rbp), %eax
	cmpl	-40(%rbp), %eax
	jge	LBB1_8
## %bb.6:                               ##   in Loop: Header=BB1_5 Depth=2
	movq	-16(%rbp), %rax
	movw	$32, (%rax)
	movq	-16(%rbp), %rax
	addq	$2, %rax
	movq	%rax, -16(%rbp)
## %bb.7:                               ##   in Loop: Header=BB1_5 Depth=2
	movl	-44(%rbp), %eax
	addl	$1, %eax
	movl	%eax, -44(%rbp)
	jmp	LBB1_5
LBB1_8:                                 ##   in Loop: Header=BB1_1 Depth=1
	movq	-8(%rbp), %rax
	addq	$1, %rax
	movq	%rax, -8(%rbp)
LBB1_9:                                 ##   in Loop: Header=BB1_1 Depth=1
	jmp	LBB1_11
LBB1_10:                                ##   in Loop: Header=BB1_1 Depth=1
	movq	-8(%rbp), %rax
	movsbl	(%rax), %eax
	movw	%ax, %cx
	movq	-16(%rbp), %rax
	movw	%cx, (%rax)
	movl	-32(%rbp), %eax
	addl	$1, %eax
	movl	%eax, -32(%rbp)
	movq	-8(%rbp), %rax
	addq	$1, %rax
	movq	%rax, -8(%rbp)
	movq	-16(%rbp), %rax
	addq	$2, %rax
	movq	%rax, -16(%rbp)
LBB1_11:                                ##   in Loop: Header=BB1_1 Depth=1
	jmp	LBB1_1
LBB1_12:
	movq	-24(%rbp), %rax
	addq	$48, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_fcreate                        ## -- Begin function fcreate
	.p2align	4, 0x90
_fcreate:                               ## @fcreate
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$64, %rsp
	movq	%rdi, -8(%rbp)
	movq	%rsi, -16(%rbp)
	movl	$1564, %eax                     ## imm = 0x61C
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -32(%rbp)
	movq	-32(%rbp), %rax
	movq	$1279677254, (%rax)             ## imm = 0x4C465346
	movq	-32(%rbp), %rax
	addq	$32, %rax
	movq	%rax, -32(%rbp)
	movq	-8(%rbp), %rdi
	movq	-32(%rbp), %rsi
	callq	_strcpy
	movq	-8(%rbp), %rdi
	callq	_strlen
	cltq
	movq	%rax, -40(%rbp)
	movq	-32(%rbp), %rax
	movq	-40(%rbp), %rcx
	shlq	$3, %rcx
	addq	%rcx, %rax
	movq	%rax, -32(%rbp)
	movq	-16(%rbp), %rcx
	movq	-32(%rbp), %rax
	movq	%rcx, (%rax)
	movq	-32(%rbp), %rax
	addq	$32, %rax
	movq	%rax, -32(%rbp)
	movq	$0, -48(%rbp)
LBB2_1:                                 ## =>This Inner Loop Header: Depth=1
	movq	-48(%rbp), %rax
	cmpq	-16(%rbp), %rax
	jge	LBB2_4
## %bb.2:                               ##   in Loop: Header=BB2_1 Depth=1
	movq	-32(%rbp), %rax
	movq	$0, (%rax)
	movq	-32(%rbp), %rax
	addq	$8, %rax
	movq	%rax, -32(%rbp)
## %bb.3:                               ##   in Loop: Header=BB2_1 Depth=1
	movq	-48(%rbp), %rax
	addq	$1, %rax
	movq	%rax, -48(%rbp)
	jmp	LBB2_1
LBB2_4:
	movq	-32(%rbp), %rax
	movl	$512, %ecx                      ## imm = 0x200
	cqto
	idivq	%rcx
	movq	%rax, -56(%rbp)
	movq	-56(%rbp), %rdi
	callq	_save_sector
	movq	-24(%rbp), %rax
	movq	(%rax), %rcx
	movq	-16(%rbp), %rax
	shlq	$3, %rax
	addq	%rax, %rcx
	movq	-24(%rbp), %rax
	movq	%rcx, (%rax)
	movl	$3, %edi
	callq	_save_sector
	addq	$64, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_find_file                      ## -- Begin function find_file
	.p2align	4, 0x90
_find_file:                             ## @find_file
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$48, %rsp
	movq	%rdi, -16(%rbp)
	movl	$1560, %eax                     ## imm = 0x618
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -32(%rbp)
LBB3_1:                                 ## =>This Inner Loop Header: Depth=1
	movq	-32(%rbp), %rax
	cmpq	$1279677254, (%rax)             ## imm = 0x4C465346
	je	LBB3_3
## %bb.2:
	jmp	LBB3_7
LBB3_3:                                 ##   in Loop: Header=BB3_1 Depth=1
	movq	-32(%rbp), %rax
	addq	$32, %rax
	movq	%rax, -32(%rbp)
	movq	-16(%rbp), %rdi
	movq	-32(%rbp), %rsi
	callq	_strcmp
	cmpl	$1, %eax
	jne	LBB3_5
## %bb.4:
	movq	-32(%rbp), %rax
	addq	$160, %rax
	movq	%rax, -32(%rbp)
	movq	-32(%rbp), %rax
	movq	%rax, -8(%rbp)
	jmp	LBB3_8
LBB3_5:                                 ##   in Loop: Header=BB3_1 Depth=1
	movq	-32(%rbp), %rax
	addq	$128, %rax
	movq	%rax, -32(%rbp)
	movq	-32(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -40(%rbp)
	movq	-32(%rbp), %rax
	addq	$32, %rax
	movq	%rax, -32(%rbp)
	movq	-32(%rbp), %rax
	movq	-40(%rbp), %rcx
	shlq	$3, %rcx
	addq	%rcx, %rax
	movq	%rax, -32(%rbp)
## %bb.6:                               ##   in Loop: Header=BB3_1 Depth=1
	jmp	LBB3_1
LBB3_7:
	movq	$0, -8(%rbp)
LBB3_8:
	movq	-8(%rbp), %rax
	addq	$48, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_fopen                          ## -- Begin function fopen
	.p2align	4, 0x90
_fopen:                                 ## @fopen
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$48, %rsp
	movw	%si, %ax
	movq	%rdi, -16(%rbp)
	movw	%ax, -18(%rbp)
	movq	-16(%rbp), %rdi
	callq	_find_file
	movq	%rax, -32(%rbp)
	cmpq	$0, -32(%rbp)
	jne	LBB4_4
## %bb.1:
	movswl	-18(%rbp), %eax
	cmpl	$1, %eax
	jne	LBB4_3
## %bb.2:
	leaq	L_.str(%rip), %rdi
	movl	$233, %esi
	xorl	%edx, %edx
	callq	_puts32
	movq	-16(%rbp), %rdi
	callq	_ffnt
	movq	%rax, %rdi
	movl	$233, %esi
	xorl	%edx, %edx
	callq	_puts32
	leaq	L_.str.1(%rip), %rdi
	movl	$233, %esi
	xorl	%edx, %edx
	callq	_puts32
	leaq	-48(%rbp), %rax
	movq	%rax, -8(%rbp)
	jmp	LBB4_5
LBB4_3:
	jmp	LBB4_4
LBB4_4:
	movq	-32(%rbp), %rax
	movq	%rax, -48(%rbp)
	movq	-16(%rbp), %rdi
	callq	_ffnt
	movq	%rax, -40(%rbp)
	leaq	-48(%rbp), %rax
	movq	%rax, -8(%rbp)
LBB4_5:
	movq	-8(%rbp), %rax
	addq	$48, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_fgetsize                       ## -- Begin function fgetsize
	.p2align	4, 0x90
_fgetsize:                              ## @fgetsize
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$32, %rsp
	movq	%rdi, -8(%rbp)
	movq	-8(%rbp), %rdi
	movl	$1, %esi
	callq	_fopen
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	addq	$-32, %rax
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	movq	(%rax), %rax
	addq	$32, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_flist                          ## -- Begin function flist
	.p2align	4, 0x90
_flist:                                 ## @flist
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$32, %rsp
	movl	$1560, %eax                     ## imm = 0x618
	movq	%rax, -8(%rbp)
	movq	-8(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -16(%rbp)
LBB6_1:                                 ## =>This Inner Loop Header: Depth=1
	movq	-16(%rbp), %rax
	cmpq	$1279677254, (%rax)             ## imm = 0x4C465346
	je	LBB6_3
## %bb.2:
	jmp	LBB6_4
LBB6_3:                                 ##   in Loop: Header=BB6_1 Depth=1
	movq	-16(%rbp), %rax
	addq	$32, %rax
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rdi
	callq	_ffnt
	movq	%rax, %rdi
	movl	$255, %esi
	xorl	%edx, %edx
	callq	_puts32
	leaq	L_.str.2(%rip), %rdi
	movl	$255, %esi
	xorl	%edx, %edx
	callq	_puts32
	movq	-16(%rbp), %rax
	addq	$128, %rax
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -24(%rbp)
	movq	-16(%rbp), %rax
	addq	$32, %rax
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rax
	movq	-24(%rbp), %rcx
	shlq	$3, %rcx
	addq	%rcx, %rax
	movq	%rax, -16(%rbp)
	jmp	LBB6_1
LBB6_4:
	addq	$32, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_fwrite                         ## -- Begin function fwrite
	.p2align	4, 0x90
_fwrite:                                ## @fwrite
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$48, %rsp
	movq	%rdi, -8(%rbp)
	movq	%rsi, -16(%rbp)
	movq	-8(%rbp), %rdi
	movl	$1, %esi
	callq	_fopen
	movq	%rax, -24(%rbp)
	movq	-24(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -32(%rbp)
	cmpq	$0, -32(%rbp)
	jne	LBB7_2
## %bb.1:
	jmp	LBB7_3
LBB7_2:
	movq	-16(%rbp), %rdi
	movq	-32(%rbp), %rsi
	callq	_strcpy
	movq	-32(%rbp), %rax
	movl	$512, %ecx                      ## imm = 0x200
	cqto
	idivq	%rcx
	movq	%rax, -40(%rbp)
	movq	-40(%rbp), %rdi
	callq	_save_sector
	movq	-40(%rbp), %rdi
	subq	$1, %rdi
	callq	_save_sector
	movq	-40(%rbp), %rdi
	addq	$1, %rdi
	callq	_save_sector
LBB7_3:
	addq	$48, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_fstrap                         ## -- Begin function fstrap
	.p2align	4, 0x90
_fstrap:                                ## @fstrap
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	subq	$48, %rsp
	movl	$1560, %eax                     ## imm = 0x618
	movq	%rax, -8(%rbp)
	movq	-8(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -16(%rbp)
LBB8_1:                                 ## =>This Inner Loop Header: Depth=1
	movq	-16(%rbp), %rax
	cmpq	$1279677254, (%rax)             ## imm = 0x4C465346
	je	LBB8_3
## %bb.2:
	jmp	LBB8_4
LBB8_3:                                 ##   in Loop: Header=BB8_1 Depth=1
	movq	-16(%rbp), %rax
	movl	$512, %ecx                      ## imm = 0x200
	cqto
	idivq	%rcx
	movq	%rax, -24(%rbp)
	movq	-16(%rbp), %rax
	addq	$160, %rax
	movq	%rax, -16(%rbp)
	movq	-16(%rbp), %rax
	movq	(%rax), %rax
	movq	%rax, -32(%rbp)
	movq	-32(%rbp), %rax
	movl	$512, %ecx                      ## imm = 0x200
	cqto
	idivq	%rcx
	movq	%rax, -40(%rbp)
	movq	-40(%rbp), %rdi
	movq	-24(%rbp), %rsi
	callq	_offset_sec_load
	movq	-8(%rbp), %rax
	movq	-32(%rbp), %rcx
	shlq	$3, %rcx
	addq	%rcx, %rax
	movq	%rax, -8(%rbp)
	jmp	LBB8_1
LBB8_4:
	addq	$48, %rsp
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.section	__TEXT,__cstring,cstring_literals
L_.str:                                 ## @.str
	.asciz	"File '"

L_.str.1:                               ## @.str.1
	.asciz	"' not found!\n"

L_.str.2:                               ## @.str.2
	.asciz	"\n"

.subsections_via_symbols
