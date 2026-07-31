#include "xsdvalidate_libxml2.h"

#include <errno.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <sys/time.h>
#include <libxml/parser.h>

#define LIBXML_STATIC

typedef struct _errCtx {
    char* errBuf;
    size_t len;
    size_t cap;
} errCtx;

static errArray initErrArray(void) {
    errArray errArr = {
        .data = calloc(2, sizeof(struct simpleXmlError)), .len = 0, .cap = 2};
    return errArr;
}

static char* copyStringOrEmpty(const char* src) {
    if (src == NULL) {
        src = "";
    }

    size_t len = strlen(src) + 1;
    char* dst = malloc(len);
    if (dst == NULL) {
        return NULL;
    }

    memcpy(dst, src, len);
    return dst;
}

static void freeSimpleXmlError(struct simpleXmlError* sErr) {
    free(sErr->message);
    free(sErr->node);
    free(sErr->nodePath);
}

void xsdValidateFreeErrArray(errArray* errArr) {
    for (int i = 0; i < errArr->len; i++) {
        freeSimpleXmlError(&errArr->data[i]);
    }
    free(errArr->data);
}

static errCtx initErrCtx(size_t len, size_t cap) {
    errCtx ectx = {.errBuf = malloc(cap), .len = len, .cap = cap};
    memset(ectx.errBuf, '\0', len);
    return ectx;
}

static void freeErrCtx(errCtx ectx) {
    free(ectx.errBuf);
    ectx.len=0;
    ectx.cap=0;
}

static void appendErrCtxErrBuff(errCtx* ectx, const char* buffStr) {
    size_t buffStrLen = strlen(buffStr);
    size_t capWanted = ectx->len + buffStrLen;

    if (capWanted > ectx->cap) {
        size_t newCap = capWanted + GO_ERR_INIT;
        char* tmp = malloc(newCap);
        memcpy(tmp, ectx->errBuf, ectx->len);
        free(ectx->errBuf);
        ectx->errBuf = tmp;
        ectx->cap = newCap;
    }

    size_t newLen = ectx->len + buffStrLen;
    char* tmp = malloc(newLen);
    if (ectx->len > 1) {
        memcpy(tmp, ectx->errBuf, ectx->len);
    }
    memcpy(&tmp[ectx->len - 1], buffStr, buffStrLen + 1);
    free(ectx->errBuf);
    ectx->errBuf = tmp;
    ectx->len = newLen;
}

static void noOutputCallback(void* ctx, const char* message, ...) {}

void xsdValidateInit(void) {
    xmlInitParser();
}

void xsdValidateCleanup(void) {
#if LIBXML_VERSION < 21000
    xmlSchemaCleanupTypes();
#endif
    xmlCleanupParser();
}

static void genErrorCallback(void* ctx, const char* message, ...) {
    errCtx* ectx = ctx;
    char* newLine = malloc(GO_ERR_INIT);

    va_list varArgs;
    va_start(varArgs, message);

    size_t lineLen = 1 + vsnprintf(newLine, GO_ERR_INIT, message, varArgs);

    if (lineLen > GO_ERR_INIT) {
        va_end(varArgs);
        va_start(varArgs, message);
        free(newLine);
        newLine = malloc(lineLen);
        vsnprintf(newLine, lineLen, message, varArgs);
        va_end(varArgs);
    } else {
        va_end(varArgs);
    }

    appendErrCtxErrBuff(ectx, newLine);
    free(newLine);
}

static void simpleStructErrorCallback(
    void* ctx,
#if LIBXML_VERSION >= 21200
    const xmlError *p
#else
    xmlErrorPtr p
#endif
) {
    if (p == NULL) {
        return;
    }

    errArray* sErrArr = ctx;
    if (sErrArr == NULL) {
        return;
    }

    struct simpleXmlError sErr = {0};
    sErr.message = copyStringOrEmpty(p->message);
    sErr.node = copyStringOrEmpty(NULL);
    sErr.nodePath = copyStringOrEmpty(NULL);
    if (sErr.message == NULL || sErr.node == NULL || sErr.nodePath == NULL) {
        freeSimpleXmlError(&sErr);
        return;
    }

    sErr.type = VALIDATION_ERROR;
    sErr.code = p->code;
    sErr.level = p->level;
    sErr.line = p->line;

    if (p->node != NULL) {
        const xmlChar* nodeName = ((xmlNodePtr)p->node)->name;
        if (nodeName != NULL) {
            char* node = copyStringOrEmpty((const char*)nodeName);
            if (node == NULL) {
                freeSimpleXmlError(&sErr);
                return;
            }
            free(sErr.node);
            sErr.node = node;
        }

        xmlChar* nodePath = xmlGetNodePath((xmlNodePtr)p->node);
        if (nodePath != NULL) {
            char* nodePathCopy = copyStringOrEmpty((const char*)nodePath);
            xmlFree(nodePath);
            if (nodePathCopy == NULL) {
                freeSimpleXmlError(&sErr);
                return;
            }
            free(sErr.nodePath);
            sErr.nodePath = nodePathCopy;
        }
    }
    if (sErrArr->len >= sErrArr->cap) {
        size_t newCap = sErrArr->cap * 2;
        if (newCap <= sErrArr->cap) {
            freeSimpleXmlError(&sErr);
            return;
        }

        struct simpleXmlError* tmp = calloc(newCap, sizeof(*tmp));
        if (tmp == NULL) {
            freeSimpleXmlError(&sErr);
            return;
        }
        memcpy(tmp, sErrArr->data, sErrArr->len * sizeof(*tmp));
        free(sErrArr->data);
        sErrArr->data = tmp;
        sErrArr->cap = newCap;
    }
    sErrArr->data[sErrArr->len] = sErr;
    sErrArr->len++;
}

