package store_test

import (
	"github.com/cloudfoundry-community/pat/config"
	"github.com/cloudfoundry-community/pat/laboratory"
	"github.com/cloudfoundry-community/pat/redis"
	. "github.com/cloudfoundry-community/pat/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config.WithStore", func() {
	var (
		csvStoreDir     string
		csvStore        laboratory.Store
		connFromFactory redis.Conn
		redisConn       redis.Conn
		redisStore      laboratory.Store
		flags           config.Config
		args            []string
	)

	BeforeEach(func() {
		flags = config.NewConfig()
		DescribeParameters(flags)
		args = []string{}

		csvStore = NewCsvStore("/tmp/fakecsvstore")
		redisStore = NewCsvStore("/tmp/fakeredisstore")
		CsvStoreFactory = func(dir string) laboratory.Store {
			csvStoreDir = dir
			return csvStore
		}

		connFromFactory = &dummyConn{}
		WithRedisConnection = func(fn func(conn redis.Conn) error) error {
			return fn(connFromFactory)
		}

		RedisStoreFactory = func(conn redis.Conn) (laboratory.Store, error) {
			redisConn = conn
			return redisStore, nil
		}
	})

	JustBeforeEach(func() {
		flags.Parse(args)
	})

	Context("When useRedis is false", func() {
		BeforeEach(func() {
			args = []string{"-use-redis=false", "-csv-dir", "foo/bar/baz"}
		})

		It("Uses the csvDir paramter to configure a CSV store", func() {
			var s laboratory.Store = nil
			WithStore(func(store laboratory.Store) error {
				s = store
				return nil
			})

			Ω(s).Should(Equal(csvStore))
			Ω(csvStoreDir).Should(Equal("foo/bar/baz"))
		})
	})

	Context("When useRedis is true", func() {
		BeforeEach(func() {
			args = []string{"-use-redis"}
		})

		It("Creates a Redis store using the connection from redis.WithRedisConnection(", func() {
			var s laboratory.Store = nil
			WithStore(func(store laboratory.Store) error {
				s = store
				return nil
			})

			Ω(s).Should(Equal(redisStore))
			Ω(redisConn).Should(Equal(connFromFactory))
		})
	})
})

type dummyConn struct{}

func (dummyConn) Do(cmd string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
