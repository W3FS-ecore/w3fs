
#ifndef _TSK_TYPE_H
#define _TSK_TYPE_H

#include "TSKDefine.h"
#include <string.h>
#include <pthread.h>
#include <stdint.h>
#include <stdbool.h>

#define IDENTITY_MAX_CIPHER_LENGTH 128
#define IDENTITY_MAX_PLAIN_LENGTH 38

#define ECC_MAX_PLAIN_LENGTH  64
#define ECC_MAX_CIPHER_LENGTH 154


#define IDENTITYID_LENGTH 20




typedef char CHAR,*PCHAR;
//typedef wchar_t WCHAR,*PWCHAR;
typedef unsigned char BYTE,*PBYTE;
typedef unsigned char BOOLEAN,*PBOOLEAN;
typedef short INT16;
typedef unsigned short UINT16;
typedef int INT32;
typedef unsigned int UINT32;
typedef int64_t INT64;
typedef uint64_t UINT64;
typedef uint64_t QWORD,*PQWORD;
typedef void * PVOID;
typedef unsigned int  DWORD;




typedef UINT64 InfoPointer;

typedef InfoPointer JIT_ID;

typedef INT32 Storage_ID;

typedef UINT32 PROC_ID;

typedef UINT32 THRD_ID;

typedef void * IdentityObject;




#define FALSE 0
#define TRUE  1

#ifndef INVALID_HANDLE_VALUE
#define INVALID_HANDLE_VALUE ((HANDLE)-1)
#endif



#define JIT_NULL     ((JIT_ID)0)



#define DIGEST_MD5_LENGTH    16

#define DIGEST_SHA256_LENGTH  32

#define IDENTITY_NULL_KEY           0
#define IDENTITY_PUBLIC_KEY     1
#define IDENTITY_PRIVATE_KEY    2
#define IDENTITY_KEY_ID    3

#define CRYPT_NULL         0
#define CRYPT_PUBLIC_KEY   1
#define CRYPT_PRIVATE_KEY  2


#define CIPHER_TAG_LENGTH    16

#define ACTION_ENCRYPT       1
#define ACTION_DECRYPT       2



#define EVENTCLASS_MIN                        0
#define EVENTCLASS_MAX                        31
#define EVENTCLASS_COUNT                      32

#define EVENTCLASS_ACK                        0

#define EVENTCLASS_LOG                        1
#define EVENTCLASS_DS                         2

#define EVENTFLAG_Broadcast                   0x1

#define EVENTFLAG_AutoExpect                  0x2

#define EVENTFLAG_NowConsume                  0x4

#define EVENTFLAG_ControllerGlue              0x8

#define EVENTFLAG_ConsumeGlue                 0x10

#define EVENTFLAG_Important                   0x20

#define EVENTFLAG_Sync                        0x40

#define MODULE_ANY                            0

#define MODULE_Controller                     1

#define MODULE_CurrentPID                    ((PROC_ID)-1)
#define MODULE_SystemPID                     ((PROC_ID)0)

#define EVENT_ANY                             0

#define EVENT_BROADCAST_ACK_SUCCESS  0
#define EVENT_BROADCAST_ACK_ERROR    -1

#define DELETEEVENT_AckCookie   1
#define DELETEEVENT_Type        2
#define DELETEEVENT_Consume     3



#define KError_Success                         0

#define KError_NoDataHandled                   1

#define KError_CreateObject                    10

#define KError_OpenObject                      11

#define KError_OverrideObject                  12


#define KError_LengthError                     -1

#define KError_TimeOut                         -2

#define KError_CantConnectKernel               -10

#define KError_FileServing                     -11

#define KError_MultiUser                       -20

#define KError_NoUser                          -21


#define KError_SyntaxError                     -30

#define KError_CantFindObject                  -31

#define KError_MultiObject                     -32

#define KError_MaybeFindObject                 -33

#define KError_CommonLogic                     -34


#define KError_CipherError                     -40

#define KError_CipherInnerError                -41

#define KError_PermissionDenied                -42

#define KError_HeaderCrcError                  -50

#define KError_IdentitySealError               -51

#define KError_Other                           -100

typedef pthread_mutex_t TEFS_LOCK;
typedef pthread_mutex_t * PTEFS_LOCK;

inline void InitTEFSLock(PTEFS_LOCK pLock) {pthread_mutex_init(pLock,NULL);}

inline void UnInitTEFSLock(PTEFS_LOCK  pLock) {pthread_mutex_destroy(pLock);}

inline void TEFSLock(PTEFS_LOCK pLock) {pthread_mutex_lock(pLock);}

inline void TEFSUnLock(PTEFS_LOCK pLock) {pthread_mutex_unlock(pLock);}


#ifdef __APPLICATION

extern void DbgPrint(const char * strDebug,...);
#define __dbg_print DbgPrint
#endif
#ifdef __KERNEL

#define __dbg_print DbgPrint
#endif



#define FlagOn(_F,_SF)        ((_F) & (_SF))

#define BooleanFlagOn(F,SF)   ((BOOLEAN)(((F) & (SF)) != 0))

#define SetFlag(_F,_SF)       ((_F) |= (_SF))

