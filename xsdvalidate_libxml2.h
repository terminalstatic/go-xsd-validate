#ifndef XSDVALIDATE_LIBXML2_H
#define XSDVALIDATE_LIBXML2_H

#include <stddef.h>
#include <stdlib.h>
#include <libxml/xmlschemastypes.h>

#define GO_ERR_INIT 1024
#define P_ERR_DEFAULT 1
#define P_ERR_VERBOSE 2

struct xsdParserResult {
    xmlSchemaPtr schemaPtr;
    char* errorStr;
};

struct xmlParserResult {
    xmlDocPtr docPtr;
    char* errorStr;
};

typedef enum {
    NO_ERROR = 0,
    LIBXML2_ERROR = 1,
    XSD_PARSER_ERROR = 2,
    XML_PARSER_ERROR = 3,
    VALIDATION_ERROR = 4
} errorType;

struct simpleXmlError {
    errorType type;
    int code;
    char* message;
    int level;
    int line;
    char* node;
    char* nodePath;
};

typedef struct _errArray {
    struct simpleXmlError* data;
    size_t len;
    size_t cap;
} errArray;

void xsdValidateInit(void);
void xsdValidateCleanup(void);
void xsdValidateFreeErrArray(errArray* errArr);

struct xsdParserResult xsdValidateParseUrlSchema(const char* url, const short int options);
struct xsdParserResult xsdValidateParseMemSchema(const void* xsd, const int goXsdSourceLen, const short int options);
struct xmlParserResult xsdValidateParseDoc(const void* goXmlSource, const int goXmlSourceLen, const short int options);
errArray xsdValidateDoc(const xmlDocPtr doc, const xmlSchemaPtr schema);
errArray xsdValidateBuf(const void* goXmlSource, const int goXmlSourceLen, const short int xmlParserOptions, const xmlSchemaPtr schema);

#endif
