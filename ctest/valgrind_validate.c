#include "xsdvalidate_libxml2.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libxml/parser.h>

static void fail(const char* msg) {
    fprintf(stderr, "%s\n", msg);
    exit(1);
}

static char* readFile(const char* path, int* len) {
    FILE* f = fopen(path, "rb");
    if (f == NULL) {
        fprintf(stderr, "failed to open %s\n", path);
        exit(1);
    }

    if (fseek(f, 0, SEEK_END) != 0) {
        fclose(f);
        fail("failed to seek file");
    }
    long size = ftell(f);
    if (size < 0) {
        fclose(f);
        fail("failed to tell file size");
    }
    if (fseek(f, 0, SEEK_SET) != 0) {
        fclose(f);
        fail("failed to rewind file");
    }

    char* buf = malloc((size_t)size + 1);
    if (buf == NULL) {
        fclose(f);
        fail("failed to allocate file buffer");
    }

    size_t n = fread(buf, 1, (size_t)size, f);
    if (n != (size_t)size) {
        free(buf);
        fclose(f);
        fail("failed to read complete file");
    }
    fclose(f);

    buf[size] = '\0';
    *len = (int)size;
    return buf;
}

static int hasError(const errArray* errArr, const char* node, const char* nodePath) {
    for (size_t i = 0; i < errArr->len; i++) {
        if (errArr->data[i].node != NULL && errArr->data[i].nodePath != NULL &&
            strcmp(errArr->data[i].node, node) == 0 &&
            strcmp(errArr->data[i].nodePath, nodePath) == 0) {
            return 1;
        }
    }
    return 0;
}

static void requireErrArrayType(const errArray* errArr, errorType typ) {
    if (errArr->len == 0) {
        fail("expected at least one error");
    }
    if (errArr->data[0].type != typ) {
        fprintf(stderr, "expected error type %d, got %d\n", typ, errArr->data[0].type);
        exit(1);
    }
}

static xmlSchemaPtr parseSchema(void) {
    struct xsdParserResult schema = xsdValidateParseUrlSchema("testdata/nodepath.xsd", P_ERR_DEFAULT);
    if (schema.schemaPtr == NULL) {
        fprintf(stderr, "schema parse failed: %s\n", schema.errorStr == NULL ? "" : schema.errorStr);
        free(schema.errorStr);
        exit(1);
    }
    free(schema.errorStr);
    return schema.schemaPtr;
}

static xmlDocPtr parseDoc(const char* path) {
    int len = 0;
    char* xml = readFile(path, &len);
    struct xmlParserResult doc = xsdValidateParseDoc(xml, len, P_ERR_DEFAULT);
    free(xml);
    if (doc.docPtr == NULL) {
        fprintf(stderr, "xml parse failed: %s\n", doc.errorStr == NULL ? "" : doc.errorStr);
        free(doc.errorStr);
        exit(1);
    }
    free(doc.errorStr);
    return doc.docPtr;
}

static void testParseUrlSchemaFail(void) {
    struct xsdParserResult schema = xsdValidateParseUrlSchema("testdata/does-not-exist.xsd", P_ERR_DEFAULT);
    if (schema.schemaPtr != NULL) {
        xmlSchemaFree(schema.schemaPtr);
        free(schema.errorStr);
        fail("expected missing schema URL parse to fail");
    }
    free(schema.errorStr);
}

static void testParseMemSchema(void) {
    int len = 0;
    char* xsd = readFile("testdata/nodepath.xsd", &len);
    struct xsdParserResult schema = xsdValidateParseMemSchema(xsd, len, P_ERR_DEFAULT);
    free(xsd);
    if (schema.schemaPtr == NULL) {
        fprintf(stderr, "memory schema parse failed: %s\n", schema.errorStr == NULL ? "" : schema.errorStr);
        free(schema.errorStr);
        exit(1);
    }
    xmlSchemaFree(schema.schemaPtr);
    free(schema.errorStr);

    const char badXsd[] = "<xs:schema xmlns:xs=\"http://www.w3.org/2001/XMLSchema\"><xs:element name=\"broken\"";
    schema = xsdValidateParseMemSchema(badXsd, (int)strlen(badXsd), P_ERR_DEFAULT);
    if (schema.schemaPtr != NULL) {
        xmlSchemaFree(schema.schemaPtr);
        free(schema.errorStr);
        fail("expected malformed memory schema parse to fail");
    }
    free(schema.errorStr);
}

