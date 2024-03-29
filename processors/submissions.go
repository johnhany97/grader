package processors

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// SubmissionsProcessor struct which attached to it
// are the methods used by the test tasks to execute
// compilation/running tasks that interface with Docker
type SubmissionsProcessor struct{}

// Language type
type Language string

// FileExtension type
type FileExtension string

// LanguageProperties is a struct containing features of a language
// which is used currently by the inner map of the languages to
// reflect the support for the languages currently hy the SubmissionsProcessor
type LanguageProperties struct {
	Name        Language      `json:"name"`
	Extension   FileExtension `json:"extension"`
	DockerImage string        `json:"dockerImage"`
}

// Inner map of the languages supported by dexec
// and the respective file extensions, language names
// and docker images
var languages = map[string]LanguageProperties{
	"c":      {Name: "C", Extension: "c", DockerImage: "dexec/lang-c"},
	"clj":    {Name: "Clojure", Extension: "clj", DockerImage: "dexec/lang-clojure"},
	"coffee": {Name: "CoffeeScript", Extension: "coffee", DockerImage: "dexec/lang-coffee"},
	"cpp":    {Name: "C++", Extension: "cpp", DockerImage: "dexec/lang-cpp"},
	"cs":     {Name: "C#", Extension: "cs", DockerImage: "dexec/lang-csharp"},
	"d":      {Name: "D", Extension: "d", DockerImage: "dexec/lang-d"},
	"erl":    {Name: "Erlang", Extension: "erl", DockerImage: "dexec/lang-erlang"},
	"fs":     {Name: "F#", Extension: "fs", DockerImage: "dexec/lang-fsharp"},
	"go":     {Name: "Go", Extension: "go", DockerImage: "dexec/lang-go"},
	"groovy": {Name: "Groovy", Extension: "groovy", DockerImage: "dexec/lang-groovy"},
	"hs":     {Name: "Haskell", Extension: "hs", DockerImage: "dexec/lang-haskell"},
	"java":   {Name: "Java", Extension: "java", DockerImage: "dexec/lang-java"},
	"lisp":   {Name: "Lisp", Extension: "lisp", DockerImage: "dexec/lang-lisp"},
	"lua":    {Name: "Lua", Extension: "lua", DockerImage: "dexec/lang-lua"},
	"js":     {Name: "JavaScript", Extension: "js", DockerImage: "dexec/lang-node"},
	"nim":    {Name: "Nim", Extension: "nim", DockerImage: "dexec/lang-nim"},
	"m":      {Name: "Objective C", Extension: "m", DockerImage: "dexec/lang-objc"},
	"ml":     {Name: "OCaml", Extension: "ml", DockerImage: "dexec/lang-ocaml"},
	"p6":     {Name: "Perl 6", Extension: "p6", DockerImage: "dexec/lang-perl6"},
	"pl":     {Name: "Perl", Extension: "pl", DockerImage: "dexec/lang-perl"},
	"php":    {Name: "PHP", Extension: "php", DockerImage: "dexec/lang-php"},
	"py":     {Name: "Python", Extension: "py", DockerImage: "dexec/lang-python"},
	"r":      {Name: "R", Extension: "r", DockerImage: "dexec/lang-r"},
	"rky":    {Name: "Racket", Extension: "rkt", DockerImage: "dexec/lang-racket"},
	"rb":     {Name: "Ruby", Extension: "rb", DockerImage: "dexec/lang-ruby"},
	"rs":     {Name: "Rust", Extension: "rs", DockerImage: "dexec/lang-rust"},
	"scala":  {Name: "Scala", Extension: "scala", DockerImage: "dexec/lang-scala"},
	"sh":     {Name: "Bash", Extension: "sh", DockerImage: "dexec/lang-bash"},
}

// Constants relating to timing out test tasks
// and the errors that may be produced. Should be set
// according to the system's configurations.
const timePerTask = 10 // time.Second
const timeoutErrorExistStatusCode = 124
const timeoutErrorMessage = "timeout"

