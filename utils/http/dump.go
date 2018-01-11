package http

import (
    "bytes"
    "io"
    "io/ioutil"
    "net/http"
    "fmt"
)

//DumpRequestBody ...
func DumpRequestBody(req *http.Request) ([]byte, error) {
    if req == nil {
        return nil, fmt.Errorf("request is nil")
    }
    body, save, err := drainBody(req.Body)
    if err != nil {
        return nil, err
    }
    req.Body = save
    return body, err
}

//DumpRequestBody ...
func DumpResponseBody(resp *http.Response) ([]byte, error){
    if resp == nil {
        return nil, fmt.Errorf("reqsponse is nil")
    }
    body, save, err := drainBody(resp.Body)
    if err != nil {
        return nil, err
    }
    resp.Body = save
    return body, nil
}

func drainBody(b io.ReadCloser) (body []byte, r io.ReadCloser, err error) {
    if b == http.NoBody {
        // No copying needed. Preserve the magic sentinel meaning of NoBody.
        return nil, http.NoBody, nil
    }
    var buf bytes.Buffer
    if _, err = buf.ReadFrom(b); err != nil {
        return nil, b, err
    }
    if err = b.Close(); err != nil {
        return nil, b, err
    }
    return buf.Bytes(), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}