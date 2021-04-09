#include <stdio.h>
#include "shim.h"

#include "MemoryUtils.h"

using namespace SourceHook;

SH_DECL_HOOK1_void(IServerGameClients, ClientDisconnect, SH_NOATTRIB, 0, edict_t *);
SH_DECL_HOOK2_void(IServerGameClients, ClientPutInServer, SH_NOATTRIB, 0, edict_t *, char const *);
SH_DECL_HOOK2_void(IServerGameClients, ClientCommand, SH_NOATTRIB, 0, edict_t *, const CCommand &);

Shim g_shim;
IVEngineServer *engine = NULL;
IServerGameClients *gameclients = NULL;

PLUGIN_EXPOSE(Shim, g_shim);

class IClient;
class CCLCMsg_VoiceData;

DETOUR_DECL_STATIC4(BroadcastVoiceData, void, IClient *, client, int, bytes, char *, data, long long, xuid)
{
    META_LOG(g_PLAPI, ">>> SV_BroadcastVoiceData(%p, %d, %p, %lld)", client, bytes, data, xuid);

	DETOUR_STATIC_CALL(BroadcastVoiceData)(client, bytes, data, xuid);

    g_shim.BroadcastVoiceData_Callback(bytes, data);
}

#ifdef _WIN32
// This function has been LTCG'd to __fastcall.
DETOUR_DECL_STATIC0(BroadcastVoiceData_Protobuf, void)
{
	IClient *client;
	CCLCMsg_VoiceData *message;
	__asm
	{
		mov client, ecx
		mov message, edx
	}

	// Call the original func before logging to try and avoid the registers being overwritten.
	DETOUR_STATIC_CALL(BroadcastVoiceData_Protobuf)();
#else
DETOUR_DECL_STATIC2(BroadcastVoiceData_Protobuf, void, IClient *, client, CCLCMsg_VoiceData *, message)
{
	DETOUR_STATIC_CALL(BroadcastVoiceData_Protobuf)(client, message);
#endif

	META_LOG(g_PLAPI, ">>> SV_BroadcastVoiceData(%p, %p)", client, message);

	// If this breaks on Linux only, check the libstdc++ ABI in use, see comment in AMBuilder.
	std::string *voiceData = *(std::string **)((uintptr_t)message + 8);

    g_shim.BroadcastVoiceData_Callback(voiceData->size(), voiceData->data());
}

bool Shim::Load(PluginId id, ISmmAPI *ismm, char *error, size_t maxlen, bool late)
{
	PLUGIN_SAVEVARS()

    GET_V_IFACE_CURRENT(GetEngineFactory, engine, IVEngineServer, INTERFACEVERSION_VENGINESERVER);
    GET_V_IFACE_ANY(GetServerFactory, gameclients, IServerGameClients, INTERFACEVERSION_SERVERGAMECLIENTS);

    SH_ADD_HOOK_MEMFUNC(IServerGameClients, ClientDisconnect, gameclients, this, &Shim::ClientDisconnect, true);
	SH_ADD_HOOK_MEMFUNC(IServerGameClients, ClientPutInServer, gameclients, this, &Shim::ClientPutInServer, true);
    SH_ADD_HOOK_MEMFUNC(IServerGameClients, ClientCommand, gameclients, this, &Shim::ClientCommand, false);

    void *engineFactory = (void *)g_SMAPI->GetEngineFactory(false);

	int engineVersion = g_SMAPI->GetSourceEngineBuild();

	void *adrVoiceData = NULL;

    switch (engineVersion)
    {
        case SOURCE_ENGINE_CSGO:
            #ifdef _WIN32
            adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x2A\x2A\x2A\x2A\x2A\x2A\x2A\x2A\xE4\x00\x00\x00\x53\x56\x57\x8B\xD9\x8B\xF2", 19);
            #else
            adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x2A\x2A\x2A\x2A\x2A\x2A\x2A\x2A\x53\x81\xEC\xEC\x00\x00\x00\x89\x04\x24\x8B\x5D\x0C\xC7\x44\x24\x04", 25);
            #endif

            break;

        default:
            #ifdef _WIN32
            adrVoiceData = g_MemUtils.FindPattern(engineFactory, "\x2A\x2A\x2A\x2A\x2A\x2A\x2A\x2A\x83\xEC\x50\x83\x78\x30\x00\x0F\x84\x2A\x2A\x2A\x2A\x53", 22);
            #else
            adrVoiceData = g_MemUtils.ResolveSymbol(engineFactory, "_Z21SV_BroadcastVoiceDataP7IClientiPcx");
            #endif
    }

	if (!adrVoiceData)
	{
		g_SMAPI->Format(error, maxlen, "SV_BroadcastVoiceData sigscan failed.");
		return false;
	}

    if (engineVersion == SOURCE_ENGINE_CSGO)
    {
        m_VoiceDetour = DETOUR_CREATE_STATIC(BroadcastVoiceData_Protobuf, adrVoiceData);
    }
    else
    {
        m_VoiceDetour = DETOUR_CREATE_STATIC(BroadcastVoiceData, adrVoiceData);
    }

    if (!m_VoiceDetour)
	{
		g_SMAPI->Format(error, maxlen, "SV_BroadcastVoiceData detour failed.");
		return false;
	}

    m_VoiceDetour->EnableDetour();

	return true;
}

void Shim::BroadcastVoiceData_Callback(int bytes, const char *data)
{
    if (!data || bytes <= 0) {
        return;
    }

#if 0
	// This is useful for dumping voice data for debugging.
	static int packet = 0;
    char filename[64];
	sprintf(filename, "voice_%lld_%02d.bin", steamId, packet++);
	FILE *file = fopen(filename, "wb");
	fwrite(data, bytes, 1, file);
	fclose(file);
#endif

    receive_audio(bytes, data);
}

void Shim::ClientPutInServer(edict_t *pEntity, char const *playername)
{
    if (!pEntity || pEntity->IsFree())
        return;

    const CSteamID *steamid = engine->GetClientSteamID(pEntity);

    if (!steamid)
        return;

    client_put_in_server(steamid->ConvertToUint64(), playername);
}

void Shim::ClientDisconnect(edict_t *pEntity)
{
    if (!pEntity || pEntity->IsFree())
        return;

    const CSteamID *steamid = engine->GetClientSteamID(pEntity);

    if (!steamid)
        return;

    client_disconnect(steamid->ConvertToUint64());
}

void Shim::ClientCommand(edict_t *pEntity, const CCommand &args)
{
}

void *Shim::OnMetamodQuery(const char* iface, int *ret)
{
    if (strcmp(iface, SOURCEMOD_NOTICE_EXTENSIONS) == 0) {
        BindToSourcemod();
    }

    if (ret != NULL) {
        *ret = IFACE_OK;
    }

    return NULL;
}

bool Shim::Unload(char *error, size_t maxlen)
{
    SM_UnloadExtension();

    SH_REMOVE_HOOK_MEMFUNC(IServerGameClients, ClientDisconnect, gameclients, this, &Shim::ClientDisconnect, true);
    SH_REMOVE_HOOK_MEMFUNC(IServerGameClients, ClientPutInServer, gameclients, this, &Shim::ClientPutInServer, true);
    SH_REMOVE_HOOK_MEMFUNC(IServerGameClients, ClientCommand, gameclients, this, &Shim::ClientCommand, false);

    if (this->m_VoiceDetour)
    {
        this->m_VoiceDetour->Destroy();
        this->m_VoiceDetour = NULL;
    }

	return true;
}

void Shim::BindToSourcemod()
{
    char error[256];

	if (!SM_LoadExtension(error, sizeof(error))) {
		char message[512];
		snprintf(message, sizeof(message), "Could not load as a SourceMod extension: %s\n", error);
		engine->LogPrint(message);
	}
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
