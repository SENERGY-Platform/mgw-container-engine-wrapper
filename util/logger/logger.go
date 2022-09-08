/*
 * Copyright 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var level = WarningLvl
var logger *log.Logger

func InitLogger(lvl Level, prefix string, microSec bool, useUtc bool) {
	level = lvl
	flags := log.Ldate | log.Ltime | log.Lmsgprefix
	if microSec {
		flags = flags | log.Lmicroseconds
	}
	if useUtc {
		flags = flags | log.LUTC
	}
	logger = log.New(os.Stderr, prefix, flags)
}

func ParseLevel(v string) (Level, error) {
	if v == debugStr {
		return Debug, nil
	}
	for i := 0; i < len(lvlStrings); i++ {
		if lvlStrings[i] == v {
			return Level(i), nil
		}
	}
	return defaultLvl, errors.New(fmt.Sprintf("unknown logging level '%s'", v))
}

func GetLevel() Level {
	return level
}

func genMessage(lvl string, val string) string {
	return fmt.Sprintf("%s %s", strings.ToUpper(lvl), val)
}

func logMessage(lvl Level, val string) {
	switch lvl {
	case PanicLvl:
		logger.Panicln(genMessage(lvl.String(), val))
	case FatalLvl:
		logger.Fatalln(genMessage(lvl.String(), val))
	default:
		if lvl <= level {
			logger.Println(genMessage(lvl.String(), val))
		}
	}
}

func Panic(v ...any) {
	logMessage(PanicLvl, fmt.Sprint(v...))
}

func PanicF(format string, v ...any) {
	logMessage(PanicLvl, fmt.Sprintf(format, v...))
}

func Fatal(v ...any) {
	logMessage(FatalLvl, fmt.Sprint(v...))
}

func Fatalf(format string, v ...any) {
	logMessage(FatalLvl, fmt.Sprintf(format, v...))
}

func Error(v ...any) {
	logMessage(ErrorLvl, fmt.Sprint(v...))
}

func Errorf(format string, v ...any) {
	logMessage(ErrorLvl, fmt.Sprintf(format, v...))
}

func Warning(v ...any) {
	logMessage(WarningLvl, fmt.Sprint(v...))
}

func WarningF(format string, v ...any) {
	logMessage(WarningLvl, fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	logMessage(InfoLvl, fmt.Sprint(v...))
}

func InfoF(format string, v ...any) {
	logMessage(InfoLvl, fmt.Sprintf(format, v...))
}

func DebugL1(v ...any) {
	logMessage(DebugL1Lvl, fmt.Sprint(v...))
}

func DebugL1F(format string, v ...any) {
	logMessage(DebugL1Lvl, fmt.Sprintf(format, v...))
}

func DebugL2(v ...any) {
	logMessage(DebugL2Lvl, fmt.Sprint(v...))
}

func DebugL2F(format string, v ...any) {
	logMessage(DebugL2Lvl, fmt.Sprintf(format, v...))
}

func DebugL3(v ...any) {
	logMessage(DebugL3Lvl, fmt.Sprint(v...))
}

func DebugL3F(format string, v ...any) {
	logMessage(DebugL3Lvl, fmt.Sprintf(format, v...))
}
