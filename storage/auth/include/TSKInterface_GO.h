
#if !defined(_TSK_INTERFACE_GO_H)
#define _TSK_INTERFACE_GO_H
#include "TSKType.h"
#ifdef _WIN32
#define TSK_DECLSPEC __declspec(dllexport)
#else
#define TSK_DECLSPEC  __attribute__ ((visibility("default")))
#define __stdcall
#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_Init();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC void __stdcall TSK_GO_UnInit();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_DigestCRC32(unsigned int nFlowLen, unsigned char * pFlow, void* * pCookie, unsigned int * pCRC32);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_DigestMD5(unsigned int nFlowLen, unsigned char * pFlow, void* * pCookie, unsigned char * pMD5);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_DigestSha256(unsigned int nFlowLen, unsigned char * pFlow, void* * pCookie, unsigned char * pSHA256);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityIssue(void* * pIdentity, unsigned char* pSeed, unsigned int nSeedLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityIssueEx(unsigned char* pSeed, unsigned int nSeedLen, unsigned char* pPublicKeyBuf, unsigned int* nPublicKeyLen, unsigned char* pPrivateKeyBuf, unsigned int* nPrivateKeyLen,unsigned char* pKeyID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityFree(void* identity);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityExport(void* identity, unsigned char Action, unsigned int BufLen, unsigned char * Buf, unsigned int * pWrittenLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityImport(unsigned char Action, unsigned int BufLen, unsigned char * Buf, void* * pIdentity);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityEncrypt(void* identity, unsigned char cryptAction,
	unsigned int nSrcFlowLen, unsigned char * pSrcFlow,
	unsigned int nTarFlowLen, unsigned char * pTarFlow,
	unsigned int * pTarFlowReturnLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_IdentityDecrypt(void* identity, unsigned char cryptAction,
	unsigned int nSrcFlowLen, unsigned char * pSrcFlow,
	unsigned int nTarFlowLen, unsigned char * pTarFlow,
	unsigned int * pTarFlowReturnLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_GetCipherDesc(unsigned char * pDescBuf, int nDescBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_GetCipherInfo(unsigned char * tagCipher, unsigned short * pKeyLength, unsigned short * pBlockSize);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_SetDefaultCipher(unsigned char * tagCipher);

#ifdef __cplusplus
extern "C"
#endif
#ifdef _WIN32
TSK_DECLSPEC int __stdcall TSK_GO_SetDefaultZone(signed long long nLen);
#else
TSK_DECLSPEC int  __stdcall TSK_GO_SetDefaultZone(int64_t nLen);
#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_CreateCipherCell(unsigned char * tagCipher, unsigned short nKeyLen, unsigned char * pKey);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_DeleteCipherCell(int nCipherCell);

#ifdef __cplusplus
extern "C"
#endif
#ifdef _WIN32
TSK_DECLSPEC int __stdcall TSK_GO_ProcessCipherCell(int nCipherCell, unsigned char nAction, signed long long nOffset, unsigned int nLength, unsigned char * pFlow);
#else
TSK_DECLSPEC int __stdcall TSK_GO_ProcessCipherCell(int nCipherCell, unsigned char nAction, int64_t nOffset, unsigned int nLength, unsigned char * pFlow);
#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_LoginUser(unsigned char* nIdentityID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_LogoutUser();


#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_SetHoldIdentity(unsigned int nListVerb,
	unsigned char* nIdentityID, unsigned short nPermission, int nLastTime, unsigned char * pKeyBuf, int nKeyLen);

#ifdef SDK_FLOWINTERFACE

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_CreateByDefault();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_CreateByFlow(unsigned char * pFlow);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_ModifyForNew(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_ModifyForSeal(int nFlowID, bool bDelAdd,
	unsigned char* nIID, unsigned short nPermission, int nLastTime, unsigned char* pKeyBuf, unsigned int nKeyBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_ModifyForSign(int nFlowID, bool bDelAdd,
	unsigned char * pDigest);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_RecountSeal(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_DigestFlow(int nFlowID,
	INT64 nOffset, unsigned int nLength, unsigned char * pFlow,
	void* * pCookie,
	unsigned char * pMD5);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_ProcessFlow(int nFlowID,
	unsigned char nAction,
	INT64 nOffset, unsigned int nLength, unsigned char * pFlow,
	INT64* pVDL);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_FreeAFlow(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_GetEFSLength(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_GetFillLength(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_SetFillLength(int nFlowID, int filllength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_GetEFSFlow(int nFlowID, unsigned char * pTEFSInfoBuf);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_GetEFSSealInfo(int nFlowID, unsigned char* pIdentityIDArray, int nArrayLength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC unsigned char* __stdcall TSKHelp_GO_Flow_GetFileOwner(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSKHelp_GO_Flow_GetLastUseTime(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC unsigned short __stdcall TSKHelp_GO_Flow_GetLastPermissionInfo(int nFlowID);
#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_Init(int nAction, int nThreadCount);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_PushAAction(int nTaskID, bool bDelAdd,
	unsigned char * IID, unsigned short nPermission, int nLastTime, unsigned char* pKeyBuf, int nKeyBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_SetShareSign(int nTaskID, unsigned char nSignAction);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_StartATask(int nTaskID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_CancelATask(int nTaskID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_UnInitATask(int nTaskID);

#ifdef __cplusplus
extern "C"
#endif
#ifdef _WIN32
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_GetATaskState(int nTaskID, int* nState, long long * nTotalLength, long long * nCurrentLength);
#else
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_GetATaskState(int nTaskID, int* nState, int64_t * nTotalLength, int64_t * nCurrentLength);
#endif

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_GetATaskReport(int nTaskID, int nReportMode, unsigned char* pReportBuf, int* nReportBufLength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_FileOpTask_CanExit();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_Go_FileOp_AdjustByFlow(unsigned char* pHeadSrc, int nHeadFlowSrcLen, unsigned char* pHeadFlowDst, int* nHeadFlowDstLen,
	bool bDelAdd, unsigned char* nIID, unsigned short nPermission, int nLastTime, unsigned char* pKeyBuf, int nKeyBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC int __stdcall TSK_GO_LocalServer_Start(int nPort);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC void __stdcall TSK_GO_LocalServer_Stop();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC void __stdcall TSK_GO_LocalServer_SetSystemTime(int nTime);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC bool __stdcall TSK_GO_LocalServer_GetSessionKey(char pSessionKey[20]);


#endif
