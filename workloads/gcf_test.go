package workloads_test

import (
	"crypto/md5"
	"io/ioutil"
	"os"

	. "github.com/cloudfoundry-community/pat/workloads"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GCF Workloads", func() {
	var (
		salt              string
		srcDir            string
		dstDir            string
	)

	BeforeEach(func() {
		salt = "1234"
		srcDir = "../assets/hello-world/"
		dstDir = "../assets/"+salt + "/"
	})

	Describe("Generating and Pushing an app", func() {
		Context("MoveAndSalt", func() {
			It("Creates a new directory and adds a comment to each file to change the hash", func() {
				MoveAndSalt(srcDir, dstDir, salt)
				
				files, _ := ioutil.ReadDir(srcDir)
				for i := 0; i < len(files); i++ {
					fileInfo, _ := os.Stat(srcDir + files[i].Name())
					if fileInfo.Mode().IsRegular() {
						input, _ := ioutil.ReadFile(srcDir+files[i].Name())
						output, err := ioutil.ReadFile(dstDir+files[i].Name())
						oldHash := md5.Sum(input)
						newHash := md5.Sum(output)
				
						Ω(newHash).ShouldNot(Equal(oldHash))
						Ω(err).ShouldNot(HaveOccured())
					}
				}
				os.RemoveAll(dstDir)
			})
		})
	})
})
