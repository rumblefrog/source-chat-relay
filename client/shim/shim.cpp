#include <stdio.h>
#include "shim.h"

#include "MemoryUtils.h"

using namespace SourceHook;

Shim g_shim;

PLUGIN_EXPOSE(Shim, g_shim);

class IClient;
class CCLCMsg_VoiceData;

DETOUR_DECL_STATIC4(BroadcastVoiceData, void, IClient *, client, int, bytes, char *, data, long long, xuid)
{
    META_LOG(g_PLAPI, ">>> SV_BroadcastVoiceData(%p, %d, %p, %lld)", client, bytes, data, xuid);

	DETOUR_STATIC_CALL(BroadcastVoiceData)(client, bytes, data, xuid);

#if 0
	// This is useful for getting the correct m_SteamID offset.
	char filename[64];
	sprintf(filename, "voice_%p_client.bin", client);
	FILE *file = fopen(filename, "wb");
	fwrite(client, 1024, 1, file);
	fclose(file);
#endif

	// xuid isn't populated pre-CS:GO/Protobuf, so get the SteamID from the client instead.
    #ifdef _WIN32
        uint64_t steamId = *(uint64_t *)((uintptr_t)client + 92);
    #else
        uint64_t steamId = *(uint64_t *)((uintptr_t)client + 96);
    #endif

    g_shim.BroadcastVoiceData_Callback(steamId, bytes, data, false);
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

	// TODO: Gamedata this.

	// If this breaks on Linux only, check the libstdc++ ABI in use, see comment in AMBuilder.
	std::string *voiceData = *(std::string **)((uintptr_t)message + 8);

	// The xuid in the message is helpfully set to the steamid on CS:GO/Protobuf, which makes finding these offsets quite easy.
	uint64_t steamId = *(uint64_t *)((uintptr_t)message + 12);

	// CS:GO/Protobuf allows individual clients to choose between Steam or Engine voice encoding, so this can differ per-packet.
	// Note: this is different from our enum, here 0 = steam, 1 = engine (sv_voicecodec).
	int voiceFormat = *(int *)((uintptr_t)message + 20);

	bool forceSteamVoice = (voiceFormat == 0);
    g_shim.BroadcastVoiceData_Callback(steamId, voiceData->size(), voiceData->data(), forceSteamVoice);
}

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

    m_Client = new_client();

	return true;
}

void Shim::BroadcastVoiceData_Callback(uint64_t steamId, int bytes, const char *data, bool forceSteamVoice)
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

    receive_audio(this->m_Client, steamId, bytes, data, forceSteamVoice);
}

bool Shim::Unload(char *error, size_t maxlen)
{
    if (this->m_VoiceDetour)
    {
        this->m_VoiceDetour->Destroy();
        this->m_VoiceDetour = NULL;
    }

    if (this->m_Client)
    {
        free_client(this->m_Client);
        this->m_Client = NULL;
    }

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
