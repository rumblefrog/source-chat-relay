#ifndef _INCLUDE_SHIM_SOURCEMOD_EXTENSION_
#define _INCLUDE_SHIM_SOURCEMOD_EXTENSION_

#include "sm_sdk_config.h"
#include "shim.h"

using namespace SourceMod;

class ShimExtension : public IExtensionInterface
{
public:
	virtual bool OnExtensionLoad(IExtension *me,
		IShareSys *sys,
		char *error,
		size_t maxlength,
		bool late);
	virtual void OnExtensionUnload();
	virtual void OnExtensionsAllLoaded();
	virtual void OnExtensionPauseChange(bool pause);
	virtual bool QueryRunning(char *error, size_t maxlength);
	virtual bool IsMetamodExtension();
	virtual const char *GetExtensionName();
	virtual const char *GetExtensionURL();
	virtual const char *GetExtensionTag();
	virtual const char *GetExtensionAuthor();
	virtual const char *GetExtensionVerString();
	virtual const char *GetExtensionDescription();
	virtual const char *GetExtensionDateString();
};

bool SM_LoadExtension(char *error, size_t maxlength);
void SM_UnloadExtension();

extern IShareSys *sharesys;
extern IExtension *myself;
extern ShimExtension g_SMExt;

#endif //_INCLUDE_SHIM_SOURCEMOD_EXTENSION_
