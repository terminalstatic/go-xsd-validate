package xsdvalidate

/*
#cgo CFLAGS: -std=c99
#cgo pkg-config: libxml-2.0
#include <string.h>
#include <libxml/xmlschemastypes.h>
#include <errno.h>
#include <stdlib.h>
#include <stdbool.h>
#define GO_ERR_INIT 1024
#define P_ERR_DEFAULT 1
#define P_ERR_VERBOSE 2
#define LIBXML_STATIC

struct xsdParserResult {
	xmlSchemaPtr schemaPtr;
	char *errorStr;
};

struct xmlParserResult {
	xmlDocPtr docPtr;
	char *errorStr;
};

struct errCtx {
	char *errBuf;
};

typedef enum {
	NO_ERROR 		= 0,
	LIBXML2_ERROR		= 1,
	XSD_PARSER_ERROR	= 2,
	XML_PARSER_ERROR 	= 3,
	VALIDATION_ERROR 	= 4
} errorType;

struct simpleXmlError {
	errorType 	type;
	int		code;
	char*		message;
	int 		level;
	int		line;
	char*		node;
};

typedef struct _errArray {
	struct 	simpleXmlError *data;
	size_t 	len;
	size_t 	cap;
} errArray;

static void freeErrArray(errArray *errArr) {
	for (int i=0;i<errArr->len;i++) {
		free(errArr->data[i].message);
		free(errArr->data[i].node);

	}
	free(errArr->data);
	free(errArr);
}

static void noOutputCallback(void *ctx, const char *message, ...) {}

static void init() {
	xmlInitParser();
}

static void cleanup() {
	xmlSchemaCleanupTypes();
	xmlCleanupParser();
}

static void genErrorCallback(void *ctx, const char *message, ...) {
	struct errCtx *ectx = ctx;
	char *newLine = malloc(GO_ERR_INIT);

	va_list varArgs;
	va_start(varArgs, message);

	int oldLen = strlen(ectx->errBuf) + 1;
	int lineLen = 1 + vsnprintf(newLine, GO_ERR_INIT, message, varArgs);

	if (lineLen  > GO_ERR_INIT) {
		va_end(varArgs);
		va_start(varArgs, message);
		free(newLine);
		newLine = malloc(lineLen);
		vsnprintf(newLine, lineLen, message, varArgs);
		va_end(varArgs);
	} else {
		va_end(varArgs);
	}

	char *tmp = malloc(oldLen + lineLen);
	memcpy(tmp, ectx->errBuf, oldLen);
	strcat(tmp, newLine);
	free(newLine);
	free(ectx->errBuf);
	ectx->errBuf = tmp;
}

static void simpleStructErrorCallback(void *ctx, xmlErrorPtr p) {
	errArray *sErrArr = ctx;

	struct simpleXmlError sErr;
	sErr.message = calloc(GO_ERR_INIT, sizeof(char));
	sErr.node = calloc(GO_ERR_INIT, sizeof(char));

	sErr.type = VALIDATION_ERROR;
	sErr.code = p->code;
	sErr.level = p->level;
	sErr.line = p->line;


	int cpyLen = 1 + snprintf(sErr.message, GO_ERR_INIT, "%s", p->message);
	if (cpyLen > GO_ERR_INIT) {
		free(sErr.message);
		sErr.message = malloc(cpyLen);
		snprintf(sErr.message, cpyLen, "%s", p->message);
	}

	if (p->node !=NULL) {
		cpyLen = 1 + snprintf(sErr.node, GO_ERR_INIT, "%s", (((xmlNodePtr) p->node)->name));
		if (cpyLen > GO_ERR_INIT) {
			free(sErr.node);
			sErr.node= malloc(cpyLen);
			snprintf(sErr.node, cpyLen, "%s", (((xmlNodePtr) p->node)->name));
		}
	}
	if (sErrArr->len >= sErrArr->cap) {
		sErrArr->cap = sErrArr->cap * 2;
		struct simpleXmlError *tmp = calloc(sErrArr->cap, sizeof(*tmp));
		memcpy(tmp, sErrArr->data, sErrArr->len * sizeof(*tmp));
		free(sErrArr->data);
		sErrArr->data = tmp;
	}
	sErrArr->data[sErrArr->len] = sErr;
	sErrArr->len++;
}

static struct xsdParserResult cParseUrlSchema(const char *url, const short int options) {
	xmlLineNumbersDefault(1);
	bool err = false;
	struct xsdParserResult parserResult;
	char *errBuf=NULL;
	struct errCtx ectx;
	ectx.errBuf=calloc(GO_ERR_INIT, sizeof(char));

	xmlSchemaPtr schema = NULL;
	xmlSchemaParserCtxtPtr schemaParserCtxt = NULL;

	schemaParserCtxt = xmlSchemaNewParserCtxt(url);

	if (schemaParserCtxt == NULL) {
		err = true;
		strcpy(ectx.errBuf, "Xsd parser internal error");
	}
	else
	{
		if (options & P_ERR_VERBOSE) {
			xmlSchemaSetParserErrors(schemaParserCtxt, noOutputCallback, noOutputCallback, NULL);
			xmlSetGenericErrorFunc(&ectx, genErrorCallback);
		} else {
			xmlSetGenericErrorFunc(NULL, noOutputCallback);
			xmlSchemaSetParserErrors(schemaParserCtxt, genErrorCallback, noOutputCallback, &ectx);
		}

		schema = xmlSchemaParse(schemaParserCtxt);

		xmlSchemaFreeParserCtxt(schemaParserCtxt);
		if (schema == NULL) {
			err = true;
			char *tmp = malloc(strlen(ectx.errBuf) + 1);
			memcpy(tmp, ectx.errBuf, strlen(ectx.errBuf) + 1);
			free(ectx.errBuf);
			ectx.errBuf = tmp;
		}
	}

	errBuf=malloc(strlen(ectx.errBuf)+1);
	memcpy(errBuf,  ectx.errBuf, strlen(ectx.errBuf)+1);
	free(ectx.errBuf);
	parserResult.schemaPtr=schema;
	parserResult.errorStr=errBuf;
	errno = err ? -1 : 0;
	return parserResult;
}

static struct xmlParserResult cParseDoc(const void *goXmlSource, const int goXmlSourceLen, const short int options) {
	xmlLineNumbersDefault(1);
	bool err = false;
	struct xmlParserResult parserResult;
	char *errBuf=NULL;
	struct errCtx ectx;
	ectx.errBuf=calloc(GO_ERR_INIT, sizeof(char));;

	xmlDocPtr doc=NULL;
	xmlParserCtxtPtr xmlParserCtxt=NULL;

	if (goXmlSourceLen == 0) {
		err = true;
		if (options & P_ERR_VERBOSE) {
			strcpy(ectx.errBuf, "parser error : Document is empty");
		} else {
			strcpy(ectx.errBuf, "Malformed xml document");
		}
	} else {
		xmlParserCtxt = xmlNewParserCtxt();

		if (xmlParserCtxt == NULL) {
			err = true;
			strcpy(ectx.errBuf, "Xml parser internal error");
		}
		else
		{
			if (options & P_ERR_VERBOSE) {
				xmlSetGenericErrorFunc(&ectx, genErrorCallback);
			} else {
				xmlSetGenericErrorFunc(NULL, noOutputCallback);
			}

			doc = xmlParseMemory(goXmlSource, goXmlSourceLen);

			xmlFreeParserCtxt(xmlParserCtxt);
			if (doc == NULL) {
				err = true;
				if (options & P_ERR_VERBOSE) {
					char *tmp = malloc(strlen(ectx.errBuf) + 1);
					memcpy(tmp, ectx.errBuf, strlen(ectx.errBuf) + 1);
					free(ectx.errBuf);
					ectx.errBuf = tmp;
				} else {
					strcpy(ectx.errBuf, "Malformed xml document");
				}
			}
		}
	}

	errBuf=malloc(strlen(ectx.errBuf)+1);
	memcpy(errBuf,  ectx.errBuf, strlen(ectx.errBuf)+1);
	free(ectx.errBuf);
	parserResult.docPtr=doc;
	parserResult.errorStr=errBuf;
	errno = err ? -1 : 0;
	return parserResult;
}

static errArray *cValidate(const xmlDocPtr doc, const xmlSchemaPtr schema) {
	xmlLineNumbersDefault(1);

	errArray *errArr = malloc(sizeof(*errArr));
	errArr->data = calloc(2, sizeof(struct simpleXmlError));
	errArr->len=0;
	errArr->cap=2;

	struct simpleXmlError simpleError;
	simpleError.message = calloc(GO_ERR_INIT, sizeof(char));
	simpleError.node = calloc(GO_ERR_INIT, sizeof(char));


	if (schema == NULL) {
		simpleError.type = LIBXML2_ERROR;
		strcpy(simpleError.message, "Xsd schema null pointer");
		errArr->data[errArr->len] = simpleError;
		errArr->len++;
	}
	else if (doc == NULL) {
		simpleError.type = LIBXML2_ERROR;
		errArr->len++;
	}
	else
	{
		xmlSchemaValidCtxtPtr schemaCtxt;
		schemaCtxt = xmlSchemaNewValidCtxt(schema);

		if (schemaCtxt == NULL) {
			simpleError.type = LIBXML2_ERROR;
			strcpy(simpleError.message, "Xml validation internal error");
			errArr->data[errArr->len] = simpleError;
			errArr->len++;
		}
		else
		{
			xmlSchemaSetValidStructuredErrors(schemaCtxt, simpleStructErrorCallback, errArr);
			int schemaErr = xmlSchemaValidateDoc(schemaCtxt, doc);
			xmlSchemaFreeValidCtxt(schemaCtxt);

			if (schemaErr < 0 && errArr->len == 0)
			{
				simpleError.type = LIBXML2_ERROR;
				strcpy(simpleError.message, "Xml validation internal error");
				errArr->data[errArr->len] = simpleError;
				errArr->len++;
			}
			else {
				simpleError.type = NO_ERROR;
			}
		}
	}

	free(simpleError.node);
	free(simpleError.message);

	errno = errArr->len == NO_ERROR ? 0 : -1;
	return errArr;
}

static errArray *cValidateBuf(const void *goXmlSource, const int goXmlSourceLen, const short int xmlParserOptions, const xmlSchemaPtr schema) {
	xmlLineNumbersDefault(1);

	errArray *errArr = malloc(sizeof(*errArr));
	errArr->data = calloc(2, sizeof(struct simpleXmlError));
	errArr->len=0;
	errArr->cap=2;

	struct simpleXmlError simpleError;
	simpleError.message = calloc(GO_ERR_INIT, sizeof(char));

	struct xmlParserResult parserResult = cParseDoc(goXmlSource, goXmlSourceLen, xmlParserOptions);

	if (schema == NULL) {
		simpleError.type = LIBXML2_ERROR;
		strcpy(simpleError.message, "Xsd schema null pointer");
		errArr->data[errArr->len] = simpleError;
		errArr->len++;
		xmlFreeDoc(parserResult.docPtr);
	}
	else if (parserResult.docPtr == NULL) {
		simpleError.type = XML_PARSER_ERROR;
		free(simpleError.message);
		simpleError.message = malloc(strlen(parserResult.errorStr) + 1);
		strcpy(simpleError.message, parserResult.errorStr);
		errArr->data[errArr->len] = simpleError;
		errArr->len++;
	}
	else
	{
		free(simpleError.message);
		free(errArr->data);
		free(errArr);

		errArr = cValidate(parserResult.docPtr, schema);
		xmlFreeDoc(parserResult.docPtr);
	}

	free(parserResult.errorStr);

	errno = errArr->len == NO_ERROR ? 0 : -1;
	return errArr;
}
*/
import "C"
import (
	"runtime"
	"strings"
	"time"
	"unsafe"
)

