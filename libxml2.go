package xsdvalidate

/*
#cgo pkg-config: libxml-2.0
#include <string.h>
#include <libxml/xmlschemastypes.h>
#include <errno.h>
#include <malloc.h>
#define GO_ERR_INIT 256
#define P_ERR_EXT 2
#define LIBXML_STATIC

struct xsdParserResult {
	xmlSchemaPtr schemaPtr;
	char *errorStr;
};

struct xmlParserResult {
	xmlDocPtr docPtr;
	char *errorStr;
} xmlParserResult;

struct errCtx {
	char *errBuf;
};

static void noOutputCallback(void *ctx, const char *message, ...) {
}

static void init() {
	xmlInitParser();
	xmlLineNumbersDefault(1);
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
		free(newLine);
		newLine = malloc(lineLen);
		vsnprintf(newLine, lineLen, message, varArgs);
	}
	va_end(varArgs);

	char *tmp = malloc(oldLen + lineLen);
	memcpy(tmp, ectx->errBuf, oldLen);
	strcat(tmp, newLine);
	free(newLine);
	free(ectx->errBuf);
	ectx->errBuf = tmp;
}

static struct xsdParserResult cParseUrlSchema(const char *url, const short int options) {
	struct xsdParserResult parserResult;
	char *errBuf=NULL;
	struct errCtx *ectx=malloc(sizeof(struct errCtx));
	ectx->errBuf=calloc(GO_ERR_INIT, sizeof(char));
	struct errCtx *genEctx=malloc(sizeof(struct errCtx));;
	genEctx->errBuf=calloc(GO_ERR_INIT, sizeof(char));

	xmlSchemaPtr schema = NULL;
	xmlSchemaParserCtxtPtr schemaParserCtxt = NULL;

	schemaParserCtxt = xmlSchemaNewParserCtxt(url);

	if (schemaParserCtxt == NULL) {
		errno = -1;
		strcpy(ectx->errBuf, "Xsd parser internal error");
	}
	else
	{
		if (options & P_ERR_EXT) {
			xmlSetGenericErrorFunc(genEctx, genErrorCallback);
		} else {
			xmlSetGenericErrorFunc(NULL, noOutputCallback);
		}

		xmlSchemaSetParserErrors(schemaParserCtxt, genErrorCallback, noOutputCallback, ectx);

		schema = xmlSchemaParse(schemaParserCtxt);
		errno = 0;

		xmlSchemaFreeParserCtxt(schemaParserCtxt);
		if (schema == NULL) {
			errno = -1;
			char *tmp = NULL;
			if (options & P_ERR_EXT) {
				tmp = (char *) malloc(strlen(ectx->errBuf) + strlen(genEctx->errBuf) + 1);
				memcpy(tmp, ectx->errBuf, strlen(ectx->errBuf) + 1);
				strcat(tmp, genEctx->errBuf);
			} else {
				tmp = (char *) malloc(strlen(ectx->errBuf) + 1);
				memcpy(tmp, ectx->errBuf, strlen(ectx->errBuf) + 1);
			}
			free(ectx->errBuf);
			ectx->errBuf = tmp;
		}
	}
	errBuf=malloc(strlen(ectx->errBuf)+1);
	memcpy(errBuf,  ectx->errBuf, strlen(ectx->errBuf)+1);

	free(ectx->errBuf);
	free(ectx);
	free(genEctx->errBuf);
	free(genEctx);
	parserResult.schemaPtr=schema;
	parserResult.errorStr=errBuf;
	return parserResult;
}

static struct xmlParserResult cParseDoc(const char *goXmlSource, const int goXmlSourceLen, const short int options) {
	struct xmlParserResult parserResult;
	char *errBuf=NULL;
	struct errCtx *ectx=malloc(sizeof(struct errCtx));
	ectx->errBuf=calloc(GO_ERR_INIT, sizeof(char));;


	xmlDocPtr doc=NULL;
	xmlParserCtxtPtr xmlParserCtxt=NULL;

	xmlParserCtxt = xmlNewParserCtxt();

	if (xmlParserCtxt == NULL) {
		errno = -1;
		strcpy(ectx->errBuf, "Xml parser internal error");
	}
	else
	{
		if (options & P_ERR_EXT) {
			xmlSetGenericErrorFunc(ectx, genErrorCallback);
		} else {
			xmlSetGenericErrorFunc(NULL, noOutputCallback);
		}

		doc = xmlParseMemory(goXmlSource, goXmlSourceLen);

		xmlFreeParserCtxt(xmlParserCtxt);

		if (doc == NULL) {
			if (options & P_ERR_EXT) {
				errno = -1;
				char *tmp = malloc(strlen(ectx->errBuf) + 1);
				memcpy(tmp, ectx->errBuf, strlen(ectx->errBuf) + 1);
				free(ectx->errBuf);
				ectx->errBuf = tmp;
			} else {
				errno = -1;
				strcpy(ectx->errBuf, "Malformed xml document");
			}
		}
	}

	errBuf=malloc(strlen(ectx->errBuf)+1);
	memcpy(errBuf,  ectx->errBuf, strlen(ectx->errBuf)+1);
	free(ectx->errBuf);
	free(ectx);
	parserResult.docPtr=doc;
	parserResult.errorStr=errBuf;
	return parserResult;
}

static char *cValidate(const xmlDocPtr doc, const xmlSchemaPtr schema) {
	char *errBuf=NULL;
	struct errCtx *ectx=malloc(sizeof(struct errCtx));
	ectx->errBuf=calloc(GO_ERR_INIT, sizeof(char));
	int schemaErr=0;

	if (schema == NULL) {
		errno = -1;
		strcpy(ectx->errBuf, "Xsd schema null pointer");
	}
	else if (doc == NULL) {
		errno = -1;
		strcpy(ectx->errBuf, "Xml schema null pointer");
	}
	else
	{
		xmlSchemaValidCtxtPtr schemaCtxt;
		schemaCtxt = xmlSchemaNewValidCtxt(schema);

		if (schemaCtxt == NULL) {
			errno = -1;
			strcpy(ectx->errBuf, "Xml validation internal error");
		}
		else
		{

			xmlSchemaSetValidErrors(schemaCtxt, genErrorCallback, noOutputCallback, ectx);
			schemaErr = xmlSchemaValidateDoc(schemaCtxt, doc);

			xmlSchemaFreeValidCtxt(schemaCtxt);

			if (schemaErr == 0)
			{
				errno = 0;
			}
			else if (schemaErr > 0)
			{
				errno = -1;
			}
			else
			{
				errno = -1;
				strcpy(ectx->errBuf, "Xml validation internal error");
			}
		}
	}

	errBuf=malloc(strlen(ectx->errBuf)+1);
	memcpy(errBuf,  ectx->errBuf, strlen(ectx->errBuf)+1);
	free(ectx->errBuf);
	free(ectx);
	return errBuf;
}
*/
import "C"
import (
	"errors"
	"strings"
	"unsafe"
)

