#if !defined(_TSK_HELP_H)
#define _TSK_HELP_H

#include "TSKType.h"
#include "TSKInterface.h"
#include "KList.h"
#include "TSKInfo.h"

#ifdef _WIN32
#include <windows.h>

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
#include <stdlib.h>

#define TSK_DECLSPEC
#endif


#define BIGFILE_LEAST (8*1024*1024)
#define BIGFILE_BASE (32*1024*1024)
#define BIGFILE_BASE_INTERVAL (2*1024*1024)
#define BIGFILE_BASE_INTERVAL_BITS 21
#define BIGFILE_BASE_INTERVAL_MASK 0x001FFFFF
#define BIGFILE_MIDDLE (256*1024*1024)
#define BIGFILE_MIDDLE_INTERVAL (8*1024*1024)
#define BIGFILE_MIDDLE_INTERVAL_BITS 23
#define BIGFILE_MIDDLE_INTERVAL_MASK 0x007FFFFF
#define BIGFILE_LARGE_INTERVAL (32*1024*1024)
#define BIGFILE_LARGE_INTERVAL_BITS 25
#define BIGFILE_LARGE_INTERVAL_MASK 0x01FFFFFF


typedef struct _TSKIdentity
{
	KLIST_ENTRY linkfield;
	BYTE nIID[IDENTITYID_LENGTH];
	PermissionInfo nPermission;
	int nLastTime;
	IdentityObject object;

	_TSKIdentity() 
	{
		memset(nIID,0,IDENTITYID_LENGTH);
		object=NULL;
	}
	~_TSKIdentity() 
	{
		if (object) TSK_IdentityFree(object);
	}

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	void Clear()
	{
		memset(nIID,0,IDENTITYID_LENGTH);
		ClearPermission(&nPermission);
		nLastTime = 0;
		if (object) {TSK_IdentityFree(object);object=NULL;}
	}

	BOOLEAN BuildFromNone(BYTE nArguIID[IDENTITYID_LENGTH], PermissionInfo nArguPermission,int nArguLastTime)
	{
		memcpy(nIID,nArguIID,IDENTITYID_LENGTH);
		CopyPermission(&nPermission, &nArguPermission);
		nLastTime = nArguLastTime;
		if (object) {TSK_IdentityFree(object);object=NULL;}
		return TRUE;
	}

	BOOLEAN BuildFromNew(BYTE nArguIID[IDENTITYID_LENGTH], PermissionInfo nArguPermission,int nArguLastTime)
	{
		memcpy(nIID, nArguIID, IDENTITYID_LENGTH);
		CopyPermission(&nPermission, &nArguPermission);
		nLastTime = nArguLastTime;
		if (object) {TSK_IdentityFree(object);object=NULL;}
		if (TSK_IdentityIssue(&object,NULL,0)) {Clear();return FALSE;}
		else return TRUE;
		
	}

	BOOLEAN BuildFromPublic(BYTE nArguIID[IDENTITYID_LENGTH], PermissionInfo nArguPermission, int nArguLastTime,BYTE * pArguPublic,INT32 nArguPublicLen)
	{
		memcpy(nIID, nArguIID, IDENTITYID_LENGTH);
		CopyPermission(&nPermission, &nArguPermission);
		nLastTime = nArguLastTime;
		if (object) {TSK_IdentityFree(object);object=NULL;}
		if (pArguPublic == NULL) {Clear();return FALSE;}
		if (TSK_IdentityImport(IDENTITY_PUBLIC_KEY,nArguPublicLen,pArguPublic,&object)) {Clear();return FALSE;}
		else return TRUE;
	}

	BOOLEAN BuildFromPrivate(BYTE nArguIID[IDENTITYID_LENGTH], int nArguLastTime,BYTE * pArguPrivate,INT32 nArguPrivateLen)
	{
		memcpy(nIID, nArguIID, IDENTITYID_LENGTH);
		SetPermission(&nPermission);
		nLastTime = nArguLastTime;
		if (object) {TSK_IdentityFree(object);object=NULL;}
		if (pArguPrivate == NULL) {Clear();return FALSE;}
		if (TSK_IdentityImport(IDENTITY_PRIVATE_KEY,nArguPrivateLen,pArguPrivate,&object)) {Clear();return FALSE;}
		else return TRUE;
	}
}TSKIdentity,*PTSKIdentity;

