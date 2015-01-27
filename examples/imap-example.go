// This file is brainless porting of examples/imap-example.go
// This code haven't anything common or related with IMAP, I haven't any idea
// why gmime developers used this name.

package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path"

    "github.com/sendgrid/go-gmime/gmime"
)


type EnvelopeT struct {
    Date  string `json:"date"`
    Subject string `json:"subject"`
    From string `json:"from"`
    Sender string `json:"sender"`
    ReplyTo string `json:"reply_to"`
    To string `json:"to"`
    Cc string `json:"cc"`
    Bcc string `json:"bcc"`
    InReplyTo string `json:"in_reply_to"`
    MessageId string `json:"message_id"`
}

type ParamT struct {
    Name string `json:"name"`
    Value string `json:"name"`
}

type ContentTypeT struct {
    Type string `json:"type"`
    SubType string `json:"type"`
    Params []ParamT `json:"params"`
}

type ContentDispositionT struct {
    ContentDisposition string `json:"content_disposition"`
    Params []ParamT `json:"params"`
}

type Bodystruct struct {
    ContentType ContentTypeT `json:"content_type"`
    ContentDisposition ContentDispositionT `json:"content_disposition"`
    Encoding string `json:"encoding"`
    Envelope EnvelopeT `json:"envelope"`
    Subparts []Bodystruct `json:"subparts"`
}

var file_name string
var scan_from bool

func init() {
    flag.StringVar(&file_name, "filename", "file.eml", "email file to parse")
    flag.BoolVar(&scan_from, "from", false, "scan from")
}

func main() {
    flag.Parse()
    fs := gmime.NewFileStreamForPath(file_name, "r")
    if fs == nil {
       panic("can't open " + file_name)
    }

    parser := gmime.NewParserWithStream(fs)
    parser.SetScanFrom(scan_from)
    message := parser.ConstructMessage()
    if message != nil {
        uid := message.MessageId()
        // FIXME: uid should be "Maybe uid" here
        if uid == "" {
            uid = path.Base(file_name)
        }
        os.Mkdir(uid, os.FileMode(0755))
        write_message(message, uid)
        reconstruct_message(uid)
    }
}

func write_message(m gmime.Message, uid string) {
    write_header(m, uid)
    write_bodystructure(m, uid)
    write_part(m.MimePart(), uid, "1")
}

func write_header(m gmime.Message, uid string) {
    fn := path.Join(uid, "HEADER")
    ioutil.WriteFile(fn, []byte(m.Headers()), os.FileMode(0644))
}

func write_part(ob gmime.Object, uid string, spec string) {
    {
        fn := path.Join(uid, fmt.Sprintf("%s.HEADER", spec))
        ioutil.WriteFile(fn, []byte(ob.Headers()), os.FileMode(0644))
    }

    if mp, ok := ob.(gmime.Multipart); ok {
        n := mp.Count()
        for i := 0; i < n; i++ {
            subpart := mp.GetPart(i)
            id := fmt.Sprintf("%s.%d", spec, i + 1)
            write_part(subpart, uid, id)
        }
    } else if mp, ok := ob.(gmime.MessagePart); ok {
        fn := path.Join(uid, fmt.Sprintf("%s.TEXT", spec))
        ostream := gmime.NewFileStreamForPath(fn, "wt")
        mp.Message().WriteToStream(ostream)
        defer ostream.Close()
    } else if p, ok := ob.(gmime.Part); ok {
        fn := path.Join(uid, fmt.Sprintf("%s.TEXT", spec))
        ostream := gmime.NewFileStreamForPath(fn, "wt")
        defer ostream.Close()
        dataWrapper := p.ContentObject()
        dataWrapper.WriteToStream(ostream)
    }
}

func write_bodystructure(m gmime.Message, uid string) {
    fn := path.Join(uid, "BODYSTRUCTURE")
    bs :=  serialize_part_bodystructure(m.MimePart())
    data, _ := json.Marshal(bs)
    ioutil.WriteFile(fn, data, os.FileMode(0644))
}

