#ifndef _INCLUDE_METAMOD_SOURCE_STUB_PLUGIN_H_
#define _INCLUDE_METAMOD_SOURCE_STUB_PLUGIN_H_

#include <ISmmPlugin.h>
#include <eiface.h>
#include <bindings.h>
#include <CDetour/detours.h>

class Shim : public ISmmPlugin
{
    public:
        Shim();
        void BroadcastVoiceData_Callback(uint64_t steamId, int bytes, const char *data, bool forceSteamVoice);
	public:
		bool Load(PluginId id, ISmmAPI *ismm, char *error, size_t maxlength, bool late);
		bool Unload(char *error, size_t maxlen);
	public:
		const char *GetAuthor();
		const char *GetName();
		const char *GetDescription();
		const char *GetURL();
		const char *GetVersion();
		const char *GetDate();
		const char *GetLicense();
		const char *GetLogTag();
    private:
        Client *m_Client;
        CDetour *m_VoiceDetour;
};

extern Shim g_shim;

PLUGIN_GLOBALVARS();

#endif //_INCLUDE_METAMOD_SOURCE_STUB_PLUGIN_H_
