// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/x86"
)

var (
	AX               = RtoB(x86.REG_AX)
	BX               = RtoB(x86.REG_BX)
	CX               = RtoB(x86.REG_CX)
	DX               = RtoB(x86.REG_DX)
	DI               = RtoB(x86.REG_DI)
	SI               = RtoB(x86.REG_SI)
	LeftRdwr  uint32 = gc.LeftRead | gc.LeftWrite
	RightRdwr uint32 = gc.RightRead | gc.RightWrite
)

// This table gives the basic information about instruction
// generated by the compiler and processed in the optimizer.
// See opt.h for bit definitions.
//
// Instructions not generated need not be listed.
// As an exception to that rule, we typically write down all the
// size variants of an operation even if we just use a subset.
//
// The table is formatted for 8-space tabs.
var progtable = [x86.ALAST]obj.ProgInfo{
	obj.ATYPE:     {Flags: gc.Pseudo | gc.Skip},
	obj.ATEXT:     {Flags: gc.Pseudo},
	obj.AFUNCDATA: {Flags: gc.Pseudo},
	obj.APCDATA:   {Flags: gc.Pseudo},
	obj.AUNDEF:    {Flags: gc.Break},
	obj.AUSEFIELD: {Flags: gc.OK},
	obj.ACHECKNIL: {Flags: gc.LeftRead},
	obj.AVARDEF:   {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARKILL:  {Flags: gc.Pseudo | gc.RightWrite},

	// NOP is an internal no-op that also stands
	// for USED and SET annotations, not the Intel opcode.
	obj.ANOP:       {Flags: gc.LeftRead | gc.RightWrite},
	x86.AADCL:      {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry | gc.UseCarry},
	x86.AADCW:      {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry | gc.UseCarry},
	x86.AADDB:      {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AADDL:      {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AADDW:      {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AADDSD:     {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	x86.AADDSS:     {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	x86.AANDB:      {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AANDL:      {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AANDW:      {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry},
	obj.ACALL:      {Flags: gc.RightAddr | gc.Call | gc.KillCarry},
	x86.ACDQ:       {Flags: gc.OK, Reguse: AX, Regset: AX | DX},
	x86.ACWD:       {Flags: gc.OK, Reguse: AX, Regset: AX | DX},
	x86.ACLD:       {Flags: gc.OK},
	x86.ASTD:       {Flags: gc.OK},
	x86.ACMPB:      {Flags: gc.SizeB | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ACMPL:      {Flags: gc.SizeL | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ACMPW:      {Flags: gc.SizeW | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ACOMISD:    {Flags: gc.SizeD | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ACOMISS:    {Flags: gc.SizeF | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ACVTSD2SL:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTSD2SS:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTSL2SD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTSL2SS:  {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTSS2SD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTSS2SL:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTTSD2SL: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ACVTTSS2SL: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.ADECB:      {Flags: gc.SizeB | RightRdwr},
	x86.ADECL:      {Flags: gc.SizeL | RightRdwr},
	x86.ADECW:      {Flags: gc.SizeW | RightRdwr},
	x86.ADIVB:      {Flags: gc.SizeB | gc.LeftRead | gc.SetCarry, Reguse: AX, Regset: AX},
	x86.ADIVL:      {Flags: gc.SizeL | gc.LeftRead | gc.SetCarry, Reguse: AX | DX, Regset: AX | DX},
	x86.ADIVW:      {Flags: gc.SizeW | gc.LeftRead | gc.SetCarry, Reguse: AX | DX, Regset: AX | DX},
	x86.ADIVSD:     {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	x86.ADIVSS:     {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	x86.AFLDCW:     {Flags: gc.SizeW | gc.LeftAddr},
	x86.AFSTCW:     {Flags: gc.SizeW | gc.RightAddr},
	x86.AFSTSW:     {Flags: gc.SizeW | gc.RightAddr | gc.RightWrite},
	x86.AFADDD:     {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFADDDP:    {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFADDF:     {Flags: gc.SizeF | gc.LeftAddr | RightRdwr},
	x86.AFCOMD:     {Flags: gc.SizeD | gc.LeftAddr | gc.RightRead},
	x86.AFCOMDP:    {Flags: gc.SizeD | gc.LeftAddr | gc.RightRead},
	x86.AFCOMDPP:   {Flags: gc.SizeD | gc.LeftAddr | gc.RightRead},
	x86.AFCOMF:     {Flags: gc.SizeF | gc.LeftAddr | gc.RightRead},
	x86.AFCOMFP:    {Flags: gc.SizeF | gc.LeftAddr | gc.RightRead},
	x86.AFUCOMIP:   {Flags: gc.SizeF | gc.LeftAddr | gc.RightRead},
	x86.AFCHS:      {Flags: gc.SizeD | RightRdwr}, // also SizeF

	x86.AFDIVDP:  {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFDIVF:   {Flags: gc.SizeF | gc.LeftAddr | RightRdwr},
	x86.AFDIVD:   {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFDIVRDP: {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFDIVRF:  {Flags: gc.SizeF | gc.LeftAddr | RightRdwr},
	x86.AFDIVRD:  {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFXCHD:   {Flags: gc.SizeD | LeftRdwr | RightRdwr},
	x86.AFSUBD:   {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFSUBDP:  {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFSUBF:   {Flags: gc.SizeF | gc.LeftAddr | RightRdwr},
	x86.AFSUBRD:  {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFSUBRDP: {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFSUBRF:  {Flags: gc.SizeF | gc.LeftAddr | RightRdwr},
	x86.AFMOVD:   {Flags: gc.SizeD | gc.LeftAddr | gc.RightWrite},
	x86.AFMOVF:   {Flags: gc.SizeF | gc.LeftAddr | gc.RightWrite},
	x86.AFMOVL:   {Flags: gc.SizeL | gc.LeftAddr | gc.RightWrite},
	x86.AFMOVW:   {Flags: gc.SizeW | gc.LeftAddr | gc.RightWrite},
	x86.AFMOVV:   {Flags: gc.SizeQ | gc.LeftAddr | gc.RightWrite},

	// These instructions are marked as RightAddr
	// so that the register optimizer does not try to replace the
	// memory references with integer register references.
	// But they do not use the previous value at the address, so
	// we also mark them RightWrite.
	x86.AFMOVDP:   {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.RightAddr},
	x86.AFMOVFP:   {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.RightAddr},
	x86.AFMOVLP:   {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.RightAddr},
	x86.AFMOVWP:   {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.RightAddr},
	x86.AFMOVVP:   {Flags: gc.SizeQ | gc.LeftRead | gc.RightWrite | gc.RightAddr},
	x86.AFMULD:    {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFMULDP:   {Flags: gc.SizeD | gc.LeftAddr | RightRdwr},
	x86.AFMULF:    {Flags: gc.SizeF | gc.LeftAddr | RightRdwr},
	x86.AIDIVB:    {Flags: gc.SizeB | gc.LeftRead | gc.SetCarry, Reguse: AX, Regset: AX},
	x86.AIDIVL:    {Flags: gc.SizeL | gc.LeftRead | gc.SetCarry, Reguse: AX | DX, Regset: AX | DX},
	x86.AIDIVW:    {Flags: gc.SizeW | gc.LeftRead | gc.SetCarry, Reguse: AX | DX, Regset: AX | DX},
	x86.AIMULB:    {Flags: gc.SizeB | gc.LeftRead | gc.SetCarry, Reguse: AX, Regset: AX},
	x86.AIMULL:    {Flags: gc.SizeL | gc.LeftRead | gc.ImulAXDX | gc.SetCarry},
	x86.AIMULW:    {Flags: gc.SizeW | gc.LeftRead | gc.ImulAXDX | gc.SetCarry},
	x86.AINCB:     {Flags: gc.SizeB | RightRdwr},
	x86.AINCL:     {Flags: gc.SizeL | RightRdwr},
	x86.AINCW:     {Flags: gc.SizeW | RightRdwr},
	x86.AJCC:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJCS:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJEQ:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJGE:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJGT:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJHI:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJLE:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJLS:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJLT:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJMI:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJNE:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJOC:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJOS:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJPC:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJPL:      {Flags: gc.Cjmp | gc.UseCarry},
	x86.AJPS:      {Flags: gc.Cjmp | gc.UseCarry},
	obj.AJMP:      {Flags: gc.Jump | gc.Break | gc.KillCarry},
	x86.ALEAW:     {Flags: gc.LeftAddr | gc.RightWrite},
	x86.ALEAL:     {Flags: gc.LeftAddr | gc.RightWrite},
	x86.AMOVBLSX:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.AMOVBLZX:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.AMOVBWSX:  {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.AMOVBWZX:  {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.AMOVWLSX:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.AMOVWLZX:  {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	x86.AMOVB:     {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move},
	x86.AMOVL:     {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move},
	x86.AMOVW:     {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move},
	x86.AMOVSB:    {Flags: gc.OK, Reguse: DI | SI, Regset: DI | SI},
	x86.AMOVSL:    {Flags: gc.OK, Reguse: DI | SI, Regset: DI | SI},
	x86.AMOVSW:    {Flags: gc.OK, Reguse: DI | SI, Regset: DI | SI},
	obj.ADUFFCOPY: {Flags: gc.OK, Reguse: DI | SI, Regset: DI | SI | CX},
	x86.AMOVSD:    {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Move},
	x86.AMOVSS:    {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Move},

	// We use MOVAPD as a faster synonym for MOVSD.
	x86.AMOVAPD:   {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Move},
	x86.AMULB:     {Flags: gc.SizeB | gc.LeftRead | gc.SetCarry, Reguse: AX, Regset: AX},
	x86.AMULL:     {Flags: gc.SizeL | gc.LeftRead | gc.SetCarry, Reguse: AX, Regset: AX | DX},
	x86.AMULW:     {Flags: gc.SizeW | gc.LeftRead | gc.SetCarry, Reguse: AX, Regset: AX | DX},
	x86.AMULSD:    {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	x86.AMULSS:    {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	x86.ANEGB:     {Flags: gc.SizeB | RightRdwr | gc.SetCarry},
	x86.ANEGL:     {Flags: gc.SizeL | RightRdwr | gc.SetCarry},
	x86.ANEGW:     {Flags: gc.SizeW | RightRdwr | gc.SetCarry},
	x86.ANOTB:     {Flags: gc.SizeB | RightRdwr},
	x86.ANOTL:     {Flags: gc.SizeL | RightRdwr},
	x86.ANOTW:     {Flags: gc.SizeW | RightRdwr},
	x86.AORB:      {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AORL:      {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AORW:      {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.APOPL:     {Flags: gc.SizeL | gc.RightWrite},
	x86.APUSHL:    {Flags: gc.SizeL | gc.LeftRead},
	x86.APXOR:     {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	x86.ARCLB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry | gc.UseCarry},
	x86.ARCLL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry | gc.UseCarry},
	x86.ARCLW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry | gc.UseCarry},
	x86.ARCRB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry | gc.UseCarry},
	x86.ARCRL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry | gc.UseCarry},
	x86.ARCRW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry | gc.UseCarry},
	x86.AREP:      {Flags: gc.OK, Reguse: CX, Regset: CX},
	x86.AREPN:     {Flags: gc.OK, Reguse: CX, Regset: CX},
	obj.ARET:      {Flags: gc.Break | gc.KillCarry},
	x86.AROLB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.AROLL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.AROLW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ARORB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ARORL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ARORW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASAHF:     {Flags: gc.OK, Reguse: AX, Regset: AX},
	x86.ASALB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASALL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASALW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASARB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASARL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASARW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASBBB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.SetCarry | gc.UseCarry},
	x86.ASBBL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry | gc.UseCarry},
	x86.ASBBW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry | gc.UseCarry},
	x86.ASETCC:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETCS:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETEQ:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETGE:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETGT:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETHI:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETLE:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETLS:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETLT:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETMI:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETNE:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETOC:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETOS:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETPC:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETPL:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASETPS:    {Flags: gc.SizeB | RightRdwr | gc.UseCarry},
	x86.ASHLB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASHLL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASHLW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASHRB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASHRL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASHRW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.ShiftCX | gc.SetCarry},
	x86.ASTOSB:    {Flags: gc.OK, Reguse: AX | DI, Regset: DI},
	x86.ASTOSL:    {Flags: gc.OK, Reguse: AX | DI, Regset: DI},
	x86.ASTOSW:    {Flags: gc.OK, Reguse: AX | DI, Regset: DI},
	obj.ADUFFZERO: {Flags: gc.OK, Reguse: AX | DI, Regset: DI},
	x86.ASUBB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.ASUBL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.ASUBW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.ASUBSD:    {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	x86.ASUBSS:    {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	x86.ATESTB:    {Flags: gc.SizeB | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ATESTL:    {Flags: gc.SizeL | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.ATESTW:    {Flags: gc.SizeW | gc.LeftRead | gc.RightRead | gc.SetCarry},
	x86.AUCOMISD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightRead},
	x86.AUCOMISS:  {Flags: gc.SizeF | gc.LeftRead | gc.RightRead},
	x86.AXCHGB:    {Flags: gc.SizeB | LeftRdwr | RightRdwr},
	x86.AXCHGL:    {Flags: gc.SizeL | LeftRdwr | RightRdwr},
	x86.AXCHGW:    {Flags: gc.SizeW | LeftRdwr | RightRdwr},
	x86.AXORB:     {Flags: gc.SizeB | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AXORL:     {Flags: gc.SizeL | gc.LeftRead | RightRdwr | gc.SetCarry},
	x86.AXORW:     {Flags: gc.SizeW | gc.LeftRead | RightRdwr | gc.SetCarry},
}

func proginfo(p *obj.Prog) {
	info := &p.Info
	*info = progtable[p.As]
	if info.Flags == 0 {
		gc.Fatal("unknown instruction %v", p)
	}

	if (info.Flags&gc.ShiftCX != 0) && p.From.Type != obj.TYPE_CONST {
		info.Reguse |= CX
	}

	if info.Flags&gc.ImulAXDX != 0 {
		if p.To.Type == obj.TYPE_NONE {
			info.Reguse |= AX
			info.Regset |= AX | DX
		} else {
			info.Flags |= RightRdwr
		}
	}

	// Addressing makes some registers used.
	if p.From.Type == obj.TYPE_MEM && p.From.Name == obj.NAME_NONE {
		info.Regindex |= RtoB(int(p.From.Reg))
	}
	if p.From.Index != x86.REG_NONE {
		info.Regindex |= RtoB(int(p.From.Index))
	}
	if p.To.Type == obj.TYPE_MEM && p.To.Name == obj.NAME_NONE {
		info.Regindex |= RtoB(int(p.To.Reg))
	}
	if p.To.Index != x86.REG_NONE {
		info.Regindex |= RtoB(int(p.To.Index))
	}
}