static void testParseDocFailures(void) {
    const char malformedXml[] = "<order>";
    struct xmlParserResult doc = xsdValidateParseDoc(malformedXml, (int)strlen(malformedXml), P_ERR_DEFAULT);
    if (doc.docPtr != NULL) {
        xmlFreeDoc(doc.docPtr);
        free(doc.errorStr);
        fail("expected direct malformed XML parse to fail");
    }
    free(doc.errorStr);

    doc = xsdValidateParseDoc("", 0, P_ERR_DEFAULT);
    if (doc.docPtr != NULL) {
        xmlFreeDoc(doc.docPtr);
        free(doc.errorStr);
        fail("expected direct empty XML parse to fail");
    }
    free(doc.errorStr);
}

static void testValidateDocPass(xmlSchemaPtr schema) {
    xmlDocPtr doc = parseDoc("testdata/nodepath_valid.xml");
    errArray errArr = xsdValidateDoc(doc, schema);
    if (errArr.len != 0) {
        xsdValidateFreeErrArray(&errArr);
        xmlFreeDoc(doc);
        fail("expected valid document to produce no validation errors");
    }
    xsdValidateFreeErrArray(&errArr);
    xmlFreeDoc(doc);
}

static void testValidateBufPass(xmlSchemaPtr schema) {
    int len = 0;
    char* xml = readFile("testdata/nodepath_valid.xml", &len);
    errArray errArr = xsdValidateBuf(xml, len, P_ERR_DEFAULT, schema);
    free(xml);

    if (errArr.len != 0) {
        xsdValidateFreeErrArray(&errArr);
        fail("expected valid buffer to produce no validation errors");
    }
    xsdValidateFreeErrArray(&errArr);
}

static void testValidateBufNodePath(xmlSchemaPtr schema) {
    int len = 0;
    char* xml = readFile("testdata/nodepath_bad_nested_value.xml", &len);
    errArray errArr = xsdValidateBuf(xml, len, P_ERR_DEFAULT, schema);
    free(xml);

    requireErrArrayType(&errArr, VALIDATION_ERROR);
    if (!hasError(&errArr, "quantity", "/order/items/item/quantity")) {
        xsdValidateFreeErrArray(&errArr);
        fail("expected quantity node path validation error");
    }
    xsdValidateFreeErrArray(&errArr);
}

static void testValidateBufMultipleErrors(xmlSchemaPtr schema) {
    int len = 0;
    char* xml = readFile("testdata/nodepath_multiple_errors.xml", &len);
    errArray errArr = xsdValidateBuf(xml, len, P_ERR_DEFAULT, schema);
    free(xml);

    requireErrArrayType(&errArr, VALIDATION_ERROR);
    if (!hasError(&errArr, "order", "/order") ||
        !hasError(&errArr, "quantity", "/order/items/item/quantity")) {
        xsdValidateFreeErrArray(&errArr);
        fail("expected multiple node path validation errors");
    }
    xsdValidateFreeErrArray(&errArr);
}

static char* makeHugeXml(int* len) {
    const int depth = 300;
    const size_t textLen = 10 * 1024 * 1024 + 1;
    const char prefix[] = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root>";
    const char dataOpen[] = "<data>";
    const char dataClose[] = "</data>";
    const char suffix[] = "</root>";
    const size_t total = strlen(prefix) + (size_t)depth * strlen("<a>") +
                         strlen(dataOpen) + textLen + strlen(dataClose) +
                         (size_t)depth * strlen("</a>") + strlen(suffix);

    char* xml = malloc(total + 1);
    if (xml == NULL) {
        fail("failed to allocate huge XML buffer");
    }

    char* p = xml;
    memcpy(p, prefix, strlen(prefix));
    p += strlen(prefix);
    for (int i = 0; i < depth; i++) {
        memcpy(p, "<a>", strlen("<a>"));
        p += strlen("<a>");
    }
    memcpy(p, dataOpen, strlen(dataOpen));
    p += strlen(dataOpen);
    memset(p, 'A', textLen);
    p += textLen;
    memcpy(p, dataClose, strlen(dataClose));
    p += strlen(dataClose);
    for (int i = 0; i < depth; i++) {
        memcpy(p, "</a>", strlen("</a>"));
        p += strlen("</a>");
    }
    memcpy(p, suffix, strlen(suffix));
    p += strlen(suffix);
    *p = '\0';

    *len = (int)total;
    return xml;
}

