#if !defined(_TSK_INFO_H)
#define _TSK_INFO_H

#include "TSKType.h"
#include "TSKInterface.h"


#define TEFS_INFO_VERSION  0

#define SEAL_INFO_VERSION  0
#define TEFS_INFO_UNIT 8192
#define TEFS_INFO_BLANK 1024
#define TEFS_INFO_MAXLENGTH 0x100000
#define TEFS_INFO_TAG_LENGTH 32
#define UFID_LENGTH 16
#define IDENTITY_CHECK_LENGTH 8
#define IDENTITY_CHECK_VALUE  "19781124"
#define MULTI_ENCRYPT_NULL          0
#define MULTI_ENCRYPT_DATA_YES      1
#define MULTI_ENCRYPT_DATA_NO       2

extern BYTE g_TEFStag[];

typedef struct _PerIdentitySealInfoFake
{
	UINT16 nSealVersion;
	BYTE   pCheck[IDENTITY_CHECK_LENGTH];
	UINT16   nSessionKeyLength;
	BYTE * pSessionKey;
	PermissionInfo nPermissionInfo;
	int nLastTime;
}PerIdentitySealInfoFake,*PPerIdentitySealInfoFake;

typedef struct _IdentitySealInfo
{
	UINT16 nCount;
	UINT16 nPerLen;
	BYTE * pFlow;

	_IdentitySealInfo() {pFlow=NULL;}
	~_IdentitySealInfo() {if (pFlow) __mem_free(pFlow);}

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_IdentitySealInfo & si);
	void Packet(BYTE * pBuf);
	BOOLEAN UnPacket(BYTE * pBuf,UINT16 nVersion=TEFS_INFO_VERSION);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{
		return (2*sizeof(UINT16)+((UINT32)nCount)*((UINT32)nPerLen));
	}
}IdentitySealInfo,*PIdentitySealInfo;

typedef struct _PerIdentityViewInfo
{
	BYTE IID[IDENTITYID_LENGTH];
	PermissionInfo nPermissionInfo;
	int nLastTime;

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_PerIdentityViewInfo & si);
	void Packet(BYTE * pBuf);
	BOOLEAN UnPacket(BYTE * pBuf,UINT16 nVersion=TEFS_INFO_VERSION);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{
		return (IDENTITYID_LENGTH +sizeof(int)+sizeof(PermissionInfo));
	}
}PerIdentityViewInfo,*PPerIdentityViewInfo;
typedef struct _IdentityViewInfo
{
	UINT16 nCount;
	UINT16 nPerLen;
	PPerIdentityViewInfo arrayInfo;

	_IdentityViewInfo() {arrayInfo=NULL;}
	~_IdentityViewInfo() {if (arrayInfo) __mem_free(arrayInfo);}

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_IdentityViewInfo & si);
	void Packet(BYTE * pBuf);
	BOOLEAN UnPacket(BYTE * pBuf,UINT16 nVersion=TEFS_INFO_VERSION);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{
		return (2*sizeof(UINT16)+((UINT32)nCount)*((UINT32)nPerLen));
	}
}IdentityViewInfo,*PIdentityViewInfo;


typedef struct _CipherInfo
{
	BYTE tag[CIPHER_TAG_LENGTH];
	UINT16 nKeyLength;
	UINT16 nBlockSize;
	BYTE nCounterMode;
	BYTE * pCipherKey;

	_CipherInfo() {pCipherKey=NULL;}
	~_CipherInfo() {if (pCipherKey) __mem_free(pCipherKey);}

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_CipherInfo & si);
	void Packet(BYTE * pBuf);
	BOOLEAN UnPacket(BYTE * pBuf,UINT16 nVersion=TEFS_INFO_VERSION);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{
		UINT32 temp=(UINT32)RoundUp(nKeyLength,nBlockSize);
		return (CIPHER_TAG_LENGTH+sizeof(UINT16)+sizeof(UINT16)+sizeof(BYTE)+temp);
	}
}CipherInfo,*PCipherInfo;

typedef struct _ZoneItem
{
	uint64_t start;
	uint64_t length;

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}
}ZoneItem,*PZoneItem;
typedef struct _ZoneInfo
{
	UINT32   nCount;
	ZoneItem zone;
	                        	
	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_ZoneInfo & si);
	void Packet(BYTE * pBuf);
	BOOLEAN UnPacket(BYTE * pBuf,UINT16 nVersion=TEFS_INFO_VERSION);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{
		return (sizeof(UINT32)+(sizeof(uint64_t)+sizeof(uint64_t)));
	}
}ZoneInfo,* PZoneInfo;