static struct xsdParserResult parseSchema(
                                           xmlSchemaParserCtxtPtr schemaParserCtxt,
                                           const short int options) {
    bool err = false;
    struct xsdParserResult parserResult;
    errCtx ectx = initErrCtx(1, GO_ERR_INIT);
    errCtx ectxParse = initErrCtx(1, GO_ERR_INIT);

    xmlSchemaPtr schema = NULL;

    if (schemaParserCtxt == NULL) {
        err = true;
        const char msg[] = "Xsd parser internal error";
        freeErrCtx(ectxParse);
        appendErrCtxErrBuff(&ectx, msg);
    } else {
        if (options & P_ERR_VERBOSE) {
            xmlSchemaSetParserErrors(schemaParserCtxt, genErrorCallback, noOutputCallback, &ectxParse);
            xmlSetGenericErrorFunc(&ectx, genErrorCallback);
        } else {
            xmlSetGenericErrorFunc(NULL, noOutputCallback);
            xmlSchemaSetParserErrors(schemaParserCtxt, genErrorCallback, noOutputCallback, &ectx);
        }

        schema = xmlSchemaParse(schemaParserCtxt);

        xmlSchemaFreeParserCtxt(schemaParserCtxt);
        if (schema == NULL) {
            freeErrCtx(ectx);
            ectx = ectxParse;
            err = true;
        } else {
            freeErrCtx(ectxParse);
        }
    }

    parserResult.errorStr = malloc(ectx.len);
    memcpy(parserResult.errorStr, ectx.errBuf, ectx.len);
    freeErrCtx(ectx);
    parserResult.schemaPtr = schema;
    errno = err ? -1 : 0;
    return parserResult;
}

struct xsdParserResult xsdValidateParseUrlSchema(const char* url,
                                                const short int options) {
    xmlSchemaParserCtxtPtr schemaParserCtxt = NULL;
    schemaParserCtxt = xmlSchemaNewParserCtxt(url);
    return parseSchema(schemaParserCtxt, options);
}

struct xsdParserResult xsdValidateParseMemSchema(const void* xsd,
                                                const int goXsdSourceLen,
                                                const short int options) {
    xmlSchemaParserCtxtPtr schemaParserCtxt = NULL;
    schemaParserCtxt = xmlSchemaNewMemParserCtxt(xsd, goXsdSourceLen);

    return parseSchema(schemaParserCtxt, options);
}

struct xmlParserResult xsdValidateParseDoc(const void* goXmlSource,
                                          const int goXmlSourceLen,
                                          const short int options) {
    bool err = false;
    struct xmlParserResult parserResult;
    errCtx ectx = initErrCtx(1, GO_ERR_INIT);

    xmlDocPtr doc = NULL;
    xmlParserCtxtPtr xmlParserCtxt = NULL;

    if (goXmlSourceLen == 0) {
        err = true;
        if (options & P_ERR_VERBOSE) {
            const char msg[] = "parser error : Document is empty";
            appendErrCtxErrBuff(&ectx, msg);
        } else {
            const char msg[] = "Malformed xml document";
            appendErrCtxErrBuff(&ectx, msg);
        }
    } else {
        xmlParserCtxt = xmlNewParserCtxt();

        if (xmlParserCtxt == NULL) {
            err = true;
            const char msg[] = "Xml parser internal error";
            appendErrCtxErrBuff(&ectx, msg);
        } else {
            if (options & P_ERR_VERBOSE) {
                xmlSetGenericErrorFunc(&ectx, genErrorCallback);
            } else {
                xmlSetGenericErrorFunc(NULL, noOutputCallback);
            }

            doc = xmlReadMemory(goXmlSource, goXmlSourceLen, NULL, NULL, 0);

            xmlFreeParserCtxt(xmlParserCtxt);
            if (doc == NULL) {
                err = true;
                if (!(options & P_ERR_VERBOSE)) {
                    const char msg[] = "Malformed xml document";
                    appendErrCtxErrBuff(&ectx, msg);
                }
            }
        }
    }

    parserResult.errorStr = malloc(ectx.len);
    memcpy(parserResult.errorStr, ectx.errBuf, ectx.len);
    freeErrCtx(ectx);
    parserResult.docPtr = doc;
    errno = err ? -1 : 0;
    return parserResult;
}

