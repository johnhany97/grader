package processors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type SubmissionsProcessor struct{}

func (p SubmissionsProcessor) Execute(file string, folder string) (string, string) {
	var cmd *exec.Cmd
	cmd = exec.Command("dexec", file)
	cmd.Dir = folder
	var out bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	cmd.Run()
	return out.String(), e.String()
}

func (p SubmissionsProcessor) ExecuteWithInput(file string, folder string, input string) (string, string) {
	// compile command
	var cmd *exec.Cmd
	cmd = exec.Command("dexec", file)
	cmd.Dir = folder
	// provide stdin
	cmd.Stdin = strings.NewReader(input)
	// take back stdout
	var out bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	cmd.Run()
	return out.String(), e.String()
}

func (p SubmissionsProcessor) ExecuteJUnitTests(className string, folder string, junitTests string) (string, error) {
	fileName := className + "Test.java"
	path := folder + fileName
	err := writeShellFile(fileName, path, junitTests)
	if err != nil {
		return "", err
	}
	// delete when done
	defer deletePath(path)
	junitCmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker run -t --rm -v $(pwd -P)/%v.java:/tmp/dexec/build/%v.java -v $(pwd -P)/%v:/tmp/dexec/build/%v johnhany97/grader-junit %v.java %v", className, className, fileName, fileName, className, fileName))
	junitCmd.Dir = folder
	var out bytes.Buffer
	junitCmd.Stdout = &out
	if err = junitCmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (p SubmissionsProcessor) ExecutePyUnitTests(file string, className string, folder string, pyUnitTests string) (string, error) {
	fileName := "test_" + file
	path := folder + fileName
	err := writeShellFile(fileName, path, pyUnitTests)
	if err != nil {
		return "", err
	}
	// delete when done
	defer deletePath(path)
	var cmd *exec.Cmd
	cmd = exec.Command("dexec", fileName, "-i", file)
	cmd.Dir = folder
	var out bytes.Buffer
	cmd.Stderr = &out
	if err = cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (p SubmissionsProcessor) ExecuteJavaStyle(file string, folder string) (string, string) {
	path := folder + "google_checks.xml"
	var cmd *exec.Cmd
	err := writeJavaStyleChecks(path)
	if err != nil {
		return "", "Error Checking Style. Please contact administrators."
	}
	// defer deletePath(path)
	cmd = exec.Command("checkstyle", "-c", "google_checks.xml", file)
	cmd.Dir = folder
	var out bytes.Buffer
	var e bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &e
	cmd.Run()
	return out.String(), e.String()
}

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

func deletePath(path string) {
	var err = os.Remove(path)
	if err != nil {
		return
	}
}

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
          "-//Checkstyle//DTD Checkstyle Configuration 1.3//EN"
          "https://checkstyle.org/dtds/configuration_1_3.dtd">

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
    <!-- Excludes all 'module-info.java' files              -->
    <!-- See https://checkstyle.org/config_filefilters.html -->
    <module name="BeforeExecutionExclusionFileFilter">
        <property name="fileNamePattern" value="module\-info\.java$"/>
    </module>
    <!-- Checks for whitespace                               -->
    <!-- See http://checkstyle.sf.net/config_whitespace.html -->
    <module name="FileTabCharacter">
        <property name="eachLine" value="true"/>
    </module>

    <module name="TreeWalker">
        <module name="OuterTypeFilename"/>
        <module name="IllegalTokenText">
            <property name="tokens" value="STRING_LITERAL, CHAR_LITERAL"/>
            <property name="format"
             value="\\u00(09|0(a|A)|0(c|C)|0(d|D)|22|27|5(C|c))|\\(0(10|11|12|14|15|42|47)|134)"/>
            <property name="message"
             value="Consider using special escape sequence instead of octal value or Unicode escaped value."/>
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
            <property name="tokens"
             value="LITERAL_TRY, LITERAL_FINALLY, LITERAL_IF, LITERAL_ELSE, LITERAL_SWITCH"/>
        </module>
        <module name="NeedBraces"/>
        <module name="LeftCurly"/>
        <module name="RightCurly">
            <property name="id" value="RightCurlySame"/>
            <property name="tokens"
             value="LITERAL_TRY, LITERAL_CATCH, LITERAL_FINALLY, LITERAL_IF, LITERAL_ELSE,
                    LITERAL_DO"/>
        </module>
        <module name="RightCurly">
            <property name="id" value="RightCurlyAlone"/>
            <property name="option" value="alone"/>
            <property name="tokens"
             value="CLASS_DEF, METHOD_DEF, CTOR_DEF, LITERAL_FOR, LITERAL_WHILE, STATIC_INIT,
                    INSTANCE_INIT"/>
        </module>
        <module name="WhitespaceAround">
            <property name="allowEmptyConstructors" value="true"/>
            <property name="allowEmptyLambdas" value="true"/>
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
        <module name="LambdaParameterName">
            <property name="format" value="^[a-z]([a-z0-9][a-zA-Z0-9]*)?$"/>
            <message key="name.invalidPattern"
                     value="Lambda parameter name ''{0}'' must match pattern ''{1}''."/>
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
            <property name="tokens"
             value="COMMA, SEMI, POST_INC, POST_DEC, DOT, ELLIPSIS, METHOD_REF"/>
            <property name="allowLineBreaks" value="true"/>
        </module>
        <module name="ParenPad"/>
        <module name="OperatorWrap">
            <property name="option" value="NL"/>
            <property name="tokens"
             value="BAND, BOR, BSR, BXOR, DIV, EQUAL, GE, GT, LAND, LE, LITERAL_INSTANCEOF, LOR,
                    LT, MINUS, MOD, NOT_EQUAL, PLUS, QUESTION, SL, SR, STAR, METHOD_REF "/>
        </module>
        <module name="AnnotationLocation">
            <property name="id" value="AnnotationLocationMostCases"/>
            <property name="tokens"
             value="CLASS_DEF, INTERFACE_DEF, ENUM_DEF, METHOD_DEF, CTOR_DEF"/>
        </module>
        <module name="AnnotationLocation">
            <property name="id" value="AnnotationLocationVariables"/>
            <property name="tokens" value="VARIABLE_DEF"/>
            <property name="allowSamelineMultipleAnnotations" value="true"/>
        </module>
        <module name="NonEmptyAtclauseDescription"/>
        <module name="JavadocTagContinuationIndentation"/>
        <module name="SummaryJavadoc">
            <property name="forbiddenSummaryFragments"
             value="^@return the *|^This method returns |^A [{]@code [a-zA-Z0-9]+[}]( is a )"/>
        </module>
        <module name="JavadocParagraph"/>
        <module name="AtclauseOrder">
            <property name="tagOrder" value="@param, @return, @throws, @deprecated"/>
            <property name="target"
             value="CLASS_DEF, INTERFACE_DEF, ENUM_DEF, METHOD_DEF, CTOR_DEF, VARIABLE_DEF"/>
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
</module>
`
