#include <stdio.h>
#include "shim.h"
#include "bindings.h"

Shim g_shim;

PLUGIN_EXPOSE(Shim, g_shim);

bool Shim::Load(PluginId id, ISmmAPI *ismm, char *error, size_t maxlen, bool late)
{
	PLUGIN_SAVEVARS()

	return true;
}

bool Shim::Unload(char *error, size_t maxlen)
{
	return true;
}

const char *Shim::GetLicense()
{
	return extension_license();
}

const char *Shim::GetName()
{
	return extension_name();
}

const char *Shim::GetAuthor()
{
	return extension_author();
}

const char *Shim::GetDescription()
{
	return extension_description();
}

const char *Shim::GetURL()
{
	return extension_url();
}

const char *Shim::GetVersion()
{
	return extension_version();
}

const char *Shim::GetLogTag()
{
	return extension_log_tag();
}

const char *Shim::GetDate()
{
	return extension_date();
}
