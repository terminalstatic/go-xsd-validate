package xsdvalidate

/*
#cgo pkg-config: libxml-2.0
#include <string.h>
#include <libxml/xmlschemastypes.h>
#include <errno.h>
#include <malloc.h>
#define GO_ERR_INIT 256
#define LIBXML_STATIC

typedef struct {
	xmlSchemaPtr schemaPtr;
	char *errorStr;
} xsdParserResult;

typedef struct {
	xmlDocPtr docPtr;
	char *errorStr;
} xmlParserResult;

typedef struct {
	char *errBuf;
} errCtx;

void noOutputCallback(void *ctx, const char *message, ...) {
}

void genErrorCallback(void *ctx, const char *message, ...) {
	errCtx *tctx=(errCtx *) ctx;
	char *newLine = malloc(GO_ERR_INIT);

	va_list varArgs;
        va_start(varArgs, message);

	int oldLen = strlen(tctx->errBuf) + 1;
	int lineLen = 1 + vsnprintf(newLine, GO_ERR_INIT, message, varArgs);

	if (lineLen  > GO_ERR_INIT) {
		free(newLine);
		newLine = malloc(lineLen);
		vsnprintf(newLine, GO_ERR_INIT, message, varArgs);
	}
	va_end(varArgs);

	char *tmp = malloc(oldLen + lineLen);
	strcpy(tmp, tctx->errBuf);
	strcat(tmp, newLine);
	free(tctx->errBuf);
	tctx->errBuf = tmp;
	free(newLine);
}

static xsdParserResult cParseUrlSchema(const char *url) {
	xsdParserResult parserResult;
	char *errBuf=NULL;
	errCtx *ectx=malloc(sizeof(errCtx));
	ectx->errBuf=calloc(sizeof(char), GO_ERR_INIT);;

	xmlSchemaPtr schema = NULL;
	xmlSchemaParserCtxtPtr schemaParserCtxt = NULL;

	//xmlSetGenericErrorFunc(ectx, genErrorCallback);
	xmlSetGenericErrorFunc(noOutputCallback, noOutputCallback);

	schemaParserCtxt = xmlSchemaNewParserCtxt(url);
	xmlSchemaSetParserErrors(schemaParserCtxt, genErrorCallback, noOutputCallback, ectx);

	schema = xmlSchemaParse(schemaParserCtxt);

	xmlSchemaFreeParserCtxt(schemaParserCtxt);

	if (schema == NULL) errno = -1;

	errBuf=calloc(sizeof(char), strlen(ectx->errBuf)+1);
	strcpy(errBuf, ectx->errBuf);
	free(ectx->errBuf);
	free(ectx);
	parserResult.schemaPtr=schema;
	parserResult.errorStr=errBuf;
	return parserResult;

}

static xmlParserResult cParseDoc(char *goXmlSource, int goXmlSourceLen) {
	xmlParserResult parserResult;
	char *errBuf=NULL;
	errCtx *ectx=malloc(sizeof(errCtx));
	ectx->errBuf=calloc(sizeof(char), GO_ERR_INIT);;

	xmlDocPtr doc=NULL;
	xmlParserCtxtPtr xmlParserCtxt=NULL;

	//xmlSetGenericErrorFunc(ectx, genErrorCallback);
	xmlSetGenericErrorFunc(noOutputCallback, noOutputCallback);

	xmlParserCtxt = xmlNewParserCtxt();

	doc = xmlParseMemory(goXmlSource, goXmlSourceLen);

	xmlFreeParserCtxt(xmlParserCtxt);

	if (doc == NULL) errno = -1;

	errBuf=calloc(sizeof(char), strlen(ectx->errBuf)+1);
	strcpy(errBuf, ectx->errBuf);
	free(ectx->errBuf);
	free(ectx);
	parserResult.docPtr=doc;
	parserResult.errorStr=errBuf;
	return parserResult;
}

static char *cValidate(xmlDocPtr doc, xmlSchemaPtr schema) {
	char *errBuf=NULL;
	errCtx *ectx=malloc(sizeof(errCtx));
	ectx->errBuf=calloc(sizeof(char), GO_ERR_INIT);;
	int schemaErr=0;

	if (schema == NULL) {
		errno = -1;
		strcpy(ectx->errBuf, "Error xsd schema: null pointer");
	}
	else if (doc == NULL) {
		errno = -1;
		strcpy(ectx->errBuf, "Error xml schema: null pointer");
	}
	else
	{
		xmlSchemaValidCtxtPtr schemaCtxt;
		schemaCtxt = xmlSchemaNewValidCtxt(schema);

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
			strcpy(ectx->errBuf, "Xml validation generated an internal error\n");
		}

	}

	errBuf=calloc(sizeof(char), strlen(ectx->errBuf)+1);
	strcpy(errBuf, ectx->errBuf);
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
	C.xmlInitParser()
	C.xmlLineNumbersDefault(1)
}

// Cleans up the libxml2 parser
func libXml2Cleanup() {
	C.xmlSchemaCleanupTypes()
	C.xmlCleanupParser()
}

// The helper function for parsing xml
func parseXmlMem(inXml []byte) (C.xmlDocPtr, error) {

	strXml := C.CString(string(inXml))
	defer C.free(unsafe.Pointer(strXml))
	pRes, err := C.cParseDoc(strXml, C.int(len(inXml)))

	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, errors.New(strings.Trim(rStr, "\n"))
	}
	return pRes.docPtr, nil
}

// The helper function for parsing the schema
func parseUrlSchema(url string) (C.xmlSchemaPtr, error) {
	strUrl := C.CString(url)
	defer C.free(unsafe.Pointer(strUrl))

	pRes, err := C.cParseUrlSchema(strUrl)
	defer C.free(unsafe.Pointer(pRes.errorStr))
	if err != nil {
		rStr := C.GoString(pRes.errorStr)
		return nil, errors.New(strings.Trim(rStr, "\n"))
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
		return errors.New(strings.Trim(rStr, "\n"))
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
