#ifndef _TSK_DEFINE_H
#define _TSK_DEFINE_H


#define __mem_alloc(MemType,MemLen,MemTag)  (malloc(MemLen))
#define __mem_free(MemPointer) free(MemPointer)
#define __mem_realloc(p,Type,Len,Tag) ((p)?__mem_free((p)):0,__mem_alloc((Type),(Len),(Tag)))

#endif