#define ClearFlag(_F,_SF)     ((_F) &= ~(_SF))

#define RoundUp(x, align) (((INT32) (x) + (align - 1)) & ~(align - 1)) 

#define RoundDown(x, align) ((INT32)(x) & ~(align - 1)) 

#define AlignedOn(x, align) (((INT32)(x) & (align - 1)) == 0) 

#define BooleanAlignedOn(x,align) ((BOOLEAN)(((INT32)(x) & (align - 1)) == 0))


#define StringLen(p) (strlen(p)+1)
#define WStringLen(p) ((wcslen(p)+1)*sizeof(WCHAR))
#define StringFlowLen(flow,pos) (strlen((PCHAR)((BYTE*)flow+pos))+1)
#define WStringFlowLen(flow,pos) ((wcslen((PWCHAR)((BYTE*)flow+pos))+1)*sizeof(WCHAR))

#define LocateFlow_SimplePointer(flow,pos,type,p) {p=(type*)((BYTE*)(flow)+pos);pos+=sizeof(type);}

#define LocateFlow_String(flow,pos,p) {p=(PCHAR)((BYTE*)(flow)+pos);pos+=(strlen(p)+1);} 

#define LocateFlow_WString(flow,pos,p) {p=(PWCHAR)((BYTE*)(flow)+pos);pos+=(wcslen(p)+1)*sizeof(WCHAR);}

#define LocateFlow_Array(flow,pos,type,p,len) {p=(type*)((BYTE*)(flow)+pos);pos+=((len)*sizeof(type));}

#define RestoreFlow_Simple(flow,pos,type,d) {d=*(type*)((BYTE*)(flow)+pos);pos+=sizeof(type);}

#define RestoreFlow_SimplePointer(flow,pos,type,p) {*(type*)(p)=*(type*)((BYTE*)(flow)+pos);pos+=sizeof(type);}

#define RestoreFlow_String(flow,pos,p) {memcpy(p,(BYTE*)(flow)+pos,StringFlowLen(flow,pos));pos+=StringFlowLen(flow,pos);}

#define RestoreFlow_WString(flow,pos,p) {memcpy(p,(BYTE*)(flow)+pos,WStringFlowLen(flow,pos));pos+=WStringFlowLen(flow,pos);}

#define RestoreFlow_Array(flow,pos,type,p,len) {memcpy(p,(BYTE*)(flow)+pos,(len)*sizeof(type));pos+=((len)*sizeof(type));}

#define SaveFlow_Simple(flow,pos,type,d) {*(type*)((BYTE*)(flow)+pos)=(d);pos+=sizeof(type);}

#define SaveFlow_SimplePointer(flow,pos,type,p) {*(type*)((BYTE*)(flow)+pos)=*(type*)(p);pos+=sizeof(type);}

#define SaveFlow_String(flow,pos,p) {memcpy((BYTE*)(flow)+pos,p,StringLen(p));pos+=StringLen(p);}

#define SaveFlow_WString(flow,pos,p) {memcpy((BYTE*)(flow)+pos,p,WStringLen(p));pos+=WStringLen(p);}

#define SaveFlow_Array(flow,pos,type,p,len) {memcpy((BYTE*)(flow)+pos,p,(len)*sizeof(type));pos+=((len)*sizeof(type));}

#define SetFlow_Zero(flow,pos,len) {memset((BYTE*)(flow)+pos,0,len);pos+=(len);}

#define SetFlow_Nega(flow,pos,len) {memset((BYTE*)(flow)+pos,-1,len);pos+=(len);}


typedef struct _CipherDescReportItem
{
	BYTE *  tagCipher;
	char * strCipherName;
	char * strCipherText;
}CipherDescReportItem,* PCipherDescReportItem;

typedef INT32 INSTANCE_ID;

#define INSTANCE_NULL  ((INSTANCE_ID)0)
#define Instance_Windows       0
#define Instance_Linux         1
#define Instance_Mac           2
#define Instance_IOS           3
#define Instance_Android       4
#define Instance_Type(n) ((UINT32)(n)>>24)

#define Instance_GetRealId(n) ((UINT32)(n)&0x00FFFFFF) 



#define LIST_VERB_DEFAULT    1

#define LIST_VERB_CLEAR      2

#define LIST_VERB_ADD        3

#define LIST_VERB_DELETE     4

#define LIST_VERB_MODIFY     5

#define PERMISSION_INFO_LENGTH     2


typedef struct _PermissionInfo
{
	UINT16    bOwner : 1;
	UINT16    bAdjust : 1;
}PermissionInfo, *PPermissionInfo;


#define SetPermission(a) *(UINT16*)(a)=(UINT16)-1;
#define ClearPermission(a) *(UINT16*)(a)=(UINT16)0;
#define CopyPermission(a,b) *(UINT16*)(a)=*(UINT16*)(b);
#define AndPermission(a,b) *(UINT16*)(a)=(*(UINT16*)(a)) & (*(UINT16*)(b))
#define OrPermission(a,b) *(UINT16*)(a)=(*(UINT16*)(a)) | (*(UINT16*)(b))

#endif
