
#if !defined(_TSK_INTERFACE_H)
#define _TSK_INTERFACE_H

#include "TSKType.h"



#ifdef _WIN32
#define TSK_EXPORT  1

#ifdef __APPLICATION
	#ifdef TSK_EXPORT
		#define TSK_DECLSPEC __declspec(dllexport)
	#else
		#define TSK_DECLSPEC __declspec(dllimport)
	#endif
#endif
#ifdef __KERNEL
	#define TSK_DECLSPEC 
#endif
#else
#define TSK_DECLSPEC  __attribute__ ((visibility("default")))
#define __stdcall
#endif


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_Init();


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC void __stdcall TSK_UnInit();

#ifdef _WIN32

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetController(BOOLEAN b);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_Exit(JIT_ID * pJitID,UINT32 * pJitInfoLen);


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_Debug(UINT32 nAction,UINT32 nObject,UINT32 nArgu1,InfoPointer pArgu2);
#endif

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_GetRandomKey(UINT32 nKeyLen,BYTE * pKey);


#ifdef _WIN32

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_GetJITInfo(JIT_ID JitID,UINT32 nJitInfoLen,BYTE * pJitInfo);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetJITInfo(UINT32 nJitInfoLen,BYTE * pJitInfo,JIT_ID * pJitID);

#endif


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_DigestCRC32(UINT32 nFlowLen,BYTE * pFlow,PVOID * pCookie,UINT32 * pCRC32);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_DigestMD5(UINT32 nFlowLen,BYTE * pFlow,PVOID * pCookie,BYTE * pMD5);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_DigestSha256(UINT32 nFlowLen,BYTE * pFlow,PVOID * pCookie,BYTE * pSHA256);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityIssue(IdentityObject * pIdentity,BYTE* pSeed,UINT32 nSeedLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityIssueEx(BYTE* pSeed, UINT32 nSeedLen,BYTE* pPublicKeyBuf, UINT32& nPublicKeyLen,BYTE* pPrivateKenBuf, UINT32& nPrivateKeyLen,BYTE* pKeyID);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityFree(IdentityObject identity);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityExport(IdentityObject identity,BYTE Action,UINT32 BufLen,BYTE * Buf,UINT32 * pWrittenLen);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityImport(BYTE Action,UINT32 BufLen,BYTE * Buf,IdentityObject * pIdentity);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityEncrypt(IdentityObject identity,BYTE cryptAction,
					UINT32 nSrcFlowLen,BYTE * pSrcFlow,
					UINT32 nTarFlowLen,BYTE * pTarFlow,
					UINT32 * pTarFlowReturnLen);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_IdentityDecrypt(IdentityObject identity,BYTE cryptAction,
					UINT32 nSrcFlowLen,BYTE * pSrcFlow,
					UINT32 nTarFlowLen,BYTE * pTarFlow,
					UINT32 * pTarFlowReturnLen);


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_GetCipherDesc(BYTE * pDescBuf,INT32 nDescBufLen);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_GetCipherInfo(BYTE * tagCipher,UINT16 * pKeyLength,UINT16 * pBlockSize);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetDefaultCipher(BYTE * tagCipher);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetDefaultZone(INT64 nLen);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_CreateCipherCell(BYTE * tagCipher,UINT32 nKeyLen,BYTE * pKey);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_DeleteCipherCell(INT32 nCipherCell);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_ProcessCipherCell(INT32 nCipherCell,BYTE nAction,INT64 nOffset,UINT32 nLength,BYTE * pFlow);


#ifdef _WIN32

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_CreateEventClass(UINT32 nEventClassID,UINT32 nEventFlag);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_DeleteEventClass(UINT32 nEventClassID);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_RegisterModule(UINT32 nModuleID,PROC_ID nPID);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_UnRegisterModule(UINT32 nModuleID,PROC_ID nPID);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_QueryEvent(UINT32 nEventClassID,UINT32 nThisModuleID,UINT32 nEventTypeLeft,UINT32 nEventTypeRight,
					UINT32 nWaitTime,UINT32 * pAckCookie,JIT_ID * pJitID,UINT32 * pJitInfoLen);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_PostEvent(UINT32 nEventClassID,UINT32 nEventFlag,UINT32 nThatModuleID,
					UINT32 nEventLen,BYTE * pEvent);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SendEvent(UINT32 nEventClassID,UINT32 nEventFlag,
					UINT32 nThatModuleID,UINT32 nWaitTime,
					UINT32 nEventLen,BYTE * pEvent,JIT_ID * pJitID,UINT32 * pJitInfoLen);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_AckEvent(UINT32 nAckCookie,UINT32 nResultLen,BYTE * pResult);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_DeleteEvent(UINT32 nEventClassID,UINT32 nDeleteLogic,UINT32 nDeleteArguLeft,UINT32 nDeleteArguRight);

#endif

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_LoginUser(BYTE* nIdentityID);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_LogoutUser();


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetHoldIdentity(UINT32 nListVerb,
					BYTE* nIdentityID, PermissionInfo nPermission, int nLastTime,BYTE * pKeyBuf,INT32 nKeyLen);

#ifdef _WIN32

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetHideFolder(UINT32 nListVerb,BOOLEAN bSelfHide,WCHAR * strUserHide);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_GetHideFolder(JIT_ID * pJitID,UINT32 * pJitInfoLen);


#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetCMGlobalIgnore(UINT32 nListVerb,WCHAR * strTemplate);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetCMGlobalMust(UINT32 nListVerb,WCHAR * strTemplate);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetCMConfigProcess(UINT32 nListVerb,
					BYTE nProcessType,BYTE nCheatLevel,BOOLEAN bInherit,
					WCHAR * strProcessName);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_SetCMRunProcess(UINT32 nListVerb,
					PROC_ID nPID,BYTE nProcessType,BYTE nCheatLevel,BOOLEAN bInherit,
					WCHAR * strProcessName);

