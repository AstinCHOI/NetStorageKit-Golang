package netstorage_test


import (
    "fmt"
    "io/ioutil"
    "os"
    "testing"
    "time"
    "github.com/AstinCHOI/NetStorageKit-GoLang/akamai/netstorage"
    "github.com/stretchr/testify/suite"
    "./spike/secrets"
    // "./akamai/netstorage"
)

var NS_HOSTNAME string = "astin-nsu.akamaihd.net"
var NS_KEYNAME  string = "astinastin"
var NS_KEY string = secrets.KEY
var NS_CPCODE string = "360949"


type NetstorageTestSuite struct {
    suite.Suite
    ns *netstorage.Netstorage
    temp_ns_dir string
    temp_file string
    temp_ns_file string
}

func check(err error, exit bool) {
    if (err != nil) {
        if exit {
            panic(err)
        } else {
            fmt.Println(err)
        }
    }
}

func (suite *NetstorageTestSuite) SetupSuite() {
    suite.ns = netstorage.NewNetstorage(NS_HOSTNAME, NS_KEYNAME, NS_KEY, false)
    suite.temp_ns_dir = fmt.Sprintf("/%s/nst_%d", NS_CPCODE, time.Now().Unix())
    suite.temp_file = fmt.Sprintf("nst_%d.txt", time.Now().Unix())
    suite.temp_ns_file = fmt.Sprintf("%s/%s", suite.temp_ns_dir, suite.temp_file)
}

func (suite *NetstorageTestSuite) TearDownSuite() {
    // delete temp files for local
    if _, err := os.Stat(suite.temp_file); err == nil {
        err = os.Remove(suite.temp_file)
        check(err, false)
        fmt.Printf("[TEARDOWN] remove %s from local done\n", suite.temp_file)    
    }

    if _, err := os.Stat(suite.temp_file + "_rename"); err == nil {
        err = os.Remove(suite.temp_file + "_rename")
        check(err, false)
        fmt.Printf("[TEARDOWN] remove %s from local done\n", suite.temp_file + "_rename")    
    }

    // delete temp files for netstorage
    if res, _, err := suite.ns.Delete(suite.temp_ns_file); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] delete %s done\n", suite.temp_ns_file)
    }

    if res, _, err := suite.ns.Delete(suite.temp_ns_file + "_lnk"); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] delete %s done\n", suite.temp_ns_file + "_lnk")
    }

    if res, _, err := suite.ns.Delete(suite.temp_ns_file + "_rename"); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] delete %s done\n", suite.temp_ns_file + "_rename")
    }

    if res, _, err := suite.ns.Rmdir(suite.temp_ns_dir); res.StatusCode == 200 && err == nil {
        fmt.Printf("[TEARDOWN] rmdir %s done\n", suite.temp_ns_dir)
    }
}

func (suite *NetstorageTestSuite) TestNetstorage() {
    // Dir
    res, _, err := suite.ns.Dir("/" + NS_CPCODE)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[dir] StatusCode should be 200 OK")
    fmt.Printf("[TEST] dir /%s done\n", NS_CPCODE)
    
    
    // Mkdir
    res, _, err = suite.ns.Mkdir(suite.temp_ns_dir)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[mkdir] StatusCode should be 200 OK")
    fmt.Printf("[TEST] mkdir %s done\n", suite.temp_ns_dir)

    // Upload
    err = ioutil.WriteFile(suite.temp_file, []byte("Hello, Netstorage API World!"), 0666)
    check(err, true)
    res, _, err = suite.ns.Upload(suite.temp_file, suite.temp_ns_file)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[upload] StatusCode should be 200 OK")
    fmt.Printf("[TEST] upload %s to %s done\n", suite.temp_file, suite.temp_ns_file)

    // Du
    res, _, err = suite.ns.Du(suite.temp_ns_dir)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[du] StatusCode should be 200 OK")
    fmt.Printf("[TEST] du %s done\n", suite.temp_ns_dir)

    // Mtime
    current_time := time.Now().Unix()
    res, _, err = suite.ns.Mtime(suite.temp_ns_file, current_time)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[mtime] StatusCode should be 200 OK")
    fmt.Printf("[TEST] mtime %s done\n", suite.temp_ns_file)

    // Stat
    res, _, err = suite.ns.Stat(suite.temp_ns_file)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[stat] StatusCode should be 200 OK")
    fmt.Printf("[TEST] stat %s done\n", suite.temp_ns_file)

    // Symlink
    res, _, err = suite.ns.Symlink(suite.temp_ns_file, suite.temp_ns_file + "_lnk")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[symlink] StatusCode should be 200 OK")
    fmt.Printf("[TEST] symlink %s to %s done\n", suite.temp_ns_file, suite.temp_ns_file + "_lnk")

    // Rename
    res, _, err = suite.ns.Rename(suite.temp_ns_file, suite.temp_ns_file + "_rename")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[rename] StatusCode should be 200 OK")
    fmt.Printf("[TEST] rename %s to %s done\n", suite.temp_ns_file, suite.temp_ns_file + "_rename")

    // Download
    res, _, err = suite.ns.Download(suite.temp_ns_file + "_rename")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[download] StatusCode should be 200 OK")
    fmt.Printf("[TEST] download %s done\n", suite.temp_ns_file)

    // Delete
    res, _, err = suite.ns.Delete(suite.temp_ns_file + "_rename")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[delete] StatusCode should be 200 OK")
    fmt.Printf("[TEST] delete %s done\n", suite.temp_ns_file + "_rename")
    res, _, err = suite.ns.Delete(suite.temp_ns_file + "_lnk")
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[delete] StatusCode should be 200 OK")
    fmt.Printf("[TEST] delete %s done\n", suite.temp_ns_file + "_lnk")

    // Rmdir
    res, _, err = suite.ns.Rmdir(suite.temp_ns_dir)
    check(err, true)
    suite.Equal(res.StatusCode, 200, "[rmdir] StatusCode should be 200 OK")
    fmt.Printf("[TEST] rmdir %s done\n", suite.temp_ns_dir)
}


func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(NetstorageTestSuite))
}