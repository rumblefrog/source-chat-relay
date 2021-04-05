#ifndef _INCLUDE_METAMOD_SOURCE_STUB_PLUGIN_H_
#define _INCLUDE_METAMOD_SOURCE_STUB_PLUGIN_H_

#include <ISmmPlugin.h>
#include <eiface.h>
#include <steam/steamclientpublic.h>
#include <bindings.h>
#include <CDetour/detours.h>

class Shim : public ISmmPlugin
{
    public:
        void BroadcastVoiceData_Callback(int bytes, const char *data);
        void ClientDisconnect(edict_t *pEntity);
        void ClientPutInServer(edict_t *pEntity, char const *playername);
        void ClientCommand(edict_t *pEntity, const CCommand &args);
	public:
		bool Load(PluginId id, ISmmAPI *ismm, char *error, size_t maxlength, bool late);
		bool Unload(char *error, size_t maxlen);
        void *OnMetamodQuery(const char *iface, int *ret);
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
        CDetour *m_VoiceDetour;
};

extern Shim g_shim;

PLUGIN_GLOBALVARS();

#endif //_INCLUDE_METAMOD_SOURCE_STUB_PLUGIN_H_
