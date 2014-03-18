package workloads

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/nu7hatch/gouuid"
	. "github.com/pivotal-cf-experimental/cf-test-helpers/cf"
)

//Todo(simon) Remove, for dev testing only
func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	r := min + rand.Intn(max-min)
	return r
}

func Dummy() error {
	time.Sleep(time.Duration(random(1, 5)) * time.Second)
	return nil
}

func DummyWithErrors() error {
	Dummy()
	if random(0, 10) > 8 {
		return errors.New("Random (dummy) error")
	}
	return nil
}

func Push() error {
	guid, _ := uuid.NewV4()
	err := Cf("push", "pats-"+guid.String(), "patsapp", "-m", "64M", "-p", "assets/hello-world").ExpectOutput("App started")
	return err
}

func MoveAndSalt(srcDir string, dstDir string, salt string){
	os.Mkdir(dstDir, 0777)
	files, _ := ioutil.ReadDir(srcDir)	          
	for i := 0; i < len(files); i++ {
		fileInfo, _ := os.Stat(srcDir + files[i].Name())
		if fileInfo.Mode().IsDir() {
			newDstDir := dstDir + "/" + fileInfo.Name()
			MoveAndSalt(srcDir + fileInfo.Name()+"/", newDstDir+"/", salt)
		}
		if fileInfo.Mode().IsRegular() {
			input, _ := ioutil.ReadFile(srcDir + files[i].Name())
			output, _ := os.Create(dstDir+"/"+files[i].Name())
			output.Write(input)
			output.Write([]byte("\n#" + salt))
			output.Close()
		}        
	}
}  

func GenerateAndPush() error {
	srcDir := "assets/hello-world/"
	guid, _ := uuid.NewV4()
	rand.Seed(time.Now().UTC().UnixNano())
	salt := strconv.FormatInt(rand.Int63(), 10)
	dstDir := "assets/"+salt
	MoveAndSalt(srcDir, dstDir, salt)

	err := Cf("push", "pats-"+guid.String(), "patsapp", "-m", "64M", "-p", dstDir).ExpectOutput("App started")
	os.RemoveAll(dstDir)
	return err
}
