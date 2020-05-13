package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)


import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonB4ad3f7dDecodeCourseraGoHw3Struct(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonB4ad3f7dEncodeCourseraGoHw3Struct(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"browsers\":"
		out.RawString(prefix)
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB4ad3f7dEncodeCourseraGoHw3Struct(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB4ad3f7dEncodeCourseraGoHw3Struct(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB4ad3f7dDecodeCourseraGoHw3Struct(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB4ad3f7dDecodeCourseraGoHw3Struct(l, v)
}


type User struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Browsers []string `json:"browsers"`
}


// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	seenBrowsers := make(map[string]interface{})
	uniqueBrowsers := 0
	email := ""
	foundUsers := ""
	isAndroid := false
	isMSIE := false

	user := &User{}

	fmt.Fprintln(out, "found users:")

	reader := bufio.NewReader(file)
	for i := 0;;i++{
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}
		// fmt.Printf("%v %v\n", err, line)
		err = user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false
		//fmt.Println(user.Name, user.Email, user.Browsers)
		for _, browser := range user.Browsers {

			if strings.Contains(browser, "Android")  {
				isAndroid = true
				if _, ok := seenBrowsers[browser]; !ok{
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers[browser] = struct {}{}
					uniqueBrowsers++
				}
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok{
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers[browser] = struct {}{}
					uniqueBrowsers++
				}
			}
		}

		if (isAndroid && isMSIE) {
			// log.Println("Android and MSIE user:", user["name"], user["email"])
			email = strings.ReplaceAll(user.Email, "@", " [at] ")
			foundUsers = fmt.Sprintf("[%d] %s <%s>", i, user.Name, email)
			fmt.Fprintln(out, foundUsers)
		}
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Total unique browsers", uniqueBrowsers)
}
