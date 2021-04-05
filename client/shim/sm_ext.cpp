#include "sm_ext.h"
#include "bindings.h"

using namespace SourceMod;

ShimExtension g_SMExt;

bool ShimExtension::OnExtensionLoad(
    IExtension *me,
	IShareSys *sys,
	char *error,
	size_t maxlength,
	bool late)
{
    sharesys = sys;
    myself = me;

    /* Get the default interfaces from our configured SDK header */
    if (!SM_AcquireInterfaces(error, maxlength)) {
        return false;
    }

    sharesys->RegisterLibrary(myself, "SCR");

    return true;
}

void ShimExtension::OnExtensionUnload()
{
    SM_UnsetInterfaces();
}

void ShimExtension::OnExtensionsAllLoaded() {}
void ShimExtension::OnExtensionPauseChange() {}

bool ShimExtension::QueryRunning(char *error, size_t maxlength)
{
    return true;
}

bool ShimExtension::IsMetamodExtension()
{
    return false;
}

const char *ShimExtension::GetExtensionName()
{
    return extension_name();
}

const char *ShimExtension::GetExtensionURL()
{
    return extension_url();
}

const char *ShimExtension::GetExtensionTag()
{
    return extension_log_tag();
}

const char *ShimExtension::GetExtensionAuthor()
{
    return extension_author();
}

const char *ShimExtension::GetExtensionVerString()
{
    return extension_version();
}

const char *ShimExtension::GetExtensionDescription()
{
    return extension_description();
}

const char *ShimExtension::GetExtensionDateString()
{
    return extension_date();
}

bool SM_LoadExtension(char *error, size_t maxlength) {
	if ((smexts = (IExtensionManager *)
			g_SMAPI->MetaFactory(SOURCEMOD_INTERFACE_EXTENSIONS, NULL, NULL)) == NULL) {
		if (error && maxlength) {
			snprintf(error, maxlength, SOURCEMOD_INTERFACE_EXTENSIONS " interface not found");
		}
		return false;
	}

	/* This could be more dynamic */
	char path[256];
	g_SMAPI->PathFormat(path, sizeof(path),  "addons/scr/bin/SCR%s",
#if defined __linux__
		"_i486.so"
#else
		".dll"
#endif
	);

	if ((myself = smexts->LoadExternal(&g_RCBotSourceMod, path, "scr.ext", error, maxlength))
			== NULL) {
		SM_UnsetInterfaces();
		return false;
	}
	return true;
}

void SM_UnloadExtension() {
	smexts->UnloadExtension(myself);
}
