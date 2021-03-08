/**
 * vim: set ts=4 :
 * =============================================================================
 * Source Dedicated Server Wrapper for Mac OS X
 * Copyright (C) 2011 Scott "DS" Ehlert.  All rights reserved.
 * =============================================================================
 *
 * This program is free software; you can redistribute it and/or modify it under
 * the terms of the GNU General Public License, version 3.0, as published by the
 * Free Software Foundation.
 * 
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#ifndef _INCLUDE_SRCDS_OSX_SH_INCLUDE_H_
#define _INCLUDE_SRCDS_OSX_SH_INCLUDE_H_

class ISourceHook;
class IHookManagerAutoGen {};
class HookManagerPubFunc {};
class CProto {};
struct IntPassInfo;
struct ProtoInfo;
struct PassInfo
{
    struct V2Info;
};

/* Define a bunch of things for SourceHook headers */
#include "metamod_oslink.h"
#include "sourcehook.h"
#include "sourcehook_hookmangen.h"
#include "sourcehook_hookmangen_x86.h"

using namespace SourceHook;
using namespace SourceHook::Impl;

/* Extra code generation helpers not present in SourceHook headers */
inline void IA32_Mov_ESP_Disp8_Reg(JitWriter *jit, jit_int8_t disp8, jit_uint8_t src)
{
	jit->write_ubyte(IA32_MOV_RM_REG);
	jit->write_ubyte(ia32_modrm(MOD_DISP8, src, REG_ESP));
	jit->write_ubyte(ia32_sib(NOSCALE, REG_NOIDX, REG_ESP));
	jit->write_byte(disp8);
}

inline void IA32_Write_Jump32_Abs(JitWriter *jit, jitoffs_t jmp, void *target)
{
	jit_int32_t disp = (int32_t)target - ((int32_t)jit->GetData() + jmp + 4);
	jit->rewrite(jmp, disp);
}

inline void IA32_Write_Jump32(JitWriter *jit, jitoffs_t jmp, jitoffs_t target)
{
	jit_int32_t disp = (target - (jmp + 4));
	jit->rewrite(jmp, disp);
}

inline void IA32_Send_Jump32_Here(JitWriter *jit, jitoffs_t jmp)
{
	jitoffs_t curptr = jit->get_outputpos();
	IA32_Write_Jump32(jit, jmp, curptr);
}

#endif // _INCLUDE_SRCDS_OSX_SH_INCLUDE_H_