func serialize_part_bodystructure(ob gmime.Object) Bodystruct {
    ct := ob.ContentType()
    contentType := ContentTypeT{
            Type: ct.MediaType(),
            SubType: ct.MediaSubtype(),
            Params: serialize_params(ct),
        }
    bs := Bodystruct{
        ContentType: contentType,
    }

    if mp, ok := ob.(gmime.Multipart); ok {
        n := mp.Count()
        for i := 0; i < n; i++ {
            sp := mp.GetPart(i)
            ssp := serialize_part_bodystructure(sp)
            bs.Subparts = append(bs.Subparts, ssp)
        }
    } else if mp, ok := ob.(gmime.MessagePart); ok {
        inner := mp.Message()
        bs.Envelope = EnvelopeT{
            Date: inner.DateAsString(),
            Subject: inner.Header("Subject"),
            To: inner.Header("To"),
            From: inner.Header("From"),
            Sender: inner.Header("Sender"),
            ReplyTo: inner.Header("Reply-To"),
            InReplyTo: inner.Header("In-Reply-To"),
            Cc: inner.Header("Cc"),
            Bcc: inner.Header("Bcc"),
            MessageId: inner.MessageId(),
        }
        ssp := serialize_part_bodystructure(mp)
        bs.Subparts = append(bs.Subparts, ssp)
    } else if p, ok := ob.(gmime.Part); ok {
        cd := ob.ContentDisposition()
        bs.ContentDisposition = ContentDispositionT{
            ContentDisposition: cd.Disposition(),
            Params: serialize_params(cd),
        }
        bs.Encoding = p.ContentEncoding()
    }

    return bs
}


func serialize_params(p gmime.Parametrized) []ParamT {
    var params []ParamT
    p.ForEachParam(func (name string, value string) {
        pp := ParamT{
            Name: name,
            Value: value,
        }
        params = append(params, pp)
    })
    return params
}

func reconstruct_message(uid string) {
    fn := path.Join(uid, "HEADER")
    data, err := ioutil.ReadFile(fn)
    if err != nil {
       panic("Can't open " + fn)
    }
    
    stream := gmime.NewMemStreamWithBuffer(string(data))
    parser := gmime.NewParserWithStream(stream)
    parser.SetScanFrom(false)
    message := parser.ConstructMessage()

    mimePart := message.MimePart()
    if mp, ok := mimePart.(gmime.Multipart); ok {
        fn := path.Join(uid, "BODYSTRUCTURE")
        bs_data, err := ioutil.ReadFile(fn)
        if err != nil {
           panic("Can't open " + fn)
       }
        bs := Bodystruct{}
        json.Unmarshal(bs_data, &bs)
        reconstruct_multipart(mp, bs, uid, "1")
    } else if p, ok := message.(gmime.Part); ok {
        reconstruct_part_content(p, uid, "1", "8bit") // 8bit is bogus here
    }


    result := path.Join(uid, "MESSAGE")
    ostream := gmime.NewFileStreamForPath(result, "wt")
    message.WriteToStream(ostream)
}

func reconstruct_multipart(part gmime.Multipart, bs Bodystruct, uid string, spec string) {
    for i, sub := range bs.Subparts {
        subspec := fmt.Sprintf("%s.%d", spec, i)
        fmt.Println("reconstructing a %s/%s part (%s)", sub.ContentType.Type, sub.ContentType, subspec)
        
        fn := path.Join(uid, fmt.Sprintf("%s.HEADER", subspec))
        data, err := ioutil.ReadFile(fn)
        if err != nil {
           panic("Can't open " + fn)
        }

        stream := gmime.NewMemStreamWithBuffer(string(data))
        parser := gmime.NewParserWithStream(stream)
        parser.SetScanFrom(false)
        subpart := parser.ConstructPart()

        if mp, ok := subpart.(gmime.Multipart); ok {
            reconstruct_multipart(mp, sub, uid, subspec)
        } else if mp, ok := subpart.(gmime.MessagePart); ok {
            reconstruct_message_part(mp, uid, subspec)
        } else if p, ok := subpart.(gmime.Part); ok {
            reconstruct_part_content(p, uid, spec + ".1", bs.Encoding)
        }
        part.AddPart(subpart)
    }
}

func reconstruct_message_part(msgpart gmime.MessagePart, uid string, spec string) {
    fn := path.Join(uid, fmt.Sprintf("%s.TEXT", spec))
    data, err := ioutil.ReadFile(fn)
    if err != nil {
        panic("can't open " + fn)
    }
    stream := gmime.NewMemStreamWithBuffer(string(data))
    parser := gmime.NewParserWithStream(stream)
    parser.SetScanFrom(false)
    message := parser.ConstructMessage()
    msgpart.SetMessage(message)
}

func reconstruct_part_content(part gmime.Part, uid string, spec string, encoding string) {
    fn := path.Join(uid, fmt.Sprintf("%s.TEXT", spec))
    data, err := ioutil.ReadFile(fn)
    if err != nil {
        panic("can't open " + fn)
    }
    stream := gmime.NewMemStreamWithBuffer(string(data))
    content := gmime.NewDataWrapperWithStream(stream, encoding)
    part.SetContentObject(content)
}