// XsdHandler handles schema parsing and validation and wraps a pointer to libxml2's xmlSchemaPtr.
type XsdHandler struct {
	schemaPtr C.xmlSchemaPtr
}

// XmlHandler handles xml parsing and wraps a pointer to libxml2's xmlDocPtr.
type XmlHandler struct {
	docPtr C.xmlDocPtr
}

// Initializes the libxml2 parser, suggested for multithreading
func libXml2Init() {
	C.init()
}

// Cleans up the libxml2 parser
func libXml2Cleanup() {
	C.cleanup()
}

// The helper function for parsing xml
func parseXmlMem(inXml []byte, options Options) (C.xmlDocPtr, error) {
	strXml := C.CBytes(inXml)
	defer C.free(unsafe.Pointer(strXml))
	pRes, err := C.cParseDoc(strXml, C.int(len(inXml)), C.short(options))

	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, XmlParserError{errorMessage{strings.Trim(rStr, "\n")}}
	}
	return pRes.docPtr, nil
}

// The helper function for parsing the schema
func parseUrlSchema(url string, options Options) (C.xmlSchemaPtr, error) {
	strUrl := C.CString(url)
	defer C.free(unsafe.Pointer(strUrl))

	pRes, err := C.cParseUrlSchema(strUrl, C.short(options))
	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, XsdParserError{errorMessage{strings.Trim(rStr, "\n")}}
	}
	return pRes.schemaPtr, nil
}

