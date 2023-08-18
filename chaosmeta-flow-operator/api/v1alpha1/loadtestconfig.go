/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import "encoding/xml"

type JMeterTestPlan struct {
	XMLName    xml.Name `xml:"jmeterTestPlan"`
	Version    string   `xml:"version,attr"`
	Properties string   `xml:"properties,attr"`
	JMeter     string   `xml:"jmeter,attr"`
	HashTree   HashTree `xml:"hashTree"`
}

type HashTree struct {
	XMLName          xml.Name         `xml:"hashTree"`
	TestPlan         TestPlan         `xml:"TestPlan"`
	ThreadGroup      ThreadGroup      `xml:"ThreadGroup"`
	HeaderManager    HeaderManager    `xml:"HeaderManager"`
	HTTPSamplerProxy HTTPSamplerProxy `xml:"HTTPSamplerProxy"`
	HashTree         *HashTree        `xml:"hashTree"`
}

type TestPlan struct {
	XMLName     xml.Name    `xml:"TestPlan"`
	GuiClass    string      `xml:"guiclass,attr"`
	TestClass   string      `xml:"testclass,attr"`
	TestName    string      `xml:"testname,attr"`
	ElementProp ElementProp `xml:"elementProp"`
}

type ElementProp struct {
	XMLName        xml.Name       `xml:"elementProp"`
	Name           string         `xml:"name,attr"`
	ElementType    string         `xml:"elementType,attr"`
	GuiClass       string         `xml:"guiclass,attr"`
	TestClass      string         `xml:"testclass,attr"`
	TestName       string         `xml:"testname,attr"`
	CollectionProp CollectionProp `xml:"collectionProp"`
}

type CollectionProp struct {
	XMLName   xml.Name   `xml:"collectionProp"`
	Name      string     `xml:"name,attr"`
	Arguments []Argument `xml:"elementProp"`
}

type Argument struct {
	XMLName        xml.Name       `xml:"elementProp"`
	Name           string         `xml:"name,attr"`
	ElementType    string         `xml:"elementType,attr"`
	CollectionProp CollectionProp `xml:"collectionProp"`
	StringProp     StringProp     `xml:"stringProp"`
}

type StringProp struct {
	XMLName xml.Name `xml:"stringProp"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

type ThreadGroup struct {
	XMLName     xml.Name     `xml:"ThreadGroup"`
	GuiClass    string       `xml:"guiclass,attr"`
	TestClass   string       `xml:"testclass,attr"`
	TestName    string       `xml:"testname,attr"`
	ElementProp ElementProp  `xml:"elementProp"`
	StringProps []StringProp `xml:"stringProp"`
	BoolProps   []BoolProp   `xml:"boolProp"`
}

type BoolProp struct {
	XMLName xml.Name `xml:"boolProp"`
	Name    string   `xml:"name,attr"`
	Value   bool     `xml:",chardata"`
}

type HTTPSamplerProxy struct {
	XMLName     xml.Name     `xml:"HTTPSamplerProxy"`
	GuiClass    string       `xml:"guiclass,attr"`
	TestClass   string       `xml:"testclass,attr"`
	TestName    string       `xml:"testname,attr"`
	BoolProps   []BoolProp   `xml:"boolProp"`
	ElementProp ElementProp  `xml:"elementProp"`
	StringProps []StringProp `xml:"stringProp"`
}

type HeaderManager struct {
	XMLName        xml.Name       `xml:"HeaderManager"`
	GuiClass       string         `xml:"guiclass,attr"`
	TestClass      string         `xml:"testclass,attr"`
	TestName       string         `xml:"testname,attr"`
	CollectionProp CollectionProp `xml:"collectionProp"`
}

type Header struct {
	XMLName     xml.Name     `xml:"elementProp"`
	Name        string       `xml:"name,attr"`
	ElementType string       `xml:"elementType,attr"`
	StringProps []StringProp `xml:"stringProp"`
}

type HTTPArgument struct {
	XMLName     xml.Name     `xml:"elementProp"`
	Name        string       `xml:"name,attr"`
	ElementType string       `xml:"elementType,attr"`
	StringProps []StringProp `xml:"stringProp"`
}
