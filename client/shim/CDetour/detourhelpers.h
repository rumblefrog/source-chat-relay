/**
 * vim: set ts=4 :
 * =============================================================================
 * SourceMod
 * Copyright (C) 2004-2007 AlliedModders LLC.  All rights reserved.
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
 * As a special exception, AlliedModders LLC gives you permission to link the
 * code of this program (as well as its derivative works) to "Half-Life 2," the
 * "Source Engine," the "SourcePawn JIT," and any Game MODs that run on software
 * by the Valve Corporation.  You must obey the GNU General Public License in
 * all respects for all other code used.  Additionally, AlliedModders LLC grants
 * this exception to all derivative works.  AlliedModders LLC defines further
 * exceptions, found in LICENSE.txt (as of this writing, version JULY-31-2007),
 * or <http://www.sourcemod.net/license.php>.
 *
 * Version: $Id$
 */

#ifndef _INCLUDE_SOURCEMOD_DETOURHELPERS_H_
#define _INCLUDE_SOURCEMOD_DETOURHELPERS_H_

struct patch_t
{
	patch_t()
	{
		patch[0] = 0;
		bytes = 0;
	}
	unsigned char patch[20];
	size_t bytes;
};

inline void SetMemPatchable(void *address, size_t size)
{
	SetMemAccess(address, size, SH_MEM_READ|SH_MEM_WRITE|SH_MEM_EXEC);
}

inline void SetMemExec(void *address, size_t size)
{
	SetMemAccess(address, size, SH_MEM_READ|SH_MEM_EXEC);
}

inline void DoGatePatch(unsigned char *target, void *callback)
{
	SetMemPatchable(target, 5);

	target[0] = IA32_JMP_IMM32;
	*(int32_t *)(&target[1]) = int32_t((unsigned char *)callback - (target + 5));
	
	SetMemExec(target, 5);
}

inline void ApplyPatch(void *address, int offset, const patch_t *patch, patch_t *restore)
{
	SetMemPatchable(address, sizeof(patch->patch));

	unsigned char *addr = (unsigned char *)address + offset;
	if (restore)
	{
		for (size_t i=0; i<patch->bytes; i++)
		{
			restore->patch[i] = addr[i];
		}
		restore->bytes = patch->bytes;
	}

	memcpy(addr, patch->patch, patch->bytes);

	SetMemExec(address, sizeof(patch->patch));
}

#endif //_INCLUDE_SOURCEMOD_DETOURHELPERS_H_