// Handles schema parsing and validation, wraps a pointer to libxml2's xmlSchemaPtr.
type XsdHandler struct {
	schemaPtr C.xmlSchemaPtr
}

// Handles xml parsing, wraps a pointer to libxml2's xmlDocPtr.
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

	strXml := C.CString(string(inXml))
	defer C.free(unsafe.Pointer(strXml))
	pRes, err := C.cParseDoc(strXml, C.int(len(inXml)), C.short(options))

	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, errors.New("Xml parser error:\n" + strings.Trim(rStr, "\n"))
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
		return nil, errors.New("Xsd parser error:\n" + strings.Trim(rStr, "\n"))
	}
	return pRes.schemaPtr, nil
}

// Helper function for validating given an xml document
func validateWithXsd(xmlHandler *XmlHandler, xsdHandler *XsdHandler) error {
	defer C.malloc_trim(0)

	sPtr, err := C.cValidate(xmlHandler.docPtr, xsdHandler.schemaPtr)
	defer C.free(unsafe.Pointer(sPtr))
	if err != nil {
		rStr := C.GoString(sPtr)
		return errors.New("Validation error:\n" + strings.Trim(rStr, "\n"))
	}
	return nil
}

//Wrapper for the xmlSchemaFree function
func freeSchemaPtr(xsdHandler *XsdHandler) {
	if xsdHandler.schemaPtr != nil {
		C.xmlSchemaFree(xsdHandler.schemaPtr)
	}
}

//Wrapper for the xmlFreeDoc function
func freeDocPtr(xmlHandler *XmlHandler) {
	if xmlHandler.docPtr != nil {
		C.xmlFreeDoc(xmlHandler.docPtr)
	}
}
