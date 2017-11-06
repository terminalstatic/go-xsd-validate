package xsdvalidate

/*
#cgo pkg-config: libxml-2.0
#include <string.h>
#include <libxml/xmlschemastypes.h>
#include <errno.h>
#include <malloc.h>
#include <stdbool.h>
#define GO_ERR_INIT 256
#define P_ERR_DEFAULT 1
#define P_ERR_EXT 2
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

struct simpleXmlError {
	int	code;
	char*	message;
	int 	level;
	int	line;
	char*	node;
};


static void noOutputCallback(void *ctx, const char *message, ...) {
}

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

void structErrorCallback(void *ctx, xmlErrorPtr p) {
	struct errCtx *ectx = ctx;
	char *newLine = malloc(GO_ERR_INIT);


	//pid_t pid = syscall(__NR_gettid);

	//printf("threadId: %li, code: %d, level: %d, node: %s, line: %d, message: %s", pid, p->code, p->level, ((xmlNodePtr) p->node)->name, p->line, p->message);
	int oldLen = strlen(ectx->errBuf) + 1;
	int lineLen = 1 + snprintf(newLine, GO_ERR_INIT, "%s", p->message);

	if (lineLen  > GO_ERR_INIT) {
		free(newLine);
		newLine = malloc(lineLen);
		snprintf(newLine, lineLen, "%s", p->message);
	}


	char *tmp = malloc(oldLen + lineLen);
	memcpy(tmp, ectx->errBuf, oldLen);
	strcat(tmp, newLine);
	free(newLine);
	free(ectx->errBuf);
	ectx->errBuf = tmp;
}


void simpleStructErrorCallback(void *ctx, xmlErrorPtr p) {
	struct simpleXmlError *sErr = ctx;
	sErr->code = p->code;
	sErr->level = p->level;
	sErr->line = p->line;

        int cpyLen = 1 + snprintf(sErr->message, GO_ERR_INIT, "%s", p->message);
	if (cpyLen > GO_ERR_INIT) {
		free(sErr->message);
		sErr->message = malloc(cpyLen);
		snprintf(sErr->message, cpyLen, "%s", p->message);
	}

	if (p->node !=NULL) {
		cpyLen = 1 + snprintf(sErr->node, GO_ERR_INIT, "%s", (((xmlNodePtr) p->node)->name));
		if (cpyLen > GO_ERR_INIT) {
			free(sErr->node);
			sErr->node= malloc(cpyLen);
			snprintf(sErr->node, cpyLen, "%s", (((xmlNodePtr) p->node)->name));
		}
	}

}

static struct xsdParserResult cParseUrlSchema(const char *url, const short int options) {
	bool err = false;
	struct xsdParserResult parserResult;
	char *errBuf=NULL;
	struct errCtx *ectx=malloc(sizeof(struct errCtx));
	ectx->errBuf=calloc(GO_ERR_INIT, sizeof(char));
	struct errCtx *genEctx=malloc(sizeof(struct errCtx));;
	genEctx->errBuf=calloc(GO_ERR_INIT, sizeof(char));

	xmlSchemaPtr schema = NULL;
	xmlSchemaParserCtxtPtr schemaParserCtxt = NULL;

	xmlLineNumbersDefault(1);

	schemaParserCtxt = xmlSchemaNewParserCtxt(url);

	if (schemaParserCtxt == NULL) {
		err = true;
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

		xmlSchemaFreeParserCtxt(schemaParserCtxt);
		if (schema == NULL) {
			err = true;
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
	errno = err ? -1 : 0;
	return parserResult;
}

static struct xmlParserResult cParseDoc(const char *goXmlSource, const int goXmlSourceLen, const short int options) {
	bool err = false;
	struct xmlParserResult parserResult;
	char *errBuf=NULL;
	struct errCtx *ectx=malloc(sizeof(struct errCtx));
	ectx->errBuf=calloc(GO_ERR_INIT, sizeof(char));;

	xmlLineNumbersDefault(1);

	xmlDocPtr doc=NULL;
	xmlParserCtxtPtr xmlParserCtxt=NULL;

	xmlParserCtxt = xmlNewParserCtxt();

	if (xmlParserCtxt == NULL) {
		err = true;
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
				err = true;
				char *tmp = malloc(strlen(ectx->errBuf) + 1);
				memcpy(tmp, ectx->errBuf, strlen(ectx->errBuf) + 1);
				free(ectx->errBuf);
				ectx->errBuf = tmp;
			} else {
				err = true;
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
	errno = err ? -1 : 0;
	return parserResult;
}

static struct simpleXmlError *cValidate(const xmlDocPtr doc, const xmlSchemaPtr schema) {
	bool err = false;
	int schemaErr=0;

	struct simpleXmlError *simpleError = malloc(sizeof(struct simpleXmlError));
	simpleError->message = calloc(GO_ERR_INIT, sizeof(char));
	simpleError->node = calloc(GO_ERR_INIT, sizeof(char));

	//xmlErrorPtr errPtr= malloc(sizeof(xmlError));

	xmlLineNumbersDefault(1);

	if (schema == NULL) {
		err = true;
		strcpy(simpleError->message, "Xsd schema null pointer");
	}
	else if (doc == NULL) {
		err = true;
		strcpy(simpleError->message, "Xml schema null pointer");
	}
	else
	{
		xmlSchemaValidCtxtPtr schemaCtxt;
		schemaCtxt = xmlSchemaNewValidCtxt(schema);

		if (schemaCtxt == NULL) {
			err = true;
			strcpy(simpleError->message, "Xml validation internal error");
		}
		else
		{

			xmlSchemaSetValidStructuredErrors(schemaCtxt, simpleStructErrorCallback, simpleError);
			schemaErr = xmlSchemaValidateDoc(schemaCtxt, doc);
			xmlSchemaFreeValidCtxt(schemaCtxt);

			if (schemaErr > 0)
			{
				err = true;
			}
			else if (schemaErr < 0)
			{
				err = true;
				strcpy(simpleError->message, "Xml validation internal error");
			}
		}
	}

	errno = err ? -1 : 0;
	return simpleError;
}

*/
import "C"
import (
	"log"
	"runtime"
	"strings"
	"time"
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

// Helper function for validating given an xml document
func validateWithXsd(xmlHandler *XmlHandler, xsdHandler *XsdHandler) error {
	sErr, err := C.cValidate(xmlHandler.docPtr, xsdHandler.schemaPtr)
	defer freeSimpleXmlError(sErr)
	if err != nil {
		return ValidationError{
			Code:     int(sErr.code),
			Message:  strings.Trim(C.GoString(sErr.message), "\n"),
			Level:    int(sErr.level),
			Line:     int(sErr.line),
			NodeName: C.GoString(sErr.node),
		}
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

// Free C struct
func freeSimpleXmlError(sxe *C.struct_simpleXmlError) {
	C.free(unsafe.Pointer(sxe.message))
	C.free(unsafe.Pointer(sxe.node))
	C.free(unsafe.Pointer(sxe))
}

// Ticker for gc and malloc_trim
func gcTicker(d time.Duration, quit chan struct{}) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			log.Println("Running GC and malloc_trim(0)")
			runtime.GC()
			C.malloc_trim(0)
		case <-quit:
			log.Println("GC ticker stopped")
			ticker.Stop()
			return
		}
	}
}