INT32 TSKHelp_CreateByDefault(PTEFSInfoRef * ppTEFSInfoRef);

INT32 TSKHelp_CreateByFlow(BYTE * pFlow,PTEFSInfoRef * ppTEFSInfoRef);

INT32 TSKHelp_ModifyForNew(PTEFSInfoRef pTEFSInfoRef);

INT32 TSKHelp_ModifyForSeal(PTEFSInfoRef pTEFSInfoRef,
								BOOLEAN bDelAdd,
								PTSKIdentity pIdentity);

INT32 TSKHelp_ModifyForSign(PTEFSInfoRef pTEFSInfoRef,
								BOOLEAN bDelAdd,
								BYTE * pDigest);


INT32 TSKHelp_ModifyForWholeSeal(PTEFSInfoRef pTEFSInfoRef,
								PKLIST_ENTRY pDistributeList);

INT32 TSKHelp_RecountSeal(PTEFSInfoRef pTEFSInfoRef);

INT32 TSKHelp_DigestFlow(PTEFSInfoRef pTEFSInfoRef,
					INT64 nOffset,UINT32 nLength,BYTE * pFlow,
					PVOID * pCookie,
					BYTE * pMD5);

#ifdef _WIN32
INT32 TSKHelp_ProcessFlow(PTEFSInfoRef pTEFSInfoRef,
					BYTE nAction,
					INT64 nOffset,UINT32 nLength,BYTE * pFlow,
					PLONGLONG pVDL);
#else
INT32  TSKHelp_ProcessFlow(PTEFSInfoRef pTEFSInfoRef,
	BYTE nAction,
	INT64 nOffset, UINT32 nLength, BYTE * pFlow,
	INT64* pVDL);
#endif

void TSKHelp_Free(PTEFSInfoRef pTEFSInfoRef);


#ifdef _WIN32
INT32 TSKHelp_GetFileVDL(HANDLE hFile,PLONGLONG pVDL);
#endif

INT32 TSKHelp_GetPathInfo(PWCHAR strPath,BOOLEAN bFolder,BYTE * pInfo);

INT32 TSKHelp_GetKernelName(PWCHAR strAppPath,BOOLEAN bFolder,PWCHAR strKernelPath,UINT32 nKernelPathLen);

#ifdef SDK_FLOWINTERFACE

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_CreateByDefault();

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_CreateByFlow(BYTE * pFlow);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_ModifyForNew(INT32 nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_ModifyForSeal(int nFlowID, BOOLEAN bDelAdd,
	BYTE* nIID, PermissionInfo nPermission, int nLastTime, BYTE* pKeyBuf, DWORD nKeyBufLen);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_ModifyForSign(int nFlowID, BOOLEAN bDelAdd,
	BYTE * pDigest);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_RecountSeal(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_DigestFlow(int nFlowID,
	INT64 nOffset, UINT32 nLength, BYTE * pFlow,
	PVOID * pCookie,
	BYTE * pMD5);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_ProcessFlow(int nFlowID,
	BYTE nAction,
	INT64 nOffset, UINT32 nLength, BYTE * pFlow,
	INT64* pVDL);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_FreeAFlow(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_GetEFSLength(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_GetFillLength(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_SetFillLength(int nFlowID, int filllength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_GetEFSFlow(int nFlowID, BYTE * pTEFSInfoBuf);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_GetEFSSealInfo(int nFlowID, BYTE* pIdentityIDArray, INT32 nArrayLength);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC BYTE* __stdcall TSKHelp_Flow_GetFileOwner(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC INT32 __stdcall TSKHelp_Flow_GetLastUseTime(int nFlowID);

#ifdef __cplusplus
extern "C"
#endif
TSK_DECLSPEC PermissionInfo __stdcall TSKHelp_Flow_GetLastPermissionInfo(int nFlowID);
#endif

#endif