static void testParseHugeXml(void) {
    int len = 0;
    char* xml = makeHugeXml(&len);

    struct xmlParserResult defaultDoc = xsdValidateParseDoc(xml, len, P_ERR_DEFAULT);
    if (defaultDoc.docPtr != NULL) {
        xmlFreeDoc(defaultDoc.docPtr);
        free(defaultDoc.errorStr);
        free(xml);
        fail("expected huge XML parsing without P_XML_HUGE to fail");
    }
    free(defaultDoc.errorStr);

    struct xmlParserResult hugeDoc = xsdValidateParseDoc(xml, len, P_XML_HUGE);
    free(xml);
    if (hugeDoc.docPtr == NULL) {
        fprintf(stderr, "huge XML parse failed: %s\n", hugeDoc.errorStr == NULL ? "" : hugeDoc.errorStr);
        free(hugeDoc.errorStr);
        exit(1);
    }
    xmlFreeDoc(hugeDoc.docPtr);
    free(hugeDoc.errorStr);
}

static void testMalformedXml(xmlSchemaPtr schema) {
    const char xml[] = "<order>";
    errArray errArr = xsdValidateBuf(xml, (int)strlen(xml), P_ERR_DEFAULT, schema);
    requireErrArrayType(&errArr, XML_PARSER_ERROR);
    xsdValidateFreeErrArray(&errArr);

    errArr = xsdValidateBuf("", 0, P_ERR_DEFAULT, schema);
    requireErrArrayType(&errArr, XML_PARSER_ERROR);
    xsdValidateFreeErrArray(&errArr);
}

static void testValidateBufNullSchema(void) {
    int len = 0;
    char* xml = readFile("testdata/nodepath_valid.xml", &len);
    errArray errArr = xsdValidateBuf(xml, len, P_ERR_DEFAULT, NULL);
    free(xml);

    requireErrArrayType(&errArr, LIBXML2_ERROR);
    xsdValidateFreeErrArray(&errArr);
}

static void testNullBranches(xmlSchemaPtr schema) {
    xmlDocPtr doc = parseDoc("testdata/nodepath_valid.xml");

    errArray missingDoc = xsdValidateDoc(NULL, schema);
    requireErrArrayType(&missingDoc, LIBXML2_ERROR);
    xsdValidateFreeErrArray(&missingDoc);

    errArray missingSchema = xsdValidateDoc(doc, NULL);
    requireErrArrayType(&missingSchema, LIBXML2_ERROR);
    xsdValidateFreeErrArray(&missingSchema);

    xmlFreeDoc(doc);
}

static void testRepeatedValidation(xmlSchemaPtr schema) {
    int len = 0;
    char* xml = readFile("testdata/nodepath_bad_nested_value.xml", &len);

    for (int i = 0; i < 1000; i++) {
        errArray errArr = xsdValidateBuf(xml, len, P_ERR_DEFAULT, schema);
        requireErrArrayType(&errArr, VALIDATION_ERROR);
        xsdValidateFreeErrArray(&errArr);
    }

    free(xml);
}

int main(void) {
    xsdValidateInit();

    testParseUrlSchemaFail();
    testParseMemSchema();
    testParseDocFailures();

    xmlSchemaPtr schema = parseSchema();
    testValidateDocPass(schema);
    testValidateBufPass(schema);
    testValidateBufNodePath(schema);
    testValidateBufMultipleErrors(schema);
    testMalformedXml(schema);
    testValidateBufNullSchema();
    testParseHugeXml();
    testNullBranches(schema);
    testRepeatedValidation(schema);
    xmlSchemaFree(schema);

    xsdValidateCleanup();
    return 0;
}