typedef struct _MainInfo
{
	BYTE    nUFID[UFID_LENGTH];

	BYTE OwnerIID[IDENTITYID_LENGTH];

	BYTE VerifyIID[IDENTITYID_LENGTH];
	BYTE  nDigestVerify[IDENTITY_MAX_CIPHER_LENGTH];
	
	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_MainInfo & si);
	void Packet(BYTE * pBuf);
	BOOLEAN UnPacket(BYTE * pBuf,UINT16 nVersion=TEFS_INFO_VERSION);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{
		return (UFID_LENGTH+ IDENTITYID_LENGTH *2+IDENTITY_MAX_CIPHER_LENGTH);
	}
}MainInfo,* PMainInfo;


class TEFSInfoRef;
typedef struct _TEFSInfo 
{
	BYTE   tag[TEFS_INFO_TAG_LENGTH];
	BYTE   nMultiEncryptInfo;
	UINT32  nCRC32Check;
 	UINT16   nVersion;
	UINT32  nContentLen;
	UINT32  nFillLen;

	MainInfo mainInfo;
	ZoneInfo zoneInfo;
	CipherInfo cipherInfo;
	IdentitySealInfo idSealInfo;
	IdentityViewInfo idViewInfo;

	void* operator new(size_t n) {return __mem_alloc(PagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}

	BOOLEAN Copy(_TEFSInfo & si);
	void Packet(BYTE * pBuf,TEFSInfoRef * pRef);
	INT32 UnPacket(BYTE * pBuf,TEFSInfoRef * pRef);
	UINT32 GetPacketLen(UINT16 nVersion=TEFS_INFO_VERSION)
	{	
		return (TEFS_INFO_TAG_LENGTH+sizeof(BYTE)+sizeof(UINT16)+sizeof(UINT32)*3+
				mainInfo.GetPacketLen(nVersion)+zoneInfo.GetPacketLen(nVersion)+cipherInfo.GetPacketLen(nVersion)+
				idSealInfo.GetPacketLen(nVersion)+idViewInfo.GetPacketLen(nVersion));
	}
}TEFSInfo,* PTEFSInfo;

class TEFSInfoRef
{

public:

	PTEFSInfo pTEFSInfo;


	BYTE   nMultiEncryptInfo;
	

	UINT16   nVersionOrigin;
	UINT32  nContentLenOrigin;
	UINT32  nFillLenOrigin;
	UINT16   nVersionMem;
	UINT32  nContentLenMem;
	UINT32  nFillLenMem;
	BOOLEAN bCheckMem;


	INT64 nEncryptZone;

	BOOLEAN bSealValid;


	UINT16   nKeyLength;
	UINT16   nBlockSize;
	BYTE   nCounterMode;
	BYTE * pSessionKey;
	BYTE * pCipherKey;
	INT32    nCipherCell;

	PermissionInfo nPermissionInfo;
	int nLastTime;

	TEFSInfoRef() 
	{
		pTEFSInfo=NULL;
		pSessionKey=NULL;
		pCipherKey=NULL;
		nCipherCell=-1;
	}
	~TEFSInfoRef() 
	{
		if (pTEFSInfo) delete pTEFSInfo;
		if (pSessionKey) __mem_free(pSessionKey);
		if (pCipherKey) __mem_free(pCipherKey);
		if (nCipherCell>=0) TSK_DeleteCipherCell(nCipherCell);
	}

	void* operator new(size_t n) {return __mem_alloc(NonPagedPool,n,POOLTAG_TSK);}
	void operator delete(void* p) {__mem_free(p);}



public:

	BOOLEAN SetTEFSInfo()
	{
		if (pTEFSInfo==NULL)
		{
			pTEFSInfo=new TEFSInfo;
			return (pTEFSInfo!=NULL);
		}
		return FALSE;
	}
	void ClearTEFSInfo()
	{
		if (pTEFSInfo)
		{
			delete pTEFSInfo;
			pTEFSInfo=NULL;
		}
	}


	void Packet(BYTE * pBuf) {pTEFSInfo->Packet(pBuf,this);}
	INT32 UnPacket(BYTE * pBuf) {return pTEFSInfo->UnPacket(pBuf,this);}
	UINT32 GetPacketLen() {return pTEFSInfo->GetPacketLen(TEFS_INFO_VERSION);}


	void ClearCheckMem() {bCheckMem=FALSE;}


	UINT32 GetHeaderLenMem() {return (nContentLenMem+nFillLenMem);}

	UINT32 GetHeaderLenOrigin() {return (nContentLenOrigin+nFillLenOrigin);}

	void MemVersion2OriginVersion() 
	{
		nVersionOrigin=nVersionMem;
		nContentLenOrigin=nContentLenMem;
		nFillLenOrigin=nFillLenMem;
	}
};
typedef class TEFSInfoRef * PTEFSInfoRef;


#define CHECK_TEFSINFO_OK             0
#define CHECK_TEFSINFO_NOT            1
#define CHECK_TEFSINFO_VERSION        2
#define CHECK_TEFSINFO_FAIL           3
extern BYTE TEFSInfoCheckKey(BYTE * pBuf,UINT16 * pVersion,UINT32 * pContentLen,UINT32 * pFillLen);

#endif
