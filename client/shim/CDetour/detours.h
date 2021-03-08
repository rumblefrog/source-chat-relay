/**
* vim: set ts=4 :
* =============================================================================
* SourceMod
* Copyright (C) 2004-2008 AlliedModders LLC.  All rights reserved.
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
*/

#ifndef _INCLUDE_SOURCEMOD_DETOURS_H_
#define _INCLUDE_SOURCEMOD_DETOURS_H_

#include <sh_include.h>
#include "detourhelpers.h"

/**
 * CDetours class for SourceMod Extensions by pRED*
 * Modified by DS to make use of SourceHook's GenBuffer and code generation.
 *
 * detourhelpers.h entirely stolen from CSS:DM and were written by BAILOPAN (I assume).
 * asm.h/c from devmaster.net (thanks cybermind) edited by pRED* to handle gcc -fPIC thunks correctly
 * Concept by Nephyrin Zey (http://www.doublezen.net/) and Windows Detour Library (http://research.microsoft.com/sn/detours/)
 * Member function pointer ideas by Don Clugston (http://www.codeproject.com/cpp/FastDelegate.asp)
 */

#define DETOUR_MEMBER_CALL(name) (this->*name##_Actual)
#define DETOUR_STATIC_CALL(name) (name##_Actual)

#define DETOUR_DECL_STATIC0(name, ret) \
ret (*name##_Actual)(void) = NULL; \
ret name(void)

#define DETOUR_DECL_STATIC1(name, ret, p1type, p1name) \
ret (*name##_Actual)(p1type) = NULL; \
ret name(p1type p1name)

#define DETOUR_DECL_STATIC2(name, ret, p1type, p1name, p2type, p2name) \
ret (*name##_Actual)(p1type, p2type) = NULL; \
ret name(p1type p1name, p2type p2name)

#define DETOUR_DECL_STATIC3(name, ret, p1type, p1name, p2type, p2name, p3type, p3name) \
ret (*name##_Actual)(p1type, p2type, p3type) = NULL; \
ret name(p1type p1name, p2type p2name, p3type p3name)

#define DETOUR_DECL_STATIC4(name, ret, p1type, p1name, p2type, p2name, p3type, p3name, p4type, p4name) \
ret (*name##_Actual)(p1type, p2type, p3type, p4type) = NULL; \
ret name(p1type p1name, p2type p2name, p3type p3name, p4type p4name)

#define DETOUR_DECL_MEMBER0(name, ret) \
class name##Class \
{ \
public: \
	ret name(); \
	static ret (name##Class::* name##_Actual)(void); \
}; \
ret (name##Class::* name##Class::name##_Actual)(void) = NULL; \
ret name##Class::name()

#define DETOUR_DECL_MEMBER1(name, ret, p1type, p1name) \
class name##Class \
{ \
public: \
	ret name(p1type p1name); \
	static ret (name##Class::* name##_Actual)(p1type); \
}; \
ret (name##Class::* name##Class::name##_Actual)(p1type) = NULL; \
ret name##Class::name(p1type p1name)

#define DETOUR_DECL_MEMBER2(name, ret, p1type, p1name, p2type, p2name) \
class name##Class \
{ \
public: \
	ret name(p1type p1name, p2type p2name); \
	static ret (name##Class::* name##_Actual)(p1type, p2type); \
}; \
ret (name##Class::* name##Class::name##_Actual)(p1type, p2type) = NULL; \
ret name##Class::name(p1type p1name, p2type p2name)

#define DETOUR_DECL_MEMBER3(name, ret, p1type, p1name, p2type, p2name, p3type, p3name) \
class name##Class \
{ \
public: \
	ret name(p1type p1name, p2type p2name, p3type p3name); \
	static ret (name##Class::* name##_Actual)(p1type, p2type, p3type); \
}; \
ret (name##Class::* name##Class::name##_Actual)(p1type, p2type, p3type) = NULL; \
ret name##Class::name(p1type p1name, p2type p2name, p3type p3name)

#define GET_MEMBER_CALLBACK(name) (void *)GetCodeAddress(&name##Class::name)
#define GET_MEMBER_TRAMPOLINE(name) (void **)(&name##Class::name##_Actual)

#define GET_STATIC_CALLBACK(name) (void *)&name
#define GET_STATIC_TRAMPOLINE(name) (void **)&name##_Actual

#define DETOUR_CREATE_MEMBER(name, addr) CDetourManager::CreateDetour(GET_MEMBER_CALLBACK(name), GET_MEMBER_TRAMPOLINE(name), addr);
#define DETOUR_CREATE_STATIC(name, addr) CDetourManager::CreateDetour(GET_STATIC_CALLBACK(name), GET_STATIC_TRAMPOLINE(name), addr);

class GenericClass {};
typedef void (GenericClass::*VoidFunc)();

inline void *GetCodeAddr(VoidFunc mfp)
{
	return *(void **)&mfp;
}

/**
 * Converts a member function pointer to a void pointer.
 * This relies on the assumption that the code address lies at mfp+0
 * This is the case for both g++ and later MSVC versions on non virtual functions but may be different for other compilers
 * Based on research by Don Clugston : http://www.codeproject.com/cpp/FastDelegate.asp
 */
#define GetCodeAddress(mfp) GetCodeAddr(reinterpret_cast<VoidFunc>(mfp))

class CDetourManager;

class CDetour
{
public:

	bool IsEnabled();

	/**
	 * These would be somewhat self-explanatory I hope
	 */
	void EnableDetour();
	void DisableDetour();

	void *GetTargetAddr();
	void SetTargetAddr(void *addr);

	void Destroy();

	friend class CDetourManager;

protected:
	CDetour(void *callbackfunction, void **trampoline);

	bool Init(void *addr);
private:

	/* These create/delete the allocated memory */
	bool CreateDetour();
	void DeleteDetour();

	bool enabled;
	bool detoured;

	patch_t detour_restore;
	/* Address of the detoured function */
	void *detour_address;
	/* Address of the allocated trampoline function */
	void *detour_trampoline;
	/* Address of the callback handler */
	void *detour_callback;
	/* The function pointer used to call our trampoline */
	void **trampoline;

	GenBuffer codegen;
};

class CDetourManager
{
public:

	/**
	 * Creates a new detour.
	 *
	 * @param callbackfunction			Void pointer to your detour callback function.
	 * @param trampoline				Address of the trampoline pointer
	 * @param signame					Section name containing a signature to fetch from the gamedata file.
	 * @return							A new CDetour pointer to control your detour.
	 *
	 * Example:
	 *
	 * CBaseServer::ConnectClient(netadr_s &, int, int, int, char  const*, char  const*, char  const*, int)
	 *
	 * Define a new class with the required function and a member function pointer to the same type:
	 *
	 * class CBaseServerDetour
	 * {
	 * public:
	 *		 bool ConnectClient(void *netaddr_s, int, int, int, char  const*, char  const*, char  const*, int);
	 *		 static bool (CBaseServerDetour::* ConnectClient_Actual)(void *netaddr_s, int, int, int, char  const*, char  const*, char  const*, int);
	 * }
	 *
	 *	void *callbackfunc = GetCodeAddress(&CBaseServerDetour::ConnectClient);
	 *	void **trampoline = (void **)(&CBaseServerDetour::ConnectClient_Actual);
	 *
	 * Creation:
	 * CDetourManager::CreateDetour(callbackfunc,  trampoline, "ConnectClient");
	 *
	 * Usage:
	 *
	 * CBaseServerDetour::ConnectClient(void *netaddr_s, int, int, int, char  const*, char  const*, char  const*, int)
	 * {
	 *			//pre hook code
	 *			bool result = (this->*ConnectClient_Actual)(netaddr_s, rest of params);
	 *			//post hook code
	 *			return result;
	 * }
	 *
	 * Note we changed the netadr_s reference into a void* to avoid needing to define the type
	 */
	//static CDetour *CreateDetour(void *callbackfunction, void **trampoline, const char *signame);

	/**
	 * Creates a new detour.
	 *
	 * @param callbackfunction			Void pointer to your detour callback function.
	 * @param trampoline				Address of the trampoline pointer.
	 * @param addr						Address of function to detour.
	 * @return							A new CDetour pointer to control your detour.
	 */
	static CDetour *CreateDetour(void *callbackfunction, void **trampoline, void *addr);

	friend class CDetour;
};

#endif // _INCLUDE_SOURCEMOD_DETOURS_H_