errArray xsdValidateDoc(const xmlDocPtr doc, const xmlSchemaPtr schema) {
    errArray errArr = initErrArray();

    struct simpleXmlError simpleError;
    simpleError.message = calloc(GO_ERR_INIT, sizeof(char));
    simpleError.node = calloc(GO_ERR_INIT, sizeof(char));
    simpleError.nodePath = calloc(GO_ERR_INIT, sizeof(char));

    if (schema == NULL) {
        simpleError.type = LIBXML2_ERROR;
        strcpy(simpleError.message, "Xsd schema null pointer");
        errArr.data[errArr.len] = simpleError;
        errArr.len++;
    } else if (doc == NULL) {
        simpleError.type = LIBXML2_ERROR;
        strcpy(simpleError.message, "Xml doc null pointer");
        errArr.data[errArr.len] = simpleError;
        errArr.len++;
    } else {
        xmlSchemaValidCtxtPtr schemaCtxt;
        schemaCtxt = xmlSchemaNewValidCtxt(schema);

        if (schemaCtxt == NULL) {
            simpleError.type = LIBXML2_ERROR;
            strcpy(simpleError.message, "Xml validation internal error");
            errArr.data[errArr.len] = simpleError;
            errArr.len++;
        } else {
            xmlSchemaSetValidStructuredErrors(schemaCtxt, simpleStructErrorCallback,
                                              &errArr);
            int schemaErr = xmlSchemaValidateDoc(schemaCtxt, doc);
            xmlSchemaFreeValidCtxt(schemaCtxt);

            if (schemaErr < 0 && errArr.len == 0) {
                simpleError.type = LIBXML2_ERROR;
                strcpy(simpleError.message, "Xml validation internal error");
                errArr.data[errArr.len] = simpleError;
                errArr.len++;
            } else {
                freeSimpleXmlError(&simpleError);
            }
        }
    }

    errno = errArr.len == NO_ERROR ? 0 : -1;
    return errArr;
}

errArray xsdValidateBuf(const void* goXmlSource,
                        const int goXmlSourceLen,
                        const short int xmlParserOptions,
                        const xmlSchemaPtr schema) {
    errArray errArr = initErrArray();

    struct simpleXmlError simpleError;
    simpleError.message = calloc(GO_ERR_INIT, sizeof(char));
    simpleError.node = calloc(GO_ERR_INIT, sizeof(char));
    simpleError.nodePath = calloc(GO_ERR_INIT, sizeof(char));

    struct xmlParserResult parserResult =
    xsdValidateParseDoc(goXmlSource, goXmlSourceLen, xmlParserOptions);

    if (schema == NULL) {
        simpleError.type = LIBXML2_ERROR;
        const char msg[] = "Xsd schema null pointer";
        strcpy(simpleError.message, msg);
        errArr.data[errArr.len] = simpleError;
        errArr.len++;

        xmlFreeDoc(parserResult.docPtr);
        free(parserResult.errorStr);
        errno = -1;
        return errArr;
    } else if (parserResult.docPtr == NULL) {
        simpleError.type = XML_PARSER_ERROR;
        free(simpleError.message);
        simpleError.message = malloc(strlen(parserResult.errorStr) + 1);
        strcpy(simpleError.message, parserResult.errorStr);
        errArr.data[errArr.len] = simpleError;
        errArr.len++;

        xmlFreeDoc(parserResult.docPtr);
        free(parserResult.errorStr);
        errno = -1;
        return errArr;
    }
    freeSimpleXmlError(&simpleError);
    xsdValidateFreeErrArray(&errArr);
    free(parserResult.errorStr);

    errArray valErrArr = xsdValidateDoc(parserResult.docPtr, schema);

    xmlFreeDoc(parserResult.docPtr);

    errno = valErrArr.len == NO_ERROR ? 0 : -1;
    return valErrArr;
}