// timedExec executes a given command normally but
// detects if the exit signal was that of a timeout
// and reports so to by returning a unique kind of error
func timedExec(cmd *exec.Cmd) (bool, error) {
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == timeoutErrorExistStatusCode {
				return false, errors.New(timeoutErrorMessage)
			}
		}
		return false, err
	}
	return true, nil
}

// Execute is a method that given a file and the folder within which it exists,
// runs the file and just sends back any stdout and any errors produced.
func (p SubmissionsProcessor) Execute(file string, folder string) (string, string, error) {
	var cmd *exec.Cmd

	// standard output and error buffers
	var out bytes.Buffer
	var e bytes.Buffer

	cmd = exec.Command("dexec", "-t", strconv.Itoa(timePerTask), file)

	cmd.Dir = folder
	cmd.Stdout = &out
	cmd.Stderr = &e

	ok, err := timedExec(cmd)
	if !ok {
		return out.String(), e.String(), err
	}
	return out.String(), e.String(), nil
}

// ExecuteWithInput is a method given a file, a folder and some input runs
// that file and provides it with the input and returns any stdout or errors produced
func (p SubmissionsProcessor) ExecuteWithInput(file string, folder string, input string) (string, string, error) {
	var cmd *exec.Cmd

	// standard output and error buffers
	var out bytes.Buffer
	var e bytes.Buffer

	cmd = exec.Command("dexec", "-t", strconv.Itoa(timePerTask), file)

	cmd.Dir = folder

	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &e

	ok, err := timedExec(cmd)
	if !ok {
		return out.String(), e.String(), err
	}
	return out.String(), e.String(), nil
}

// ExecuteJUnitTests is a method used to execute JUnit Tests against a given file.
// It starts by writing the JUnit Test file then running a specially created docker image
// with the required files and returns back any output or errors produced.
func (p SubmissionsProcessor) ExecuteJUnitTests(className string, folder string, junitTests string) (string, error) {
	fileName := className + "Test.java"
	path := folder + fileName
	err := writeShellFile(fileName, path, junitTests)
	if err != nil {
		return "", err
	}
	defer deletePath(path)

	// standard output buffer
	var out bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker run -t --rm -v $(pwd -P)/%v.java:/tmp/dexec/build/%v.java -v $(pwd -P)/%v:/tmp/dexec/build/%v johnhany97/grader-junit %v.java %v", className, className, fileName, fileName, className, fileName))

	cmd.Dir = folder

	cmd.Stdout = &out

	ok, err := timedExec(cmd)
	if !ok {
		return out.String(), err
	}
	return out.String(), nil
}

// ExecutePyUnitTests is a method used to execute Python unit tests against a given file.
// It starts by writing the Python unit test file then running the python docker image
// with the required files and returns back any output or errors produced.
func (p SubmissionsProcessor) ExecutePyUnitTests(file string, className string, folder string, pyUnitTests string) (string, error) {
	fileName := "test_" + file
	path := folder + fileName
	err := writeShellFile(fileName, path, pyUnitTests)
	if err != nil {
		return "", err
	}
	defer deletePath(path)

	var cmd *exec.Cmd

	// standard output buffer
	var out bytes.Buffer

	cmd = exec.Command("dexec", "-t", strconv.Itoa(timePerTask), fileName, "-i", file)

	cmd.Dir = folder

	cmd.Stderr = &out // BUG: For some reason python unit tests output to stderr

	ok, err := timedExec(cmd)
	if !ok {
		return out.String(), err
	}
	return out.String(), nil
}

// ExecuteJavaStyle is a method used to run style checks using the Google CheckStyle XML
// against a given Java Class.
func (p SubmissionsProcessor) ExecuteJavaStyle(file string, folder string) (string, string, error) {
	path := folder + "google_checks.xml"

	err := writeJavaStyleChecks(path)
	if err != nil {
		return "", "", errors.New("error checking style. please contact administrators")
	}
	defer deletePath(path)

	var cmd *exec.Cmd

	// standard output and error buffers
	var out bytes.Buffer
	var e bytes.Buffer

	cmd = exec.Command("checkstyle", "-c", "google_checks.xml", file)

	cmd.Dir = folder

	cmd.Stdout = &out
	cmd.Stderr = &e

	ok, err := timedExec(cmd)
	if !ok {
		return out.String(), e.String(), err
	}
	return out.String(), e.String(), nil
}