func handleErrArray(errSlice []C.struct_simpleXmlError) ValidationError {
	ve := ValidationError{make([]StructError, len(errSlice))}
	for i := 0; i < len(errSlice); i++ {
		ve.Errors[i] = StructError{
			Code:     int(errSlice[i].code),
			Message:  strings.Trim(C.GoString(errSlice[i].message), "\n"),
			Level:    int(errSlice[i].level),
			Line:     int(errSlice[i].line),
			NodeName: C.GoString(errSlice[i].node)}
	}
	return ve

}

// Helper function for validating given an xml document
func validateWithXsd(xmlHandler *XmlHandler, xsdHandler *XsdHandler) error {
	sErr, err := C.cValidate(xmlHandler.docPtr, xsdHandler.schemaPtr)
	defer C.freeErrArray(sErr)
	if err != nil {
		errSlice := (*[1 << 30]C.struct_simpleXmlError)(unsafe.Pointer(sErr.data))[:sErr.len:sErr.len]
		return handleErrArray(errSlice)
	}
	return nil
}

// Helper function for validating given an xml byte slice
func validateBufWithXsd(inXml []byte, options Options, xsdHandler *XsdHandler) error {
	strXml := C.CBytes(inXml)
	defer C.free(unsafe.Pointer(strXml))
	sErr, err := C.cValidateBuf(strXml, C.int(len(inXml)), C.short(options), xsdHandler.schemaPtr)
	defer C.freeErrArray(sErr)
	if err != nil {
		errSlice := (*[1 << 30]C.struct_simpleXmlError)(unsafe.Pointer(sErr.data))[:sErr.len:sErr.len]
		switch errSlice[0]._type {
		case C.VALIDATION_ERROR:
			return handleErrArray(errSlice)
		case C.XML_PARSER_ERROR:
			return XmlParserError{errorMessage{strings.Trim(C.GoString(errSlice[0].message), "\n")}}
		case C.LIBXML2_ERROR:
			return Libxml2Error{errorMessage{strings.Trim(C.GoString(errSlice[0].message), "\n")}}
		case C.XSD_PARSER_ERROR:
			return XsdParserError{errorMessage{strings.Trim(C.GoString(errSlice[0].message), "\n")}}
		default:
			return Libxml2Error{errorMessage{"Unknown error"}}
		}
		return ValidationError{}
	}
	return nil
}

// Wrapper for the xmlSchemaFree function
func freeSchemaPtr(xsdHandler *XsdHandler) {
	if xsdHandler.schemaPtr != nil {
		C.xmlSchemaFree(xsdHandler.schemaPtr)
	}
}

// Wrapper for the xmlFreeDoc function
func freeDocPtr(xmlHandler *XmlHandler) {
	if xmlHandler.docPtr != nil {
		C.xmlFreeDoc(xmlHandler.docPtr)
	}
}

// Ticker for gc
func gcTicker(d time.Duration, quit chan struct{}) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			runtime.GC()
			//C.malloc_trim(0)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
