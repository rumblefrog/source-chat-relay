#include <stdio.h>
#include "shim.h"

#include "MemoryUtils.h"
#include <CDetour/detours.h>

using namespace SourceHook;

Shim g_shim;

PLUGIN_EXPOSE(Shim, g_shim);

Shim::Shim()
{
    m_Client = NULL;
}

bool Shim::Load(PluginId id, ISmmAPI *ismm, char *error, size_t maxlen, bool late)
{
	PLUGIN_SAVEVARS()

    void *engineFactory = (void *)g_SMAPI->GetEngineFactory(false);
	int engineVersion = g_SMAPI->GetSourceEngineBuild();
	void *adrVoiceData = NULL;
	int offsPlayerSlot = 0;

    switch (engineVersion)
	{
		case SOURCE_ENGINE_CSGO:
#ifdef _WIN32
			adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x55\x8B\xEC\x81\xEC\xD0\x00\x00\x00\x53\x56\x57", 12);
			offsPlayerSlot = 15;
#else
			adrVoiceData = g_MemUtils.ResolveSymbol(engineFactory, "_Z21SV_BroadcastVoiceDataP7IClientiPcx");
			offsPlayerSlot = 16;
#endif
			break;

		case SOURCE_ENGINE_LEFT4DEAD2:
#ifdef _WIN32
			adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x55\x8B\xEC\x83\xEC\x70\xA1\x2A\x2A\x2A\x2A\x33\xC5\x89\x45\xFC\xA1\x2A\x2A\x2A\x2A\x53\x56", 23);
			offsPlayerSlot = 14;
#else
			adrVoiceData = g_MemUtils.ResolveSymbol(engineFactory, "_Z21SV_BroadcastVoiceDataP7IClientiPcx");
			offsPlayerSlot = 15;
#endif
			break;

		case SOURCE_ENGINE_NUCLEARDAWN:
#ifdef _WIN32
			adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x55\x8B\xEC\xA1\x2A\x2A\x2A\x2A\x83\xEC\x58\x57\x33\xFF", 14);
			offsPlayerSlot = 14;
#else
			adrVoiceData = g_MemUtils.ResolveSymbol(engineFactory, "_Z21SV_BroadcastVoiceDataP7IClientiPcx");
			offsPlayerSlot = 15;
#endif
			break;

		case SOURCE_ENGINE_INSURGENCY:
#ifdef _WIN32
			adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x55\x8B\xEC\x83\xEC\x74\x68\x2A\x2A\x2A\x2A\x8D\x4D\xE4\xE8", 15);
			offsPlayerSlot = 14;
#else
			adrVoiceData = g_MemUtils.ResolveSymbol(engineFactory, "_Z21SV_BroadcastVoiceDataP7IClientiPcx");
			offsPlayerSlot = 15;
#endif
			break;

		case SOURCE_ENGINE_TF2:
		case SOURCE_ENGINE_CSS:
		case SOURCE_ENGINE_HL2DM:
		case SOURCE_ENGINE_DODS:
		case SOURCE_ENGINE_SDK2013:
#ifdef _WIN32
			adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x55\x8B\xEC\xA1\x2A\x2A\x2A\x2A\x83\xEC\x50\x83\x78\x30", 14);
			offsPlayerSlot = 14;
#else
			adrVoiceData = g_MemUtils.ResolveSymbol(engineFactory, "_Z21SV_BroadcastVoiceDataP7IClientiPcx");
			offsPlayerSlot = 15;
#endif
			break;

		default:
			g_SMAPI->Format(error, maxlen, "Unsupported game.");
			return false;
	}

	if (!adrVoiceData)
	{
		g_SMAPI->Format(error, maxlen, "SV_BroadcastVoiceData sigscan failed.");
		return false;
	}

    m_Client = new_client();

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
