package main

import (
	"fmt"
	"strings"
)

var fmtWhole = `namespace NConfig
{
    public static class Error
    {
        public static readonly string[] Keys;

        public static string Desc(int errCode)
        {
            if (errCode >= 0 && errCode < Keys.Length)
            {
                return LanguageReader.Error[Keys[errCode]].CN;
            }
            return "ErrCode" + errCode.ToString();
        }

        static Error()
        {
            Keys = new string[]
            {{ErrorKey}
            };
        }
    }
}
`

var fmtErrorKey = `
                "%v",`

func tomlToCSharp(errReasonToml *errReasonToml) string {
	keys := ""

	for _, one := range errReasonToml.Error {
		keys += fmt.Sprintf(fmtErrorKey, one.ErrCode)
	}

	return strings.Replace(fmtWhole, "{ErrorKey}", keys, -1)
}