#ifdef __cplusplus
extern "C" 
#endif
TSK_DECLSPEC INT32 __stdcall TSK_GetCMRunProcess(JIT_ID * pJitID,UINT32 * pJitInfoLen);

#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_Init(INT32 nAction,INT32 nThreadCount);

#ifdef __cplusplus
extern "C"
#endif
#ifdef _WIN32
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_PushAFileToTask(INT32 nTaskID, WCHAR* pSrcFullPath,WCHAR* pDstFullPath);
#else
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_PushAFileToTask(INT32 nTaskID, char* pSrcFullPath, char* pDstFullPath);
#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_PushAAction(INT32 nTaskID, BOOLEAN bDelAdd,
	BYTE* IID, PermissionInfo nPermission, int nLastTime, BYTE* pKeyBuf, DWORD nKeyBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_SetShareSign(INT32 nTaskID, BYTE nSignAction);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_StartATask(INT32 nTaskID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_CancelATask(INT32 nTaskID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_UnInitATask(INT32 nTaskID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_GetATaskState(INT32 nTaskID,INT32* nState,INT64* nTotalLength,INT64* nCurrentLength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_GetATaskReport(INT32 nTaskID, INT32 nReportMode,BYTE* pReportBuf,INT32* nReportBufLength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOpTask_CanExit();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_FileOp_AdjustByFlow(BYTE* pHeadSrc, INT32 nHeadFlowSrcLen, BYTE* pHeadFlowDst, INT32* nHeadFlowDstLen,
	BOOLEAN bDelAdd, BYTE* nIID, PermissionInfo nPermission, int nLastTime, BYTE* pKeyBuf, DWORD nKeyBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSK_LocalServer_Start(INT32 nPort);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC void __stdcall TSK_LocalServer_Stop();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC void __stdcall TSK_LocalServer_SetSystemTime(INT32 nTime);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC bool __stdcall TSK_LocalServer_GetSessionKey(char pSessionKey[20]);

#ifdef __cplusplus
extern "C"
#endif
#ifdef _WIN32
TSK_DECLSPEC void __stdcall TSK_LocalServer_SetCurrentWebSiteParam(WCHAR* strRootDir,WCHAR* strDynamicFlag,WCHAR* strRemoteServer,WCHAR* strIndexFileName,WCHAR* strPemFilePath);
#else
TSK_DECLSPEC void __stdcall TSK_LocalServer_SetCurrentWebSiteParam(char* strRootDir, char* strDynamicFlag, char* strRemoteServer,char* strIndexFileName,char* strPemFilePath);
#endif

#endif