// Method that given a name for the file, a path to where it'll exist
// and the data to write within it creates that file and returns
// any errors found, if any
func writeShellFile(name string, path string, data string) error {
	// detect if file exists
	_, err := os.Stat(path)
	// delete file if exists
	if os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

// deletePath is used to delete a given path
func deletePath(path string) {
	var err = os.Remove(path)
	if err != nil {
		return
	}
}

// writeJavaStyleChecks is used to write the google_checks.xml
// used by checkstyle which is relevant to the javaStyle test task
func writeJavaStyleChecks(path string) error {
	// detect if file exists
	_, err := os.Stat(path)
	// delete file if exists
	if os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(path, []byte(googleChecks), 0644)
	if err != nil {
		return err
	}
	return nil
}

const googleChecks = `<?xml version="1.0"?>
<!DOCTYPE module PUBLIC
          "-//Puppy Crawl//DTD Check Configuration 1.3//EN"
          "http://checkstyle.sourceforge.net/dtds/configuration_1_3.dtd">

<!--
    Checkstyle configuration that checks the Google coding conventions from Google Java Style
    that can be found at https://google.github.io/styleguide/javaguide.html.

    Checkstyle is very configurable. Be sure to read the documentation at
    http://checkstyle.sf.net (or in your downloaded distribution).

    To completely disable a check, just comment it out or delete it from the file.

    Authors: Max Vetrenko, Ruslan Diachenko, Roman Ivanov.
 -->

<module name = "Checker">
    <property name="charset" value="UTF-8"/>

    <property name="severity" value="warning"/>

    <property name="fileExtensions" value="java, properties, xml"/>
    <!-- Checks for whitespace                               -->
    <!-- See http://checkstyle.sf.net/config_whitespace.html -->
    <module name="FileTabCharacter">
        <property name="eachLine" value="true"/>
    </module>

    <module name="TreeWalker">
        <module name="OuterTypeFilename"/>
        <module name="IllegalTokenText">
            <property name="tokens" value="STRING_LITERAL, CHAR_LITERAL"/>
            <property name="format" value="\\u00(09|0(a|A)|0(c|C)|0(d|D)|22|27|5(C|c))|\\(0(10|11|12|14|15|42|47)|134)"/>
            <property name="message" value="Consider using special escape sequence instead of octal value or Unicode escaped value."/>
        </module>
        <module name="AvoidEscapedUnicodeCharacters">
            <property name="allowEscapesForControlCharacters" value="true"/>
            <property name="allowByTailComment" value="true"/>
            <property name="allowNonPrintableEscapes" value="true"/>
        </module>
        <module name="LineLength">
            <property name="max" value="100"/>
            <property name="ignorePattern" value="^package.*|^import.*|a href|href|http://|https://|ftp://"/>
        </module>
        <module name="AvoidStarImport"/>
        <module name="OneTopLevelClass"/>
        <module name="NoLineWrap"/>
        <module name="EmptyBlock">
            <property name="option" value="TEXT"/>
            <property name="tokens" value="LITERAL_TRY, LITERAL_FINALLY, LITERAL_IF, LITERAL_ELSE, LITERAL_SWITCH"/>
        </module>
        <module name="NeedBraces"/>
        <module name="LeftCurly"/>
        <module name="RightCurly">
            <property name="id" value="RightCurlySame"/>
            <property name="tokens" value="LITERAL_TRY, LITERAL_CATCH, LITERAL_FINALLY, LITERAL_IF, LITERAL_ELSE, LITERAL_DO"/>
        </module>
        <module name="RightCurly">
            <property name="id" value="RightCurlyAlone"/>
            <property name="option" value="alone"/>
            <property name="tokens" value="CLASS_DEF, METHOD_DEF, CTOR_DEF, LITERAL_FOR, LITERAL_WHILE, STATIC_INIT, INSTANCE_INIT"/>
        </module>
        <module name="WhitespaceAround">
            <property name="allowEmptyConstructors" value="true"/>
            <property name="allowEmptyMethods" value="true"/>
            <property name="allowEmptyTypes" value="true"/>
            <property name="allowEmptyLoops" value="true"/>
            <message key="ws.notFollowed"
             value="WhitespaceAround: ''{0}'' is not followed by whitespace. Empty blocks may only be represented as '{}' when not part of a multi-block statement (4.1.3)"/>
             <message key="ws.notPreceded"
             value="WhitespaceAround: ''{0}'' is not preceded with whitespace."/>
        </module>
        <module name="OneStatementPerLine"/>
        <module name="MultipleVariableDeclarations"/>
        <module name="ArrayTypeStyle"/>
        <module name="MissingSwitchDefault"/>
        <module name="FallThrough"/>
        <module name="UpperEll"/>
        <module name="ModifierOrder"/>
        <module name="EmptyLineSeparator">
            <property name="allowNoEmptyLineBetweenFields" value="true"/>
        </module>
        <module name="SeparatorWrap">
            <property name="id" value="SeparatorWrapDot"/>
            <property name="tokens" value="DOT"/>
            <property name="option" value="nl"/>
        </module>
        <module name="SeparatorWrap">
            <property name="id" value="SeparatorWrapComma"/>
            <property name="tokens" value="COMMA"/>
            <property name="option" value="EOL"/>
        </module>
        <module name="SeparatorWrap">
            <!-- ELLIPSIS is EOL until https://github.com/google/styleguide/issues/258 -->
            <property name="id" value="SeparatorWrapEllipsis"/>
            <property name="tokens" value="ELLIPSIS"/>
            <property name="option" value="EOL"/>
        </module>
        <module name="SeparatorWrap">
            <!-- ARRAY_DECLARATOR is EOL until https://github.com/google/styleguide/issues/259 -->
            <property name="id" value="SeparatorWrapArrayDeclarator"/>
            <property name="tokens" value="ARRAY_DECLARATOR"/>
            <property name="option" value="EOL"/>
        </module>
        <module name="SeparatorWrap">
            <property name="id" value="SeparatorWrapMethodRef"/>
            <property name="tokens" value="METHOD_REF"/>
            <property name="option" value="nl"/>
        </module>
        <module name="PackageName">
            <property name="format" value="^[a-z]+(\.[a-z][a-z0-9]*)*$"/>
            <message key="name.invalidPattern"
             value="Package name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="TypeName">
            <message key="name.invalidPattern"
             value="Type name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="MemberName">
            <property name="format" value="^[a-z][a-z0-9][a-zA-Z0-9]*$"/>
            <message key="name.invalidPattern"
             value="Member name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="ParameterName">
            <property name="format" value="^[a-z]([a-z0-9][a-zA-Z0-9]*)?$"/>
            <message key="name.invalidPattern"
             value="Parameter name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="CatchParameterName">
            <property name="format" value="^[a-z]([a-z0-9][a-zA-Z0-9]*)?$"/>
            <message key="name.invalidPattern"
             value="Catch parameter name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="LocalVariableName">
            <property name="tokens" value="VARIABLE_DEF"/>
            <property name="format" value="^[a-z]([a-z0-9][a-zA-Z0-9]*)?$"/>
            <message key="name.invalidPattern"
             value="Local variable name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="ClassTypeParameterName">
            <property name="format" value="(^[A-Z][0-9]?)$|([A-Z][a-zA-Z0-9]*[T]$)"/>
            <message key="name.invalidPattern"
             value="Class type name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="MethodTypeParameterName">
            <property name="format" value="(^[A-Z][0-9]?)$|([A-Z][a-zA-Z0-9]*[T]$)"/>
            <message key="name.invalidPattern"
             value="Method type name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="InterfaceTypeParameterName">
            <property name="format" value="(^[A-Z][0-9]?)$|([A-Z][a-zA-Z0-9]*[T]$)"/>
            <message key="name.invalidPattern"
             value="Interface type name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="NoFinalizer"/>
        <module name="GenericWhitespace">
            <message key="ws.followed"
             value="GenericWhitespace ''{0}'' is followed by whitespace."/>
             <message key="ws.preceded"
             value="GenericWhitespace ''{0}'' is preceded with whitespace."/>
             <message key="ws.illegalFollow"
             value="GenericWhitespace ''{0}'' should followed by whitespace."/>
             <message key="ws.notPreceded"
             value="GenericWhitespace ''{0}'' is not preceded with whitespace."/>
        </module>
        <module name="Indentation">
            <property name="basicOffset" value="2"/>
            <property name="braceAdjustment" value="0"/>
            <property name="caseIndent" value="2"/>
            <property name="throwsIndent" value="4"/>
            <property name="lineWrappingIndentation" value="4"/>
            <property name="arrayInitIndent" value="2"/>
        </module>
        <module name="AbbreviationAsWordInName">
            <property name="ignoreFinal" value="false"/>
            <property name="allowedAbbreviationLength" value="1"/>
        </module>
        <module name="OverloadMethodsDeclarationOrder"/>
        <module name="VariableDeclarationUsageDistance"/>
        <module name="CustomImportOrder">
            <property name="sortImportsInGroupAlphabetically" value="true"/>
            <property name="separateLineBetweenGroups" value="true"/>
            <property name="customImportOrderRules" value="STATIC###THIRD_PARTY_PACKAGE"/>
        </module>
        <module name="MethodParamPad"/>
        <module name="NoWhitespaceBefore">
          <property name="tokens" value="COMMA, SEMI, POST_INC, POST_DEC, DOT, ELLIPSIS, METHOD_REF"/>
            <property name="allowLineBreaks" value="true"/>
        </module>
        <module name="ParenPad"/>
        <module name="OperatorWrap">
            <property name="option" value="NL"/>
            <property name="tokens" value="BAND, BOR, BSR, BXOR, DIV, EQUAL, GE, GT, LAND, LE, LITERAL_INSTANCEOF, LOR, LT, MINUS, MOD, NOT_EQUAL, PLUS, QUESTION, SL, SR, STAR, METHOD_REF "/>
        </module>
        <module name="AnnotationLocation">
            <property name="id" value="AnnotationLocationMostCases"/>
            <property name="tokens" value="CLASS_DEF, INTERFACE_DEF, ENUM_DEF, METHOD_DEF, CTOR_DEF"/>
        </module>
        <module name="AnnotationLocation">
            <property name="id" value="AnnotationLocationVariables"/>
            <property name="tokens" value="VARIABLE_DEF"/>
            <property name="allowSamelineMultipleAnnotations" value="true"/>
        </module>
        <module name="NonEmptyAtclauseDescription"/>
        <module name="JavadocTagContinuationIndentation"/>
        <module name="SummaryJavadoc">
            <property name="forbiddenSummaryFragments" value="^@return the *|^This method returns |^A [{]@code [a-zA-Z0-9]+[}]( is a )"/>
        </module>
        <module name="JavadocParagraph"/>
        <module name="AtclauseOrder">
            <property name="tagOrder" value="@param, @return, @throws, @deprecated"/>
            <property name="target" value="CLASS_DEF, INTERFACE_DEF, ENUM_DEF, METHOD_DEF, CTOR_DEF, VARIABLE_DEF"/>
        </module>
        <module name="JavadocMethod">
            <property name="scope" value="public"/>
            <property name="allowMissingParamTags" value="true"/>
            <property name="allowMissingThrowsTags" value="true"/>
            <property name="allowMissingReturnTag" value="true"/>
            <property name="minLineCount" value="2"/>
            <property name="allowedAnnotations" value="Override, Test"/>
            <property name="allowThrowsTagsForSubclasses" value="true"/>
        </module>
        <module name="MethodName">
            <property name="format" value="^[a-z][a-z0-9][a-zA-Z0-9_]*$"/>
            <message key="name.invalidPattern"
             value="Method name ''{0}'' must match pattern ''{1}''."/>
        </module>
        <module name="SingleLineJavadoc">
            <property name="ignoreInlineTags" value="false"/>
        </module>
        <module name="EmptyCatchBlock">
            <property name="exceptionVariableName" value="expected"/>
        </module>
        <module name="CommentsIndentation"/>
    </module>
</module>`
