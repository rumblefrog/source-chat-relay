#pragma semicolon 1

#define PLUGIN_AUTHOR "Fishy"
#define PLUGIN_VERSION "0.0.8"

#include <sourcemod>
#include <morecolors>
#include <socket>
#include <smlib>

#pragma newdecls required

#define HEADER_LEN 161

enum RelayFrame {
	Authenticate,
	Message,
	Unknown
}

char sHostname[64];
char sHost[64] = "127.0.0.1";
char sToken[64];
char sPrefix[8];

// Randomly selected port
int iPort = 57452;
int iFlag;

bool bFlag;

ConVar cHost;
ConVar cPort;
ConVar cPrefix;
ConVar cFlag;

Handle hSocket;

public Plugin myinfo = 
{
	name = "Source Chat Relay",
	author = PLUGIN_AUTHOR,
	description = "Prototype Stage",
	version = PLUGIN_VERSION,
	url = "https://keybase.io/RumbleFrog"
};

public void OnPluginStart()
{
	CreateConVar("rf_scr_version", PLUGIN_VERSION, "Source Chat Relay Version", FCVAR_REPLICATED | FCVAR_SPONLY | FCVAR_DONTRECORD | FCVAR_NOTIFY);

	cHost = CreateConVar("rf_scr_host", "127.0.0.1", "Relay Server Host", FCVAR_PROTECTED);

	cPort = CreateConVar("rf_scr_port", "57452", "Relay Server Port", FCVAR_PROTECTED);
	
	cPrefix = CreateConVar("rf_scr_prefix", "", "Prefix required to send message to Discord", FCVAR_NONE);
	
	cFlag = CreateConVar("rf_scr_flag", "", "If prefix is enabled, this admin flag is required to send message using the prefix", FCVAR_PROTECTED);

	AutoExecConfig(true, "Source-Server-Relay");

	hSocket = SocketCreate(SOCKET_TCP, OnSocketError);

	SocketSetOption(hSocket, SocketReuseAddr, 1);
	SocketSetOption(hSocket, SocketKeepAlive, 1);
	
	#if defined DEBUG
		SocketSetOption(hSocket, DebugMode, 1);
	#endif
}

public void OnConfigsExecuted()
{
	GetConVarString(FindConVar("hostname"), sHostname, sizeof sHostname);

	cHost.GetString(sHost, sizeof sHost);
	
	cPrefix.GetString(sPrefix, sizeof sPrefix);
	
	iPort = cPort.IntValue;
	
	char sFlag[8];
	
	cFlag.GetString(sFlag, sizeof sFlag);
	
	if (!StrEqual(sFlag, ""))
	{
		AdminFlag aFlag;
		
		bFlag = FindFlagByChar(sFlag[0], aFlag);
		
		iFlag = FlagToBit(aFlag);
	}
		
	char sPath[PLATFORM_MAX_PATH], sIP[64];
	
	Server_GetIPString(sIP, sizeof sIP);
	
	BuildPath(Path_SM, sPath, sizeof sPath, "data/%s:%d.data", sIP, Server_GetPort());
	
	if (FileExists(sPath, false))
	{
		File tFile = OpenFile(sPath, "r", false);
		
		tFile.ReadString(sToken, sizeof sToken, -1);
		
		delete tFile;
	} else
	{
		File tFile = OpenFile(sPath, "w", false);
	
		GenerateRandomChars(sToken, sizeof sToken, 32);
	
		tFile.WriteString(sToken, true);
	
		delete tFile;
	}

	if (!SocketIsConnected(hSocket))
		ConnectRelay();
}

void ConnectRelay()
{	
	if (!SocketIsConnected(hSocket))
		SocketConnect(hSocket, OnSocketConnected, OnSocketReceive, OnSocketDisconnected, sHost, iPort);
	else
		PrintToServer("Socket is already connected?");
}

public Action Timer_Reconnect(Handle timer)
{
	ConnectRelay();
}

void StartReconnectTimer()
{
	if (SocketIsConnected(hSocket))
		SocketDisconnect(hSocket);
		
	CreateTimer(10.0, Timer_Reconnect);
}

public int OnSocketDisconnected(Handle socket, any arg)
{	
	StartReconnectTimer();
	
	PrintToServer("Socket disconnected");
}

public int OnSocketError(Handle socket, int errorType, int errorNum, any ary)
{
	StartReconnectTimer();
	
	LogError("Socket error %i (errno %i)", errorType, errorNum);
}

public int OnSocketConnected(Handle socket, any arg)
{
	PackFrame(Authenticate, sToken);
	
	PrintToServer("Successfully Connected");
}

public int OnSocketReceive(Handle socket, const char[] receiveData, int dataSize, any arg)
{
	PrintToServer(receiveData);
	
	ParseMessageFrame(receiveData);
}

public void OnClientSayCommand_Post(int client, const char[] command, const char[] sArgs)
{
	if (!Client_IsValid(client))
		return;
		
	if (!SocketIsConnected(hSocket))
		return;
		
	if (StrEqual(sPrefix, ""))
		PackMessage(client, sArgs);
	else
	{
		if (bFlag && !CheckCommandAccess(client, "arandomcommandthatsnotregistered", iFlag, true))
			return;
		
		int iLen = strlen(sPrefix);
		
		bool bMatch = true;
		
		for (int i = 0; i < iLen; i++)
		{
			if (sPrefix[i] != sArgs[i])
			{
				bMatch = false;
				break;
			}
		}
		
		char sBuffer[MAX_MESSAGE_LENGTH];
		
		for (int i = iLen; i < strlen(sArgs); i++)
			Format(sBuffer, sizeof sBuffer, "%s%c", sBuffer, sArgs[i]);
		
		if (bMatch)
			PackMessage(client, sBuffer);
	}
}

void PackMessage(int client, const char[] message)
{
	// 64 - Hostname
	// 64 - ClientID (Snowflake / SteamID64)
	// 32 - ClientName
	// remaining
	
	int iMessageLen = strlen(message);
	
	int iFrameLen = HEADER_LEN + iMessageLen;
	
	char[] sFrame = new char[iFrameLen];
	
	char sName[32], sID64[64];
	
	Format(sFrame, iFrameLen, "%s%-64s", sFrame, sHostname);
	
	GetClientAuthId(client, AuthId_SteamID64, sID64, sizeof sID64);
	
	GetClientName(client, sName, sizeof sName);
	
	Format(sFrame, iFrameLen, "%s%-64s", sFrame, sID64);
	
	Format(sFrame, iFrameLen, "%s%-32s", sFrame, sName);
	
	Format(sFrame, iFrameLen, "%s%s", sFrame, message);
	
	PackFrame(Message, sFrame);
}

void PackFrame(RelayFrame opcode, const char[] payload)
{
	int iPayloadLen = strlen(payload);
	int iLen = iPayloadLen + 4;
	
	char[] sFrame = new char[iLen];
	
	switch (opcode)
	{
		case Authenticate:
		{
			sFrame[0] = '0';
			
			Format(sFrame, iLen, "%s%s", sFrame, payload);
		}
		case Message:
		{
			sFrame[0] = '1';
				
			Format(sFrame, iLen, "%s%s", sFrame, payload);
		}
	}
	
	#if defined DEBUG
		PrintToServer("%s", sFrame);
	#endif
	
	SendFrame(sFrame);
}

void SendFrame(const char[] frame)
{
	SocketSend(hSocket, frame);
	
	#if defined DEBUG
		PrintToConsoleAll(frame);
	#endif
}

void ParseMessageFrame(const char[] frame)
{
	if (frame[0] != '1')
		return;
	
	char b1[64], b2[64], b3[32], hostname[64], id64[64], name[32];
	
	int iLen = strlen(frame);
	
	int iOffset = 1;
	
	for (int i = 0; i < 64; i++)
	{
		Format(b1, sizeof b1, "%s%c", b1, frame[iOffset]);
		iOffset++;
	}
	
	String_Trim(b1, hostname, sizeof hostname);
	
	for (int i = 0; i < 64; i++)
	{
		Format(b2, sizeof b2, "%s%c", b2, frame[iOffset]);
		iOffset++;
	}
	
	
	String_Trim(b2, id64, sizeof id64);
	
	for (int i = 0; i < 32; i++)
	{
		Format(b3, sizeof b3, "%s%c", b3, frame[iOffset]);
		iOffset++;
	}
	
	String_Trim(b3, name, sizeof name);
	
	int iContentLen = iLen - iOffset;
	
	char[] sContent = new char[iContentLen];
	
	for (int i = 0; i < iContentLen; i++)
	{
		Format(sContent, iContentLen + 1, "%s%c", sContent, frame[iOffset]);
		iOffset++;
	}
	
	#if defined DEBUG
		PrintToConsoleAll("===== PARSING =====");
	
		PrintToConsoleAll("hostname: %s", hostname);
	
		PrintToConsoleAll("id64: %s", id64);
	
		PrintToConsoleAll("name: %s", name);
	
		PrintToConsoleAll("sContentLen: %d", iContentLen);
	
		PrintToConsoleAll("sContent: %s", sContent);
	
		PrintToConsoleAll("===================");
	#endif
	
	CPrintToChatAll("{purple}[%s] {lightgreen}%s{white}: {grey}%s", hostname, name, sContent);
}

stock void GenerateRandomChars(char[] buffer, int buffersize, int len)
{
	char charset[] = "adefghijstuv6789!@#$%^klmwxyz01bc2345nopqr&+=";
	
	for (int i = 0; i < len; i++)
		Format(buffer, buffersize, "%s%c", buffer, charset[GetRandomInt(0, sizeof charset)]);
}
